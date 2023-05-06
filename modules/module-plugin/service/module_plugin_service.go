package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/initialize"
	audit_model "github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	"github.com/eolinker/apinto-dashboard/modules/core"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/dto"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/entry"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/store"
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"os"
	"time"
)

type modulePluginService struct {
	pluginStore        store.IModulePluginStore
	pluginEnableStore  store.IModulePluginEnableStore
	pluginPackageStore store.IModulePluginPackageStore

	coreService            core.ICore
	installedCache         IInstalledCache
	navigationModulesCache module_plugin.INavigationModulesCache
	lockService            locker_service.IAsynLockService
}

func newModulePluginService() module_plugin.IModulePluginService {
	s := &modulePluginService{}
	bean.Autowired(&s.pluginStore)
	bean.Autowired(&s.pluginEnableStore)
	bean.Autowired(&s.pluginPackageStore)

	bean.Autowired(&s.coreService)
	bean.Autowired(&s.installedCache)
	bean.Autowired(&s.navigationModulesCache)
	bean.Autowired(&s.lockService)
	return s
}

func (m *modulePluginService) GetPlugins(ctx context.Context, groupID, searchName string) ([]*model.ModulePluginItem, error) {
	var pluginEntries []*entry.PluginListItem
	var err error
	//判断groupID是不是其它分组
	if groupID == pluginGroupOther {
		groupList := initialize.GetModulePluginGroups()
		groups := make([]string, 0, len(groupList))
		for _, group := range groupList {
			groups = append(groups, group.ID)
		}
		pluginEntries, err = m.pluginStore.GetOtherGroupPlugins(ctx, groups, searchName)
	} else {
		pluginEntries, err = m.pluginStore.GetPluginList(ctx, groupID, searchName)
	}

	if err != nil {
		return nil, err
	}
	plugins := make([]*model.ModulePluginItem, 0, len(pluginEntries))
	for _, pluginEntry := range pluginEntries {
		//框架插件不返回
		if pluginEntry.Type == pluginTypeFrame {
			continue
		}
		plugin := &model.ModulePluginItem{
			PluginListItem: pluginEntry,
			IsEnable:       false,
			IsInner:        true,
		}
		//若为非内置
		if !IsInnerPlugin(pluginEntry.Type) {
			plugin.IsInner = false
		}

		//若插件已启用
		if pluginEntry.IsEnable == statusPluginEnable {
			plugin.IsEnable = true
		}
		plugins = append(plugins, plugin)
	}
	return plugins, nil
}

func (m *modulePluginService) GetPluginInfo(ctx context.Context, pluginUUID string) (*model.ModulePluginInfo, error) {
	plugin, err := m.pluginStore.GetPluginInfo(ctx, pluginUUID)
	if err != nil {
		return nil, err
	}

	info := &model.ModulePluginInfo{
		ModulePlugin: plugin,
		Enable:       false,
		CanDisable:   false,
		Uninstall:    false,
	}

	enableEntry, err := m.pluginEnableStore.Get(ctx, plugin.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		if err == gorm.ErrRecordNotFound {
			return info, nil
		}
		return nil, err
	}
	info.CanDisable = enableEntry.IsCanDisable

	//若插件可卸载，且为停用，才能返回可卸载
	if enableEntry.IsEnable == statusPluginDisable && enableEntry.IsCanUninstall {
		info.Uninstall = true
	}

	//若插件已启用
	if enableEntry.IsEnable == statusPluginEnable {
		info.Enable = true
	}
	return info, nil
}

