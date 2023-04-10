package warn_timer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/api"
	api_model "github.com/eolinker/apinto-dashboard/modules/api/model"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-model"
	"github.com/eolinker/apinto-dashboard/modules/monitor"
	"github.com/eolinker/apinto-dashboard/modules/monitor/model"
	"github.com/eolinker/apinto-dashboard/modules/namespace"
	"github.com/eolinker/apinto-dashboard/modules/notice"
	"github.com/eolinker/apinto-dashboard/modules/notice/notice-model"
	"github.com/eolinker/apinto-dashboard/modules/upstream"
	"github.com/eolinker/apinto-dashboard/modules/upstream/model"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/apinto-dashboard/modules/user/user-model"
	"github.com/eolinker/apinto-dashboard/modules/warn"
	"github.com/eolinker/apinto-dashboard/modules/warn/warn-model"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/go-basic/uuid"
	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/errgroup"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type IMonitorWarn interface {
	monitorWarn()
}

type monitorWarn struct {
	warnStrategyService  warn.IWarnStrategyService
	monitorStatistics    monitor.IMonitorStatistics
	namespaceService     namespace.INamespaceService
	monitorService       monitor.IMonitorService
	userService          user.IUserInfoService
	clusterService       cluster.IClusterService
	warnHistoryService   warn.IWarnHistoryService
	apiService           api.IAPIService
	service              upstream.IService
	commonCache          cache.ICommonCache
	noticeChannelService notice.INoticeChannelService
	noticeChannelDriver  notice.INoticeChannelDriverManager
}

func newMonitorWarn() IMonitorWarn {
	mon := &monitorWarn{}
	bean.Autowired(&mon.warnStrategyService)
	bean.Autowired(&mon.monitorStatistics)
	bean.Autowired(&mon.namespaceService)
	bean.Autowired(&mon.userService)
	bean.Autowired(&mon.monitorService)
	bean.Autowired(&mon.warnHistoryService)
	bean.Autowired(&mon.apiService)
	bean.Autowired(&mon.noticeChannelDriver)
	bean.Autowired(&mon.service)
	bean.Autowired(&mon.clusterService)
	bean.Autowired(&mon.noticeChannelService)
	bean.Autowired(&mon.commonCache)
	return mon
}

func (mon *monitorWarn) monitorWarn() {
	ctx := context.Background()
	//当前时间前一分钟
	t := time.Now().Add(-time.Minute)
	endTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)

	namespaces, err := mon.namespaceService.GetAll()
	if err != nil {
		log.Errorf("monitorWarn-monitorWarn namespaceService.GetAll error:%s", err.Error())
		return
	}

	//获取所有用户信息
	userInfos, err := mon.userService.GetUserInfoMaps(ctx)
	if err != nil {
		log.Errorf("monitorWarn-monitorWarn userService.userInfos error:%s", err.Error())
		return
	}

	for _, ns := range namespaces {
		//1.获取所有告警策略
		strategiesAll, err := mon.warnStrategyService.WarnStrategyAll(ctx, ns.Id, 1)
		if err != nil {
			log.Errorf("monitorWarn-monitorWarn WarnStrategyAll error:%s", err.Error())
			return
		}

		//以分区ID+统计粒度+指标+统计时间为分组生成策略
		group := common.SliceToMapArray(strategiesAll, func(t *warn_model.WarnStrategy) string {
			return fmt.Sprintf("%s:%s:%s:%d", t.PartitionUUID, t.Dimension, t.Quota, t.Every)
		})

		//命名空间下的所有api/service/noticeChannel信息
		apiMaps, serviceMaps, noticeChannelMaps, err := mon.getSourceInfo(ctx, ns.Id)
		if err != nil {
			log.Errorf("monitorWarn-monitorWarn getSourceInfo error:%s", err.Error())
			return
		}

		for key, strategies := range group {
			mon.task(ctx, ns.Id, key, endTime, strategies, userInfos, apiMaps, serviceMaps, noticeChannelMaps)
		}
	}

}

