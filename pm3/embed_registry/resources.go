package embed_registry

import (
	"embed"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/model"
	"path"
)

type embedPluginResources struct {
	icon   string
	readme string
	root   string
	Fs     *embed.FS
}

func NewEmbedPluginResources(fs *embed.FS, root string, icon string, readme string) *model.PluginResources {
	resources := &embedPluginResources{icon: icon, readme: readme, root: path.Join(root), Fs: fs}
	return resources.toResources()
}
func (e *embedPluginResources) toResources() *model.PluginResources {

	files := make(map[string][]byte)

	readDir(e.Fs, e.root, "", files)
	return model.NewPluginResources(e.icon, e.readme, files)
}
func readDir(fs *embed.FS, root, name string, out map[string][]byte) {
	dir := path.Join(root, name)
	dirEntries, err := fs.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range dirEntries {

		cn := path.Join(name, e.Name())
		if e.IsDir() {
			readDir(fs, root, cn, out)
		} else {
			file, err := fs.ReadFile(path.Join(root, cn))
			if err != nil {
				continue
			}
			out[cn] = file
		}
	}
}

//
//func (e *embedPluginResources) ICon() ([]byte, bool) {
//	data, err := e.Fs.ReadFile(path.Join(e.root, e.icon))
//	if err != nil {
//		return nil, false
//	}
//	return data, true
//}
//
//func (e *embedPluginResources) RM() ([]byte, bool) {
//	data, err := e.Fs.ReadFile(path.Join(e.root, e.readme))
//	if err != nil {
//		return nil, false
//	}
//	return data, true
//}
//
//func (e *embedPluginResources) ReadMe(name string) ([]byte, bool) {
//	data, err := e.Fs.ReadFile(path.Join(e.root, name))
//	if err != nil {
//		return nil, false
//	}
//	return data, true
//}
//
//func (e *embedPluginResources) Resources(file string) ([]byte, bool) {
//	data, err := e.Fs.ReadFile(path.Join(e.root, "resources", file))
//	if err != nil {
//		return nil, false
//	}
//	return data, true
//}