func (m *modulePluginService) GetPluginGroups(ctx context.Context) ([]*model.PluginGroup, error) {
	pluginEntries, err := m.pluginStore.GetPluginList(ctx, "", "")
	if err != nil {
		return nil, err
	}

	groupCountMap := make(map[string]int, 7) //记录分组名对应的插件数
	hasGroupPluginsCount := 0                //记录有具体分组名的插件数量
	total := 0
	for _, item := range pluginEntries {
		if item.Type != pluginTypeFrame {
			groupCountMap[item.Group]++
			total++
		}
	}

	list := initialize.GetModulePluginGroups()
	groups := make([]*model.PluginGroup, 0, len(list)+1)
	for _, item := range list {
		hasGroupPluginsCount = hasGroupPluginsCount + groupCountMap[item.ID]
		groups = append(groups, &model.PluginGroup{
			UUID:  item.ID,
			Name:  item.Name,
			Count: groupCountMap[item.ID],
		})
	}
	groups = append(groups, &model.PluginGroup{
		UUID:  pluginGroupOther,
		Name:  "其它",
		Count: total - hasGroupPluginsCount,
	})
	return groups, nil
}

func (m *modulePluginService) GetInnerPluginList(ctx context.Context) ([]*model.ModulePluginInfo, error) {
	pluginList, err := m.pluginStore.GetInnerPluginList(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*model.ModulePluginInfo, 0, len(pluginList))
	for _, item := range pluginList {
		items = append(items, &model.ModulePluginInfo{
			ModulePlugin: item,
		})
	}
	return items, nil
}

func (m *modulePluginService) GetPluginEnableInfo(ctx context.Context, pluginUUID string) (*model.PluginEnableInfo, error) {
	pluginInfo, err := m.pluginStore.GetPluginInfo(ctx, pluginUUID)
	if err != nil {
		return nil, err
	}
	enableEntry, err := m.pluginEnableStore.Get(ctx, pluginInfo.Id)
	if err != nil {
		return nil, err
	}

	enableCfg := new(model.PluginEnableCfg)
	_ = json.Unmarshal(enableEntry.Config, enableCfg)

	info := &model.PluginEnableInfo{
		Name:       enableEntry.Name,
		Server:     enableCfg.Server,
		Header:     enableCfg.Header,
		Query:      enableCfg.Query,
		Initialize: enableCfg.Initialize,
	}

	return info, nil
}