func (mon *monitorWarn) task(ctx context.Context, namespaceId int, key string, endTime time.Time, strategies []*warn_model.WarnStrategy, userMaps map[int]*user_model.UserInfo, apiMaps map[string]*api_model.APIInfo, serviceMaps map[string]*upstream_model.ServiceListItem, noticeChannelMaps map[string]*notice_model.NoticeChannel) {

	//key+时间戳
	lockKey := fmt.Sprintf("%s_%d", key, endTime.Unix())
	if err := mon.lock(ctx, lockKey); err != nil {
		log.Warnf("lockKey=%s error=%s", lockKey, err.Error())
		return
	}
	//因为key是由时间戳组成，所以不需要释放锁
	//defer mon.unlock(ctx, lockKey)

	//拿到维度、指标和周期获取指标值
	split := strings.Split(key, ":")
	if len(split) < 4 {
		log.Errorf("split 长度不对")
		return
	}
	everyStr := split[3]
	every, _ := strconv.Atoi(everyStr)
	if every == 0 {
		log.Errorf("every is 0")
		return

	}
	// 开始时间为当前时间的前一分钟+统计的时间粒度
	startTime := endTime.Add(-time.Minute * time.Duration(every))

	groupStr := split[1]
	if groupStr == warn_model.DimensionTypePartition {
		//维度为分区是查询该分区下集群的聚合数据
		groupStr = warn_model.DimensionTypeCluster
	}

	if groupStr == warn_model.DimensionTypeService {
		//apinto上游对应的字符串为upstream
		groupStr = "upstream"
	}

	partitionUUID := split[0]
	quotaType := warn_model.QuotaType(split[2])

	//查询监控数据
	statistics, err := mon.monitorStatistics.WarnStatistics(ctx, namespaceId, partitionUUID, startTime, endTime, groupStr, quotaType, nil)
	if err != nil {
		log.Errorf("monitorWarn-WarnStatistics error:%s", err.Error())
		return
	}
	log.DebugF("partitionUUID=%s startTime=%v endTime=%v groupStr=%s quotaType=%s statistics=%v", partitionUUID, startTime, endTime, groupStr, quotaType, statistics)

	//获取当前分区下的集群信息
	clusterNameMaps, clusterUuidMaps, err := mon.getClustersByPartitionUuid(ctx, namespaceId, partitionUUID)
	if err != nil {
		log.Errorf("monitorWarn-getClustersByPartitionUuid error:%s", err.Error())
		return
	}

	for _, strategy := range strategies {

		//获取用户的邮箱和通知渠道ID
		userEmailStr := make([]string, 0)
		noticeUserId := make([]string, 0)
		for _, userId := range strategy.WarnStrategyConfig.Users {
			if u, ok := userMaps[userId]; ok {
				if len(strings.TrimSpace(u.Email)) > 0 {
					userEmailStr = append(userEmailStr, u.Email)
				}

				if len(strings.TrimSpace(u.NoticeUserId)) > 0 {
					noticeUserId = append(noticeUserId, u.NoticeUserId)
				}
			}
		}

		target := strategy.WarnStrategyConfig.Target
		//获取实际的告警目录
		values := mon.getRealTargetValues(apiMaps, serviceMaps, clusterNameMaps, strategy.Dimension, target.Rule, target.Values)

		for _, rule := range strategy.WarnStrategyConfig.Rule {
			warnList := make([]warn_model.NoticeChannelWarn, 0)

			if strategy.Dimension == warn_model.DimensionTypePartition {
				//分区是把分区下的集群聚合查询
				targetValue := 0.0
				for _, value := range values {
					targetValue += statistics[value]
				}
				switch strategy.Quota {
				case warn_model.QuotaTypeReqFailRate, warn_model.QuotaTypeProxyFailRate:
					//请求失败率和转发失败率计算,因是计算的该分区下集群的聚合数据，所以需要求平均数
					targetValue *= 100
					targetValue = targetValue / float64(len(values))
				case warn_model.QuotaTypeAvgResp:
					//平均响应时间计算
					targetValue = targetValue / float64(len(values))
				}

				//计算告警指标和规则是否触发告警
				warns := mon.warnCount(ctx, namespaceId, startTime, endTime,
					targetValue, strategy, rule, values, apiMaps, clusterUuidMaps, statistics)
				if len(warns) > 0 {
					//需要告警的信息
					warnList = append(warnList, warns...)
				}

			} else {
				for _, value := range values {
					targetValue := statistics[value]
					switch strategy.Quota {
					case warn_model.QuotaTypeReqFailRate, warn_model.QuotaTypeProxyFailRate:
						//请求失败率和转发失败率计算
						targetValue *= 100
					}

					//计算告警指标和规则是否触发告警
					warns := mon.warnCount(ctx, namespaceId, startTime, endTime,
						targetValue, strategy, rule, []string{value}, apiMaps, clusterUuidMaps, statistics)
					if len(warns) > 0 {
						//需要告警的信息
						warnList = append(warnList, warns...)
					}
				}
			}

			//warnList > 0 说明触发告警策略了
			if len(warnList) > 0 {

				historyInfo := getWarnHistoryInfo(strategy, warnList, endTime)
				//是否需要发送告警消息，不需要直接跳过
				if !mon.isSendTo(ctx, endTime, strategy) {
					//设置告警但未发送状态
					mon.setWarnMinuteStatus(ctx, endTime, strategy, warn_model.WarnStatusTrigger)
					//告警历史 插入告警历史
					_ = mon.warnHistoryService.Create(ctx, namespaceId, strategy.PartitionId, historyInfo)
					continue
				}

				//写入告警次数到缓存
				mon.setWarnNum(ctx, endTime, strategy)
				mon.setWarnMinuteStatus(ctx, endTime, strategy, warn_model.WarnStatusSendTrigger)

				//发送失败的次数和需要发送的次数做对比
				var sendFail = new(int64)

				noticeErrGroup, _ := errgroup.WithContext(ctx)
				sendMsgErrors := make([]*warn_model.SendMsgError, 0)

				for _, uid := range rule.ChannelUuids {
					channelUuid := uid

					//利用协程快速发送通知消息
					noticeErrGroup.Go(func() error {
						noticeChannelDriver := mon.noticeChannelDriver.GetDriver(channelUuid)
						if noticeChannelDriver == nil {
							log.Errorf("获取不到通知渠道 渠道uuid：%s", channelUuid)
							return errors.New("渠道通知获取失败")
						}
						sendMsgErrorUuid := uuid.New()
						sendMsgError := &warn_model.SendMsgError{
							UUID:              sendMsgErrorUuid,
							NoticeChannelUUID: channelUuid,
						}

						if channel, ok := noticeChannelMaps[channelUuid]; ok {
							sends := make([]string, 0)
							msg := ""
							if channel.Type == 2 {
								if len(userEmailStr) == 0 {
									atomic.AddInt64(sendFail, 1)
									sendMsgError.Msg = "收件人邮箱为空"
									sendMsgErrors = append(sendMsgErrors, sendMsgError)
									return errors.New(sendMsgError.Msg)
								}
								//获取邮箱msg
								msg = formatWarnEmailMsg(strategy.Title, strategy.Dimension, endTime, warnList)
								sends = userEmailStr
							} else {
								if len(noticeUserId) == 0 {
									atomic.AddInt64(sendFail, 1)
									sendMsgError.Msg = "通知用户ID为空"
									sendMsgErrors = append(sendMsgErrors, sendMsgError)
									return errors.New(sendMsgError.Msg)
								}
								msg = formatWarnWebhookMsg(strategy.Title, strategy.Dimension, endTime, warnList)
								sends = noticeUserId
							}
							if err = noticeChannelDriver.SendTo(sends, strategy.Title, msg); err != nil {
								sendMsgError.Msg = err.Error()
								sendMsgErrors = append(sendMsgErrors, sendMsgError)
								log.Errorf("告警消息发送失败 sendMsgErrorUuid=%s channelUuid=%s users=%v err=%s", sendMsgErrorUuid, channelUuid, sends, err.Error())
								atomic.AddInt64(sendFail, 1)
								return err
							}
						}

						return nil
					})

				}

				sendStatus := 1
				//发生错误 发送失败
				if err = noticeErrGroup.Wait(); err != nil {
					//部分发送成功
					sendStatus = 3
					//全部发送失败
					if *sendFail == int64(len(rule.ChannelUuids)) {
						sendStatus = 2
					}
				}
				errMsg, _ := json.Marshal(sendMsgErrors)

				historyInfo.Status = sendStatus
				historyInfo.ErrMsg = string(errMsg)
				//告警历史 插入告警历史
				_ = mon.warnHistoryService.Create(ctx, namespaceId, strategy.PartitionId, historyInfo)

			} else {
				mon.setWarnMinuteStatus(ctx, endTime, strategy, warn_model.WarnStatusNotTrigger)
			}
		}

	}

	return
}

