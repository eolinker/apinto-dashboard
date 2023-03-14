package monitor_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/model/monitor-model"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"sort"
	"strings"
	"time"
)

type IMonitorStatisticsCache interface {
	GetStatisticsCache(ctx context.Context, partitionId string, start, end time.Time, groupBy string, wheres []monitor_model.MonWhereItem, limit int) (map[string]monitor_model.MonCommonData, error)
	SetStatisticsCache(ctx context.Context, partitionId string, start, end time.Time, groupBy string, wheres []monitor_model.MonWhereItem, limit int, values map[string]monitor_model.MonCommonData) error

	GetTrendCache(ctx context.Context, partitionId string, start, end time.Time, wheres []monitor_model.MonWhereItem) (*monitor_model.MonCallCountInfo, error)
	SetTrendCache(ctx context.Context, partitionId string, start, end time.Time, wheres []monitor_model.MonWhereItem, value *monitor_model.MonCallCountInfo) error

	GetCircularMap(ctx context.Context, partitionId string, start, end time.Time, wheres []monitor_model.MonWhereItem) (request, proxy *monitor_model.CircularDate, err error)
	SetCircularMap(ctx context.Context, partitionId string, start, end time.Time, wheres []monitor_model.MonWhereItem, request, proxy *monitor_model.CircularDate) error

	GetMessageTrend(ctx context.Context, partitionId string, start, end time.Time, wheres []monitor_model.MonWhereItem) (*monitor_model.MessageTrend, error)
	SetMessageTrend(ctx context.Context, partitionId string, start, end time.Time, wheres []monitor_model.MonWhereItem, val *monitor_model.MessageTrend) error
}
type monitorStatisticsCache struct {
	commonCache cache.ICommonCache
}

func newMonitorStatisticsCache() IMonitorStatisticsCache {
	val := new(monitorStatisticsCache)
	bean.Autowired(&val.commonCache)
	return val
}

func (m *monitorStatisticsCache) GetStatisticsCache(ctx context.Context, partitionId string, start, end time.Time, groupBy string, wheres []monitor_model.MonWhereItem, limit int) (map[string]monitor_model.MonCommonData, error) {

	key := fmt.Sprintf("monitor:statistics:%s:%d_%d:%s:%s:%d", partitionId, start.Unix(), end.Unix(), groupBy, m.formatWhereKey(wheres), limit)

	maps, err := m.commonCache.HGetAll(ctx, key)
	if err != nil {
		log.Errorf("GetStatisticsCache cache.HGetAll key=%s err=%s", key, err.Error())
		return nil, err
	}
	valMap := make(map[string]monitor_model.MonCommonData)
	for k, v := range maps {
		commonData := &monitor_model.MonCommonData{}
		if err = json.Unmarshal([]byte(v), commonData); err != nil {
			log.Errorf("GetStatisticsCache json.Unmarshal err=%s", err.Error())
			return nil, err
		}
		valMap[k] = *commonData
	}

	return valMap, nil
}

func (m *monitorStatisticsCache) SetStatisticsCache(ctx context.Context, partitionId string, start, end time.Time, groupBy string, wheres []monitor_model.MonWhereItem, limit int, values map[string]monitor_model.MonCommonData) error {
	key := fmt.Sprintf("monitor:statistics:%s:%d_%d:%s:%s:%d", partitionId, start.Unix(), end.Unix(), groupBy, m.formatWhereKey(wheres), limit)

	maps := make(map[string][]byte)
	for k, data := range values {
		bytes, err := json.Marshal(data)
		if err != nil {
			log.Errorf("SetStatisticsCache json.Marshal key=%s err=%s", key, err.Error())
			return err
		}
		maps[k] = bytes
	}

	return m.commonCache.HMSet(ctx, key, maps, 5*time.Minute)
}

func (m *monitorStatisticsCache) GetTrendCache(ctx context.Context, partitionId string, start, end time.Time, wheres []monitor_model.MonWhereItem) (*monitor_model.MonCallCountInfo, error) {
	key := fmt.Sprintf("monitor:trend:%s:%d_%d:%s", partitionId, start.Unix(), end.Unix(), m.formatWhereKey(wheres))

	bytes, err := m.commonCache.Get(ctx, key)
	if err != nil {
		log.Errorf("GetTrendCache cache.Get key=%s err=%s", key, err.Error())
		return nil, err
	}

	val := new(monitor_model.MonCallCountInfo)

	if err = json.Unmarshal(bytes, val); err != nil {
		log.Errorf("GetTrendCache json.Unmarshal key=%s bytes=%v err=%s", key, bytes, err.Error())
		return nil, err
	}

	return val, nil
}

func (m *monitorStatisticsCache) SetTrendCache(ctx context.Context, partitionId string, start, end time.Time, wheres []monitor_model.MonWhereItem, value *monitor_model.MonCallCountInfo) error {
	key := fmt.Sprintf("monitor:trend:%s:%d_%d:%s", partitionId, start.Unix(), end.Unix(), m.formatWhereKey(wheres))

	bytes, err := json.Marshal(value)
	if err != nil {
		log.Errorf("SetTrendCache json.Marshal key=%s val=%v err=%s", key, value, err.Error())
		return err
	}

	return m.commonCache.Set(ctx, key, bytes, 5*time.Minute)
}

