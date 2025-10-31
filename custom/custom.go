package custom

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Assets       map[string]string `yaml:"assets"`
	Power        string            `yaml:"power"`
	Product      string            `yaml:"product"`
	Guide        *bool             `yaml:"guide"`
	PluginIgnore []string          `yaml:"ignore"`
}

var (
	replaceAccess map[string][]byte = make(map[string][]byte)
	pluginIgnores                   = make(map[string]struct{})
)

func init() {
	data, err := os.ReadFile("customized/custom.yml")
	if err != nil {
		return
	}
	c := new(Config)
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return
	}
	if c.Assets != nil {
		for n, p := range c.Assets {
			if p != "" {
				file, err := os.ReadFile(filepath.Join("customized", p))
				if err != nil {
					continue
				}
				replaceAccess[n] = file
			}
		}
	}

	if c.Power != "" {
		powered = c.Power
	}
	if c.Product != "" {
		product = c.Product
	}
	if c.Guide != nil {
		guide = *c.Guide
	}
	for _, p := range c.PluginIgnore {
		pluginIgnores[p] = struct{}{}
	}
}
func AssetsReplace(name string) ([]byte, bool) {

	path, has := replaceAccess[name]
	return path, has
}

func Powered() string {
	return powered
}
func Produce() string {
	return product
}
func Guide() bool {
	return guide
}

func IgnorePlugin(name string) bool {
	_, ig := pluginIgnores[name]
	return ig
}