func getWarnHistoryInfo(strategy *warn_model.WarnStrategy, warnList []warn_model.NoticeChannelWarn, t time.Time) *warn_model.WarnHistoryInfo {
	historyTargets := make([]string, 0)
	contentMaps := make(map[string]string)
	for _, warnInfo := range warnList {
		contentMaps[warnInfo.Name] = getContent(warnInfo, false)
		historyTargets = append(historyTargets, warnInfo.Name)
	}

	content, _ := json.Marshal(contentMaps)
	content, _ = unescapeUnicode(content)

	historyInfo := &warn_model.WarnHistoryInfo{
		StrategyTitle: strategy.Title,
		Dimension:     strategy.Dimension,
		Quota:         string(strategy.Quota),
		Target:        strings.Join(historyTargets, ","),
		Content:       string(content),
		CreateTime:    t,
	}

	return historyInfo
}

// 告警计算 判断是否触发告警
func (mon *monitorWarn) warnCount(ctx context.Context, namespaceId int, startTime, endTime time.Time, targetValue float64,
	strategy *warn_model.WarnStrategy, rule *warn_model.WarnStrategyConfigRule, values []string, apiMaps map[string]*api_model.APIInfo, clusterMaps map[string]*cluster_model.Cluster, statistics map[string]float64) []warn_model.NoticeChannelWarn {
	warnCount := 0

	ratio := 0.0
	yearBasis := 0.0

	isQueryRing := false
	isQueryBasis := false
	for _, condition := range rule.Condition {
		switch condition.Compare {
		case "ring_ratio_add", "ring_ratio_reduce": //环比增加
			isQueryRing = true
		case "year_basis_add", "year_basis_reduce": //同比增加
			isQueryBasis = true
		}
	}
	log.DebugF("monitorWarn-warnCount targetValue=%v strategy=%v values=%v statistics=%v", targetValue, strategy, values, statistics)

	//查询同比或环比的数据
	if isQueryRing {
		ratio, _ = mon.getCompare(ctx, namespaceId, values, startTime, endTime, strategy, true)
		switch strategy.Quota {
		case warn_model.QuotaTypeReqFailRate, warn_model.QuotaTypeProxyFailRate:
			ratio *= 100
		}

	}
	if isQueryBasis {
		yearBasis, _ = mon.getCompare(ctx, namespaceId, values, startTime, endTime, strategy, false)
		switch strategy.Quota {
		case warn_model.QuotaTypeReqFailRate, warn_model.QuotaTypeProxyFailRate:
			yearBasis *= 100
		}
	}

	msgConditions := make([]*warn_model.MsgCondition, 0)
	for _, condition := range rule.Condition {
		realValue := 0.0
		switch condition.Compare {
		case ">":
			realValue = targetValue
			if targetValue > condition.Value {
				warnCount++
			}
		case ">=":
			realValue = targetValue
			if targetValue >= condition.Value {
				warnCount++
			}
		case "<":
			realValue = targetValue
			if targetValue < condition.Value {
				warnCount++
			}
		case "<=":
			realValue = targetValue
			if targetValue <= condition.Value {
				warnCount++
			}
		case "==":
			realValue = targetValue
			if condition.Value == targetValue {
				warnCount++
			}
		case "!=":
			realValue = targetValue
			if condition.Value != targetValue {
				warnCount++
			}
		case "ring_ratio_add":
			if ratio == 0 {
				continue
			}
			f := ((targetValue - ratio) / ratio) * 100
			realValue = f
			if f > condition.Value {
				warnCount++
			}
		case "ring_ratio_reduce":
			if ratio == 0 {
				continue
			}
			f := ((ratio - targetValue) / ratio) * 100
			realValue = f
			if f < condition.Value {
				warnCount++
			}
		case "year_basis_add": //同比增加
			if yearBasis > 0 {
				continue
			}
			f := ((targetValue - yearBasis) / yearBasis) * 100
			realValue = f
			if f > condition.Value {
				warnCount++
			}
		case "year_basis_reduce": //同比减少
			if yearBasis > 0 {
				continue
			}
			f := ((targetValue - yearBasis) / yearBasis) * 100
			realValue = f
			if f < condition.Value {
				warnCount++
			}
		}
		msgConditions = append(msgConditions, &warn_model.MsgCondition{
			RealValue: realValue,
			Compare:   condition.Compare,
			Unit:      condition.Unit,
			Value:     condition.Value,
		})
	}

	if isQueryBasis || isQueryRing {
		go func() {
			redisKey := fmt.Sprintf("%s:%s:%s:%d:%d:%d", strategy.PartitionUUID, strategy.Dimension, strategy.Quota, strategy.Every, startTime.Unix(), endTime.Unix())
			redisValue := make(map[string][]byte, 0)
			for _, key := range values {
				if val, ok := statistics[key]; ok {
					redisValue[key] = []byte(strconv.FormatFloat(val, 'g', -1, 64))
				}
			}

			if err := mon.commonCache.HMSet(ctx, redisKey, redisValue, time.Hour*25); err != nil {
				log.Errorf("环比同比数据插入redis失败 error=%s", err.Error())
			}
		}()

	}

	warnList := make([]warn_model.NoticeChannelWarn, 0)
	if warnCount == len(rule.Condition) {
		url := ""
		name := ""

		if strategy.Dimension == warn_model.DimensionTypeApi {
			if apiInfo, ok := apiMaps[values[0]]; ok {
				url = apiInfo.RequestPathLabel
				name = apiInfo.Name
			}
		} else if strategy.Dimension == warn_model.DimensionTypePartition {
			clusters := make([]string, 0)
			for _, value := range values {
				if clusterInfo, ok := clusterMaps[value]; ok {
					clusters = append(clusters, clusterInfo.Name)
				}
			}
			name = strings.Join(clusters, ",")
		} else if strategy.Dimension == warn_model.DimensionTypeCluster {
			if clusterInfo, ok := clusterMaps[values[0]]; ok {
				name = clusterInfo.Name
			} else {
				name = "未知集群名"
			}

		} else {
			name = values[0]
		}

		warnList = append(warnList, warn_model.NoticeChannelWarn{
			Url:       url,
			Name:      name,
			Every:     strategy.Every,
			Quota:     strategy.Quota,
			Condition: msgConditions,
		})

	}
	return warnList
}

