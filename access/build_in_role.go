package access

import (
	_ "embed"
	"gopkg.in/yaml.v3"
)

//go:embed config/role.yml
var roleData []byte
var (
	roleList  []*BuildInRole
	roleMap   map[string]*BuildInRole
	roleTitle map[string]struct{}
)

type BuildInRoleConf struct {
	ID        int      `yaml:"id"`
	Title     string   `yaml:"title"`
	Uuid      string   `yaml:"uuid"`
	Desc      string   `yaml:"desc"`
	IsAddable bool     `yaml:"is_addable"`
	Module    string   `yaml:"module"`
	Access    []string `yaml:"access"`
}

type BuildInRole struct {
	ID        int
	Title     string
	Uuid      string
	Desc      string
	IsAddable bool
	Module    string
	AccessID  []Access
	Access    []string
}

func initBuildInRoleConfig() {
	conf := make([]*BuildInRoleConf, 1)
	err := yaml.Unmarshal(roleData, &conf)
	if err != nil {
		panic(err)
	}
	rm := make(map[string]*BuildInRole, len(conf))
	rl := make([]*BuildInRole, len(conf))
	rt := make(map[string]struct{}, len(conf))
	for i, role := range conf {
		if role.Uuid == "" {
			panic("build-in role's uuid can't be nil. ")
		}
		if role.Title == "" {
			panic("build-in role's uuid can't be nil. ")
		}

		accessID := make([]Access, len(role.Access))
		for idx, key := range role.Access {
			id, err := Parse(key)
			if err != nil {
				panic(err)
			}
			accessID[idx] = id
		}

		buildInRole := &BuildInRole{
			ID:        role.ID,
			Title:     role.Title,
			Uuid:      role.Uuid,
			Desc:      role.Desc,
			IsAddable: role.IsAddable,
			Module:    role.Module,
			AccessID:  accessID,
			Access:    role.Access,
		}

		rm[role.Uuid] = buildInRole
		rl[i] = buildInRole
		rt[role.Title] = struct{}{}
	}

	roleMap = rm
	roleList = rl
	roleTitle = rt
}

func IsBuildInRole(roleUUID string) bool {
	_, has := roleMap[roleUUID]
	return has
}

//IsBuildInRoleTitle 判断有无和内置角色title重复
func IsBuildInRoleTitle(title string) bool {
	_, has := roleTitle[title]
	return has
}

func GetBuildInRoleList() []*BuildInRole {
	return roleList
}

func GetBuildInRoleMap() map[string]*BuildInRole {
	return roleMap
}

func GetBuildInRole(roleUUID string) *BuildInRole {
	return roleMap[roleUUID]
}

func GetAdminRole() *BuildInRole {
	return roleMap["super-admin"]
}
