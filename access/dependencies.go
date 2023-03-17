package access

import (
	_ "embed"
	"gopkg.in/yaml.v3"
)

//go:embed config/dependencies.yml
var dependenciesData []byte
var (
	dependenciesMap map[string][]string
)

func initDependencies() {
	dependenciesMap = make(map[string][]string)
	err := yaml.Unmarshal(dependenciesData, &dependenciesMap)
	if err != nil {
		panic(err)
	}

}