func (mon *monitorWarn) lock(ctx context.Context, key string) error {
	log.DebugF("monitorWarn-lock key=%s", key)
	b, err := mon.commonCache.SetNX(ctx, key, "1", time.Minute)
	if err != nil {
		return err
	}
	if b {
		return nil
	}
	return errors.New("锁已被占用")
}

func (mon *monitorWarn) unlock(ctx context.Context, key string) {
	log.DebugF("monitorWarn-unlock key=%s", key)
	_ = mon.commonCache.Del(ctx, key)
}

func getHourMaxKey(uuid string, t time.Time) string {
	return fmt.Sprintf("warn_notice_hour_limit:%s:%s", uuid, t.Format("2006-01-02"))
}

func getContinuityKey(uuid string, t time.Time) string {
	return fmt.Sprintf("warn_notice_minute_num:%s:%s", uuid, t.Add(-time.Minute).Format("2006-01-02 15:04"))
}

// 是否需要发送告警消息
func (mon *monitorWarn) isSendTo(ctx context.Context, t time.Time, strategy *warn_model.WarnStrategy) bool {
	hourMax := strategy.WarnStrategyConfig.HourMax
	if hourMax > 0 {
		key := getHourMaxKey(strategy.Uuid, t)
		//是否发送告警通知
		num, err := mon.commonCache.GetInt(ctx, key)
		if err != nil && err != redis.Nil {
			log.Errorf("isSendTo-commonCache.GetInt err=%s", err.Error())
			return false
		}

		if int(num) >= strategy.WarnStrategyConfig.HourMax {
			log.Infof("当前策略达到小时最大发送次数 策略UUID=%s 已发送次数=%d  最大限制次数=%s", strategy.Uuid, num, hourMax)
			return false
		}
	}

	continuity := strategy.WarnStrategyConfig.Continuity
	if continuity > 0 {
		//5
		//54

		//53
		//y

		//53
		//y

		//53
		//t

		//51 52 53 54 55
		//f  y  t  t  t
		//51 52 53 54 55
		//t  y  t  t  t

		//
		//先查询前一分钟告警状态
		frontMinuteKey := getContinuityKey(strategy.Uuid, t.Add(-time.Minute))
		frontMinuteStatus, err := mon.commonCache.Get(ctx, frontMinuteKey)
		if err != nil && err != redis.Nil {
			log.Errorf("isSendTo-commonCache.Get err=%s", err.Error())
			return false
		}

		switch string(frontMinuteStatus) {
		case warn_model.WarnStatusNotTrigger:
			//没触发则表示当前分钟可以触发告警
			return true
		case warn_model.WarnStatusTrigger:
			if continuity == 1 { //如果只配了1分钟 那么直接返回true 可以触发告警
				return true
			}
			//这种情况要判断前 continuity 分钟的数据是否都是WarnStatusTrigger
			for i := continuity; i > 0; i-- {
				if i == 1 {
					continue
				}

				key := getContinuityKey(strategy.Uuid, t.Add(-time.Minute*time.Duration(i)))
				status, err := mon.commonCache.Get(ctx, key)
				if err != nil && err != redis.Nil {
					log.Errorf("isSendTo-commonCache.Get err=%s", err.Error())
					return false
				}
				if string(status) == warn_model.WarnStatusNotTrigger || string(status) == warn_model.WarnStatusSendTrigger {
					return false
				}

			}
		case warn_model.WarnStatusSendTrigger:
			//已触发则表示当前分钟不可以再触发告警了
			return false
		}

	}

	return true
}

