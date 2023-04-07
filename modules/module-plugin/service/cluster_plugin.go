package service

import (
	"context"
	"errors"
	"github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/dto"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/store"
	"github.com/eolinker/eosc/common/bean"
)

var (
	modulePluginNotFound = errors.New("plugin doesn't exist. ")
)

type modulePluginService struct {
	pluginStore        store.IModulePluginStore
	pluginEnableStore  store.IModulePluginEnableStore
	pluginPackageStore store.IModulePluginPackageStore
	lockService        locker_service.IAsynLockService
}

func newModulePluginService() module_plugin.IModulePluginService {

	s := &modulePluginService{}
	bean.Autowired(&s.pluginStore)
	bean.Autowired(&s.pluginEnableStore)
	bean.Autowired(&s.pluginPackageStore)
	bean.Autowired(&s.lockService)
	return s
}

func (m *modulePluginService) GetPlugins(ctx context.Context, groupUUID, searchName string) ([]*model.ModulePluginItem, []*model.PluginGroup, error) {
	//TODO implement me
	panic("implement me")
}

func (m *modulePluginService) GetPluginInfo(ctx context.Context, pluginUUID string) (*model.ModulePluginInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (m *modulePluginService) GetPluginGroupsEnum(ctx context.Context) ([]*model.PluginGroup, error) {
	//TODO implement me
	panic("implement me")
}

func (m *modulePluginService) GetPluginEnableInfo(ctx context.Context, pluginUUID string) (*model.PluginEnableInfo, *model.PluginEnableRender, error) {
	//TODO implement me
	panic("implement me")
}

func (m *modulePluginService) InstallPlugin(ctx context.Context, groupName string, packageContent []byte) error {
	//TODO implement me
	panic("implement me")
}

func (m *modulePluginService) EnablePlugin(ctx context.Context, pluginUUID string, enableInfo *dto.PluginEnableInfo) error {
	//TODO implement me
	panic("implement me")
}

func (m *modulePluginService) DisablePlugin(ctx context.Context, pluginUUID string) error {
	//TODO implement me
	panic("implement me")
}
