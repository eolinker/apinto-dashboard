package audit_service

import (
	"context"
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/audit"
	"github.com/eolinker/apinto-dashboard/modules/audit/audit-entry"
	"github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-dashboard/modules/audit/audit-store"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"time"
)

type auditLogService struct {
	auditLogStore   audit_store.IAuditLogStore
	clusterService  cluster.IClusterService
	userInfoService user.IUserInfoService
}

func newAuditLogService() audit.IAuditLogService {
	s := &auditLogService{}

	bean.Autowired(&s.auditLogStore)
	bean.Autowired(&s.clusterService)
	bean.Autowired(&s.userInfoService)
	return s
}

func (a *auditLogService) GetLogsList(ctx context.Context, namespaceID, operateType int, kind, keyword string, start, end int64, pageNum, pageSize int) ([]*audit_model.LogListItem, int, error) {
	list, total, err := a.auditLogStore.GetLogsByCondition(ctx, namespaceID, operateType, kind, keyword, start, end, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}

	items := make([]*audit_model.LogListItem, 0, len(list))
	for _, info := range list {
		operator := &audit_model.OperatorInfo{
			Username: info.Username,
		}
		userInfo, _ := a.userInfoService.GetUserInfo(ctx, info.UserID)
		if userInfo != nil {
			operator.UserID = userInfo.Id
			operator.Nickname = userInfo.NickName
			operator.Avatar = userInfo.Avatar
		}

		items = append(items, &audit_model.LogListItem{
			ID:          info.Id,
			Operator:    operator,
			OperateType: audit_model.LogOperateType(info.OperateType),
			Kind:        info.Kind,
			Time:        common.TimeToStr(info.StartTime),
			IP:          info.IP,
		})
	}

	return items, total, nil
}

func (a *auditLogService) GetLogDetail(ctx context.Context, logID int) ([]*audit_model.LogDetailArg, error) {
	info, err := a.auditLogStore.Get(ctx, logID)
	if err != nil {
		return nil, err
	}
	args := make([]*audit_model.LogDetailArg, 0, 10)

	userInfo, _ := a.userInfoService.GetUserInfo(ctx, info.UserID)
	if userInfo != nil {
		args = append(args, &audit_model.LogDetailArg{
			Attr:  "用户名",
			Value: userInfo.UserName,
		})
	}

	args = append(args, &audit_model.LogDetailArg{
		Attr:  "操作类型",
		Value: audit_model.GetLogOperateTitle(audit_model.LogOperateType(info.OperateType)),
	})

	args = append(args, &audit_model.LogDetailArg{
		Attr:  "操作对象",
		Value: info.Kind,
	})

	args = append(args, &audit_model.LogDetailArg{
		Attr:  "操作时间",
		Value: common.TimeToStr(info.StartTime),
	})

	args = append(args, &audit_model.LogDetailArg{
		Attr:  "操作IP",
		Value: info.IP,
	})

	args = append(args, &audit_model.LogDetailArg{
		Attr:  "URL",
		Value: info.URL,
	})

	args = append(args, &audit_model.LogDetailArg{
		Attr:  "User-Agent",
		Value: info.UserAgent,
	})
	//对象 信息 object的optType为1时,为批量操作数据
	object := new(audit_model.LogObjectInfo)
	_ = json.Unmarshal([]byte(info.Object), object)
	if object.ClusterName != "" {
		args = append(args, &audit_model.LogDetailArg{
			Attr:  "操作集群",
			Value: object.ClusterName,
		})
	}
	if object.PublishType != 0 {
		args = append(args, &audit_model.LogDetailArg{
			Attr:  "发布类型",
			Value: audit_model.GetPublishTypeTitle(object.PublishType),
		})
	}
	if object.EnableOperate != 0 {
		args = append(args, &audit_model.LogDetailArg{
			Attr:  "启用禁用操作",
			Value: audit_model.GetEnableTypeTitle(object.EnableOperate),
		})
	}
	if object.Name != "" {
		args = append(args, &audit_model.LogDetailArg{
			Attr:  "对象名",
			Value: object.Name,
		})
	}
	if object.Uuid != "" {
		args = append(args, &audit_model.LogDetailArg{
			Attr:  "UUID",
			Value: object.Uuid,
		})
	}

	if info.Body != "" {
		args = append(args, &audit_model.LogDetailArg{
			Attr:  "请求内容",
			Value: info.Body,
		})
	}

	if info.Err != "" {
		args = append(args, &audit_model.LogDetailArg{
			Attr:  "错误信息",
			Value: info.Err,
		})
	}

	return args, nil
}

func (a *auditLogService) Log(namespace int, userId int, operate int, kind string, url, object, ip, userAgent, body, errInfo string, start, end time.Time) {
	ctx := context.Background()
	userInfo, _ := a.userInfoService.GetUserInfo(ctx, userId)
	logInfo := &audit_entry.AuditLog{
		NamespaceId: namespace,
		UserID:      userId,
		IP:          ip,
		OperateType: operate,
		Kind:        kind,
		Object:      object,
		URL:         url,
		StartTime:   start,
		EndTime:     end,
		UserAgent:   userAgent,
		Body:        body,
		Err:         errInfo,
	}
	if userInfo != nil {
		logInfo.Username = userInfo.UserName
	}

	err := a.auditLogStore.Insert(ctx, logInfo)
	if err != nil {
		log.Error(err)
	}
}