func (mon *monitorWarn) setWarnNum(ctx context.Context, t time.Time, strategy *warn_model.WarnStrategy) {
	if strategy.WarnStrategyConfig.HourMax <= 0 {
		return
	}
	hourKey := getHourMaxKey(strategy.Uuid, t)
	err := mon.commonCache.IncrBy(ctx, hourKey, 1, time.Minute*65)
	if err != nil {
		log.Errorf("告警次数写入缓存失败 err=%s", err.Error())
	}

}

func (mon *monitorWarn) setWarnMinuteStatus(ctx context.Context, t time.Time, strategy *warn_model.WarnStrategy, status string) {
	if strategy.WarnStrategyConfig.Continuity <= 0 {
		return
	}
	minuteKey := getContinuityKey(strategy.Uuid, t)
	if err := mon.commonCache.Set(ctx, minuteKey, []byte(status), time.Minute*60); err != nil {
		log.Errorf("告警次数写入缓存失败 err=%s", err.Error())
	}
}

// 获取%比较数据 isRingRatio==true查询环比数据  否则查询同比数据
func (mon *monitorWarn) getCompare(ctx context.Context, namespaceId int, values []string, startTime, endTime time.Time, strategy *warn_model.WarnStrategy, isRingRatio bool) (float64, error) {

	//环比取上个时间段的数据
	if isRingRatio {
		startTime = startTime.Add(-time.Minute * time.Duration(strategy.Every))
		endTime = endTime.Add(-time.Minute * time.Duration(strategy.Every))
	} else {
		//环比取昨天同时间段的数据
		startTime = startTime.Add(-time.Hour * 24)
		endTime = endTime.Add(-time.Hour * 24)
	}

	redisKey := fmt.Sprintf("%s:%s:%s:%d:%d:%d", strategy.PartitionUUID, strategy.Dimension, strategy.Quota, strategy.Every, startTime.Unix(), endTime.Unix())

	maps, _ := mon.commonCache.HGetAll(ctx, redisKey)
	if len(maps) == 0 { //重新查一次

		wheres := make([]model.MonWhereItem, 0)

		group := strategy.Dimension
		whereKey := strategy.Dimension
		//分区查该分区下的集群 所有group为cluster
		if strategy.Dimension == warn_model.DimensionTypePartition {
			group = warn_model.DimensionTypeCluster
			whereKey = warn_model.DimensionTypeCluster
		} else if strategy.Dimension == warn_model.DimensionTypeService {
			whereKey = "upstream"
		}

		wheres = append(wheres, model.MonWhereItem{
			Key:       whereKey,
			Operation: "in",
			Values:    values,
		})

		statistics, err := mon.monitorStatistics.WarnStatistics(ctx, namespaceId, strategy.PartitionUUID, startTime, endTime, group, strategy.Quota, wheres)
		if err != nil {
			return 0, err
		}

		if strategy.Dimension == warn_model.DimensionTypePartition {
			targetValue := 0.0
			count := 0
			for _, value := range values {
				if v, ok := statistics[value]; ok {
					targetValue += v
					count++
				}

			}
			switch strategy.Quota {
			case warn_model.QuotaTypeReqFailRate, warn_model.QuotaTypeProxyFailRate, warn_model.QuotaTypeAvgResp:
				targetValue = targetValue / float64(count)
			}
			return targetValue, nil
		}

		return statistics[values[0]], nil
	}

	if strategy.Dimension == warn_model.DimensionTypePartition {
		for _, value := range values {

			count := 0
			targetValue := 0.0
			if v, ok := maps[value]; ok {
				result, _ := strconv.ParseFloat(v, 64)
				targetValue += result
				count++
			}

			switch strategy.Quota {
			case warn_model.QuotaTypeReqFailRate, warn_model.QuotaTypeProxyFailRate, warn_model.QuotaTypeAvgResp:
				targetValue = targetValue / float64(count)
			}
			return targetValue, nil
		}
	} else {
		if v, ok := maps[values[0]]; ok {
			result, _ := strconv.ParseFloat(v, 64)
			return result, nil
		}
	}

	return 0, nil
}

