/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package service

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/modules/system"
	"github.com/eolinker/eosc/common/bean"
)

type systemServiceIml[T any] struct {
	key         string
	dataService system.ISystemConfigDataService
}

func NewSystemServiceIml[T any](key string) system.ISystemConfigService[T] {
	v := &systemServiceIml[T]{
		key: key,
	}
	bean.Autowired(&v.dataService)

	return v
}

func (s *systemServiceIml[T]) Get() (*T, error) {
	data, err := s.dataService.Get(s.key)

	v := new(T)
	err = json.Unmarshal(data, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (s *systemServiceIml[T]) Set(v *T) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return s.dataService.Set(s.key, data)
}