func (m *modulePluginService) GetPluginEnableRender(ctx context.Context, pluginUUID string) (*model.PluginEnableRender, error) {
	pluginInfo, err := m.pluginStore.GetPluginInfo(ctx, pluginUUID)
	if err != nil {
		return nil, err
	}

	renderCfg := &model.PluginEnableRender{
		//Invisible: true,
		Internet:     false,
		NameConflict: false,
	}

	enableEntry, err := m.pluginEnableStore.Get(ctx, pluginInfo.Id)
	if err != nil {
		return nil, err
	}
	//检测已启用的插件中是否有同名的
	enabledPlugin, err := m.pluginEnableStore.GetEnabledPluginByName(ctx, enableEntry.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if enabledPlugin != nil && enabledPlugin.Id != pluginInfo.Id {
		renderCfg.NameConflict = true
	}

	pluginDefine := new(model.PluginDefine)
	_ = json.Unmarshal(pluginInfo.Details, pluginDefine)

	switch pluginInfo.Driver {
	case pluginDriverRemote:
	case pluginDriverLocal:
		renderCfg.Headers = pluginDefine.Headers
		//renderCfg.Invisible = localDefine.Invisible
	}
	if enableEntry.IsShowServer {
		renderCfg.Internet = true
	}
	renderCfg.Querys = pluginDefine.Querys
	renderCfg.Initialize = pluginDefine.Initialize

	return renderCfg, nil
}

func (m *modulePluginService) InstallPlugin(ctx context.Context, userID int, pluginYml *model.PluginYmlCfg, packageContent []byte) error {
	//通过插件id来判断插件是否已安装
	_, err := m.pluginStore.GetPluginInfo(ctx, pluginYml.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == nil {
		return ErrModulePluginInstalled
	}

	driver, has := apinto_module.GetDriver(pluginYml.Driver)
	if !has {
		return ErrModulePluginDriverNotFound
	}
	plugin, err := driver.CreatePlugin(pluginYml.Define)
	if err != nil {
		return fmt.Errorf("创建插件 %s 失败:%s", pluginYml.Name, err.Error())
	}

	//全局异步锁
	err = m.lockService.Lock(locker_service.LockNameModulePlugin, 0)
	if err != nil {
		return errors.New("现在有人在操作,请稍后再试")
	}
	defer m.lockService.Unlock(locker_service.LockNameModulePlugin, 0)

	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: pluginYml.ID,
		Name: pluginYml.CName,
	})

	err = m.pluginStore.Transaction(ctx, func(txCtx context.Context) error {
		t := time.Now()
		details, _ := json.Marshal(pluginYml.Define)

		pluginInfo := &entry.ModulePlugin{
			UUID:       pluginYml.ID,
			Name:       pluginYml.Name,
			Version:    pluginYml.Version,
			Group:      pluginYml.GroupID,
			Navigation: pluginYml.Navigation,
			CName:      pluginYml.CName,
			Resume:     pluginYml.Resume,
			ICon:       pluginYml.ICon,
			Type:       pluginTypeNotInner,
			Driver:     pluginYml.Driver,
			Details:    details,
			Operator:   userID,
			CreateTime: t,
			UpdateTime: t,
		}
		if err = m.pluginStore.Save(txCtx, pluginInfo); err != nil {
			return err
		}
		enableInfo := &entry.ModulePluginEnable{
			Id:              pluginInfo.Id,
			Name:            pluginYml.Name,
			Navigation:      pluginYml.Navigation,
			IsEnable:        statusPluginDisable,
			IsCanDisable:    plugin.IsCanDisable(),
			IsCanUninstall:  plugin.IsCanUninstall(),
			IsShowServer:    plugin.IsShowServer(),
			IsPluginVisible: plugin.IsPluginVisible(),
			Config:          []byte{},
			Operator:        userID,
			UpdateTime:      t,
		}

		if err = m.pluginEnableStore.Save(txCtx, enableInfo); err != nil {
			return err
		}

		return m.pluginPackageStore.Save(txCtx, &entry.ModulePluginPackage{
			Id:      pluginInfo.Id,
			Package: packageContent,
		})

	})
	if err != nil {
		return err
	}
	//缓存
	_ = m.installedCache.Set(ctx, pluginYml.ID, &model.PluginInstalledStatus{
		Installed: true,
	}, time.Hour)

	return nil
}

func (m *modulePluginService) UninstallPlugin(ctx context.Context, userID int, pluginID string) error {
	//校验插件存在，且为非内置插件
	pluginInfo, err := m.pluginStore.GetPluginInfo(ctx, pluginID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("插件不存在")
		}
		return err
	}

	//全局异步锁
	err = m.lockService.Lock(locker_service.LockNameModulePlugin, 0)
	if err != nil {
		return errors.New("现在有人在操作,请稍后再试")
	}
	defer m.lockService.Unlock(locker_service.LockNameModulePlugin, 0)

	enableInfo, err := m.pluginEnableStore.Get(ctx, pluginInfo.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrModulePluginHasDisabled
		}
		return err
	}

	if !enableInfo.IsCanUninstall {
		return errors.New("内置插件不可以卸载")
	}

	//校验插件启用状态
	if enableInfo.IsEnable == statusPluginEnable {
		return errors.New("插件启用中，不可以卸载")
	}

	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: pluginID,
		Name: pluginInfo.CName,
	})

	err = m.pluginStore.Transaction(ctx, func(txCtx context.Context) error {
		//从package表，启用表，插件表中删除
		_, err = m.pluginPackageStore.Delete(txCtx, pluginInfo.Id)
		if err != nil {
			return err
		}
		_, err = m.pluginEnableStore.Delete(txCtx, pluginInfo.Id)
		if err != nil {
			return err
		}
		_, err = m.pluginStore.Delete(txCtx, pluginInfo.Id)
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return err
	}

	_ = m.installedCache.Set(ctx, pluginID, &model.PluginInstalledStatus{
		Installed: false,
	}, time.Hour)

	return nil
}

