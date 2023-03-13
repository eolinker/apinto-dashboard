package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	driver_manager "github.com/eolinker/apinto-dashboard/driver-manager"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-basic/uuid"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

var _ IMonitorService = (*monitorService)(nil)

type IMonitorService interface {
	PartitionList(ctx context.Context, namespaceId int) ([]*model.MonPartitionListItem, error)
	PartitionInfo(ctx context.Context, namespaceId int, uuid string) (*model.MonPartitionInfo, error)
	PartitionById(ctx context.Context, id int) (*model.MonPartitionListItem, error)
	CreatePartition(ctx context.Context, namespaceId, userId int, input *dto.MonitorPartitionInfoProxy) (*model.MonPartitionListItem, error)
	UpdatePartition(ctx context.Context, namespaceId, userId int, uuid string, input *dto.MonitorPartitionInfoProxy) (*model.MonPartitionListItem, error)
	DelPartition(ctx context.Context, namespaceId int, uuid string) error

	CheckInput(sourceType string, input []byte) ([]byte, error)
	GetInfluxDbConfig(ctx context.Context, namespaceId int, uuid string) (*model.MonitorInfluxV2Config, error)
}

type monitorService struct {
	monitorStore        store.IMonitorStore
	lockService         IAsynLockService
	warnStrategyService IWarnStrategyService
	userInfoService     IUserInfoService
	clusterService      IClusterService
	monSourceManager    driver_manager.IMonitorSourceManager
}

func newMonitorService() IMonitorService {
	m := &monitorService{}
	bean.Autowired(&m.monitorStore)
	bean.Autowired(&m.lockService)
	bean.Autowired(&m.userInfoService)
	bean.Autowired(&m.clusterService)
	bean.Autowired(&m.monSourceManager)
	bean.Autowired(&m.warnStrategyService)

	return m
}

func (m *monitorService) PartitionList(ctx context.Context, namespaceId int) ([]*model.MonPartitionListItem, error) {
	items, err := m.monitorStore.GetList(ctx, namespaceId)
	if err != nil {
		return nil, err
	}

	list := make([]*model.MonPartitionListItem, 0, len(items))
	for _, item := range items {
		clusterIDs := strings.Split(item.ClusterIDs, ",")
		clusterNames := make([]string, 0, len(clusterIDs))
		for _, cluId := range clusterIDs {
			id, _ := strconv.Atoi(cluId)
			info, _ := m.clusterService.GetByClusterId(ctx, id)
			if info != nil {
				clusterNames = append(clusterNames, info.Name)
			}
		}

		partition := &model.MonPartitionListItem{
			Id:           item.UUID,
			Name:         item.Name,
			ClusterNames: clusterNames,
		}
		list = append(list, partition)
	}
	return list, nil
}

func (m *monitorService) PartitionById(ctx context.Context, id int) (*model.MonPartitionListItem, error) {
	item, err := m.monitorStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.MonPartitionListItem{
		Id:   item.UUID,
		Name: item.Name,
	}, nil

}
func (m *monitorService) PartitionInfo(ctx context.Context, namespaceId int, uuid string) (*model.MonPartitionInfo, error) {
	partitionInfo, err := m.monitorStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("monitor partition not found. uuid:%s ", uuid)
		}
		return nil, err
	}

	cluIDList := strings.Split(partitionInfo.ClusterIDs, ",")
	cluNames := make([]string, 0, len(cluIDList))
	for _, clusterID := range cluIDList {
		cluID, _ := strconv.Atoi(clusterID)
		cluInfo, err := m.clusterService.GetByClusterId(ctx, cluID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			return nil, err
		}
		cluNames = append(cluNames, cluInfo.Name)
	}

	return &model.MonPartitionInfo{
		Id:           partitionInfo.Id,
		Name:         partitionInfo.Name,
		SourceType:   partitionInfo.SourceType,
		Config:       partitionInfo.Config,
		Env:          partitionInfo.Env,
		ClusterNames: cluNames,
	}, nil
}

func (m *monitorService) GetInfluxDbConfig(ctx context.Context, namespaceId int, uuid string) (*model.MonitorInfluxV2Config, error) {
	info, err := m.PartitionInfo(ctx, namespaceId, uuid)
	if err != nil {
		return nil, err
	}

	val := new(model.MonitorInfluxV2Config)

	if err = json.Unmarshal(info.Config, val); err != nil {
		return nil, err
	}
	return val, nil
}

