package warn_timer

import (
	"context"
	"fmt"
	mock_cache "github.com/eolinker/apinto-dashboard/cache/common-cache-mock"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/api/api-entry"
	api_mock "github.com/eolinker/apinto-dashboard/modules/api/mock"
	api_model "github.com/eolinker/apinto-dashboard/modules/api/model"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-entry"
	cluster_mock "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-mock"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-model"
	monitor_model "github.com/eolinker/apinto-dashboard/modules/monitor/model"
	monitor_mock "github.com/eolinker/apinto-dashboard/modules/monitor/monitor-mock"
	namespace_entry "github.com/eolinker/apinto-dashboard/modules/namespace/namespace-entry"
	namespace_mock "github.com/eolinker/apinto-dashboard/modules/namespace/namespace-mock"
	"github.com/eolinker/apinto-dashboard/modules/namespace/namespace-model"
	"github.com/eolinker/apinto-dashboard/modules/notice"
	notice_mock "github.com/eolinker/apinto-dashboard/modules/notice/notice-mock"
	"github.com/eolinker/apinto-dashboard/modules/notice/notice-model"
	notice_service "github.com/eolinker/apinto-dashboard/modules/notice/notice-service"
	upstream_mock "github.com/eolinker/apinto-dashboard/modules/upstream/mock"
	"github.com/eolinker/apinto-dashboard/modules/upstream/model"
	"github.com/eolinker/apinto-dashboard/modules/user/user-entry"
	user_mock "github.com/eolinker/apinto-dashboard/modules/user/user-mock"
	"github.com/eolinker/apinto-dashboard/modules/user/user-model"
	warn_mock "github.com/eolinker/apinto-dashboard/modules/warn/warn-mock"
	"github.com/eolinker/apinto-dashboard/modules/warn/warn-model"

	"github.com/golang/mock/gomock"
	"net/http"
	"testing"
	"time"
)

func warnStrategyAll(ctx context.Context, namespaceId, partitionId int, partitionUUID string, channelUuids, targetValues []string, userId []int) []*warn_model.WarnStrategy {

	warnStrategyConfigRule := make([]*warn_model.WarnStrategyConfigRule, 0)
	warnStrategyConfigRuleCondition := make([]*warn_model.WarnStrategyConfigRuleCondition, 0)

	warnStrategyConfigRuleCondition = append(warnStrategyConfigRuleCondition, &warn_model.WarnStrategyConfigRuleCondition{
		Compare: ">",
		Unit:    "%",
		Value:   1.01,
	})

	warnStrategyConfigRule = append(warnStrategyConfigRule, &warn_model.WarnStrategyConfigRule{
		ChannelUuids: channelUuids,
		Condition:    warnStrategyConfigRuleCondition,
	})
	warnStrategies := make([]*warn_model.WarnStrategy, 0)
	warnStrategies = append(warnStrategies, &warn_model.WarnStrategy{
		PartitionId: partitionId,
		NamespaceId: namespaceId,
		Uuid:        partitionUUID,
		Title:       "partition告警策略1",
		IsEnable:    true,
		Dimension:   "partition",
		Quota:       warn_model.QuotaTypeReqFailRate,
		Every:       3,
		WarnStrategyConfig: &warn_model.WarnStrategyConfig{
			Target: warn_model.WarnStrategyConfigTarget{
				Rule: "unlimited",
				//Values: targetValues,
			},
			Rule:       warnStrategyConfigRule,
			Continuity: 0,
			HourMax:    0,
			Users:      userId,
		},
		PartitionUUID: partitionUUID,
	})

	return warnStrategies
}

func warnStatistics(apiIds []string) map[string]float64 {
	m := make(map[string]float64)
	//for _, id := range apiIds {
	//	if id == "83a98e31-f4cd-eae4-cac3-d1031540e030" {
	//
	//	}
	//	m[id] = 10
	//}
	m["83a98e31-f4cd-eae4-cac3-d1031540e030"] = 10
	m["887a8d44-6a89-bfc8-3d5e-44f1c3928310"] = 30

	return m
}