func (m *modulePluginService) EnablePlugin(ctx context.Context, userID int, pluginUUID string, enableInfo *dto.PluginEnableInfo) error {
	pluginInfo, err := m.pluginStore.GetPluginInfo(ctx, pluginUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrModulePluginNotFound
		}
		return err
	}

	//全局异步锁
	err = m.lockService.Lock(locker_service.LockNameModulePlugin, 0)
	if err != nil {
		return errors.New("现在有人在操作,请稍后再试")
	}
	defer m.lockService.Unlock(locker_service.LockNameModulePlugin, 0)

	//判断插件是否已启用
	enable, err := m.pluginEnableStore.Get(ctx, pluginInfo.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if enable != nil && enable.IsEnable == statusPluginEnable {
		return ErrModulePluginHasEnabled
	}

	//若输入的启用模块名为空，则为默认的模块名
	if enableInfo.Name == "" {
		enableInfo.Name = pluginInfo.Name
	}

	//校验模块名, 若启用的模块有同名的则报错
	info, err := m.pluginEnableStore.GetEnabledPluginByName(ctx, enableInfo.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if info != nil && info.Id != pluginInfo.Id {
		return errors.New("已有同名的插件")
	}

	headers := make([]*model.ExtendParams, 0, len(enableInfo.Header))
	querys := make([]*model.ExtendParams, 0, len(enableInfo.Query))
	initializes := make([]*model.ExtendParams, 0, len(enableInfo.Initialize))
	for _, h := range enableInfo.Header {
		headers = append(headers, &model.ExtendParams{
			Name:  h.Name,
			Value: h.Value,
		})
	}
	for _, q := range enableInfo.Query {
		querys = append(querys, &model.ExtendParams{
			Name:  q.Name,
			Value: q.Value,
		})
	}
	for _, i := range enableInfo.Initialize {
		initializes = append(initializes, &model.ExtendParams{
			Name:  i.Name,
			Value: i.Value,
		})
	}
	enableCfg := &model.PluginEnableCfg{
		Server:     enableInfo.Server,
		Header:     headers,
		Query:      querys,
		Initialize: initializes,
	}
	config, _ := json.Marshal(enableCfg)
	checkConfig := enabledCfgListToMap(enableCfg)

	define := pluginInfo.Details

	err = m.coreService.CheckNewModule(pluginInfo.UUID, enableInfo.Name, pluginInfo.Driver, define, checkConfig)
	if err != nil {
		return err
	}

	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:          pluginUUID,
		Name:          pluginInfo.CName,
		EnableOperate: 1,
	})

	driver, has := apinto_module.GetDriver(pluginInfo.Driver)
	if !has {
		return ErrModulePluginDriverNotFound
	}
	plugin, err := driver.CreatePlugin(define)
	if err != nil {
		return fmt.Errorf("创建插件驱动 %s 失败:%s", pluginInfo.Name, err.Error())
	}

	err = m.pluginStore.Transaction(ctx, func(txCtx context.Context) error {
		enableEntry := &entry.ModulePluginEnable{
			Id:              pluginInfo.Id,
			Name:            enableInfo.Name,
			Navigation:      pluginInfo.Navigation,
			IsEnable:        statusPluginEnable,
			IsCanDisable:    plugin.IsCanDisable(),
			IsCanUninstall:  plugin.IsCanUninstall(),
			IsShowServer:    plugin.IsShowServer(),
			IsPluginVisible: plugin.IsPluginVisible(),
			Frontend:        plugin.GetPluginFrontend(enableInfo.Name),
			Config:          config,
			Operator:        userID,
			UpdateTime:      time.Now(),
		}

		return m.pluginEnableStore.Save(txCtx, enableEntry)
	})
	if err != nil {
		return err
	}

	//重新生成路由
	m.coreService.ResetVersion("")

	//缓存已启用的模块列表
	enableModules, err := m.pluginStore.GetNavigationModules(ctx)
	if err != nil {
		log.Errorf(err.Error())
		return nil
	}

	_ = m.navigationModulesCache.SetAll(ctx, enableModules, 10*time.Minute)

	return nil
}

