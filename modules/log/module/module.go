/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package module

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
)

type Module struct {
	*pm3.ModuleTool

	name    string
	routers apinto_module.RoutersInfo
}

func (m *Module) Frontend() []pm3.FrontendAsset {
	return nil
}

func (m *Module) Apis() []pm3.Api {
	return m.routers
}

func (m *Module) Middleware() []pm3.Middleware {
	return nil
}

func (m *Module) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func (m *Module) Name() string {
	return m.name
}

func NewModule(id, name string) apinto_module.Module {
	m := &Module{ModuleTool: pm3.NewModuleTool(id, name),
		name: name, routers: NewController().RoutersInfo()}
	m.InitAccess(m.routers)
	return m
}
