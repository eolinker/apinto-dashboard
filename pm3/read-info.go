package pm3

import (
	"fmt"
)

func ReadPluginAssembly(info *PluginDefine) (ms []PModule, acs []PAccess, fs []PFrontend, err error) {
	pName := info.Id
	// frontend
	fs = make([]PFrontend, 0, len(info.Frontend))
	for _, f := range info.Frontend {
		fs = append(fs, f)
	}

	// modules

	acs = make([]PAccess, 0, len(ms)*2+len(info.Access))

	ms = make([]PModule, 0, len(info.Navigations))
	for _, n := range info.Navigations {

		m := PModule{
			Navigation: n.Navigation,
			Name:       fmt.Sprintf("%s.%s", pName, n.Name),
			Cname:      n.Cname,
			Router:     n.Router,
		}

		ms = append(ms, m)
		if len(n.Access) > 0 {
			for _, ac := range n.Access {
				acs = append(acs, createAccess(m.Name, ac.Name, ac.Cname, ac.Depend...))
			}
		} else {
			acs = append(acs, typicalAccess(m.Name)...)
		}
	}

	//access

	if info.Access != nil {
		for md, acList := range info.Access {
			for _, ac := range acList {
				acs = append(acs, PAccess{
					Name:   fmt.Sprintf("%s.%s", pName, ac.Name),
					Cname:  ac.Cname,
					Module: md,
					Depend: ac.Depend,
				})
			}
		}
	}

	return
}

func typicalAccess(mName string) []PAccess {
	return []PAccess{

		createAccess(mName, fmt.Sprintf("%s.view", mName), "查看"),

		createAccess(mName, fmt.Sprintf("%s.edit", mName), "编辑", fmt.Sprintf("%s.view", mName)),
	}
}
func createAccess(mname, name, cname string, depend ...string) PAccess {
	return PAccess{
		Name:   name,
		Cname:  cname,
		Module: mname,
		Depend: depend,
	}
}
