package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
	"go/ast"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"reflect"
)

type Table interface {
	schema.Tabler
	IdValue() int
}
type IDB interface {
	DB(ctx context.Context) *gorm.DB
	IsTxCtx(ctx context.Context) bool
}
type myDB struct {
	db *gorm.DB
}

func (m *myDB) DB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txContextKey{}).(*gorm.DB); ok {
		return tx
	}
	return m.db.WithContext(ctx)
}
func (m *myDB) IsTxCtx(ctx context.Context) bool {
	if _, ok := ctx.Value(txContextKey{}).(*gorm.DB); ok {
		return ok
	}
	return false
}

type IBaseStore[T any] interface {
	IDB
	Get(ctx context.Context, id int) (*T, error)
	Save(ctx context.Context, t *T) error
	UpdateByUnique(ctx context.Context, t *T, uniques []string) error
	Delete(ctx context.Context, id int) (int, error)
	UpdateWhere(ctx context.Context, t *T, m map[string]interface{}) (int, error)
	Update(ctx context.Context, t *T) (int, error)
	DeleteWhere(ctx context.Context, m map[string]interface{}) (int, error)
	Insert(ctx context.Context, t ...*T) error
	List(ctx context.Context, m map[string]interface{}) ([]*T, error)
	ListQuery(ctx context.Context, sql string, args []interface{}, order string) ([]*T, error)
	First(ctx context.Context, m map[string]interface{}) (*T, error)
	FirstQuery(ctx context.Context, sql string, args []interface{}, order string) (*T, error)
	ListPage(ctx context.Context, sql string, pageNum, pageSize int, args []interface{}, order string) ([]*T, int, error)
	Transaction(ctx context.Context, f func(txCtx context.Context) error) error
}

type txContextKey struct{}

type baseStore[T any] struct {
	IDB
	uniqueList []string
	targetType *T
}

func createStore[T any](db IDB) *baseStore[T] {
	b := &baseStore[T]{
		IDB:        db,
		targetType: new(T),
	}
	modelType := reflect.TypeOf(new(T)).Elem()
	for i := 0; i < modelType.NumField(); i++ {
		if fieldStruct := modelType.Field(i); ast.IsExported(fieldStruct.Name) {
			tagSetting := schema.ParseTagSetting(fieldStruct.Tag.Get("gorm"), ";")
			if _, ok := tagSetting["DBUNIQUEINDEX"]; ok {
				b.uniqueList = append(b.uniqueList, tagSetting["COLUMN"])
			}
		}
	}

	return b
}

func (b *baseStore[T]) Get(ctx context.Context, id int) (*T, error) {
	value := new(T)
	err := b.DB(ctx).First(value, id).Error
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (b *baseStore[T]) Save(ctx context.Context, t *T) error {

	var v interface{} = t
	if table, ok := v.(Table); ok {

		if table.IdValue() != 0 {
			return b.DB(ctx).Save(t).Error
		}
		//没查到主键ID的数据 看看有没有唯一索引 有唯一索引 用唯一索引更新所有字段
		if len(b.uniqueList) > 0 {
			return b.UpdateByUnique(ctx, t, b.uniqueList)
		}
	}
	return b.Insert(ctx, t)
}

func (b *baseStore[T]) UpdateByUnique(ctx context.Context, t *T, uniques []string) error {
	columns := make([]clause.Column, 0, len(uniques))
	for _, unique := range uniques {
		columns = append(columns, clause.Column{
			Name: unique,
		})
	}
	return b.DB(ctx).Clauses(clause.OnConflict{
		Columns:   columns,
		UpdateAll: true,
	}).Create(t).Error
}

func (b *baseStore[T]) Delete(ctx context.Context, id int) (int, error) {

	result := b.DB(ctx).Delete(b.targetType, id)

	return int(result.RowsAffected), result.Error
}

func (b *baseStore[T]) UpdateWhere(ctx context.Context, t *T, m map[string]interface{}) (int, error) {

	result := b.DB(ctx).Model(t).Updates(m)

	return int(result.RowsAffected), result.Error
}

func (b *baseStore[T]) Update(ctx context.Context, t *T) (int, error) {

	result := b.DB(ctx).Updates(t)

	return int(result.RowsAffected), result.Error
}
func (b *baseStore[T]) DeleteWhere(ctx context.Context, m map[string]interface{}) (int, error) {
	if len(m) == 0 {
		return 0, gorm.ErrMissingWhereClause
	}
	result := b.DB(ctx).Where(m).Delete(b.targetType)

	return int(result.RowsAffected), result.Error
}

func (b *baseStore[T]) Insert(ctx context.Context, t ...*T) error {

	return b.DB(ctx).Create(t).Error
}

func (b *baseStore[T]) List(ctx context.Context, m map[string]interface{}) ([]*T, error) {
	list := make([]*T, 0)

	err := b.DB(ctx).Where(m).Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}
func (b *baseStore[T]) ListQuery(ctx context.Context, where string, args []interface{}, order string) ([]*T, error) {
	list := make([]*T, 0)
	db := b.DB(ctx)
	db = db.Where(where, args...)
	if order != "" {
		db = db.Order(order)
	}
	err := db.Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (b *baseStore[T]) First(ctx context.Context, m map[string]interface{}) (*T, error) {
	value := new(T)
	db := b.DB(ctx)

	err := db.Where(m).First(value).Error
	if err != nil {
		return nil, err
	}

	return value, nil
}
func (b *baseStore[T]) FirstQuery(ctx context.Context, where string, args []interface{}, order string) (*T, error) {
	value := new(T)
	db := b.DB(ctx)
	if order != "" {
		db = db.Order(order)
	}
	err := db.Where(where, args...).Take(value).Error
	if err != nil {
		return nil, err
	}

	return value, nil
}
func (b *baseStore[T]) ListPage(ctx context.Context, where string, pageNum, pageSize int, args []interface{}, order string) ([]*T, int, error) {
	list := make([]*T, 0, pageSize)
	db := b.DB(ctx).Where(where, args...)
	if order != "" {
		db = db.Order(order)
	}
	count := int64(0)
	err := db.Model(list).Count(&count).Limit(pageSize).Offset(entry.PageIndex(pageNum, pageSize)).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, int(count), nil
}

// Transaction 执行事务
func (b *baseStore[T]) Transaction(ctx context.Context, f func(context.Context) error) error {
	return b.DB(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txContextKey{}, tx)
		return f(txCtx)
	})
}
