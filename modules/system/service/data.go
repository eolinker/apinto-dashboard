/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/system"
	system_store "github.com/eolinker/apinto-dashboard/modules/system/system-store"
	"github.com/eolinker/eosc/common/bean"
	"time"
)

type dataServiceIml struct {
	cache cache.ICommonCache
	store system_store.ISystemInfoStore
}

func newDataServiceIml() system.ISystemConfigDataService {
	s := &dataServiceIml{}
	bean.Autowired(&s.cache)
	bean.Autowired(&s.store)
	return s
}

func (d *dataServiceIml) Get(key string) ([]byte, error) {
	data, err := d.cache.Get(context.Background(), key)
	if err == nil {
		return data, nil
	}
	data, err = d.store.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}
	_ = d.cache.Set(context.Background(), key, data, time.Minute*5)
	return data, nil
}

func (d *dataServiceIml) Set(key string, v []byte) error {
	err := d.store.Set(context.Background(), key, v)
	if err != nil {
		return err
	}
	_ = d.cache.Set(context.Background(), key, v, time.Minute*5)
	return nil
}
