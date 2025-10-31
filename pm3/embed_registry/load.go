package embed_registry

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path"

	"github.com/eolinker/apinto-dashboard/common"

	"github.com/eolinker/apinto-dashboard/modules/mpm3/model"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/apinto-dashboard/pm3/pinstall"
	"github.com/eolinker/eosc/log"
)

type PluginCfg struct {
	define    *pm3.PluginDefine
	auto      bool
	isDisable bool
	resources *model.PluginResources
}

func LoadPlugins(fs *embed.FS, dir string, target string) error {
	entries, err := fs.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		filePath := path.Join(dir, e.Name(), target)
		fileContent, err := fs.ReadFile(filePath)
		if err != nil {
			return err
		}
		pluginCfg, auto, isDisable, err := pinstall.ReadInner(fileContent)
		if err != nil {
			log.Errorf("read inert plugin file(%s) error: %v", filePath, err)
			return err
		}
		registerEmbedPlugin(&PluginCfg{
			define:    pluginCfg,
			auto:      auto,
			isDisable: isDisable,
			resources: NewEmbedPluginResources(fs, path.Join("embed", e.Name()), pluginCfg.ICon, "README.md"),
		})

	}

	return err
}

// LoadLocalPlugins 加载本地插件
func LoadLocalPlugins(dir string, target string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	ctx := context.Background()
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		filePath := path.Join(dir, e.Name())
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			log.Errorf("read inert plugin file(%s) error: %v", filePath, err)
			continue
		}
		files, err := common.UnzipFromBytes(fileContent)
		if err != nil {
			log.Errorf("解压插件安装包失败. err:%s", err.Error())
			continue
		}
		externPCfg, err := ValidPluginFiles(files)
		if err != nil {
			log.Errorf("verify plugin file(%s) error: %v", filePath, err)
			continue
		}
		resources := model.NewPluginResources(externPCfg.ICon, "README.md", files)
		err = service.Install(ctx, externPCfg, resources)
		if err != nil {
			log.Errorf("安装插件失败. err:%s", err.Error())
			continue
		}
		err = os.Remove(filePath)
		if err != nil {
			log.Errorf("删除插件安装包失败. err:%s", err.Error())
			continue
		}
	}

	return err
}

func ValidPluginFiles(files map[string][]byte) (*pm3.PluginDefine, error) {
	//校验解压目录下有没有必要的文件 plugin.yml icon README.md
	pluginYml, has := files["plugin.yml"]
	if !has {
		return nil, fmt.Errorf("安装插件失败, plugin.yml不存在")
	}

	//README.md是否存在
	_, has = files["README.md"]
	if !has {
		return nil, fmt.Errorf("安装插件失败 README.md 文件不存在")
	}

	externPCfg, err := pinstall.Read(pluginYml)

	if err != nil {
		return nil, fmt.Errorf("plugin.yml解析失败. err:%s", err.Error())
	}
	//TODO 校验plugin.yml

	//图标文件是否存在
	_, has = files[externPCfg.ICon]
	if err != nil {
		return nil, fmt.Errorf("安装插件失败 图标文件不存在")
	}
	return externPCfg, nil
}
