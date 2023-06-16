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

type Driver struct {
}

func NewDriver() *Driver {
	return &Driver{}
}

func (d *Driver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return new(Plugin), nil
}

type Plugin struct {
}

func (p *Plugin) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewModule(name), nil
}

func (p *Plugin) CheckConfig(name string, config interface{}) error {
	return nil
}

func (p *Plugin) GetPluginFrontend(moduleName string) string {
	return "/log"
}

func (p *Plugin) IsShowServer() bool {
	return false
}
