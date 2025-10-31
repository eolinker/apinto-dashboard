package store

import (
	"context"
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/entry"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/log"
)

type IPluginResources interface {
	store.IBaseStore[entry.Resources]
}

type pluginResources struct {
	store.IBaseStore[entry.Resources]
}

func newPluginResourcesStore(db store.IDB, infoBaseStore IPluginStore) IPluginResources {
	ipr := &pluginResources{IBaseStore: store.CreateStore[entry.Resources](db)}

	ctx := context.Background()
	dbg := db.DB(ctx)
	if dbg.Migrator().HasTable(&entry.Package{}) {
		//迁移旧数据至新表resources
		packageList := make([]*entry.Package, 0)
		err := dbg.Find(&packageList).Error
		if err != nil {
			panic(err)
		}
		if len(packageList) > 0 {

			ids := make([]int, 0, len(packageList))
			for _, item := range packageList {
				ids = append(ids, item.Id)
			}
			infoList, err := infoBaseStore.ListQuery(ctx, "`id` in (?)", []interface{}{ids}, "")
			if err != nil {
				panic(err)
			}
			infoMap := common.SliceToMap(infoList, func(t *entry.Plugin) int {
				return t.Id
			})
			for _, item := range packageList {
				info, has := infoMap[item.Id]
				if !has {
					continue
				}
				//解压
				files, err := common.UnzipFromBytes(item.Package)
				if err != nil {
					panic(err)
				}
				resources, _ := json.Marshal(&PluginResources{
					Icon:   info.ICon,
					Readme: "README.md",
					Files:  files,
				})
				err = ipr.Save(ctx, &entry.Resources{
					ID:        info.Id,
					Uuid:      info.UUID,
					Resources: resources,
				})
				if err != nil {
					panic(err)
				}
			}

		}
		err = dbg.Migrator().DropTable(&entry.Package{})
		if err != nil {
			log.Warn("drop talble module_plugin_package:", err)
		}
		dbg.Commit()
	}

	return ipr
}