func (m *modulePluginService) DisablePlugin(ctx context.Context, userID int, pluginUUID string) error {
	pluginInfo, err := m.pluginStore.GetPluginInfo(ctx, pluginUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrModulePluginNotFound
		}
		return err
	}

	err = m.lockService.Lock(locker_service.LockNameModulePlugin, 0)
	if err != nil {
		return err
	}
	defer m.lockService.Unlock(locker_service.LockNameModulePlugin, 0)

	enableInfo, err := m.pluginEnableStore.Get(ctx, pluginInfo.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrModulePluginHasDisabled
		}
		return err
	}

	if !enableInfo.IsCanDisable {
		return ErrModulePluginCantDisabled
	}

	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:          pluginUUID,
		Name:          pluginInfo.CName,
		EnableOperate: 2,
	})
	err = m.pluginEnableStore.Transaction(ctx, func(txCtx context.Context) error {
		enableInfo.IsEnable = statusPluginDisable
		enableInfo.Operator = userID
		enableInfo.UpdateTime = time.Now()
		_, err = m.pluginEnableStore.Update(txCtx, enableInfo)

		return err
	})
	if err != nil {
		return err
	}

	//重新生成路由
	m.coreService.ResetVersion("")

	//缓存已启用的模块列表
	enableModules, err := m.pluginStore.GetNavigationModules(ctx)
	if err != nil {
		log.Errorf(err.Error())
		return nil
	}

	_ = m.navigationModulesCache.SetAll(ctx, enableModules, 10*time.Minute)

	return nil
}

func (m *modulePluginService) InstallInnerPlugin(ctx context.Context, pluginYml *model.InnerPluginYmlCfg) error {
	driver, has := apinto_module.GetDriver(pluginYml.Driver)
	if !has {
		panic(fmt.Errorf("not find driver:%s", pluginYml.Driver))
	}

	plugin, err := driver.CreatePlugin(pluginYml.Define)
	if err != nil {
		panic(fmt.Errorf("create plugin %s error:%s", pluginYml.Name, err.Error()))
	}

	err = m.pluginStore.Transaction(ctx, func(txCtx context.Context) error {
		t := time.Now()

		details, _ := json.Marshal(pluginYml.Define)

		pluginInfo := &entry.ModulePlugin{
			UUID:       pluginYml.ID,
			Name:       pluginYml.Name,
			Version:    pluginYml.Version,
			Group:      pluginYml.GroupID,
			Navigation: pluginYml.Navigation,
			CName:      pluginYml.CName,
			Resume:     pluginYml.Resume,
			ICon:       pluginYml.ICon,
			Type:       pluginYml.Type,
			Driver:     pluginYml.Driver,
			Details:    details,
			Operator:   0,
			CreateTime: t,
			UpdateTime: t,
		}
		if err := m.pluginStore.Save(txCtx, pluginInfo); err != nil {
			return err
		}
		isEnable := statusPluginDisable
		if pluginYml.Auto {
			isEnable = statusPluginEnable
		}
		enable := &entry.ModulePluginEnable{
			Id:              pluginInfo.Id,
			Name:            pluginYml.Name,
			Navigation:      pluginYml.Navigation,
			IsEnable:        isEnable,
			IsCanDisable:    plugin.IsCanDisable(),
			IsCanUninstall:  plugin.IsCanUninstall(),
			IsShowServer:    plugin.IsShowServer(),
			IsPluginVisible: plugin.IsPluginVisible(),
			Frontend:        plugin.GetPluginFrontend(pluginYml.Name),
			Config:          []byte{},
			Operator:        0,
			UpdateTime:      t,
		}

		return m.pluginEnableStore.Save(txCtx, enable)
	})
	if err != nil {
		return err
	}
	//缓存
	_ = m.installedCache.Set(ctx, pluginYml.ID, &model.PluginInstalledStatus{
		Installed: true,
	}, time.Hour)

	return nil
}

