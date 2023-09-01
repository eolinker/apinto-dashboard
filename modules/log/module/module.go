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
)

type Module struct {
	name    string
	routers apinto_module.RoutersInfo
}

func (m *Module) RoutersInfo() apinto_module.RoutersInfo {
	return m.routers
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) Routers() (apinto_module.Routers, bool) {
	return m, true
}

func (m *Module) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func (m *Module) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func NewModule(name string) *Module {
	m := &Module{name: name, routers: NewController().RoutersInfo()}

	return m
}
