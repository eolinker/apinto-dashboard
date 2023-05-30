/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package service

import (
	"github.com/eolinker/apinto-dashboard/modules/system"
	"github.com/eolinker/eosc/common/bean"
)

type systemConfigServiceSimpleIml struct {
	key         string
	dataService system.ISystemConfigDataService
}

func (s *systemConfigServiceSimpleIml) Get() ([]byte, error) {
	return s.dataService.Get(s.key)
}

func (s *systemConfigServiceSimpleIml) Set(v []byte) error {
	return s.dataService.Set(s.key, v)
}

func NewSystemConfigServiceSimpleIml(key string) system.ISystemConfigServiceSimple {
	s := &systemConfigServiceSimpleIml{key: key}
	bean.Autowired(&s.dataService)
	return s
}
