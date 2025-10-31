package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/eolinker/apinto-dashboard/controller"
	plugin_group "github.com/eolinker/apinto-dashboard/initialize/plugin-group"
	audit_model "github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-dashboard/modules/core"
	"github.com/eolinker/apinto-dashboard/modules/mpm3"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/entry"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/model"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/store"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-basic/uuid"
	"gorm.io/gorm"
	"time"
)

var (
	_ mpm3.IPluginService = (*PluginService)(nil)
)

type PluginService struct {
	pluginStore store.IPluginStore
	enableStore store.EnableStore
	coreService core.ICore

	frontendService mpm3.IFrontendService
	moduleService   mpm3.IModuleService
	accessService   mpm3.IAccessService
}

func NewPluginService() mpm3.IPluginService {
	p := &PluginService{}
	bean.Autowired(&p.pluginStore)
	bean.Autowired(&p.enableStore)
	bean.Autowired(&p.coreService)
	bean.Autowired(&p.frontendService)
	bean.Autowired(&p.moduleService)
	bean.Autowired(&p.accessService)
	return p
}

func (p *PluginService) GetPlugin(ctx context.Context, uuid string) (*model.PluginInfo, error) {
	entry, err := p.pluginStore.GetPluginInfo(ctx, uuid)
	if err != nil {
		return nil, err
	}
	info := &model.PluginInfo{
		Plugin: model.Plugin{
			UUID:    entry.UUID,
			Name:    entry.Name,
			CName:   entry.CName,
			Resume:  entry.Resume,
			ICon:    entry.ICon,
			Enable:  false,
			IsInner: entry.IsInner,
		},
		Uninstall:  !entry.IsInner,
		CanDisable: entry.IsCanDisable,
	}
	enableInfo, err := p.enableStore.Get(ctx, entry.Id)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return info, nil
	}
	info.Enable = enableInfo.IsEnable == statusPluginEnable
	// 已启用或是内置插件,不允许卸载
	if info.Enable || entry.IsInner {
		info.Uninstall = false
	} else {
		info.Uninstall = true
	}
	return info, nil

}

func (p *PluginService) GetEnabled(ctx context.Context) ([]*model.PluginConfig, error) {

	plugins, err := p.pluginStore.GetEnabledPlugins(ctx)
	if err != nil {
		return nil, err
	}
	cfs := make([]*model.PluginConfig, 0, len(plugins))
	for _, i := range plugins {
		cfs = append(cfs, &model.PluginConfig{
			UUID:   i.UUID,
			Driver: i.Driver,
			Define: i.Details,
			Config: i.Config,
		})
	}
	return cfs, nil
}

func (p *PluginService) GetEnableRender(ctx context.Context, pluginUUID string) (*model.PluginEnableRender, error) {
	pluginInfo, err := p.pluginStore.GetPluginInfo(ctx, pluginUUID)
	if err != nil {
		return nil, err
	}

	renderCfg := &model.PluginEnableRender{}

	if pluginInfo.Details != nil || pluginInfo.Details.Define != nil {
		bd, err := json.Marshal(pluginInfo.Details.Define)
		if err == nil {
			json.Unmarshal(bd, renderCfg)
		}
	}

	return renderCfg, nil

}

func (p *PluginService) GetPluginEnableInfo(ctx context.Context, pluginUUID string) (*model.PluginEnableCfg, error) {
	pluginInfo, err := p.pluginStore.GetPluginInfo(ctx, pluginUUID)
	if err != nil {
		return nil, err
	}
	enableEntry, err := p.enableStore.Get(ctx, pluginInfo.Id)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return new(model.PluginEnableCfg), nil
	}

	enableCfg := new(model.PluginEnableCfg)
	_ = json.Unmarshal(enableEntry.Config, enableCfg)

	return enableCfg, nil
}

