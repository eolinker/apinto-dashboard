package model

import (
	"embed"
	"fmt"
	"path"
)

type PluginResources struct {
	Icon   string
	Readme string
	Files  map[string][]byte
}

func (p *PluginResources) ICon() ([]byte, bool) {
	data, has := p.Files[p.Icon]
	return data, has
}

func (p *PluginResources) RM() ([]byte, bool) {
	data, has := p.Files[p.Readme]
	return data, has
}

func (p *PluginResources) ReadMe(name string) ([]byte, bool) {
	data, has := p.Files[name]
	return data, has
}

func (p *PluginResources) Resources(path string) ([]byte, bool) {
	path = fmt.Sprintf("resources/%s", path)
	data, has := p.Files[path]
	return data, has
}

type EmbedPluginResources struct {
	icon   string
	readme string
	root   string
	Fs     embed.FS
}

func NewEmbedPluginResources(fs embed.FS, root string, icon string, readme string) *EmbedPluginResources {
	return &EmbedPluginResources{icon: icon, readme: readme, root: path.Join(root), Fs: fs}
}

func (e *EmbedPluginResources) ICon() ([]byte, bool) {
	data, err := e.Fs.ReadFile(path.Join(e.root, e.icon))
	if err != nil {
		return nil, false
	}
	return data, true
}

func (e *EmbedPluginResources) RM() ([]byte, bool) {
	data, err := e.Fs.ReadFile(path.Join(e.root, e.readme))
	if err != nil {
		return nil, false
	}
	return data, true
}

func (e *EmbedPluginResources) ReadMe(name string) ([]byte, bool) {
	data, err := e.Fs.ReadFile(path.Join(e.root, name))
	if err != nil {
		return nil, false
	}
	return data, true
}

func (e *EmbedPluginResources) Resources(file string) ([]byte, bool) {
	data, err := e.Fs.ReadFile(path.Join(e.root, "resources", file))
	if err != nil {
		return nil, false
	}
	return data, true
}
