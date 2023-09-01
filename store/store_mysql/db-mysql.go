/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package store_mysql

import (
	"context"
	"github.com/eolinker/apinto-dashboard/store"
	"gorm.io/gorm"
)

var (
	_ store.IDB = (*myDB)(nil)
)

type myDB struct {
	db   *gorm.DB
	info store.DBInfo
}

func NewMyDB(db *gorm.DB, info store.DBInfo) store.IDB {
	return &myDB{db: db, info: info}
}

func (m *myDB) Info() store.DBInfo {
	return m.info
}

func (m *myDB) DB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(store.TxContextKey).(*gorm.DB); ok {
		return tx
	}
	return m.db.WithContext(ctx)
}
func (m *myDB) IsTxCtx(ctx context.Context) bool {
	if _, ok := ctx.Value(store.TxContextKey).(*gorm.DB); ok {
		return ok
	}
	return false
}
