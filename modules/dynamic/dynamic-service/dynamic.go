package dynamic_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	v2 "github.com/eolinker/apinto-dashboard/client/v2"
	dynamic_entry "github.com/eolinker/apinto-dashboard/modules/dynamic/dynamic-entry"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/apinto-dashboard/modules/cluster"

	"github.com/eolinker/apinto-dashboard/modules/user"

	dynamic_model "github.com/eolinker/apinto-dashboard/modules/dynamic/dynamic-model"

	"github.com/eolinker/apinto-dashboard/modules/dynamic"
	dynamic_store "github.com/eolinker/apinto-dashboard/modules/dynamic/dynamic-store"
	"github.com/eolinker/eosc/common/bean"
)

const (
	statusOffline = "未发布"
	statusPre     = "待发布"
	statusOnline  = "已发布"
)

type dynamicService struct {
	userService    user.IUserInfoService
	clusterService cluster.IClusterService

	dynamicStore        dynamic_store.IDynamicStore
	publishHistoryStore dynamic_store.IDynamicPublishHistoryStore
	publishVersionStore dynamic_store.IDynamicPublishVersionStore
}

func (d *dynamicService) Info(ctx context.Context, namespaceId int, module string, name string) (*dynamic_model.DynamicBasicInfo, map[string]interface{}, error) {
	info, err := d.dynamicStore.First(ctx, map[string]interface{}{
		"namespace": namespaceId,
		"module":    module,
		"name":      name,
	})
	if err != nil {
		return nil, nil, err
	}
	tmp := make(map[string]interface{})
	err = json.Unmarshal([]byte(info.Config), &tmp)
	if err != nil {
		return nil, nil, err
	}

	return &dynamic_model.DynamicBasicInfo{
		ID:          info.Name,
		Title:       info.Title,
		Description: info.Description,
	}, tmp, nil
}

func (d *dynamicService) List(ctx context.Context, namespaceId int, module string, columns []string, keyword string, page int, pageSize int) ([]map[string]string, error) {
	list, err := d.dynamicStore.ListPageByKeyword(ctx, map[string]interface{}{
		"namespace": namespaceId,
		"module":    module,
	}, keyword, page, pageSize)
	if err != nil {
		return nil, err
	}
	items := make([]map[string]string, 0, len(list))
	for _, l := range list {
		updater := ""
		u, err := d.userService.GetUserInfo(ctx, l.Updater)
		if err == nil {
			updater = u.UserName
		}

		item := map[string]string{
			"id":          l.Name,
			"title":       l.Title,
			"updater":     updater,
			"update_time": l.UpdateTime.Format("2006-01-02 15:04:05"),
		}
		tmp := make(map[string]interface{})
		err = json.Unmarshal([]byte(l.Config), &tmp)
		if err == nil {
			for _, column := range columns {
				if _, ok := item[column]; ok {
					continue
				}
				item[column] = ""
				vv, has := tmp[column]
				if has {
					v, ok := vv.(string)
					if ok {
						item[column] = v
						continue
					}
				}
			}
		}
		items = append(items, item)
	}
	return items, nil
}

func (d *dynamicService) ClusterStatuses(ctx context.Context, namespaceId int, module string, profession string, names []string, keyword string, page int, pageSize int) (map[string]map[string]string, error) {
	list, err := d.dynamicStore.ListPageByKeyword(ctx, map[string]interface{}{
		"namespace": namespaceId,
		"module":    module,
	}, keyword, page, pageSize)
	if err != nil {
		return nil, err
	}
	clusters, err := d.clusterService.GetByNames(ctx, namespaceId, names)
	versionMap := make(map[string]map[string]string)
	for _, c := range clusters {
		client := GetClusterClient(c.Name, c.Addr)
		workers, err := client.List(profession)
		if err != nil {
			continue
		}
		for _, w := range workers {
			if _, ok := versionMap[w.Id]; !ok {
				versionMap[w.Name] = make(map[string]string)
			}
			versionMap[w.Name][c.Name] = w.Version
		}
	}
	result := make(map[string]map[string]string)
	for _, l := range list {
		hasVersion := true
		v, ok := versionMap[l.Name]
		if !ok {
			hasVersion = false
		}
		clusterStatus := map[string]string{}
		for _, name := range names {
			version := v[name]
			if !hasVersion || version == "" {
				clusterStatus[name] = statusOffline
				continue
			}
			if l.Version != version {
				clusterStatus[name] = statusPre
			} else {
				clusterStatus[name] = statusOnline
			}

		}
		result[l.Name] = clusterStatus
	}
	return result, nil
}

