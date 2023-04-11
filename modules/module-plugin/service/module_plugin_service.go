package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	"github.com/eolinker/apinto-dashboard/modules/core"
	"github.com/eolinker/apinto-dashboard/modules/group"
	group_service "github.com/eolinker/apinto-dashboard/modules/group/group-service"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/dto"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/entry"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/store"
	"github.com/eolinker/apinto-dashboard/modules/navigation"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/go-basic/uuid"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"os"
	"time"
)

var (
	ErrModulePluginNotFound    = errors.New("plugin doesn't exist. ")
	ErrModulePluginInstalled   = errors.New("plugin has installed. ")
	ErrModulePluginHasDisabled = errors.New("plugin has disabled. ")
)

type modulePluginService struct {
	pluginStore        store.IModulePluginStore
	pluginEnableStore  store.IModulePluginEnableStore
	pluginPackageStore store.IModulePluginPackageStore

	commonGroup       group.ICommonGroupService
	navigationService navigation.INavigationService
	coreService       core.ICore
	installedCache    IInstalledCache
	lockService       locker_service.IAsynLockService
}

func newModulePluginService() module_plugin.IModulePluginService {
	s := &modulePluginService{}
	bean.Autowired(&s.pluginStore)
	bean.Autowired(&s.pluginEnableStore)
	bean.Autowired(&s.pluginPackageStore)

	bean.Autowired(&s.commonGroup)
	bean.Autowired(&s.navigationService)
	bean.Autowired(&s.coreService)
	bean.Autowired(&s.installedCache)
	bean.Autowired(&s.lockService)
	return s
}