func (mon *monitorWarn) getClustersByPartitionUuid(ctx context.Context, namespaceId int, partitionUuid string) (clusterNameMaps map[string]*cluster_model.Cluster, clusterUuidMaps map[string]*cluster_model.Cluster, err error) {
	partitionInfo, err := mon.monitorService.PartitionInfo(ctx, namespaceId, partitionUuid)
	if err != nil {
		log.Errorf("monitorWarn-monitorService.PartitionInfo err=%s", err.Error())
		return nil, nil, err
	}

	clusters, err := mon.clusterService.GetByNames(ctx, namespaceId, partitionInfo.ClusterNames)
	if err != nil {
		log.Errorf("monitorWarn-clusterService.GetByNames err=%s", err.Error())
		return nil, nil, err
	}

	clusterNameMaps = common.SliceToMap(clusters, func(t *cluster_model.Cluster) string {
		return t.Name
	})

	clusterUuidMaps = common.SliceToMap(clusters, func(t *cluster_model.Cluster) string {
		return t.UUID
	})

	return
}

func (mon *monitorWarn) getSourceInfo(ctx context.Context, namespaceId int) (apiMaps map[string]*api_model.APIInfo, serviceMaps map[string]*upstream_model.ServiceListItem, noticeChannelMap map[string]*notice_model.NoticeChannel, err error) {
	errGroup, _ := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		apiInfos, err := mon.apiService.GetAPIInfoAll(ctx, namespaceId)
		if err != nil {
			log.Errorf("monitorWarn-apiService.GetAPIInfoAll err=%s", err.Error())
			return err
		}
		apiMaps = common.SliceToMap(apiInfos, func(t *api_model.APIInfo) string {
			return t.UUID
		})
		return nil
	})
	errGroup.Go(func() error {
		serviceListAll, err := mon.service.GetServiceListAll(ctx, namespaceId)
		if err != nil {
			log.Errorf("monitorWarn-service.GetServiceListAll err=%s", err.Error())
			return err
		}
		serviceMaps = common.SliceToMap(serviceListAll, func(t *upstream_model.ServiceListItem) string {
			return t.Name
		})
		return nil
	})

	errGroup.Go(func() error {
		channelList, err := mon.noticeChannelService.NoticeChannelList(ctx, namespaceId, 0)
		if err != nil {
			log.Errorf("monitorWarn-noticeChannelService.NoticeChannelList err=%s", err.Error())
			return err
		}

		noticeChannelMap = common.SliceToMap(channelList, func(t *notice_model.NoticeChannel) string {
			return t.Name
		})

		return nil
	})

	if err = errGroup.Wait(); err != nil {
		return nil, nil, nil, err
	}

	return apiMaps, serviceMaps, noticeChannelMap, nil
}

