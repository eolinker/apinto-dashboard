package common

type Set[K comparable] map[K]struct{}

func (s Set[K]) Has(k K) bool {
	_, has := s[k]
	return has
}
func (s Set[K]) Remove(k K) {
	delete(s, k)
}
func (s Set[K]) Set(k K) {
	s[k] = struct{}{}
}

type anySet[K comparable, T any] struct {
	setData map[K]T
}

func (a *anySet[K, T]) Has(k K) bool {
	_, h := a.setData[k]
	return h
}

func (a *anySet[K, T]) Remove(k K) (T, bool) {
	t, h := a.setData[k]
	if h {
		delete(a.setData, k)
	}
	return t, h
}

type AnySet[K comparable, T any] interface {
	Has(k K) bool
	Remove(k K) (T, bool)
}

func NewAnySet[K comparable, T any](setData map[K]T) AnySet[K, T] {
	return &anySet[K, T]{
		setData: setData,
	}
}
