package common

import (
	"fmt"
	"reflect"
	"sync"
)

type IPool[T any] interface {
	Get() *T
	Put(t *T)
}
type StructPool[T any] struct {
	pool *sync.Pool
}

func (s *StructPool[T]) Get() *T {
	return s.pool.Get().(*T)
}
func (s *StructPool[T]) Put(t *T) {
	s.pool.Put(t)
}

var (
	poolLocker sync.RWMutex
	pools      = make(map[string]*sync.Pool)
)

func CreatePool[T any]() IPool[T] {
	fi := new(T)
	tp := reflect.TypeOf(fi)
	tpN := fmt.Sprintf("%s:%s", tp.PkgPath(), tp.String())
	poolLocker.RLock()
	p, has := pools[tpN]
	poolLocker.RUnlock()
	if has {
		p.Put(fi)
		return &StructPool[T]{
			pool: p,
		}
	}
	poolLocker.Lock()
	defer poolLocker.Unlock()
	p, has = pools[tpN]
	if has {
		p.Put(fi)
		return &StructPool[T]{
			pool: p,
		}
	}
	p = &sync.Pool{New: func() any { return new(T) }}
	pools[tpN] = p
	p.Put(fi)
	return &StructPool[T]{
		pool: p,
	}
}
