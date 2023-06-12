package model

import (
	"embed"
	"fmt"
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
	PluginID string
	Icon     string
	Readme   string
	Fs       embed.FS
}

func (e *EmbedPluginResources) ICon() ([]byte, bool) {
	data, err := e.Fs.ReadFile(fmt.Sprintf("plugins/%s/%s", e.PluginID, e.Icon))
	if err != nil {
		return nil, false
	}
	return data, true
}

func (e *EmbedPluginResources) RM() ([]byte, bool) {
	data, err := e.Fs.ReadFile(fmt.Sprintf("plugins/%s/%s", e.PluginID, e.Readme))
	if err != nil {
		return nil, false
	}
	return data, true
}

func (e *EmbedPluginResources) ReadMe(name string) ([]byte, bool) {
	data, err := e.Fs.ReadFile(fmt.Sprintf("plugins/%s/%s", e.PluginID, name))
	if err != nil {
		return nil, false
	}
	return data, true
}

func (e *EmbedPluginResources) Resources(path string) ([]byte, bool) {
	data, err := e.Fs.ReadFile(fmt.Sprintf("plugins/%s/resources/%s", e.PluginID, path))
	if err != nil {
		return nil, false
	}
	return data, true
}