func (m *monitorService) CreatePartition(ctx context.Context, namespaceId, userId int, input *dto.MonitorPartitionInfoProxy) (*model.MonPartitionListItem, error) {
	//分区名查重
	partitions, err := m.monitorStore.GetByName(ctx, namespaceId, input.Name)
	if err != nil {
		return nil, err
	}
	if len(partitions) > 0 {
		return nil, errors.New("monitor partition name is reduplicated. ")
	}

	clusterIds := make([]string, 0, len(input.ClusterNames))
	for _, cluName := range input.ClusterNames {
		cluInfo, err := m.clusterService.GetByNamespaceByName(ctx, namespaceId, cluName)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			return nil, err
		}
		clusterIds = append(clusterIds, strconv.Itoa(cluInfo.Id))
	}

	t := time.Now()
	monitorInfo := &entry.MonitorPartition{
		UUID:       uuid.New(),
		Namespace:  namespaceId,
		Name:       input.Name,
		SourceType: input.SourceType,
		Config:     input.Config,
		Env:        input.Env,
		ClusterIDs: strings.Join(clusterIds, ","),
		Operator:   userId,
		CreateTime: t,
		UpdateTime: t,
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: monitorInfo.UUID,
		Name: input.Name,
	})

	partitionInfo := &model.MonPartitionListItem{
		Id:           monitorInfo.UUID,
		Name:         monitorInfo.Name,
		ClusterNames: input.ClusterNames,
	}

	return partitionInfo, m.monitorStore.Transaction(ctx, func(txCtx context.Context) error {
		return m.monitorStore.Save(txCtx, monitorInfo)
	})
}

func (m *monitorService) UpdatePartition(ctx context.Context, namespaceId, userId int, uuid string, input *dto.MonitorPartitionInfoProxy) (*model.MonPartitionListItem, error) {
	partitionInfo, err := m.monitorStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("monitor partition not found. id:%s ", uuid)
		}
		return nil, err
	}

	err = m.lockService.Lock(LockNameMonPartition, partitionInfo.Id)
	if err != nil {
		return nil, err
	}
	defer m.lockService.Unlock(LockNameMonPartition, partitionInfo.Id)

	partitionInfo, err = m.monitorStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("monitor partition not found. id:%s ", uuid)
		}
		return nil, err
	}

	//分区名查重
	partitions, err := m.monitorStore.GetByName(ctx, namespaceId, input.Name)
	if err != nil {
		return nil, err
	}
	if len(partitions) > 0 {
		for _, partition := range partitions {
			//若为分区本身
			if partition.Id == partitionInfo.Id {
				continue
			}
			return nil, errors.New("monitor partition name is reduplicated. ")
		}
	}

	clusterIds := make([]string, 0, len(input.ClusterNames))
	for _, cluName := range input.ClusterNames {
		cluInfo, err := m.clusterService.GetByNamespaceByName(ctx, namespaceId, cluName)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			return nil, err
		}
		clusterIds = append(clusterIds, strconv.Itoa(cluInfo.Id))
	}

	partitionInfo.Name = input.Name
	partitionInfo.SourceType = input.SourceType
	partitionInfo.Config = input.Config
	partitionInfo.Env = input.Env
	partitionInfo.ClusterIDs = strings.Join(clusterIds, ",")
	partitionInfo.Operator = userId
	partitionInfo.UpdateTime = time.Now()

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: uuid,
		Name: input.Name,
	})

	respInfo := &model.MonPartitionListItem{
		ClusterNames: input.ClusterNames,
	}

	return respInfo, m.monitorStore.Transaction(ctx, func(txCtx context.Context) error {
		return m.monitorStore.Save(txCtx, partitionInfo)
	})
}

func (m *monitorService) DelPartition(ctx context.Context, namespaceId int, uuid string) error {
	partitionInfo, err := m.monitorStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("monitor partition not found. id:%s ", uuid)
		}
		return err
	}

	err = m.lockService.Lock(LockNameMonPartition, partitionInfo.Id)
	if err != nil {
		return err
	}
	defer m.lockService.Unlock(LockNameMonPartition, partitionInfo.Id)

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: uuid,
		Name: partitionInfo.Name,
	})

	return m.monitorStore.Transaction(ctx, func(txCtx context.Context) error {
		_, err = m.monitorStore.Delete(txCtx, partitionInfo.Id)
		if err != nil {
			return err
		}
		//删除分区下的告警策略
		return m.warnStrategyService.DeleteWarnStrategyByPartitionId(txCtx, namespaceId, partitionInfo.Id)

	})
}

func (m *monitorService) CheckInput(sourceType string, input []byte) ([]byte, error) {
	driver := m.monSourceManager.GetDriver(sourceType)
	if driver == nil {
		return nil, fmt.Errorf("source type %s is invalid. ", sourceType)
	}
	return driver.CheckInput(input)
}