func (m *modulePluginService) UpdateInnerPlugin(ctx context.Context, pluginYml *model.InnerPluginYmlCfg) error {

	pluginInfo, err := m.pluginStore.GetPluginInfo(ctx, pluginYml.ID)
	if err != nil {
		return err
	}
	t := time.Now()

	pluginInfo.Name = pluginYml.Name
	pluginInfo.Version = pluginYml.Version
	pluginInfo.Navigation = pluginYml.Navigation
	pluginInfo.CName = pluginYml.CName
	pluginInfo.Resume = pluginYml.Resume
	pluginInfo.ICon = pluginYml.ICon
	pluginInfo.Type = pluginYml.Type
	pluginInfo.Driver = pluginYml.Driver
	pluginInfo.UpdateTime = t

	details, _ := json.Marshal(pluginYml.Define)
	pluginInfo.Details = details

	driver, has := apinto_module.GetDriver(pluginYml.Driver)
	if !has {
		panic(fmt.Errorf("not find driver:%s", pluginYml.Driver))
	}
	plugin, err := driver.CreatePlugin(pluginYml.Define)
	if err != nil {
		panic(fmt.Errorf("create plugin %s error:%s", pluginYml.Name, err.Error()))
	}

	return m.pluginStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = m.pluginStore.Save(txCtx, pluginInfo); err != nil {
			return err
		}
		enableInfo, err := m.pluginEnableStore.Get(txCtx, pluginInfo.Id)
		if err != nil {
			return err
		}
		//name和enable不更新
		enableInfo.Navigation = pluginYml.Navigation
		enableInfo.IsCanDisable = plugin.IsCanDisable()
		enableInfo.IsCanUninstall = plugin.IsCanUninstall()
		enableInfo.IsShowServer = plugin.IsShowServer()
		enableInfo.IsPluginVisible = plugin.IsPluginVisible()
		enableInfo.Frontend = plugin.GetPluginFrontend(enableInfo.Name)
		enableInfo.Operator = 0
		enableInfo.UpdateTime = t

		return m.pluginEnableStore.Save(txCtx, enableInfo)

	})
}

func (m *modulePluginService) CheckPluginInstalled(ctx context.Context, pluginID string) (bool, error) {
	isInstalled := false

	key := pluginID
	value, err := m.installedCache.Get(ctx, key)
	if err != nil && err != redis.Nil {
		return false, err
	}

	//若redis存在值
	if err == nil {
		isInstalled = value.Installed
	} else {
		//若redis无缓存
		_, err = m.pluginStore.GetPluginInfo(ctx, pluginID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return false, err
		} else if err == gorm.ErrRecordNotFound {
			isInstalled = false
		} else {
			isInstalled = true
		}
		//缓存
		m.installedCache.Set(ctx, key, &model.PluginInstalledStatus{
			Installed: isInstalled,
		}, time.Hour)
	}

	return isInstalled, nil
}

func (m *modulePluginService) CheckPluginISDeCompress(ctx context.Context, pluginDir string, pluginID string) error {
	pluginInfo, err := m.pluginStore.GetPluginInfo(ctx, pluginID)
	if err != nil {
		return err
	}
	//若插件已安装且为非内置插件，检查本地缓存是否存在
	if !IsInnerPlugin(pluginInfo.Type) {
		dirPath := fmt.Sprintf("%s/%s", pluginDir, pluginID)
		// 检查目录是否存在, 若不存在，则从数据库读取数据并解压
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			packageEntry, err := m.pluginPackageStore.Get(ctx, pluginInfo.Id)
			if err != nil {
				return err
			}
			//创建目录
			err = os.MkdirAll(dirPath, os.ModePerm)
			if err != nil {
				log.Error("安装插件失败, 无法创建目录:", err)
				return err
			}

			err = common.UnzipFromBytes(packageEntry.Package, dirPath)
			if err != nil {
				//删除目录
				os.RemoveAll(dirPath)
				return err
			}
		}
	}
	return nil
}
