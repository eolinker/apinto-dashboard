package strategy_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-entry"
	"github.com/eolinker/apinto-dashboard/store"
	"time"
)

type IStrategyStore interface {
	store.IBaseStore[strategy_entry.Strategy]
	GetByUUID(ctx context.Context, uuid string) (*strategy_entry.Strategy, error)
	GetListByType(ctx context.Context, clusterId int, strategyType string) ([]*strategy_entry.Strategy, error)
	GetByPriority(ctx context.Context, clusterId, priority int, strategyType string) (*strategy_entry.Strategy, error)
	GetMaxPriority(ctx context.Context, clusterId int, strategyType string) (int, error)
	UpdateStop(ctx context.Context, id int, isStop bool) error
	UpdatePriority(ctx context.Context, maps map[string]int) error
}

type strategyApiStore struct {
	*store.BaseStore[strategy_entry.Strategy]
}

func newStrategyStore(db store.IDB) IStrategyStore {
	return &strategyApiStore{BaseStore: store.CreateStore[strategy_entry.Strategy](db)}
}

func (s *strategyApiStore) GetByUUID(ctx context.Context, uuid string) (*strategy_entry.Strategy, error) {
	return s.FirstQuery(ctx, "`uuid` = ?", []interface{}{uuid}, "")
}

func (s *strategyApiStore) GetListByType(ctx context.Context, clusterId int, strategyType string) ([]*strategy_entry.Strategy, error) {
	return s.ListQuery(ctx, "`cluster` = ? and `type` = ?", []interface{}{clusterId, strategyType}, "priority asc")
}

func (s *strategyApiStore) UpdateStop(ctx context.Context, id int, isStop bool) error {
	_, err := s.UpdateWhere(ctx, &strategy_entry.Strategy{Id: id}, map[string]interface{}{"`is_stop`": isStop, "update_time": time.Now()})
	return err
}

func (s *strategyApiStore) GetByPriority(ctx context.Context, clusterId, priority int, strategyType string) (*strategy_entry.Strategy, error) {
	strategy := new(strategy_entry.Strategy)
	err := s.DB(ctx).Where("`cluster` = ? and `type` = ? and `priority` = ?", clusterId, strategyType, priority).Find(strategy).Error
	return strategy, err
}

func (s *strategyApiStore) GetMaxPriority(ctx context.Context, clusterId int, strategyType string) (int, error) {
	priority := 0
	db := s.DB(ctx).Table("strategy").Select("IFNULL(MAX(`priority`),0) AS `priority`") //IFNULL MAX 为了处理 N/A默认值的问题
	err := db.Where("`cluster` = ? and  `type` = ? ", clusterId, strategyType).Row().Scan(&priority)
	return priority, err
}

func (s *strategyApiStore) UpdatePriority(ctx context.Context, maps map[string]int) error {

	for uuid, priority := range maps {
		err := s.DB(ctx).Exec("update `strategy` set `priority` = ? where `uuid` = ?", priority, uuid).Error
		if err != nil {
			return err
		}
	}

	return nil
}