// 获取实际的告警目标
func (mon *monitorWarn) getRealTargetValues(apiMaps map[string]*api_model.APIInfo, serviceMaps map[string]*upstream_model.ServiceListItem, clusterMaps map[string]*cluster_model.Cluster, dimension, rule string, oldValues []string) []string {
	values := make([]string, 0)
	switch dimension {
	case warn_model.DimensionTypeApi:
		switch rule {
		case warn_model.RuleTypeUnlimited: //不限（查询所有的）
			for uid := range apiMaps {
				values = append(values, uid)
			}
			return values
		case warn_model.RuleTypeContain: //包含
			return oldValues
		case warn_model.RuleTypeNotContain: //不包含
			tempMaps := common.CopyMaps(apiMaps)
			for _, value := range oldValues {
				if _, ok := tempMaps[value]; ok {
					delete(tempMaps, value)
				}
			}
			for uid := range tempMaps {
				values = append(values, uid)
			}
			return values
		}
	case warn_model.DimensionTypeService:
		switch rule {
		case warn_model.RuleTypeUnlimited: //不限（查询所有的）
			for uid := range serviceMaps {
				values = append(values, uid)
			}
			return values
		case warn_model.RuleTypeContain: //包含
			return oldValues
		case warn_model.RuleTypeNotContain: //不包含
			tempMaps := common.CopyMaps(serviceMaps)
			for _, value := range oldValues {
				if _, ok := tempMaps[value]; ok {
					delete(tempMaps, value)
				}
			}
			for name := range tempMaps {
				values = append(values, name)
			}
			return values
		}
	case warn_model.DimensionTypeCluster:
		for _, value := range oldValues {
			if clusterInfo, ok := clusterMaps[value]; ok {
				values = append(values, clusterInfo.UUID)
			}
		}
		return values
	case warn_model.DimensionTypePartition:
		for _, clusterInfo := range clusterMaps {
			values = append(values, clusterInfo.UUID)
		}
		return values
	}
	return nil
}

// 格式化告警信息成字符串
func getContent(val warn_model.NoticeChannelWarn, isHtml bool) string {
	warnFrequency := fmt.Sprintf("%d分钟", val.Every)
	if val.Every == 60 {
		warnFrequency = "1小时"
	}

	contents := make([]string, 0)
	for i, condition := range val.Condition {

		unitStr := ""
		switch condition.Unit {
		case "num":
			unitStr = "次"
		case "%":
			unitStr = "%"
		case "ms":
			unitStr = "ms"
		case "kb":
			unitStr = "kb"
		}
		compare := warn_model.CompareValue[condition.Compare]
		quota := warn_model.QuotaRuleMap[val.Quota]
		conditionValue := common.FloatToString(condition.Value)
		realValue := common.FloatToString(condition.RealValue)
		if i == 0 {
			if isHtml {
				content := fmt.Sprintf(`%s/统计粒度%v %s %v%s<span class="strategy-color"> 实际值：%v%s</span>`, quota, warnFrequency, compare,
					conditionValue, unitStr, realValue, unitStr)
				contents = append(contents, content)
			} else {
				content := fmt.Sprintf(`%s/统计粒度%s %s %v%s 实际值：%v%s`, quota, warnFrequency, compare,
					conditionValue, unitStr, realValue, unitStr)
				contents = append(contents, content)
			}

			continue
		}
		if isHtml {
			content := fmt.Sprintf(` %s %v%s<span class="strategy-color"> 实际值：%v%s</span>`, compare,
				conditionValue, unitStr, realValue, unitStr)
			contents = append(contents, content)
		} else {
			content := fmt.Sprintf(` %s %s%s 实际值：%s%s`, compare,
				conditionValue, unitStr, realValue, unitStr)
			contents = append(contents, content)
		}

	}
	content := strings.Join(contents, " 且")
	return content
}