func (d *dynamicService) Online(ctx context.Context, namespaceId int, module string, profession string, name string, names []string, updater int) ([]string, []string, error) {
	info, err := d.dynamicStore.First(ctx, map[string]interface{}{
		"namespace": namespaceId,
		"module":    module,
		"name":      name,
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, fmt.Errorf("%s(%s) not found", module, name)
		}
		return nil, nil, err
	}

	clusters, err := d.clusterService.GetByNames(ctx, namespaceId, names)
	if err != nil {
		return nil, nil, err
	}
	successClusters := make([]string, 0, len(clusters))
	failClusters := make([]string, 0, len(clusters))
	now := time.Now()
	var publishConfig dynamic_entry.DynamicPublishConfig
	err = json.Unmarshal([]byte(info.Config), &publishConfig)
	if err != nil {
		return nil, nil, err
	}

	for _, c := range clusters {
		version := &dynamic_entry.DynamicPublishVersion{
			ClusterId:   c.Id,
			NamespaceId: namespaceId,
			Publish:     &publishConfig,
			Operator:    updater,
			CreateTime:  now,
		}
		history := &dynamic_entry.DynamicPublishHistory{
			VersionName: info.Version,
			ClusterId:   c.Id,
			NamespaceId: namespaceId,
			Publish:     &publishConfig,
			OptType:     1,
			Operator:    updater,
			CreateTime:  now,
		}
		err = d.saveVersion(ctx, version, history, GetClusterClient(c.Name, c.Addr))
		if err != nil {
			log.Errorf("fail to online config in cluster(%s),addr is %s,profession is %s,uuid is %s,config is %s", c.Name, c.Addr, profession, name, info.Config)
			failClusters = append(failClusters, c.Name)
			continue
		}
		successClusters = append(successClusters, c.Name)
	}
	if len(failClusters) > 0 {
		if len(successClusters) < 1 {
			err = errors.New("some clusters failed to go online")
		} else {
			err = errors.New("all clusters failed to go online")
		}
	}
	return successClusters, failClusters, nil
}

func (d *dynamicService) Offline(ctx context.Context, namespaceId int, module string, profession, name string, names []string, updater int) ([]string, []string, error) {
	info, err := d.dynamicStore.First(ctx, map[string]interface{}{
		"namespace": namespaceId,
		"module":    module,
		"name":      name,
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, fmt.Errorf("%s(%s) not found", module, name)
		}
		return nil, nil, err
	}
	clusters, err := d.clusterService.GetByNames(ctx, namespaceId, names)
	successClusters := make([]string, 0, len(clusters))
	failClusters := make([]string, 0, len(clusters))
	now := time.Now()
	var publishConfig dynamic_entry.DynamicPublishConfig
	err = json.Unmarshal([]byte(info.Config), &publishConfig)
	if err != nil {
		return nil, nil, err
	}
	for _, c := range clusters {
		version := &dynamic_entry.DynamicPublishVersion{
			ClusterId:   c.Id,
			NamespaceId: namespaceId,
			Publish:     &publishConfig,
			Operator:    updater,
			CreateTime:  now,
		}
		history := &dynamic_entry.DynamicPublishHistory{
			VersionName: info.Version,
			ClusterId:   c.Id,
			NamespaceId: namespaceId,
			Publish:     &publishConfig,
			OptType:     3,
			Operator:    updater,
			CreateTime:  now,
		}
		err = d.saveVersion(ctx, version, history, GetClusterClient(c.Name, c.Addr))
		if err != nil {
			log.Errorf("fail to online config in cluster(%s),addr is %s,profession is %s,uuid is %s,config is %s", c.Name, c.Addr, profession, name, info.Config)
			failClusters = append(failClusters, c.Name)
			continue
		}

		successClusters = append(successClusters, c.Name)
	}
	if len(failClusters) > 0 {
		if len(successClusters) < 1 {
			err = errors.New("some clusters failed to go offline")
		} else {
			err = errors.New("all clusters failed to go offline")
		}
	}
	return successClusters, failClusters, nil
}

