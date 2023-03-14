package service

import (
	"context"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/entry/monitor-entry"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/eolinker/apinto-dashboard/store/flux"
	"github.com/eolinker/eosc/common/bean"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"strings"
	"time"
)

const (
	oneHour = 3600
	oneDay  = 24 * oneHour
	tenDay  = 10 * oneDay
	oneYear = 365 * oneDay

	bucketMinuteRetention = (7 - 1) * oneDay
	bucketHourRetention   = (90 - 1) * oneDay
	bucketDayRetention    = (5*365 - 1) * oneDay

	bucketMinute = "apinto/minute"
	bucketHour   = "apinto/hour"
	bucketDay    = "apinto/day"
	bucketWeek   = "apinto/week"

	timeZone = "Asia/Shanghai"
)

type IMonitorStatistics interface {
	Statistics(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, groupBy string, wheres []model.MonWhereItem, limit int) (map[string]model.MonCommonData, error)
	ProxyStatistics(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, groupBy string, wheres []model.MonWhereItem, limit int) (map[string]model.MonCommonData, error)
	Trend(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (*model.MonCallCountInfo, string, error)
	ProxyTrend(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (*model.MonCallCountInfo, string, error)
	IPTrend(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (*model.MonCallCountInfo, error)
	// CircularMap 饼状图数据
	CircularMap(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (request, proxy *model.CircularDate, err error)
	MessageTrend(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (*model.MessageTrend, string, error)
	WarnStatistics(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, group string, quotaType model.QuotaType, wheres []model.MonWhereItem) (map[string]float64, error)
}

type monitorStatistics struct {
	fluxQuery              flux.IFluxQuery
	monitorService         IMonitorService
	monitorStatisticsCache IMonitorStatisticsCache
}

func newMonitorStatistics() IMonitorStatistics {
	statistics := &monitorStatistics{}
	bean.Autowired(&statistics.fluxQuery)
	bean.Autowired(&statistics.monitorStatisticsCache)
	bean.Autowired(&statistics.monitorService)
	return statistics
}

func (m *monitorStatistics) Statistics(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, groupBy string, wheres []model.MonWhereItem, limit int) (map[string]model.MonCommonData, error) {

	statisticsCache, _ := m.monitorStatisticsCache.GetStatisticsCache(ctx, partitionId, start, end, groupBy, wheres, limit)
	if len(statisticsCache) > 0 {
		return statisticsCache, nil
	}

	dbConfig, err := m.monitorService.GetInfluxDbConfig(ctx, namespaceId, partitionId)
	if err != nil {
		return nil, err
	}

	openAPI, err := m.getOpenApi(dbConfig)
	if err != nil {
		return nil, err
	}
	filters := m.formatFilter(wheres)
	newStartTime, _, _, bucket := m.getTimeIntervelAndBucket(start, end)

	statisticsConf := []*monitor_entry.StatisticsFilterConf{
		{
			Measurement: "request",
			AggregateFn: "sum()",
			Fields:      []string{"total", "success", "timing", "request"},
		},
		{
			Measurement: "proxy",
			AggregateFn: "sum()",
			Fields:      []string{"p_total", "p_success"},
		},
		{
			Measurement: "request",
			AggregateFn: "max()",
			Fields:      []string{"timing_max", "request_max"},
		}, {
			Measurement: "request",
			AggregateFn: "min()",
			Fields:      []string{"timing_min", "request_min"},
		},
	}

	results, err := m.fluxQuery.CommonStatistics(ctx, openAPI, newStartTime, end, bucket, groupBy, filters, statisticsConf, limit)
	if err != nil {
		return nil, err
	}

	resultMaps := make(map[string]model.MonCommonData)
	for key, result := range results {

		requestRate := 0.0
		if result.Total == 0 {
			requestRate = 0.0
		} else {
			requestRate = common.FormatFloat64(float64(result.Success)/float64(result.Total), 4)
		}

		proxyRate := 0.0
		if result.ProxyTotal == 0 {
			proxyRate = 0.0
		} else {
			proxyRate = common.FormatFloat64(float64(result.ProxySuccess)/float64(result.ProxyTotal), 4)
		}

		monCommonData := model.MonCommonData{
			RequestTotal:   result.Total,
			RequestSuccess: result.Success,
			RequestRate:    requestRate,
			ProxyTotal:     result.ProxyTotal,
			ProxySuccess:   result.ProxySuccess,
			ProxyRate:      proxyRate,
			StatusFail:     result.Total - result.Success,
			AvgResp:        float64(result.TotalTiming) / float64(result.Total),
			MaxResp:        result.MaxTiming,
			MinResp:        result.MinTiming,
			AvgTraffic:     float64(result.TotalRequest) / float64(result.Total),
			MaxTraffic:     result.RequestMax,
			MinTraffic:     result.RequestMin,
		}

		resultMaps[key] = monCommonData
	}

	_ = m.monitorStatisticsCache.SetStatisticsCache(ctx, partitionId, start, end, groupBy, wheres, limit, resultMaps)

	return resultMaps, nil

}

func (m *monitorStatistics) ProxyStatistics(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, groupBy string, wheres []model.MonWhereItem, limit int) (map[string]model.MonCommonData, error) {

	statisticsCache, _ := m.monitorStatisticsCache.GetStatisticsCache(ctx, partitionId, start, end, groupBy, wheres, limit)
	if len(statisticsCache) > 0 {
		return statisticsCache, nil
	}

	dbConfig, err := m.monitorService.GetInfluxDbConfig(ctx, namespaceId, partitionId)
	if err != nil {
		return nil, err
	}

	openAPI, err := m.getOpenApi(dbConfig)
	if err != nil {
		return nil, err
	}

	filters := m.formatFilter(wheres)
	newStartTime, _, _, bucket := m.getTimeIntervelAndBucket(start, end)

	statisticsConf := []*monitor_entry.StatisticsFilterConf{
		{
			Measurement: "proxy",
			AggregateFn: "sum()",
			Fields:      []string{"p_total", "p_success", "p_timing", "p_request"},
		},
		{
			Measurement: "proxy",
			AggregateFn: "max()",
			Fields:      []string{"p_timing_max", "p_request_max"},
		}, {
			Measurement: "proxy",
			AggregateFn: "min()",
			Fields:      []string{"p_timing_min", "p_request_min"},
		},
	}

	results, err := m.fluxQuery.CommonProxyStatistics(ctx, openAPI, newStartTime, end, bucket, groupBy, filters, statisticsConf, limit)
	if err != nil {
		return nil, err
	}
	resultMaps := make(map[string]model.MonCommonData)
	for key, result := range results {

		proxyRate := 0.0
		if result.ProxyTotal == 0 {
			proxyRate = 0.0
		} else {
			proxyRate = common.FormatFloat64(float64(result.ProxySuccess)/float64(result.ProxyTotal), 4)
		}

		monCommonData := model.MonCommonData{
			ProxyTotal:   result.ProxyTotal,
			ProxySuccess: result.ProxySuccess,
			ProxyRate:    proxyRate,
			StatusFail:   result.ProxyTotal - result.ProxySuccess,
			AvgResp:      float64(result.TotalTiming) / float64(result.ProxyTotal),
			MaxResp:      result.MaxTiming,
			MinResp:      result.MinTiming,
			AvgTraffic:   float64(result.TotalRequest) / float64(result.ProxyTotal),
			MaxTraffic:   result.RequestMax,
			MinTraffic:   result.RequestMin,
		}

		resultMaps[key] = monCommonData
	}

	_ = m.monitorStatisticsCache.SetStatisticsCache(ctx, partitionId, start, end, groupBy, wheres, limit, resultMaps)

	return resultMaps, nil
}

func (m *monitorStatistics) CircularMap(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (request, proxy *model.CircularDate, err error) {

	request, proxy, _ = m.monitorStatisticsCache.GetCircularMap(ctx, partitionId, start, end, wheres)
	if request != nil && proxy != nil {
		return request, proxy, nil
	}

	dbConfig, err := m.monitorService.GetInfluxDbConfig(ctx, namespaceId, partitionId)
	if err != nil {
		return nil, nil, err
	}

	openAPI, err := m.getOpenApi(dbConfig)
	if err != nil {
		return nil, nil, err
	}

	filters := m.formatFilter(wheres)
	newStartTime, _, _, bucket := m.getTimeIntervelAndBucket(start, end)

	requestFieldsConf := &monitor_entry.StatisticsFilterConf{
		Measurement: "request",
		AggregateFn: "sum()",
		Fields:      []string{"total", "success", "s4xx", "s5xx"},
	}
	//获取请求表的饼状图
	requestQueryOnce, err := m.fluxQuery.CommonQueryOnce(ctx, openAPI, newStartTime, end, bucket, filters, requestFieldsConf)
	if err != nil {
		return nil, nil, err
	}

	proxyFieldsConf := &monitor_entry.StatisticsFilterConf{
		Measurement: "proxy",
		AggregateFn: "sum()",
		Fields:      []string{"p_total", "p_success", "p_s4xx", "p_s5xx"},
	}
	//获取转发表的饼状图
	proxyQueryOnce, err := m.fluxQuery.CommonQueryOnce(ctx, openAPI, newStartTime, end, bucket, filters, proxyFieldsConf)
	if err != nil {
		return nil, nil, err
	}
	request = new(model.CircularDate)
	proxy = new(model.CircularDate)

	if v, ok := requestQueryOnce["s4xx"]; ok {
		request.Status4Xx = common.FmtIntFromInterface(v)
	}
	if v, ok := requestQueryOnce["s5xx"]; ok {
		request.Status5Xx = common.FmtIntFromInterface(v)
	}
	if v, ok := requestQueryOnce["success"]; ok {
		request.Success = common.FmtIntFromInterface(v)
	}
	if v, ok := requestQueryOnce["total"]; ok {
		request.Total = common.FmtIntFromInterface(v)
	}
	request.Fail = request.Total - request.Success

	if v, ok := proxyQueryOnce["p_s4xx"]; ok {
		proxy.Status4Xx = common.FmtIntFromInterface(v)
	}
	if v, ok := proxyQueryOnce["p_s5xx"]; ok {
		proxy.Status5Xx = common.FmtIntFromInterface(v)
	}
	if v, ok := proxyQueryOnce["p_success"]; ok {
		proxy.Success = common.FmtIntFromInterface(v)
	}
	if v, ok := proxyQueryOnce["p_total"]; ok {
		proxy.Total = common.FmtIntFromInterface(v)
	}
	proxy.Fail = proxy.Total - proxy.Success

	_ = m.monitorStatisticsCache.SetCircularMap(ctx, partitionId, start, end, wheres, request, proxy)

	return request, proxy, nil
}

func (m *monitorStatistics) MessageTrend(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (*model.MessageTrend, string, error) {
	newStartTime, every, windowOffset, bucket := m.getTimeIntervelAndBucket(start, end)

	messageTrend, _ := m.monitorStatisticsCache.GetMessageTrend(ctx, partitionId, start, end, wheres)
	if messageTrend != nil {
		return messageTrend, every, nil
	}

	dbConfig, err := m.monitorService.GetInfluxDbConfig(ctx, namespaceId, partitionId)
	if err != nil {
		return nil, "", err
	}

	openAPI, err := m.getOpenApi(dbConfig)
	if err != nil {
		return nil, "", err
	}

	filters := m.formatFilter(wheres)

	fieldsConditions := []string{"request", "response"}

	dates, groupValues, err := m.fluxQuery.CommonTendency(ctx, openAPI, newStartTime, end, bucket, "request", filters, fieldsConditions, every, windowOffset)
	if err != nil {
		return nil, "", err
	}

	resultVal := &model.MessageTrend{
		Dates:       dates,
		ReqMessage:  m.formatMessageTrendData(groupValues["request"]),
		RespMessage: m.formatMessageTrendData(groupValues["response"]),
	}

	_ = m.monitorStatisticsCache.SetMessageTrend(ctx, partitionId, start, end, wheres, resultVal)

	return resultVal, every, nil
}

func (m *monitorStatistics) Trend(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (*model.MonCallCountInfo, string, error) {

	newStartTime, every, windowOffset, bucket := m.getTimeIntervelAndBucket(start, end)

	trendCache, _ := m.monitorStatisticsCache.GetTrendCache(ctx, partitionId, start, end, wheres)
	if trendCache != nil {
		return trendCache, every, nil
	}

	dbConfig, err := m.monitorService.GetInfluxDbConfig(ctx, namespaceId, partitionId)
	if err != nil {
		return nil, "", err
	}

	openAPI, err := m.getOpenApi(dbConfig)
	if err != nil {
		return nil, "", err
	}

	filters := m.formatFilter(wheres)

	requestConditions := []string{"total", "success", "s4xx", "s5xx"}

	dates, requestValues, err := m.fluxQuery.CommonTendency(ctx, openAPI, newStartTime, end, bucket, "request", filters, requestConditions, every, windowOffset)
	if err != nil {
		return nil, "", err
	}
	requestRate := make([]float64, 0, len(dates))
	//计算请求成功率
	requestTotal := requestValues["total"]
	requestSuccess := requestValues["success"]
	for i, total := range requestTotal {
		if total == 0 {
			requestRate = append(requestRate, 0)
			continue
		}
		rate := common.FormatFloat64(float64(requestSuccess[i])/float64(total), 4)
		requestRate = append(requestRate, rate)
	}

	proxyConditions := []string{"p_total", "p_success"}

	_, proxyValues, err := m.fluxQuery.CommonTendency(ctx, openAPI, newStartTime, end, bucket, "proxy", filters, proxyConditions, every, windowOffset)
	if err != nil {
		return nil, "", err
	}
	//计算转发成功率
	proxyTotal := proxyValues["p_total"]
	proxySuccess := proxyValues["p_success"] //proxySuccess和proxyTotal必定等长
	proxyRate := make([]float64, 0, len(proxyTotal))
	for i, total := range proxyTotal {
		if total == 0 {
			proxyRate = append(proxyRate, 0)
			continue
		}
		rate := common.FormatFloat64(float64(proxySuccess[i])/float64(total), 4)
		proxyRate = append(proxyRate, rate)
	}
	resultVal := &model.MonCallCountInfo{
		Date:         dates,
		Status5XX:    requestValues["s5xx"],
		Status4XX:    requestValues["s4xx"],
		ProxyRate:    proxyRate,
		ProxyTotal:   proxyValues["p_total"],
		RequestRate:  requestRate,
		RequestTotal: requestValues["total"],
	}

	_ = m.monitorStatisticsCache.SetTrendCache(ctx, partitionId, start, end, wheres, resultVal)

	return resultVal, every, nil
}

func (m *monitorStatistics) ProxyTrend(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (*model.MonCallCountInfo, string, error) {
	newStartTime, every, windowOffset, bucket := m.getTimeIntervelAndBucket(start, end)

	trendCache, _ := m.monitorStatisticsCache.GetTrendCache(ctx, partitionId, start, end, wheres)
	if trendCache != nil {
		return trendCache, every, nil
	}

	dbConfig, err := m.monitorService.GetInfluxDbConfig(ctx, namespaceId, partitionId)
	if err != nil {
		return nil, "", err
	}

	openAPI, err := m.getOpenApi(dbConfig)
	if err != nil {
		return nil, "", err
	}

	filters := m.formatFilter(wheres)

	proxyConditions := []string{"p_total", "p_success", "p_s4xx", "p_s5xx"}

	dates, proxyValues, err := m.fluxQuery.CommonTendency(ctx, openAPI, newStartTime, end, bucket, "proxy", filters, proxyConditions, every, windowOffset)
	if err != nil {
		return nil, "", err
	}
	proxyRate := make([]float64, 0, len(dates))
	//计算请求成功率
	proxyTotal := proxyValues["p_total"]
	proxySuccess := proxyValues["p_success"]
	for i, total := range proxyTotal {
		if total == 0 {
			proxyRate = append(proxyRate, 0)
			continue
		}
		rate := common.FormatFloat64(float64(proxySuccess[i])/float64(total), 4)
		proxyRate = append(proxyRate, rate)
	}

	resultVal := &model.MonCallCountInfo{
		Date:       dates,
		Status5XX:  proxyValues["p_s5xx"],
		Status4XX:  proxyValues["p_s4xx"],
		ProxyRate:  proxyRate,
		ProxyTotal: proxyValues["p_total"],
	}

	_ = m.monitorStatisticsCache.SetTrendCache(ctx, partitionId, start, end, wheres, resultVal)

	return resultVal, every, nil
}

func (m *monitorStatistics) IPTrend(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, wheres []model.MonWhereItem) (*model.MonCallCountInfo, error) {

	trendCache, _ := m.monitorStatisticsCache.GetTrendCache(ctx, partitionId, start, end, wheres)
	if trendCache != nil {
		return trendCache, nil
	}

	dbConfig, err := m.monitorService.GetInfluxDbConfig(ctx, namespaceId, partitionId)
	if err != nil {
		return nil, err
	}

	openAPI, err := m.getOpenApi(dbConfig)
	if err != nil {
		return nil, err
	}

	filters := m.formatFilter(wheres)
	newStartTime, every, windowOffset, bucket := m.getTimeIntervelAndBucket(start, end)

	requestConditions := []string{"total", "success", "s4xx", "s5xx"}

	dates, requestValues, err := m.fluxQuery.CommonTendency(ctx, openAPI, newStartTime, end, bucket, "request", filters, requestConditions, every, windowOffset)
	if err != nil {
		return nil, err
	}
	requestRate := make([]float64, 0, len(dates))
	//计算请求成功率
	requestTotal := requestValues["total"]
	requestSuccess := requestValues["success"]
	for i, total := range requestTotal {
		if total == 0 {
			requestRate = append(requestRate, 0)
			continue
		}
		rate := common.FormatFloat64(float64(requestSuccess[i])/float64(total), 4)
		requestRate = append(requestRate, rate)
	}

	resultVal := &model.MonCallCountInfo{
		Date:         dates,
		Status5XX:    requestValues["s5xx"],
		Status4XX:    requestValues["s4xx"],
		RequestRate:  requestRate,
		RequestTotal: requestValues["total"],
	}

	_ = m.monitorStatisticsCache.SetTrendCache(ctx, partitionId, start, end, wheres, resultVal)

	return resultVal, nil
}

func (m *monitorStatistics) WarnStatistics(ctx context.Context, namespaceId int, partitionId string, start, end time.Time, group string, quotaType model.QuotaType, wheres []model.MonWhereItem) (map[string]float64, error) {
	/*
		当group为api，upstream, wheres为空
		当group为cluster, wheres可能为空可能不为空
	*/
	dbConfig, err := m.monitorService.GetInfluxDbConfig(ctx, namespaceId, partitionId)
	if err != nil {
		return nil, err
	}

	openAPI, err := m.getOpenApi(dbConfig)
	if err != nil {
		return nil, err
	}

	filters := m.formatFilter(wheres)
	//校验start，end， 原则上start最多距离现在一个小时
	newStartTime, _, _, bucket := m.getTimeIntervelAndBucket(start, end)
	if bucket != bucketMinute {
		return nil, fmt.Errorf("start and end is illegal. ")
	}

	statisticsConf, err := m.getStatisticsFilterByQuota(quotaType)
	if err != nil {
		return nil, err
	}

	results, err := m.fluxQuery.CommonWarnStatistics(ctx, openAPI, newStartTime, end, bucket, group, filters, statisticsConf)
	if err != nil {
		return nil, err
	}

	return m.calculateWarnStatistics(quotaType, results), nil
}

func (m *monitorStatistics) calculateWarnStatistics(quotaType model.QuotaType, data map[string]*monitor_entry.FluxWarnStatistics) map[string]float64 {
	results := make(map[string]float64)

	switch quotaType {
	case model.QuotaTypeReqFailCount:
		for key, statistic := range data {
			results[key] = float64(statistic.Total - statistic.Success)
		}
	case model.QuotaTypeReqFailRate:
		for key, statistic := range data {
			if statistic.Total != 0 {
				results[key] = float64(statistic.Total-statistic.Success) / float64(statistic.Total)
			} else {
				results[key] = 0
			}
		}
	case model.QuotaTypeReqStatus4xx:
		for key, statistic := range data {
			results[key] = float64(statistic.S4xx)
		}
	case model.QuotaTypeReqStatus5xx:
		for key, statistic := range data {
			results[key] = float64(statistic.S5xx)
		}
	case model.QuotaTypeProxyFailCount:
		for key, statistic := range data {
			results[key] = float64(statistic.ProxyTotal - statistic.ProxySuccess)
		}
	case model.QuotaTypeProxyFailRate:
		for key, statistic := range data {
			if statistic.ProxyTotal != 0 {
				results[key] = float64(statistic.ProxyTotal-statistic.ProxySuccess) / float64(statistic.ProxyTotal)
			} else {
				results[key] = 0
			}
		}
	case model.QuotaTypeProxyStatus4xx:
		for key, statistic := range data {
			results[key] = float64(statistic.ProxyS4xx)
		}
	case model.QuotaTypeProxyStatus5xx:
		for key, statistic := range data {
			results[key] = float64(statistic.ProxyS5xx)
		}
	case model.QuotaTypeReqMessage:
		for key, statistic := range data {
			results[key] = float64(statistic.TotalRequest) / 1024
		}
	case model.QuotaTypeRespMessage:
		for key, statistic := range data {
			results[key] = float64(statistic.TotalResponse) / 1024
		}
	case model.QuotaTypeAvgResp:
		for key, statistic := range data {
			if statistic.Total != 0 {
				results[key] = float64(statistic.TotalTiming) / float64(statistic.Total)
			} else {
				results[key] = 0
			}
		}
	}
	return results
}

// getClient todo 测试所用
func (m *monitorStatistics) getClient(ctx context.Context, namespaceId int, partitionId string) (influxdb2.Client, error) {
	dbConfig, err := m.monitorService.GetInfluxDbConfig(ctx, namespaceId, partitionId)
	if err != nil {
		return nil, err
	}

	//token := "LtzG9kZyxwAdH0yUP9XqBZKPXsrR04F4QaWxYzHCgnHwsWKrYy7waLySXWOhhMXv49M3OZATfqLoGUPtfVj4sw=="
	//client := influxdb2.NewClient("http://192.168.72.128:8086", token)
	//_, err := client.Ping(context.Background())
	//if err != nil {
	//	return nil, err
	//}

	client := influxdb2.NewClient(dbConfig.Addr, dbConfig.Token)
	_, err = client.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (m *monitorStatistics) getOpenApi(dbConfig *model.MonitorInfluxV2Config) (api.QueryAPI, error) {

	client := influxdb2.NewClient(dbConfig.Addr, dbConfig.Token)
	//_, err := client.Ping(context.Background())
	//if err != nil {
	//	return nil, err
	//}

	return client.QueryAPI(dbConfig.Org), nil
}

func (m *monitorStatistics) formatFilter(wheres []model.MonWhereItem) string {
	filter := ``
	if len(wheres) > 0 {
		filters := make([]string, 0, len(wheres))
		for _, where := range wheres {
			if len(where.Values) > 0 {
				wl := make([]string, 0, len(where.Values))
				for _, v := range where.Values {
					wl = append(wl, fmt.Sprintf(fmt.Sprintf(`r["%s"] == "%s"`, where.Key, v)))
				}
				filters = append(filters, fmt.Sprint("(", strings.Join(wl, " or "), ")"))
			}
		}
		filter = fmt.Sprint(`|> filter(fn:(r)=>`, strings.Join(filters, " and "), ")")
	}
	return filter
}

// getStatisticsFilterByQuota 根据quotaType获取flux脚本需要的统计过滤配置
func (m *monitorStatistics) getStatisticsFilterByQuota(quotaType model.QuotaType) (*monitor_entry.StatisticsFilterConf, error) {
	filterConf := &monitor_entry.StatisticsFilterConf{
		Measurement: "request",
		AggregateFn: "sum()",
		Fields:      nil,
	}

	switch quotaType {
	case model.QuotaTypeReqFailCount, model.QuotaTypeReqFailRate:
		filterConf.Fields = []string{"total", "success"}
		return filterConf, nil
	case model.QuotaTypeReqStatus4xx:
		filterConf.Fields = []string{"s4xx"}
		return filterConf, nil
	case model.QuotaTypeReqStatus5xx:
		filterConf.Fields = []string{"s5xx"}
		return filterConf, nil
	case model.QuotaTypeProxyFailCount, model.QuotaTypeProxyFailRate:
		filterConf.Measurement = "proxy"
		filterConf.Fields = []string{"p_total", "p_success"}
		return filterConf, nil
	case model.QuotaTypeProxyStatus4xx:
		filterConf.Measurement = "proxy"
		filterConf.Fields = []string{"p_s4xx"}
		return filterConf, nil
	case model.QuotaTypeProxyStatus5xx:
		filterConf.Measurement = "proxy"
		filterConf.Fields = []string{"p_s5xx"}
		return filterConf, nil
	case model.QuotaTypeReqMessage:
		filterConf.Fields = []string{"request"}
		return filterConf, nil
	case model.QuotaTypeRespMessage:
		filterConf.Fields = []string{"response"}
		return filterConf, nil
	case model.QuotaTypeAvgResp:
		filterConf.Fields = []string{"total", "timing"}
		return filterConf, nil
	}
	return filterConf, fmt.Errorf("quotaType %s is illegal. ", quotaType)
}

// formatMessageTrendData 格式化报文总量趋势图的数据，将B转换为KB,保留两位小数
func (m *monitorStatistics) formatMessageTrendData(data []int64) []float64 {
	floatData := make([]float64, 0, len(data))
	for _, d := range data {
		floatData = append(floatData, common.FormatFloat64(float64(d)/1024, 2))
	}
	return floatData
}

// getTimeIntervelAndBucket 根据start和end来获取窗口时间间隔，窗口偏移量offset，以及使用的bucket, 查询的startTime也会格式化
func (m *monitorStatistics) getTimeIntervelAndBucket(startTime, endTime time.Time) (time.Time, string, string, string) {
	//根据start距离现在的时长算出可使用的最小桶
	minimumBucket := ""
	startToNow := time.Now().Unix() - startTime.Unix()
	if startToNow <= bucketMinuteRetention {
		minimumBucket = bucketMinute
	} else if startToNow <= bucketHourRetention {
		minimumBucket = bucketHour
	} else if startToNow <= bucketDayRetention {
		minimumBucket = bucketDay
	} else {
		minimumBucket = bucketWeek
	}

	//结合可使用的最小桶，根据end-start时间间隔来得出合适的桶和趋势图时间间隔
	diff := endTime.Unix() - startTime.Unix()
	location, _ := time.LoadLocation(timeZone)
	if diff <= oneHour {

		switch minimumBucket {
		case bucketMinute:
			return startTime, "1m", "", bucketMinute
		case bucketHour:
			//start变成小时整
			newStart := formatStartTimeHour(startTime, location)
			return newStart, "1h", "", bucketHour
		case bucketDay:
			//start 变成一天整
			newStart := formatStartTimeDay(startTime, location)
			return newStart, "1d", "", bucketDay
		case bucketWeek:
			//将startTime往前顺延到周一
			newStart := formatStartTimeToMonday(startTime, location)
			return newStart, "1w", "-4d", bucketWeek
		}

	} else if diff <= oneDay {

		switch minimumBucket {
		case bucketMinute:
			offset := ""
			offsetTime := startTime.Minute() % 5
			if offsetTime != 0 {
				offset = fmt.Sprintf("%dm", offsetTime)
			}
			return startTime, "5m", offset, bucketMinute

		case bucketHour:
			newStart := formatStartTimeHour(startTime, location)
			return newStart, "1h", "", bucketHour
		case bucketDay:
			newStart := formatStartTimeDay(startTime, location)
			return newStart, "1d", "", bucketDay
		case bucketWeek:
			//将startTime往前顺延到周一
			newStart := formatStartTimeToMonday(startTime, location)
			return newStart, "1w", "-4d", bucketWeek
		}

	} else if diff <= tenDay {

		switch minimumBucket {
		case bucketMinute, bucketHour:
			newStart := formatStartTimeHour(startTime, location)
			return newStart, "1h", "", bucketHour
		case bucketDay:
			newStart := formatStartTimeDay(startTime, location)
			return newStart, "1d", "", bucketDay
		case bucketWeek:
			//将startTime往前顺延到周一
			newStart := formatStartTimeToMonday(startTime, location)
			return newStart, "1w", "-4d", bucketWeek
		}

	} else if diff < oneYear {

		switch minimumBucket {
		case bucketMinute, bucketHour, bucketDay:
			newStart := formatStartTimeDay(startTime, location)
			return newStart, "1d", "", bucketDay
		case bucketWeek:
			//将startTime往前顺延到周一
			newStart := formatStartTimeToMonday(startTime, location)
			return newStart, "1w", "-4d", bucketWeek
		}

	}

	//end-start大于1年 时间间隔为1周
	//将startTime往前顺延到周一
	newStart := formatStartTimeToMonday(startTime, location)
	return newStart, "1w", "-4d", bucketWeek
}

// formatStartTimeHour 将time格式化为小时整
func formatStartTimeHour(t time.Time, location *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, location)
}

// formatStartTimeDay 将time格式化为天整
func formatStartTimeDay(t time.Time, location *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, location)
}

// formatStartTimeDay 将startTime向前移到周一，因为week桶里的time是每个周一才有数据
func formatStartTimeToMonday(t time.Time, location *time.Location) time.Time {
	var dayDiff int
	switch t.Weekday() {
	case time.Monday:
		dayDiff = 0
	case time.Tuesday:
		dayDiff = -1
	case time.Wednesday:
		dayDiff = -2
	case time.Thursday:
		dayDiff = -3
	case time.Friday:
		dayDiff = -4
	case time.Saturday:
		dayDiff = -5
	case time.Sunday:
		dayDiff = -6
	}

	return time.Date(t.Year(), t.Month(), t.Day()+dayDiff, 0, 0, 0, 0, location)
}

//
//func (m *monitorStatistics) formatTrendResult(result []map[string]interface{}) ([]time.Time, map[string][]interface{}) {
//	dates := make([]time.Time, 0)
//	values := make(map[string][]interface{})
//	for _, maps := range result {
//
//		isProxyRate := false
//		proxyTotal := int64(0)
//		if v, ok := maps["proxy_total"]; ok {
//			proxyTotal = common.FmtIntFromInterface(v)
//			values["proxy_total"] = append(values["proxy_total"], proxyTotal)
//			isProxyRate = true
//		}
//		proxySuccessTotal := int64(0)
//		if v, ok := maps["proxy_success_total"]; ok {
//			proxySuccessTotal = common.FmtIntFromInterface(v)
//			values["proxy_success_total"] = append(values["proxy_success_total"], proxySuccessTotal)
//			isProxyRate = true
//		}
//		proxyRate := 0.0
//		if proxySuccessTotal == 0 {
//			proxyRate = 0.0
//		} else {
//			proxyRate, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", float64(proxyTotal)/float64(proxySuccessTotal)), 64)
//		}
//
//		if isProxyRate {
//			values["proxy_rate"] = append(values["proxy_rate"], proxyRate)
//		}
//
//		isRequestRate := false
//		requestTotal := int64(0)
//		if v, ok := maps["req_total"]; ok {
//			requestTotal = common.FmtIntFromInterface(v)
//			values["request_total"] = append(values["request_total"], requestTotal)
//			isRequestRate = true
//		}
//		requestSuccessTotal := int64(0)
//		if v, ok := maps["req_success_total"]; ok {
//			requestSuccessTotal = common.FmtIntFromInterface(v)
//			values["request_success_total"] = append(values["request_success_total"], requestSuccessTotal)
//			isRequestRate = true
//		}
//		requestRate := 0.0
//		if requestSuccessTotal == 0 {
//			requestRate = 0.0
//		} else {
//			requestRate, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", float64(requestTotal)/float64(requestSuccessTotal)), 64)
//		}
//
//		if isRequestRate {
//			values["request_rate"] = append(values["request_rate"], requestRate)
//		}
//
//		if v, ok := maps["status_4xx"]; ok {
//			values["status_4xx"] = append(values["status_4xx"], v)
//		}
//
//		if v, ok := maps["status_5xx"]; ok {
//			values["status_5xx"] = append(values["status_5xx"], v)
//		}
//		t, _ := maps["_time"].(time.Time)
//		dates = append(dates, t)
//	}
//	return dates, values
//}
//
//func (m *monitorStatistics) formatResultOrdered(groupBy string, result []map[string]interface{}) []model.MonCommonStatistics {
//	resultList := make([]model.MonCommonStatistics, 0, len(result))
//	for _, val := range result {
//		avgResp := common.FmtFloatFromInterface(val["avg_resp_timing"])
//		maxResp := common.FmtIntFromInterface(val["max_resp_timing"])
//		minResp := common.FmtIntFromInterface(val["min_resp_timing"])
//		avgTraffic := common.FmtFloatFromInterface(val["avg_req_mess"])
//		maxTraffic := common.FmtIntFromInterface(val["max_req_mess"])
//		minTraffic := common.FmtIntFromInterface(val["min_req_mess"])
//
//		res := model.MonCommonStatistics{
//			MonCommonData: model.MonCommonData{
//				AvgResp:    avgResp,
//				MaxResp:    maxResp,
//				MinResp:    minResp,
//				AvgTraffic: avgTraffic,
//				MaxTraffic: maxTraffic,
//				MinTraffic: minTraffic,
//			},
//		}
//		switch groupBy {
//		case "api", "app":
//			if groupBy == "api" {
//				res.ApiId, _ = val["api"].(string)
//			} else if groupBy == "app" {
//				res.AppId, _ = val["app"].(string)
//			}
//			res.RequestTotal = common.FmtIntFromInterface(val["request_total"])
//			res.RequestSuccess = common.FmtIntFromInterface(val["request_success_total"])
//			res.RequestRate, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", float64(res.RequestSuccess)/float64(res.RequestTotal)), 64)
//
//			res.ProxyTotal = common.FmtIntFromInterface(val["proxy_total"])
//			res.ProxySuccess = common.FmtIntFromInterface(val["proxy_success_total"])
//			if res.ProxyTotal == 0 {
//				res.ProxyRate = 0.0
//			} else {
//				res.ProxyRate, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", float64(res.ProxySuccess)/float64(res.ProxyTotal)), 64)
//			}
//			res.StatusFail = res.RequestTotal - res.RequestSuccess
//		case "upstream":
//			res.ServiceName, _ = val["upstream"].(string)
//			res.ProxyTotal = common.FmtIntFromInterface(val["proxy_total"])
//			res.ProxySuccess = common.FmtIntFromInterface(val["proxy_success_total"])
//			if res.ProxyTotal == 0 {
//				res.ProxyRate = 0.0
//			} else {
//				res.ProxyRate, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", float64(res.ProxySuccess)/float64(res.ProxyTotal)), 64)
//			}
//			res.StatusFail = res.ProxyTotal - res.ProxySuccess
//		case "ip":
//			res.Ip, _ = val["ip"].(string)
//			res.RequestTotal = common.FmtIntFromInterface(val["request_total"])
//			res.RequestSuccess = common.FmtIntFromInterface(val["request_success_total"])
//			res.RequestRate, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", float64(res.RequestSuccess)/float64(res.RequestTotal)), 64)
//			res.StatusFail = res.RequestTotal - res.RequestSuccess
//		}
//		resultList = append(resultList, res)
//	}
//	return resultList
//}