func formatWarnEmailMsg(title, dimension string, t time.Time, list []warn_model.NoticeChannelWarn) string {

	thead := ``
	tbody := ``
	if dimension == warn_model.DimensionTypeApi {
		thead = fmt.Sprintf(`<thead>
          <tr>
            <th>告警策略名称</th>
            <th>%s名称</th>
            <th>接口URL</th>
            <th>告警内容</th>
            <th>告警时间</th>
          </tr>
        </thead>`, "API")
		tr := ``
		for _, val := range list {
			content := getContent(val, true)

			tr += fmt.Sprintf(`<tr>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>
                %s
            </td>
            <td>%s</td>
          </tr>`, title, val.Name, val.Url, content, common.TimeToStr(t))
		}
		tbody = fmt.Sprintf(`<tbody>
          <!-- 多行表格内容则循环tr标签-->
          %s
        </tbody>`, tr)
	} else {
		dimensionVal := ""
		if dimension == warn_model.DimensionTypeService {
			dimensionVal = "上游"
		} else if dimension == warn_model.DimensionTypeCluster || dimension == warn_model.DimensionTypePartition {
			dimensionVal = "集群"
		}
		thead = fmt.Sprintf(`<thead>
          <tr>
            <th>告警策略名称</th>
            <th>%s名称</th>
            <th>告警内容</th>
            <th>告警时间</th>
          </tr>
        </thead>`, dimensionVal)
		tr := ``
		for _, val := range list {

			content := getContent(val, true)

			tr += fmt.Sprintf(`<tr>
            <td>%s</td>
            <td>%s</td>
            <td>
                %s
            </td>
            <td>%s</td>
          </tr>`, title, val.Name, content, common.TimeToStr(t))
		}
		tbody = fmt.Sprintf(`<tbody>
          <!-- 多行表格内容则循环tr标签-->
          %s
        </tbody>`, tr)
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>Apinto Pro</title>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <link rel="icon" type="image/x-icon" href="favicon.ico" />
    <style type="text/css">
      body,
      .alert-content {
        margin: 0;
        display: flex;
        flex-direction: column;
        align-items: center;
      }

      .alert-content {
        top: 50px;
        overflow: auto;
      }

      p.email-title {
        color: rgba(0, 0, 0, 1);
        font-size: 24px;
        font-weight: 500;
        font-family: 'PingFang TC';
        text-align: center;
        line-height: 22px;
        margin-top: 14vh;
      }

      p.email-subtitle {
        color: rgba(51, 51, 51, 1);
        font-size: 14px;
        font-weight: 500;
        font-family: 'PingFang TC';
        text-align: center;
        line-height: 22px;
        margin-top: 32px;
      }

      table {
        width: 68%%;
        font-size: 14px;
        font-weight: 400;
        line-height: 22px;
        margin-top: 32px;
      }

      thead tr {
        color: #999999;
        background-color: #f2f2f2;
      }

      thead tr th {
        height: 40px;
        font-size: 14px;
        font-weight: 400;
        padding: 0 10px;
        border: 1px solid #e8e8e8;
        border-right: none;
        border-left: none;
        white-space: nowrap;
        text-align: left;
      }

      thead tr th:first-child {
        border-left: 1px solid #e8e8e8;
      }

      thead tr th:last-child {
        border-right: 1px solid #e8e8e8;
      }

      tbody tr td {
        min-height: 40px;
        height: 40px;
        color: #333;
        padding: 0 10px;
        line-height: 22px;
        border-bottom: 1px solid #e8e8e8;
      }

      tbody tr td:first-child {
        border-left: 1px solid #e8e8e8;
      }

      tbody tr td:last-child {
        border-right: 1px solid #e8e8e8;
      }

      .strategy-color {
        color: #f04864;
      }
    </style>
  </head>
  <body>
    <div class="alert-content">
      <p class="email-title">%s</p>
      <p class="email-subtitle">告警时间:%s</p>
      <table cellspacing="0">
        %s
        %s
      </table>
    </div>
  </body>
</html>
`, title, common.TimeToStr(t), thead, tbody)
	return html
}

func formatWarnWebhookMsg(title, _ string, t time.Time, list []warn_model.NoticeChannelWarn) string {
	webhookMsg := make([]warn_model.MsgWebhook, 0)
	for _, warnInfo := range list {

		webhookMsg = append(webhookMsg, warn_model.MsgWebhook{
			Title:   title,
			Name:    warnInfo.Name,
			Url:     warnInfo.Url,
			Content: getContent(warnInfo, false),
			Time:    common.TimeToStr(t),
		})
	}
	bytes, _ := json.Marshal(webhookMsg)
	bytes, _ = unescapeUnicode(bytes)
	quote := strconv.Quote(string(bytes))
	quote = quote[1 : len(quote)-1] //去掉转义后的前后两个""
	return quote
}

func unescapeUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}