func (d *dynamicService) ClusterStatus(ctx context.Context, namespaceId int, module string, profession, name string) (*dynamic_model.DynamicBasicInfo, []*dynamic_model.DynamicCluster, error) {
	moduleInfo, err := d.dynamicStore.First(ctx, map[string]interface{}{
		"namespace": namespaceId,
		"module":    module,
		"name":      name,
	})
	if err != nil {
		return nil, nil, err
	}

	clusters, err := d.clusterService.GetAllCluster(ctx)
	if err != nil {
		return nil, nil, err
	}
	result := make([]*dynamic_model.DynamicCluster, 0, len(clusters))
	for _, c := range clusters {
		v, err := d.publishHistoryStore.First(ctx, map[string]interface{}{
			"namespace":    namespaceId,
			"cluster":      c.Cluster,
			"target":       moduleInfo.Id,
			"version_name": moduleInfo.Version,
			"kind":         "dynamic_module",
		})
		if err != nil {
			result = append(result, &dynamic_model.DynamicCluster{
				Name:   c.Name,
				Title:  c.Name,
				Status: statusOffline,
			})
			continue
		}

		client := GetClusterClient(c.Name, c.Addr)
		info, err := client.Info(profession, name)
		if err != nil {
			result = append(result, &dynamic_model.DynamicCluster{
				Name:   c.Name,
				Title:  c.Name,
				Status: statusOffline,
			})
			continue
		}
		status := statusOnline
		if info.Version != moduleInfo.Version {
			status = statusOffline
		}
		updater := ""
		u, err := d.userService.GetUserInfo(ctx, v.Operator)
		if err == nil {
			updater = u.UserName
		}
		result = append(result, &dynamic_model.DynamicCluster{
			Name:       c.Name,
			Title:      c.Name,
			Status:     status,
			Updater:    updater,
			UpdateTime: v.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return &dynamic_model.DynamicBasicInfo{
		ID:          moduleInfo.Name,
		Title:       moduleInfo.Title,
		Description: moduleInfo.Description,
	}, result, nil
}

func (d *dynamicService) Create(ctx context.Context, namespaceId int, module string, title string, name string, driver string, description string, body string, updater int) error {
	now := time.Now()
	info := &dynamic_entry.Dynamic{
		NamespaceId: namespaceId,
		Name:        name,
		Title:       title,
		Driver:      driver,
		Description: description,
		Version:     genVersion(now),
		Config:      body,
		Module:      module,
		Updater:     updater,
		CreateTime:  now,
		UpdateTime:  now,
	}
	return d.dynamicStore.Insert(ctx, info)
}

func (d *dynamicService) Save(ctx context.Context, namespaceId int, module string, title string, name string, description string, body string, updater int) error {

	info, err := d.dynamicStore.First(ctx, map[string]interface{}{
		"namespace": namespaceId,
		"module":    module,
		"name":      name,
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("%s(%s) not found", module, name)
		}
		return err
	}

	now := time.Now()

	info.Title = title
	info.Description = description
	info.Config = body
	info.Updater = updater
	info.Version = genVersion(now)
	info.UpdateTime = now

	return d.dynamicStore.Save(ctx, info)
}

func (d *dynamicService) Delete(ctx context.Context, namespaceId int, module string, name string) error {
	_, err := d.dynamicStore.DeleteWhere(ctx, map[string]interface{}{
		"namespace": namespaceId,
		"module":    module,
		"name":      name,
	})
	return err
}

func (d *dynamicService) saveVersion(ctx context.Context, version *dynamic_entry.DynamicPublishVersion, history *dynamic_entry.DynamicPublishHistory, client v2.IClient) error {
	return d.publishVersionStore.Transaction(ctx, func(txCtx context.Context) error {
		var err error
		if err = d.publishVersionStore.Save(txCtx, version); err != nil {
			return err
		}

		if history != nil {
			history.VersionId = version.Id
			if err = d.publishHistoryStore.Insert(txCtx, history); err != nil {
				return err
			}
			return client.Set(history.Publish.Profession, history.Publish.Name, &v2.WorkerInfo{
				BasicInfo: &v2.BasicInfo{
					Profession:  version.Publish.BasicInfo.Profession,
					Name:        version.Publish.BasicInfo.Name,
					Driver:      version.Publish.BasicInfo.Name,
					Description: version.Publish.BasicInfo.Description,
					Version:     history.VersionName,
				},
				Append: version.Publish.Append,
			})
		}
		return nil
	})
}

func newDynamicService() dynamic.IDynamicService {
	d := &dynamicService{}
	bean.Autowired(&d.dynamicStore)
	bean.Autowired(&d.userService)
	bean.Autowired(&d.clusterService)
	return d
}

func genVersion(now time.Time) string {
	return now.Format("20060102150405")
}