func warnPartitionStatistics() map[string]float64 {
	m := make(map[string]float64)
	//for _, id := range apiIds {
	//	if id == "83a98e31-f4cd-eae4-cac3-d1031540e030" {
	//
	//	}
	//	m[id] = 10
	//}
	m["zzy_local"] = 0.23
	m["zzy_test"] = 0.22

	return m
}

func partitionInfo(clusterNames []string) *monitor_model.MonPartitionInfo {
	return &monitor_model.MonPartitionInfo{
		Id:           1,
		Name:         "zzy",
		SourceType:   "",
		Config:       nil,
		Env:          "",
		ClusterNames: clusterNames,
	}
}

func userInfoAll(userId []int) map[int]*user_model.UserInfo {
	maps := make(map[int]*user_model.UserInfo)
	for _, id := range userId {
		maps[id] = user_model.CreateUserInfo(&user_entry.UserInfo{
			Id:           id,
			Sex:          1,
			UserName:     "zzy",
			NoticeUserId: "1324204490",
			NickName:     "张泽意",
			Email:        "1324204490@qq.com",
		})
	}
	return maps
}

func getClustersNames(names []string) []*cluster_model.Cluster {
	list := make([]*cluster_model.Cluster, 0)
	for i, name := range names {
		list = append(list, &cluster_model.Cluster{
			Cluster: &cluster_entry.Cluster{
				Id:          i + 1,
				NamespaceId: 1,
				Name:        name,
				UUID:        name,
			},
			Status: 0,
		})
	}
	return list
}

func getApiList(apiIds []string) []*api_model.APIInfo {
	list := make([]*api_model.APIInfo, 0)
	for i, id := range apiIds {
		list = append(list, &api_model.APIInfo{
			API: &api_entry.API{
				Id:               i + 1,
				NamespaceId:      0,
				UUID:             id,
				Name:             "API" + fmt.Sprintf("%d", i+1),
				RequestPath:      "",
				RequestPathLabel: "/baidu/" + fmt.Sprintf("%d", i),
			},
		})
	}
	return list
}

func getServiceList() []*upstream_model.ServiceListItem {
	return nil
}

