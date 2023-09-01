package store

import (
	"context"
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/entry"
	"github.com/eolinker/apinto-dashboard/store"
)

var (
	_ IModulePluginPackageStore = (*modulePluginPackage)(nil)
)

type IModulePluginPackageStore interface {
	store.IBaseStore[entry.ModulePluginPackage]
}

type modulePluginPackage struct {
	*store.BaseStore[entry.ModulePluginPackage]
}

func newModulePluginPackageStore(db store.IDB) IModulePluginPackageStore {
	packageBaseStore := store.CreateStore[entry.ModulePluginPackage](db)
	infoBaseStore := store.CreateStore[entry.ModulePlugin](db)
	resourcesBaseStore := store.CreateStore[entry.PluginResources](db)
	//迁移旧数据至新表resources
	ctx := context.Background()
	packageList, err := packageBaseStore.List(ctx, nil)
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
		infoMap := common.SliceToMap(infoList, func(t *entry.ModulePlugin) int {
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
			err = resourcesBaseStore.Save(ctx, &entry.PluginResources{
				ID:        info.Id,
				Uuid:      info.UUID,
				Resources: resources,
			})
			if err != nil {
				panic(err)
			}
		}
		//删除package
		_, _ = packageBaseStore.Delete(ctx, ids...)
	}

	return &modulePluginPackage{BaseStore: packageBaseStore}
}

type PluginResources struct {
	Icon   string
	Readme string
	Files  map[string][]byte
}