func (m *modulePluginService) GetPlugins(ctx context.Context, groupUUID, searchName string) ([]*model.ModulePluginItem, error) {
	groupID := -1
	if groupUUID != "" {
		groupInfo, err := m.commonGroup.GetGroupInfo(ctx, groupUUID)
		if err != nil {
			return nil, err
		}
		groupID = groupInfo.Id
	}
	pluginEntries, err := m.pluginStore.GetPluginList(ctx, groupID, searchName)
	if err != nil {
		return nil, err
	}
	plugins := make([]*model.ModulePluginItem, 0, len(pluginEntries))
	for _, pluginEntry := range pluginEntries {
		plugin := &model.ModulePluginItem{
			ModulePlugin: pluginEntry,
			IsEnable:     false,
			IsInner:      true,
		}
		//若为非内置
		if pluginEntry.Type == 2 {
			plugin.IsInner = false
		}
		enableEntry, err := m.pluginEnableStore.Get(ctx, pluginEntry.Id)
		if err != nil {
			//若启用表没有插件信息，则为未启用
			if err == gorm.ErrRecordNotFound {
				plugin.IsEnable = false
				plugins = append(plugins, plugin)
				continue
			}
			return nil, err
		}

		//若插件已启用
		if enableEntry.IsEnable == 2 {
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
		IsEnable:     false,
		Uninstall:    true,
	}
	//若为非内置
	if plugin.Type == 2 {
		info.Uninstall = false
	}

	enableEntry, err := m.pluginEnableStore.Get(ctx, plugin.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		if err == gorm.ErrRecordNotFound {
			return info, nil
		}
		return nil, err
	}

	//若插件已启用
	if enableEntry.IsEnable == 2 {
		info.IsEnable = true
	}
	return info, nil
}

func (m *modulePluginService) GetPluginGroups(ctx context.Context) ([]*model.PluginGroup, error) {
	groupEntries, err := m.commonGroup.GroupListAll(ctx, -1, group_service.ModulePlugin, group_service.ModulePlugin)
	if err != nil {
		return nil, err
	}
	groups := make([]*model.PluginGroup, 0, len(groupEntries))
	for _, entry := range groupEntries {
		groups = append(groups, &model.PluginGroup{
			UUID: entry.Uuid,
			Name: entry.Name,
		})
	}
	return groups, nil
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

	//通过导航id获取导航信息
	navigationUUID, _ := m.navigationService.GetUUIDByID(ctx, enableEntry.Navigation)

	enableCfg := new(model.PluginEnableCfg)
	_ = json.Unmarshal(enableEntry.Config, enableCfg)

	info := &model.PluginEnableInfo{
		Name:       enableEntry.Name,
		Navigation: navigationUUID,
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
		Internet:  false,
		Invisible: true,
	}
	switch pluginInfo.Driver {
	case "remote":
		remoteDefine := new(model.RemoteDefine)
		_ = json.Unmarshal(pluginInfo.Details, remoteDefine)
		if !remoteDefine.Internet {
			renderCfg.Internet = true
		}
		renderCfg.Querys = remoteDefine.Querys
		renderCfg.Initialize = remoteDefine.Initialize
	case "local":
		localDefine := new(model.LocalDefine)
		_ = json.Unmarshal(pluginInfo.Details, localDefine)
		renderCfg.Headers = localDefine.Headers
		renderCfg.Querys = localDefine.Querys
		renderCfg.Initialize = localDefine.Initialize
		renderCfg.Invisible = localDefine.Invisible
	}
	return renderCfg, nil
}

func (m *modulePluginService) InstallPlugin(ctx context.Context, userID int, groupName string, pluginYml *model.PluginYmlCfg, packageContent []byte) error {
	//通过插件id来判断插件是否已安装
	_, err := m.pluginStore.GetPluginInfo(ctx, pluginYml.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		return ErrModulePluginInstalled
	}

	//判断groupName存不存在，不存在则新建
	groupInfo, err := m.commonGroup.GetGroupByName(ctx, groupName, 0)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	groupID := 0
	return m.pluginStore.Transaction(ctx, func(txCtx context.Context) error {
		if err == gorm.ErrRecordNotFound {
			groupID, err = m.commonGroup.CreateGroup(txCtx, -1, userID, group_service.ModulePlugin, "", groupName, uuid.New(), "")
		} else {
			groupID = groupInfo.Id
		}

		t := time.Now()
		var details []byte
		switch pluginYml.Driver {
		case "remote":
			details, _ = json.Marshal(pluginYml.Remote)
		case "local":
			details, _ = json.Marshal(pluginYml.Local)
		case "profession":
			details, _ = json.Marshal(pluginYml.Profession)
		}

		pluginInfo := &entry.ModulePlugin{
			UUID:       pluginYml.ID,
			Name:       pluginYml.Name,
			Version:    pluginYml.Version,
			Group:      groupID,
			CName:      pluginYml.CName,
			Resume:     pluginYml.Resume,
			ICon:       pluginYml.ICon,
			Type:       2,
			Front:      "", //TODO
			Driver:     pluginYml.Driver,
			Details:    details,
			Operator:   userID,
			CreateTime: t,
		}
		if err = m.pluginStore.Save(txCtx, pluginInfo); err != nil {
			return err
		}

		return m.pluginPackageStore.Save(txCtx, &entry.ModulePluginPackage{
			Id:      pluginInfo.Id,
			Package: packageContent,
		})

	})
}

func (m *modulePluginService) EnablePlugin(ctx context.Context, userID int, pluginUUID string, enableInfo *dto.PluginEnableInfo) error {
	pluginInfo, err := m.pluginStore.GetPluginInfo(ctx, pluginUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrModulePluginNotFound
		}
		return err
	}

	err = m.lockService.Lock(locker_service.LockNameModulePlugin, pluginInfo.Id)
	if err != nil {
		return err
	}
	defer m.lockService.Unlock(locker_service.LockNameModulePlugin, pluginInfo.Id)

	navigationID, err := m.navigationService.GetIDByUUID(ctx, enableInfo.Navigation)
	if err != nil {
		return err
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

	checkConfig := enableCfg
	//若为内置插件
	if pluginInfo.Type == 0 || pluginInfo.Type == 1 {
		checkConfig = nil
	}

	err = m.pluginStore.Transaction(ctx, func(txCtx context.Context) error {
		enable := &entry.ModulePluginEnable{
			Id:         pluginInfo.Id,
			Name:       enableInfo.Name,
			Navigation: navigationID,
			IsEnable:   2,
			Config:     config,
			Operator:   userID,
			UpdateTime: time.Now(),
		}
		if err = m.pluginEnableStore.Save(txCtx, enable); err != nil {
			return err
		}

		return m.coreService.CheckNewModule(pluginInfo.UUID, enableInfo.Name, checkConfig)
	})
	if err != nil {
		return err
	}
	//TODO 重新生成路由
	err = m.coreService.ReloadModule()
	if err != nil {
		log.Error(err)
	}
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

	err = m.lockService.Lock(locker_service.LockNameModulePlugin, pluginInfo.Id)
	if err != nil {
		return err
	}
	defer m.lockService.Unlock(locker_service.LockNameModulePlugin, pluginInfo.Id)

	enableInfo, err := m.pluginEnableStore.Get(ctx, pluginInfo.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrModulePluginHasDisabled
		}
		return err
	}

	err = m.pluginEnableStore.Transaction(ctx, func(txCtx context.Context) error {
		enableInfo.IsEnable = 1
		enableInfo.Operator = userID
		enableInfo.UpdateTime = time.Now()
		_, err = m.pluginEnableStore.Update(txCtx, enableInfo)

		return err
	})
	if err != nil {
		return err
	}

	//TODO 重新生成路由
	err = m.coreService.ReloadModule()
	if err != nil {
		log.Error(err)
	}

	return nil
}

