package apinto_module

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	eventHandlerLocker = sync.RWMutex{}
	eventHandlers      = map[string]any{}
)

func RegisterEventHandler[K comparable](event K, handler func(event K, v any)) {
	hm := getCreateHandlerManager(event)
	hm.Register(event, handler)
}
func getCreateHandlerManager[K comparable](event K) EventHandlerManager[K] {
	t := reflect.TypeOf(event)
	tName := fmt.Sprintf("%s:%s", t.PkgPath(), t.Name())

	eventHandlerLocker.Lock()
	defer eventHandlerLocker.Unlock()
	hst, has := eventHandlers[tName]
	if !has {
		hst = newTEventHandlerManager[K]()
		eventHandlers[tName] = hst
	}
	hm, ok := hst.(EventHandlerManager[K])
	if !ok {
		hm = newTEventHandlerManager[K]()
		eventHandlers[tName] = hm
	}
	return hm
}
func DoEvent[K comparable](event K, v any) {
	hm := getCreateHandlerManager(event)
	hm.DoEvent(event, v)
}

type EventHandlerManager[K comparable] interface {
	Register(event K, handler func(event K, v any))
	DoEvent(event K, v interface{})
}

type tEventHandlerManager[K comparable] struct {
	locker   sync.RWMutex
	handlers map[K][]func(k K, v any)
}

func newTEventHandlerManager[K comparable]() *tEventHandlerManager[K] {
	return &tEventHandlerManager[K]{
		handlers: make(map[K][]func(k K, v any)),
	}
}

func (m *tEventHandlerManager[K]) Register(event K, handler func(k K, v any)) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.handlers[event] = append(m.handlers[event], handler)
}

func (m *tEventHandlerManager[K]) DoEvent(event K, v interface{}) {
	m.locker.RLock()
	hs := m.handlers[event]
	m.locker.RUnlock()
	for _, h := range hs {
		h(event, v)
	}
}
