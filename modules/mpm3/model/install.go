package model

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type PluginResources struct {
	icon   string
	readme string
	files  map[string][]byte
	hash   string
}
type pluginResourcesMarshal struct {
	Icon   string            `json:"icon"`
	ReadMe string            `json:"readme"`
	Files  map[string][]byte `json:"files"`
}

func (p *PluginResources) UnmarshalJSON(bytes []byte) error {
	p2 := new(pluginResourcesMarshal)
	err := json.Unmarshal(bytes, p2)
	if err != nil {
		return err
	}
	p.files = p2.Files
	p.readme = p2.ReadMe
	p.icon = p2.Icon
	p.hashFiles()
	return nil
}

func (p *PluginResources) MarshalJSON() ([]byte, error) {

	return json.Marshal(&pluginResourcesMarshal{
		Icon:   p.icon,
		ReadMe: p.readme,
		Files:  p.files,
	})

}

func NewPluginResources(icon string, readme string, files map[string][]byte) *PluginResources {
	p := &PluginResources{icon: icon, readme: readme, files: files}
	p.hashFiles()
	return p
}
func (p *PluginResources) hashFiles() {
	if len(p.files) == 0 {
		return
	}
	hd := md5.New()
	for n, f := range p.files {
		hd.Write([]byte(n))
		hd.Write(f)
	}

	p.hash = hex.EncodeToString(hd.Sum(make([]byte, 0, 16)))
}
func (p *PluginResources) Hash() string {
	return p.hash
}
func (p *PluginResources) ICon() ([]byte, bool) {
	data, has := p.files[p.icon]
	return data, has
}

func (p *PluginResources) RM() ([]byte, bool) {
	data, has := p.files[p.readme]
	return data, has
}

func (p *PluginResources) ReadMe(name string) ([]byte, bool) {
	data, has := p.files[name]
	return data, has
}

func (p *PluginResources) Resources(path string) ([]byte, bool) {
	path = fmt.Sprintf("resources/%s", path)
	data, has := p.files[path]
	return data, has
}

type PluginInstalledStatus struct {
	Installed bool
}
