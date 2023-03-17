package openapp_service

import (
	"context"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	"github.com/eolinker/apinto-dashboard/modules/openapp"
	"github.com/eolinker/apinto-dashboard/modules/openapp/open-app-dto"
	"github.com/eolinker/apinto-dashboard/modules/openapp/open-app-entry"
	openappModel "github.com/eolinker/apinto-dashboard/modules/openapp/open-app-model"
	"github.com/eolinker/apinto-dashboard/modules/openapp/openapp-store"
	"github.com/eolinker/apinto-dashboard/modules/random"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/eosc/common/bean"
	"gorm.io/gorm"
	"strings"
	"time"
)

var _ openapp.IExternalApplicationService = (*externalApplicationService)(nil)

type externalApplicationService struct {
	externalAppStore openapp_store.IExternalApplicationStore
	randomService    random.IRandomService
	lockService      locker_service.IAsynLockService
	userInfoService  user.IUserInfoService
}

func newExternalApplicationService() openapp.IExternalApplicationService {
	e := &externalApplicationService{}
	bean.Autowired(&e.externalAppStore)
	bean.Autowired(&e.randomService)
	bean.Autowired(&e.lockService)
	bean.Autowired(&e.userInfoService)

	return e
}

func (e *externalApplicationService) AppList(ctx context.Context, namespaceId int) ([]*openappModel.ExtAppListItem, error) {
	items, err := e.externalAppStore.GetList(ctx, namespaceId)
	if err != nil {
		return nil, err
	}

	list := make([]*openappModel.ExtAppListItem, 0, len(items))
	for _, item := range items {
		userInfo, err := e.userInfoService.GetUserInfo(ctx, item.Operator)
		if err != nil {
			return nil, err
		}

		status := 1 //1表示启用， 2表示禁用
		if item.IsDisable {
			status = 2
		}

		app := &openappModel.ExtAppListItem{
			Id:         item.UUID,
			Name:       item.Name,
			Token:      item.Token,
			Tags:       item.Tags,
			Status:     status,
			Operator:   userInfo.NickName,
			UpdateTime: common.TimeToStr(item.UpdateTime),
		}

		list = append(list, app)
	}

	return list, nil
}

func (e *externalApplicationService) AppInfo(ctx context.Context, namespaceId int, uuid string) (*openappModel.ExternalAppInfo, error) {
	appInfo, err := e.externalAppStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("external-app not found. id:%s ", uuid)
		}
		return nil, err
	}

	info := &openappModel.ExternalAppInfo{
		ExternalApplication: appInfo,
	}

	return info, nil
}

func (e *externalApplicationService) CreateApp(ctx context.Context, namespaceId, userId int, input *open_app_dto.ExternalAppInfoInput) error {
	//外部应用id查重
	_, err := e.externalAppStore.GetByUUID(ctx, namespaceId, input.Id)
	if err != gorm.ErrRecordNotFound {
		if err == nil {
			return fmt.Errorf("id %s already exist. ", input.Name)
		}
		return err
	}

	t := time.Now()
	appInfo := &open_app_entry.ExternalApplication{
		UUID:       input.Id,
		Namespace:  namespaceId,
		Name:       input.Name,
		Token:      e.randomService.RandomStr("external-app-token"),
		Desc:       input.Desc,
		Tags:       "",
		IsDisable:  false,
		IsDelete:   false,
		Operator:   userId,
		CreateTime: t,
		UpdateTime: t,
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: input.Id,
		Name: input.Name,
	})

	return e.externalAppStore.Transaction(ctx, func(txCtx context.Context) error {
		return e.externalAppStore.Save(txCtx, appInfo)
	})
}

func (e *externalApplicationService) UpdateApp(ctx context.Context, namespaceId, userId int, input *open_app_dto.ExternalAppInfoInput) error {
	appInfo, err := e.externalAppStore.GetByUUID(ctx, namespaceId, input.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("external-app not found. id:%s ", input.Id)
		}
		return err
	}

	err = e.lockService.Lock(locker_service.LockNameExtApp, appInfo.Id)
	if err != nil {
		return err
	}
	defer e.lockService.Unlock(locker_service.LockNameExtApp, appInfo.Id)

	_, err = e.externalAppStore.GetByUUID(ctx, namespaceId, input.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("external-app not found. id:%s ", input.Id)
		}
		return err
	}

	t := time.Now()
	//只更新name和desc
	newInfo := &open_app_entry.ExternalApplication{
		Id:         appInfo.Id,
		Name:       input.Name,
		Desc:       input.Desc,
		Operator:   userId,
		UpdateTime: t,
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: input.Id,
		Name: input.Name,
	})

	return e.externalAppStore.Transaction(ctx, func(txCtx context.Context) error {
		_, err = e.externalAppStore.Update(txCtx, newInfo)
		return err
	})
}

