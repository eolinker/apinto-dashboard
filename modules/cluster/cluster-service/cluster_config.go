package cluster_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-entry"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-store"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"gorm.io/gorm"
	"time"
)

type clusterConfigService struct {
	configStore     cluster_store.IClusterConfigStore
	runtimeStore    cluster_store.IClusterConfigRuntimeStore
	clusterService  cluster.IClusterService
	userInfoService user.IUserInfoService
	apintoClient    cluster.IApintoClient
	clConfigManager cluster.ICLConfigDriverManager
}

func newClusterConfigService() cluster.IClusterConfigService {
	s := &clusterConfigService{}

	bean.Autowired(&s.configStore)
	bean.Autowired(&s.runtimeStore)
	bean.Autowired(&s.clusterService)
	bean.Autowired(&s.userInfoService)
	bean.Autowired(&s.apintoClient)
	bean.Autowired(&s.clConfigManager)
	return s
}

func (c *clusterConfigService) Get(ctx context.Context, namespaceId int, clusterName, configType string) (interface{}, error) {
	//获取当前集群信息
	clusterInfo, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, err
	}
	info, err := c.configStore.GetConfigByTypeByCluster(ctx, clusterInfo.Id, configType)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	userInfo, err := c.userInfoService.GetUserInfo(ctx, info.Operator)
	if err != nil {
		return nil, err
	}

	return c.FormatOutput(configType, userInfo.UserName, info), nil
}

func (c *clusterConfigService) Edit(ctx context.Context, namespaceId, operator int, clusterName, configType string, config []byte) error {
	//获取当前集群信息
	clusterInfo, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	info, err := c.configStore.GetConfigByTypeByCluster(ctx, clusterInfo.Id, configType)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	t := time.Now()
	if err == gorm.ErrRecordNotFound {
		info = &cluster_entry.ClusterConfig{
			NamespaceId: namespaceId,
			ClusterId:   clusterInfo.Id,
			Type:        configType,
			IsEnable:    true, //新建默认是启用状态
			Data:        config,
			Operator:    operator,
			CreateTime:  t,
			UpdateTime:  t,
		}
	} else {
		info.Operator = operator
		info.UpdateTime = t
		info.Data = config
	}

	//获取apinto client
	client, err := c.apintoClient.GetClient(ctx, clusterInfo.Id)
	if err != nil {
		return err
	}

	//事务
	return c.configStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = c.configStore.Save(txCtx, info); err != nil {
			return err
		}
		//获取当前运行的版本
		runtime, err := c.runtimeStore.GetForCluster(ctx, info.Id, clusterInfo.Id)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		//若runtime为空
		if err == gorm.ErrRecordNotFound {
			runtime = &cluster_entry.ClusterConfigRuntime{
				NamespaceId: namespaceId,
				ConfigID:    info.Id,
				ClusterId:   clusterInfo.Id,
				IsOnline:    info.IsEnable,
				Operator:    operator,
				CreateTime:  t,
				UpdateTime:  t,
			}
		} else {
			runtime.IsOnline = info.IsEnable
			runtime.Operator = operator
			runtime.UpdateTime = t
		}

		if err = c.runtimeStore.Save(txCtx, runtime); err != nil {
			return err
		}

		//初始化配置
		if err = c.initConfig(configType, config); err != nil {
			return err
		}

		//启用状态则直接发布，非启用状态则先发布temp，再删除
		if info.IsEnable {
			if err = c.ToApinto(client, configType, configType, config); err != nil {
				return err
			}
		} else {
			tempName := fmt.Sprintf("%s_temp", configType)

			if err = c.ToApinto(client, tempName, configType, config); err != nil {
				return err
			}

			if err = c.OfflineApinto(client, tempName, configType); err != nil {
				return fmt.Errorf("offline %s from Apinto fail. err: %s", tempName, err)
			}
		}

		return nil
	})
}

