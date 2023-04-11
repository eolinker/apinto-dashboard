package monitor

import (
	"context"
	"github.com/eolinker/apinto-dashboard/driver"
	"github.com/eolinker/apinto-dashboard/modules/monitor/model"
	"github.com/eolinker/apinto-dashboard/modules/monitor/monitor-dto"
	"github.com/eolinker/apinto-dashboard/modules/warn/warn-model"
	"time"
)

type IMonitorService interface {
	PartitionList(ctx context.Context, namespaceId int) ([]*model.MonPartitionListItem, error)
	PartitionInfo(ctx context.Context, namespaceId int, uuid string) (*model.MonPartitionInfo, error)
	PartitionById(ctx context.Context, id int) (*model.MonPartitionListItem, error)
	CreatePartition(ctx context.Context, namespaceId, userId int, input *monitor_dto.MonitorPartitionInfoProxy) (*model.MonPartitionListItem, error)
	UpdatePartition(ctx context.Context, namespaceId, userId int, uuid string, input *monitor_dto.MonitorPartitionInfoProxy) (*model.MonPartitionListItem, error)
	DelPartition(ctx context.Context, namespaceId int, uuid string) error

	CheckInput(sourceType string, input []byte) ([]byte, error)
	GetInfluxDbConfig(ctx context.Context, namespaceId int, uuid string) (*model.MonitorInfluxV2Config, error)
}

type IMonitorStatistics interface {
	Statistics(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, groupBy string, wheres []model.MonWhereItem, limit int) (map[string]model.MonCommonData, error)
	ProxyStatistics(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, groupBy string, wheres []model.MonWhereItem, limit int) (map[string]model.MonCommonData, error)
	Trend(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (*model.MonCallCountInfo, string, error)
	ProxyTrend(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (*model.MonCallCountInfo, string, error)
	IPTrend(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (*model.MonCallCountInfo, error)
	// CircularMap 饼状图数据
	CircularMap(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (request, proxy *model.CircularDate, err error)
	MessageTrend(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (*model.MessageTrend, string, error)
	WarnStatistics(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, group string, quotaType warn_model.QuotaType, wheres []model.MonWhereItem) (map[string]float64, error)
}

type IMonitorStatisticsCache interface {
	GetStatisticsCache(ctx context.Context, partitionId string, start, end time.Time, groupBy string, wheres []model.MonWhereItem, limit int) (map[string]model.MonCommonData, error)
	SetStatisticsCache(ctx context.Context, partitionId string, start, end time.Time, groupBy string, wheres []model.MonWhereItem, limit int, values map[string]model.MonCommonData) error

	GetTrendCache(ctx context.Context, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (*model.MonCallCountInfo, error)
	SetTrendCache(ctx context.Context, partitionId string, start, end time.Time, wheres []model.MonWhereItem, value *model.MonCallCountInfo) error

	GetCircularMap(ctx context.Context, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (request, proxy *model.CircularDate, err error)
	SetCircularMap(ctx context.Context, partitionId string, start, end time.Time, wheres []model.MonWhereItem, request, proxy *model.CircularDate) error

	GetMessageTrend(ctx context.Context, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (*model.MessageTrend, error)
	SetMessageTrend(ctx context.Context, partitionId string, start, end time.Time, wheres []model.MonWhereItem, val *model.MessageTrend) error
}

// IMonitorSourceDriver 监控数据源驱动
type IMonitorSourceDriver interface {
	CheckInput(config []byte) ([]byte, error)
}

type IMonitorSourceManager interface {
	driver.IDriverManager[IMonitorSourceDriver]
	List() []string
}
