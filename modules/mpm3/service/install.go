package service

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	audit_model "github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	locker_service "github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	"github.com/eolinker/apinto-dashboard/modules/mpm3"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/entry"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/model"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/store"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/eosc/common/bean"
	"gorm.io/gorm"
)

var (
	_ mpm3.IInstallService = (*InstallService)(nil)
)

type InstallService struct {
	installedCache       IInstalledCache
	pluginStore          store.IPluginStore
	syncLockService      locker_service.ISyncLockService
	asynLockService      locker_service.IAsynLockService
	pluginResourcesStore mpm3.IResourcesService
	pluginEnableStore    store.EnableStore
	moduleService        mpm3.IModuleService
	accessService        mpm3.IAccessService
	frontendService      mpm3.IFrontendService
}

func newInstallService() mpm3.IInstallService {
	s := &InstallService{installedCache: newIInstalledCache()}
	bean.Autowired(&s.pluginStore)
	bean.Autowired(&s.syncLockService)
	bean.Autowired(&s.asynLockService)
	bean.Autowired(&s.pluginResourcesStore)
	bean.Autowired(&s.pluginEnableStore)
	bean.Autowired(&s.moduleService)
	bean.Autowired(&s.accessService)
	bean.Autowired(&s.frontendService)

	return s
}

func (s *InstallService) CheckPluginInstalled(ctx context.Context, pluginID string) (bool, error) {
	isInstalled := false

	key := pluginID
	value, err := s.installedCache.Get(ctx, key)
	if err != nil && err != redis.Nil {
		return false, err
	}

	//若redis存在值
	if err == nil {
		isInstalled = value.Installed
	} else {
		//若redis无缓存
		_, err = s.pluginStore.GetPluginInfo(ctx, pluginID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return false, err
		} else if err == gorm.ErrRecordNotFound {
			isInstalled = false
		} else {
			isInstalled = true
		}
		//缓存
		_ = s.installedCache.Set(ctx, key, &model.PluginInstalledStatus{
			Installed: isInstalled,
		})
	}

	return isInstalled, nil
}

func (s *InstallService) Uninstall(ctx context.Context, pluginID string) error {
	//校验插件存在，且为非内置插件
	pluginInfo, err := s.pluginStore.GetPluginInfo(ctx, pluginID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("插件不存在")
		}
		return err
	}
	//全局异步锁
	err = s.asynLockService.Lock(locker_service.LockNameModulePlugin, 0)
	if err != nil {
		return errors.New("现在有人在操作,请稍后再试")
	}
	defer s.asynLockService.Unlock(locker_service.LockNameModulePlugin, 0)

	if pluginInfo.IsInner {
		return errors.New("该插件不可以卸载")
	}

	enableInfo, err := s.pluginEnableStore.Get(ctx, pluginInfo.Id)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		enableInfo = &entry.PluginEnable{
			Id:         pluginInfo.Id,
			IsEnable:   statusPluginDisable,
			Config:     nil,
			Operator:   0,
			UpdateTime: time.Time{},
		}

	}

	//校验插件启用状态
	if enableInfo.IsEnable == statusPluginEnable {
		return errors.New("插件启用中，不可以卸载")
	}

	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: pluginID,
		Name: pluginInfo.CName,
	})

	err = s.pluginStore.Transaction(ctx, func(txCtx context.Context) error {
		//从插件资源表，启用表，插件表中删除
		err := s.moduleService.Save(txCtx, pluginInfo.Id, nil)
		if err != nil {
			return err
		}
		err = s.accessService.Save(txCtx, pluginInfo.Id, nil)
		if err != nil {
			return err
		}
		err = s.pluginResourcesStore.Delete(txCtx, pluginInfo.Id)
		if err != nil {
			return err
		}
		_, err = s.pluginEnableStore.Delete(txCtx, pluginInfo.Id)
		if err != nil {
			return err
		}
		_, err = s.pluginStore.Delete(txCtx, pluginInfo.Id)
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return err
	}

	_ = s.installedCache.Set(ctx, pluginID, &model.PluginInstalledStatus{
		Installed: false,
	})
	return nil
}