func (m *modulePluginService) GetEnablePluginsByNavigation(ctx context.Context, navigationID int) ([]*model.NavigationEnabledPlugin, error) {
	enabledPlugins, err := m.pluginEnableStore.GetListByNavigation(ctx, navigationID)
	if err != nil {
		return nil, err
	}
	plugins := make([]*model.NavigationEnabledPlugin, 0, len(enabledPlugins))
	for _, p := range enabledPlugins {
		pluginInfo, err := m.pluginStore.Get(ctx, p.Id)
		if err != nil {
			return nil, err
		}
		plugins = append(plugins, &model.NavigationEnabledPlugin{
			ModulePluginEnable: p,
			UUID:               pluginInfo.UUID,
		})
	}

	return plugins, nil
}

func (m *modulePluginService) InstallInnerPlugin(ctx context.Context, pluginYml *model.InnerPluginYmlCfg) error {
	//判断groupName存不存在，不存在则新建
	groupInfo, err := m.commonGroup.GetGroupByName(ctx, "内置插件", 0)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	groupID := 0
	if err == gorm.ErrRecordNotFound {
		groupID, err = m.commonGroup.CreateGroup(ctx, -1, 0, group_service.ModulePlugin, "", "内置插件", uuid.New(), "")
	} else {
		groupID = groupInfo.Id
	}

	navigationID := -1
	if pluginYml.Navigation != "" {
		navigationID, err = m.navigationService.GetIDByUUID(ctx, pluginYml.Navigation)
		if err != nil {
			return fmt.Errorf("navigation %s doesn't exist. ", pluginYml.Navigation)
		}
	}

	return m.pluginStore.Transaction(ctx, func(txCtx context.Context) error {

		t := time.Now()
		pluginInfo := &entry.ModulePlugin{
			UUID:       pluginYml.ID,
			Name:       pluginYml.Name,
			Version:    pluginYml.Version,
			Group:      groupID,
			CName:      pluginYml.CName,
			Resume:     pluginYml.Resume,
			ICon:       pluginYml.ICon,
			Type:       2,
			Front:      pluginYml.Front,
			Driver:     pluginYml.Driver,
			Details:    []byte{},
			Operator:   0,
			CreateTime: t,
		}
		if err = m.pluginStore.Save(txCtx, pluginInfo); err != nil {
			return err
		}
		isEnable := 1
		if pluginYml.Auto {
			isEnable = 2
		}
		enable := &entry.ModulePluginEnable{
			Id:         pluginInfo.Id,
			Name:       pluginYml.Name,
			Navigation: navigationID,
			IsEnable:   isEnable,
			Config:     []byte{},
			Operator:   0,
			UpdateTime: t,
		}

		return m.pluginEnableStore.Save(txCtx, enable)
	})
}

func (m *modulePluginService) CheckPluginInstalled(ctx context.Context, pluginID string) (bool, error) {
	isInstalled := false

	key := m.installedCache.Key(pluginID)
	value, err := m.installedCache.Get(ctx, key)
	if err != nil && err != redis.Nil {
		return false, err
	}

	var pluginInfo *entry.ModulePlugin
	//若redis存在值
	if err == nil {
		isInstalled = value.Installed
	} else {
		//若redis无缓存
		pluginInfo, err = m.pluginStore.GetPluginInfo(ctx, pluginID)
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
	//若未安装直接返回
	if !isInstalled {
		return false, nil
	}

	//若插件已安装，检查本地缓存是否存在
	dirPath := fmt.Sprintf("./plugin/%s", pluginID)
	// 检查目录是否存在, 若不存在，则从数据库读取数据并解压
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		packageEntry, err := m.pluginPackageStore.Get(ctx, pluginInfo.Id)
		if err != nil {
			return false, err
		}
		packageFile := bytes.NewReader(packageEntry.Package)
		err = common.DeCompress(packageFile, dirPath)
		if err != nil {
			//删除目录
			os.RemoveAll(dirPath)
			return false, err
		}
	}

	return true, nil
}