func (p *PluginService) Search(ctx context.Context, groupID, searchName string) ([]*model.Plugin, error) {
	var pluginEntries []*entry.PluginListItem
	var err error
	//判断groupID是不是其它分组
	if groupID == pluginGroupOther {
		groupList := plugin_group.GetModulePluginGroups()
		groups := make([]string, 0, len(groupList))
		for _, group := range groupList {
			groups = append(groups, group.ID)
		}
		pluginEntries, err = p.pluginStore.SearchFromOther(ctx, groups, searchName)
	} else {
		pluginEntries, err = p.pluginStore.Search(ctx, groupID, searchName)
	}

	if err != nil {
		return nil, err
	}

	plugins := make([]*model.Plugin, 0, len(pluginEntries))
	for _, pluginEntry := range pluginEntries {

		plugin := &model.Plugin{
			UUID:    pluginEntry.UUID,
			Name:    pluginEntry.Name,
			CName:   pluginEntry.CName,
			Resume:  pluginEntry.Resume,
			ICon:    pluginEntry.ICon,
			Enable:  pluginEntry.IsEnable == statusPluginEnable,
			IsInner: pluginEntry.IsInner,
		}

		plugins = append(plugins, plugin)
	}
	return plugins, nil
}

func (p *PluginService) GetGroups(ctx context.Context) ([]*model.PluginGroup, error) {
	pluginEntries, err := p.pluginStore.Search(ctx, "", "")
	if err != nil {
		return nil, err
	}

	groupCountMap := make(map[string]int, 7) //记录分组名对应的插件数
	hasGroupPluginsCount := 0                //记录有具体分组名的插件数量
	total := 0
	for _, item := range pluginEntries {

		groupCountMap[item.Group]++
		total++

	}

	list := plugin_group.GetModulePluginGroups()
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

func (p *PluginService) EnablePlugin(ctx context.Context, userID int, pluginUUID string, enableInfo *model.PluginEnableCfg) error {
	pluginInfo, err := p.pluginStore.GetPluginInfo(ctx, pluginUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrModulePluginNotFound
		}
		return err
	}
	//
	////全局异步锁
	//err = p.asynLockService.Lock(locker_service.LockNameModulePlugin, 0)
	//if err != nil {
	//	return errors.New("现在有人在操作,请稍后再试")
	//}
	//defer m.asynLockService.Unlock(locker_service.LockNameModulePlugin, 0)

	//判断插件是否已启用
	enable, err := p.enableStore.Get(ctx, pluginInfo.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if enable != nil && enable.IsEnable == statusPluginEnable {
		return ErrModulePluginHasEnabled
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

		Header:     headers,
		Query:      querys,
		Initialize: initializes,
	}
	config, _ := json.Marshal(enableCfg)

	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:          pluginUUID,
		Name:          pluginInfo.CName,
		EnableOperate: 1,
	})

	err = p.pluginStore.Transaction(ctx, func(txCtx context.Context) error {
		enableEntry := &entry.PluginEnable{
			Id:         pluginInfo.Id,
			IsEnable:   statusPluginEnable,
			Config:     config,
			Operator:   userID,
			UpdateTime: time.Now(),
		}

		return p.enableStore.Save(txCtx, enableEntry)
	})
	if err != nil {
		return err
	}

	p.frontendService.Clean(ctx)

	p.moduleService.Clean(ctx)
	p.accessService.Clean(ctx)
	//重新生成路由
	p.coreService.ResetVersion(uuid.New())
	return nil
}

func (p *PluginService) DisablePlugin(ctx context.Context, userID int, pluginUUID string) error {

	pluginInfo, err := p.pluginStore.GetPluginInfo(ctx, pluginUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrModulePluginNotFound
		}
		return err
	}
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:          pluginUUID,
		Name:          pluginInfo.CName,
		EnableOperate: 2,
	})
	if !pluginInfo.IsCanDisable {
		return ErrModulePluginCantDisabled
	}

	enableInfo, err := p.enableStore.Get(ctx, pluginInfo.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrModulePluginHasDisabled
		}
		return err
	}

	err = p.pluginStore.Transaction(ctx, func(txCtx context.Context) error {
		enableInfo.IsEnable = statusPluginDisable
		enableInfo.Operator = userID
		enableInfo.UpdateTime = time.Now()

		return p.enableStore.Save(txCtx, enableInfo)
	})
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	p.frontendService.Clean(ctx)

	p.moduleService.Clean(ctx)
	p.accessService.Clean(ctx)
	//重新生成路由
	p.coreService.ResetVersion(uuid.New())
	return nil
}