func (s *InstallService) Install(ctx context.Context, cfg *pm3.PluginDefine, resource *model.PluginResources) error {
	return s.install(ctx, cfg, resource, false, false, true)
}

func (s *InstallService) InstallInner(ctx context.Context, cfg *pm3.PluginDefine, resource *model.PluginResources, auto, canDisable bool) error {
	return s.install(ctx, cfg, resource, true, auto, canDisable)
}

func (s *InstallService) ClearInnerNotExits(ctx context.Context, ids []string) error {

	return s.pluginStore.Transaction(ctx, func(ctx context.Context) error {
		deleteIds, err := s.pluginStore.DeleteNotIn(ctx, true, ids...)
		if err != nil {
			return err
		}
		if len(deleteIds) == 0 {
			return nil
		}
		_, err = s.pluginEnableStore.Delete(ctx, deleteIds...)
		if err != nil {
			return err
		}
		err = s.pluginResourcesStore.Delete(ctx, deleteIds...)
		if err != nil {
			return err
		}
		for _, pid := range deleteIds {
			err := s.moduleService.Save(ctx, pid, nil)
			if err != nil {
				return err
			}
			err = s.accessService.Save(ctx, pid, nil)
			if err != nil {
				return err
			}
			err = s.frontendService.Save(ctx, pid, nil)
			if err != nil {
				return err
			}
		}

		return nil

	})

}

func (s *InstallService) install(ctx context.Context, cfg *pm3.PluginDefine, resource *model.PluginResources, inner, auto, canDisable bool) error {

	userId := common.GetUser(ctx)
	driver, has := apinto_module.GetDriver(cfg.Driver)
	if !has {
		return apinto_module.ErrorDriverNotExist
	}
	modules, access, frontends, err := driver.Install(cfg)
	if err != nil {
		return err
	}
	plugin, err := s.pluginStore.GetPluginInfo(ctx, cfg.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if plugin != nil {
		if plugin.Name != cfg.Name {
			return apinto_module.ErrorModuleNameConflict
		}

		if plugin.Hash == resource.Hash() {
			return nil
		}

	} else {
		plugin = &entry.Plugin{
			Id:         0,
			UUID:       cfg.Id,
			Name:       cfg.Name,
			CreateTime: time.Now(),
		}
	}

	plugin.Version = cfg.Version
	plugin.Group = cfg.GroupId
	plugin.CName = cfg.Cname
	plugin.ICon = cfg.ICon
	plugin.Driver = cfg.Driver
	plugin.Resume = cfg.Resume
	plugin.IsCanDisable = canDisable
	plugin.IsInner = inner
	plugin.Details = cfg
	plugin.Operator = userId
	plugin.Hash = resource.Hash()
	plugin.UpdateTime = time.Now()

	return s.pluginStore.Transaction(ctx, func(txCtx context.Context) error {
		//从插件资源表，启用表，插件表中删除
		err := s.pluginStore.Save(ctx, plugin)
		if err != nil {
			return err
		}

		if auto {
			err := s.pluginEnableStore.Save(ctx, &entry.PluginEnable{
				Id:         plugin.Id,
				IsEnable:   2,
				Config:     nil,
				Operator:   userId,
				UpdateTime: time.Now(),
			})
			if err != nil {
				return err
			}
		}

		err = s.accessService.Save(ctx, plugin.Id, access)
		if err != nil {
			return err
		}

		err = s.frontendService.Save(ctx, plugin.Id, frontends)
		if err != nil {
			return err
		}

		err = s.moduleService.Save(ctx, plugin.Id, modules)
		if err != nil {
			return err
		}
		if resource != nil {

			return s.pluginResourcesStore.Save(txCtx, plugin.Id, cfg.Id, resource)
		}

		return nil

	})
}