func (m *monitorStatisticsCache) GetCircularMap(ctx context.Context, partitionId string, start, end time.Time, wheres []monitor_model.MonWhereItem) (*monitor_model.CircularDate, *monitor_model.CircularDate, error) {
	key := fmt.Sprintf("monitor:circular_map:%s:%d_%d:%s", partitionId, start.Unix(), end.Unix(), m.formatWhereKey(wheres))

	maps, err := m.commonCache.HGetAll(ctx, key)
	if err != nil {
		log.Errorf("GetCircularMap cache.HGetAll key=%s err=%s", key, err.Error())
		return nil, nil, err
	}

	requestDate := new(monitor_model.CircularDate)
	if req, ok := maps["request"]; ok {
		if err = json.Unmarshal([]byte(req), requestDate); err != nil {
			log.Errorf("GetCircularMap json.Unmarshal request=%s  key=%s err=%s", req, key, err.Error())
			return nil, nil, err
		}
	} else {
		return nil, nil, errors.New("获取不到数据")
	}

	proxyDate := new(monitor_model.CircularDate)
	if proxy, ok := maps["proxy"]; ok {
		if err = json.Unmarshal([]byte(proxy), proxyDate); err != nil {
			log.Errorf("GetCircularMap json.Unmarshal proxy=%s key=%s err=%s", proxy, key, err.Error())
			return nil, nil, err
		}
	} else {
		return nil, nil, errors.New("获取不到数据")
	}

	return requestDate, proxyDate, nil

}

func (m *monitorStatisticsCache) SetCircularMap(ctx context.Context, partitionId string, start, end time.Time, wheres []monitor_model.MonWhereItem, request, proxy *monitor_model.CircularDate) error {
	key := fmt.Sprintf("monitor:circular_map:%s:%d_%d:%s", partitionId, start.Unix(), end.Unix(), m.formatWhereKey(wheres))

	maps := make(map[string][]byte)
	reqByte, err := json.Marshal(request)
	if err != nil {
		log.Errorf("SetCircularMap json.Marshal reqByte key=%s err=%s", key, err.Error())
		return err
	}

	proxyByte, err := json.Marshal(proxy)
	if err != nil {
		log.Errorf("SetStatisticsCache json.Marshal proxyByte key=%s err=%s", key, err.Error())
		return err
	}
	maps["request"] = reqByte
	maps["proxy"] = proxyByte

	return m.commonCache.HMSet(ctx, key, maps, 5*time.Minute)
}

func (m *monitorStatisticsCache) GetMessageTrend(ctx context.Context, partitionId string, start, end time.Time, wheres []monitor_model.MonWhereItem) (*monitor_model.MessageTrend, error) {
	key := fmt.Sprintf("monitor:message_trend:%s:%d_%d:%s", partitionId, start.Unix(), end.Unix(), m.formatWhereKey(wheres))

	bytes, err := m.commonCache.Get(ctx, key)
	if err != nil {
		log.Errorf("GetMessageTrend cache.Get key=%s err=%s", key, err.Error())
		return nil, err
	}

	val := new(monitor_model.MessageTrend)

	if err = json.Unmarshal(bytes, val); err != nil {
		log.Errorf("GetMessageTrend json.Unmarshal key=%s bytes=%v err=%s", key, bytes, err.Error())
		return nil, err
	}

	return val, nil
}

func (m *monitorStatisticsCache) SetMessageTrend(ctx context.Context, partitionId string, start, end time.Time, wheres []monitor_model.MonWhereItem, value *monitor_model.MessageTrend) error {
	key := fmt.Sprintf("monitor:message_trend:%s:%d_%d:%s", partitionId, start.Unix(), end.Unix(), m.formatWhereKey(wheres))

	bytes, err := json.Marshal(value)
	if err != nil {
		log.Errorf("SetMessageTrend json.Marshal key=%s val=%v err=%s", key, value, err.Error())
		return err
	}

	return m.commonCache.Set(ctx, key, bytes, 5*time.Minute)
}

// formatWhereKey
// 排序规则 addr:api:app:cluster:ip:path:upstream
func (m *monitorStatisticsCache) formatWhereKey(wheres []monitor_model.MonWhereItem) string {

	whereMap := make(map[string]monitor_model.MonWhereItem)
	keys := make([]string, 0, len(wheres))
	for _, where := range wheres {
		whereMap[where.Key] = where
		keys = append(keys, where.Key)
	}

	sort.Strings(keys)

	redisKeys := make([]string, 0)
	for _, key := range keys {
		if v, ok := whereMap[key]; ok {
			sort.Strings(v.Values)
			redisKeys = append(redisKeys, fmt.Sprintf("%v", strings.Join(v.Values, "_")))
		}
	}

	return strings.Join(redisKeys, ":")
}