func (e *externalApplicationService) DelApp(ctx context.Context, namespaceId, userId int, uuid string) error {
	appInfo, err := e.externalAppStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("external-app not found. id:%s ", uuid)
		}
		return err
	}

	err = e.lockService.Lock(locker_service.LockNameExtApp, appInfo.Id)
	if err != nil {
		return err
	}
	defer e.lockService.Unlock(locker_service.LockNameExtApp, appInfo.Id)

	appInfo, err = e.externalAppStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("external-app not found. id:%s ", uuid)
		}
		return err
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: appInfo.UUID,
		Name: appInfo.Name,
	})

	return e.externalAppStore.Transaction(ctx, func(txCtx context.Context) error {
		return e.externalAppStore.SoftDelete(txCtx, userId, appInfo.Id)
	})
}

func (e *externalApplicationService) Enable(ctx context.Context, namespaceId, userId int, uuid string) error {
	appInfo, err := e.externalAppStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("external-app not found. id:%s ", uuid)
		}
		return err
	}

	err = e.lockService.Lock(locker_service.LockNameExtApp, appInfo.Id)
	if err != nil {
		return err
	}
	defer e.lockService.Unlock(locker_service.LockNameExtApp, appInfo.Id)

	appInfo, err = e.externalAppStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("external-app not found. id:%s ", uuid)
		}
		return err
	}

	//若app已经是开启状态则返回
	if !appInfo.IsDisable {
		return nil
	}

	appInfo.IsDisable = false
	appInfo.Operator = userId
	appInfo.UpdateTime = time.Now()

	return e.externalAppStore.Transaction(ctx, func(txCtx context.Context) error {
		return e.externalAppStore.Save(txCtx, appInfo)
	})
}

func (e *externalApplicationService) Disable(ctx context.Context, namespaceId, userId int, uuid string) error {
	appInfo, err := e.externalAppStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("external-app not found. id:%s ", uuid)
		}
		return err
	}

	err = e.lockService.Lock(locker_service.LockNameExtApp, appInfo.Id)
	if err != nil {
		return err
	}
	defer e.lockService.Unlock(locker_service.LockNameExtApp, appInfo.Id)

	appInfo, err = e.externalAppStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("external-app not found. id:%s ", uuid)
		}
		return err
	}

	//若app已经是禁用状态则返回
	if appInfo.IsDisable {
		return nil
	}

	appInfo.IsDisable = true
	appInfo.Operator = userId
	appInfo.UpdateTime = time.Now()

	return e.externalAppStore.Transaction(ctx, func(txCtx context.Context) error {
		return e.externalAppStore.Save(txCtx, appInfo)
	})
}

func (e *externalApplicationService) FlushToken(ctx context.Context, namespaceId, userId int, uuid string) error {
	appInfo, err := e.externalAppStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("external-app not found. id:%s ", uuid)
		}
		return err
	}

	err = e.lockService.Lock(locker_service.LockNameExtApp, appInfo.Id)
	if err != nil {
		return err
	}
	defer e.lockService.Unlock(locker_service.LockNameExtApp, appInfo.Id)

	appInfo, err = e.externalAppStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("external-app not found. id:%s ", uuid)
		}
		return err
	}

	newToken := e.randomService.RandomStr("external-app-token")

	return e.externalAppStore.Transaction(ctx, func(txCtx context.Context) error {
		return e.externalAppStore.FlushToken(txCtx, userId, appInfo.Id, newToken)
	})
}

func (e *externalApplicationService) CheckExtAPPToken(ctx context.Context, namespaceId int, token string) (int, error) {
	appInfo, err := e.externalAppStore.GetByToken(ctx, namespaceId, token)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("ext-app not found. token is invalid. ")
		}
		return 0, err
	}
	if appInfo.IsDisable {
		return 0, fmt.Errorf("ext-app is disable. ")
	}

	return appInfo.Id, nil
}
func (e *externalApplicationService) UpdateExtAPPTags(ctx context.Context, namespaceId, appID int, label string) error {
	app, err := e.externalAppStore.Get(ctx, appID)
	if err != nil {
		return err
	}
	if label == "" {
		return nil
	}

	if app.Tags != "" {
		tagItems := strings.Split(app.Tags, ",")
		isReduplicated := false
		for _, tag := range tagItems {
			if tag == label {
				isReduplicated = true
				break
			}
		}
		if !isReduplicated {
			tagItems = append(tagItems, label)
		}
		app.Tags = strings.Join(tagItems, ",")
	} else {
		app.Tags = label
	}
	return e.externalAppStore.Save(ctx, app)
}

func (e *externalApplicationService) GetExtAppName(ctx context.Context, id int) (string, error) {
	info, err := e.externalAppStore.Get(ctx, id)
	if err != nil {
		return "", err
	}
	return info.Name, nil
}