func getNoticeChannel(channelUuids []string) []*notice_model.NoticeChannel {
	list := make([]*notice_model.NoticeChannel, 0)
	for i, uuid := range channelUuids {

		title := "email"
		if i+1 == 1 {
			title = "webhook"
		}
		list = append(list, &notice_model.NoticeChannel{
			Id:    i + 1,
			Name:  uuid,
			Title: title,
			Type:  i + 1,
		})
	}
	return list
}
func Test_newMonitorWarn(t *testing.T) {
	ctl := gomock.NewController(t)
	ctx := context.Background()
	namespaceId := 1

	partitionId := 1
	partitionUUID := "101"
	userIds := make([]int, 0)
	userIds = append(userIds, 18912, 18914)

	apiIds := make([]string, 0)
	apiIds = append(apiIds, "83a98e31-f4cd-eae4-cac3-d1031540e030", "887a8d44-6a89-bfc8-3d5e-44f1c3928310")

	clusterNames := make([]string, 0)
	clusterNames = append(clusterNames, "zzy_local", "zzy_test")

	channelUuids := make([]string, 0)
	channelUuids = append(channelUuids, "channelID1", "channelID2")

	warnStrategyService := warn_mock.NewMockIWarnStrategyService(ctl)
	strategies := warnStrategyAll(ctx, namespaceId, partitionId, partitionUUID, channelUuids, apiIds, userIds)
	warnStrategyService.EXPECT().WarnStrategyAll(ctx, namespaceId, 1).Return(strategies, nil)

	endTime, _ := time.ParseInLocation("2006-01-02 15:04", time.Now().Add(-time.Minute).Format("2006-01-02 15:04"), time.Local)
	startTime := endTime.Add(-time.Minute * time.Duration(3))

	statisticsService := monitor_mock.NewMockIMonitorStatistics(ctl)
	statisticsService.EXPECT().WarnStatistics(ctx, namespaceId, partitionUUID, startTime, endTime, "cluster", warn_model.QuotaTypeReqFailRate, nil).Return(warnPartitionStatistics(), nil)

	namespaceService := namespace_mock.NewMockINamespaceService(ctl)
	namespaces := make([]*namespace_model.Namespace, 0)
	namespaces = append(namespaces, &namespace_model.Namespace{Namespace: &namespace_entry.Namespace{
		Id:   1,
		Name: "default",
	}})
	namespaceService.EXPECT().GetAll().Return(namespaces, nil)

	monitorService := monitor_mock.NewMockIMonitorService(ctl)

	monitorService.EXPECT().PartitionInfo(ctx, namespaceId, partitionUUID).Return(partitionInfo(clusterNames), nil)

	userService := user_mock.NewMockIUserInfoService(ctl)
	userService.EXPECT().GetUserInfoMaps(ctx).Return(userInfoAll(userIds), nil)

	clusterService := cluster_mock.NewMockIClusterService(ctl)

	clusterService.EXPECT().GetByNames(ctx, namespaceId, clusterNames).Return(getClustersNames(clusterNames), nil)

	warnHistoryService := warn_mock.NewMockIWarnHistoryService(ctl)

	warnHistoryService.EXPECT().Create(ctx, namespaceId, partitionId, &warn_model.WarnHistoryInfo{}).Return(nil)

	apiService := api_mock.NewMockIAPIService(ctl)
	apiService.EXPECT().GetAPIInfoAll(ctx, namespaceId).Return(getApiList(apiIds), nil)

	service := upstream_mock.NewMockIService(ctl)
	service.EXPECT().GetServiceListAll(ctx, namespaceId).Return(getServiceList(), nil)

	noticeChannelService := notice_mock.NewMockINoticeChannelService(ctl)
	noticeChannelService.EXPECT().NoticeChannelList(ctx, namespaceId, 0).Return(getNoticeChannel(channelUuids), nil)

	commonCache := mock_cache.NewMockICommonCache(ctl)

	noticeChannelDriver := notice_service.NewNoticeChannelDriverManager()

	for i, uuid := range channelUuids {
		var driverManager notice.IDriverNoticeChannel
		if i+1 == 1 {
			driverManager = common.NewWebhook("https://open.feishu.cn/open-apis/bot/v2/hook/a7ed8efa-88ac-4721-af0e-c00d02172312", http.MethodPost, "JSON", common.SingleNotice, ",", map[string]string{"test1": "test1", "test2": "test2"}, `{
    "msg_type": "interactive",
    "card": {
        "elements": [{
                "tag": "div",
                "text": {
                        "content": "{msg}",
                        "tag": "lark_md"
                }
        }],
        "header": {
                "title": {
                        "content": "{title}",
                        "tag": "plain_text"
                }
        }
    }
}`)
		} else {
			driverManager = common.NewSmtp("smtp.qq.com", 587, "", "1324204490@qq.com", "zzeqxoubzzoababg", "1324204490@qq.com")
		}

		noticeChannelDriver.RegisterDriver(uuid, driverManager)
	}

	tests := []struct {
		name string
		want IMonitorWarn
	}{
		{
			name: "",
			want: &monitorWarn{
				warnStrategyService:  warnStrategyService,
				monitorStatistics:    statisticsService,
				namespaceService:     namespaceService,
				monitorService:       monitorService,
				userService:          userService,
				clusterService:       clusterService,
				warnHistoryService:   warnHistoryService,
				apiService:           apiService,
				service:              service,
				commonCache:          commonCache,
				noticeChannelService: noticeChannelService,
				noticeChannelDriver:  noticeChannelDriver,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.monitorWarn()
		})
	}
}