func (c *clusterConfigService) Enable(ctx context.Context, namespaceId, operator int, clusterName, configType string) error {
	//获取当前集群信息
	clusterInfo, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}
	info, err := c.configStore.GetConfigByTypeByCluster(ctx, clusterInfo.Id, configType)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("%s config doesn't exist. ", configType)
		}
		return err
	}

	//获取apinto client
	client, err := c.apintoClient.GetClient(ctx, clusterInfo.Id)
	if err != nil {
		return err
	}

	//获取当前运行的版本
	runtime, err := c.runtimeStore.GetForCluster(ctx, info.Id, clusterInfo.Id)
	if err != nil {
		return err
	}

	t := time.Now()

	runtime.IsOnline = true
	runtime.Operator = operator
	runtime.UpdateTime = t

	info.IsEnable = true
	info.Operator = operator
	info.UpdateTime = t
	//事务
	return c.configStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = c.configStore.Save(txCtx, info); err != nil {
			return err
		}

		if err = c.runtimeStore.Save(txCtx, runtime); err != nil {
			return err
		}

		return c.ToApinto(client, configType, configType, info.Data)
	})
}

func (c *clusterConfigService) Disable(ctx context.Context, namespaceId, operator int, clusterName, configType string) error {

	//获取当前集群信息
	clusterInfo, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}
	info, err := c.configStore.GetConfigByTypeByCluster(ctx, clusterInfo.Id, configType)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("%s config doesn't exist. ", configType)
		}
		return err
	}

	//获取apinto client
	client, err := c.apintoClient.GetClient(ctx, clusterInfo.Id)
	if err != nil {
		return err
	}

	//获取当前运行的版本
	runtime, err := c.runtimeStore.GetForCluster(ctx, info.Id, clusterInfo.Id)
	if err != nil {
		return err
	}

	t := time.Now()

	runtime.IsOnline = false
	runtime.Operator = operator
	runtime.UpdateTime = t

	info.IsEnable = false
	info.Operator = operator
	info.UpdateTime = t
	//事务
	return c.configStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = c.configStore.Save(txCtx, info); err != nil {
			return err
		}

		if err = c.runtimeStore.Save(txCtx, runtime); err != nil {
			return err
		}
		return c.OfflineApinto(client, configType, configType)
	})
}

func (c *clusterConfigService) IsConfigTypeExist(configType string) bool {
	if driver := c.clConfigManager.GetDriver(configType); driver == nil {
		return false
	}
	return true
}

func (c *clusterConfigService) CheckInput(configType string, config []byte) error {
	driver := c.clConfigManager.GetDriver(configType)
	return driver.CheckInput(config)
}

func (c *clusterConfigService) FormatOutput(configType string, operator string, config *cluster_entry.ClusterConfig) interface{} {
	driver := c.clConfigManager.GetDriver(configType)
	return driver.FormatOut(operator, config)
}

// initConfig 初始化配置，比如influxDB初始化bucket和定时脚本
func (c *clusterConfigService) initConfig(configType string, config []byte) error {
	driver := c.clConfigManager.GetDriver(configType)
	return driver.InitConfig(config)
}

func (c *clusterConfigService) ToApinto(client v1.IClient, name, configType string, config []byte) error {
	driver := c.clConfigManager.GetDriver(configType)
	if driver == nil {
		return fmt.Errorf("Get CLConfig Driver fail. type: %s ", configType)
	}
	apintoConfig := driver.ToApinto(name, config)

	switch configType {
	case cluster.CLConfigRedis, cluster.CLConfigInfluxV2:
		return client.ForOutput().Create(apintoConfig)
	default:
		return errors.New("configType doesn't exist. ")
	}
}

func (c *clusterConfigService) OfflineApinto(client v1.IClient, name, configType string) error {
	switch configType {
	case cluster.CLConfigRedis, cluster.CLConfigInfluxV2:
		return common.CheckWorkerNotExist(client.ForOutput().Delete(name))
	default:
		return errors.New("configType doesn't exist. ")
	}
}

func (c *clusterConfigService) ResetOnline(ctx context.Context, _, clusterId int) {
	//获取apinto client
	client, err := c.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		log.Errorf("Get Apinto Client fail. clusterId:%d ", clusterId)
		return
	}

	configList, _ := c.configStore.GetConfigsByClusterID(ctx, clusterId)
	for _, conf := range configList {
		if conf.IsEnable {
			err = c.ToApinto(client, conf.Type, conf.Type, conf.Data)
			if err != nil {
				log.Errorf("Publish Cluster Config %s to Apinto fail. clusterId: %d. err: %s ", conf.Type, clusterId, err)
			}
		}
	}
}
