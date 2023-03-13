package model

import "github.com/eolinker/apinto-dashboard/access"

type RoleAccess map[access.Access]struct{}

// IsAccess 只要符合一个传入的权限id就算通过
func (ra RoleAccess) IsAccess(access ...access.Access) bool {
	if len(ra) == 0 {
		return false
	}

	for _, a := range access {
		if _, ok := ra[a]; ok {
			return true
		}
	}
	return false
}

func CreateRoleAccess(acs ...access.Access) *RoleAccess {
	ra := make(RoleAccess, len(acs))
	for _, ac := range acs {
		ra[ac] = struct{}{}
	}
	return &ra
}
