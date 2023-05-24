package common

func MapToSlice[K comparable, T any, D any](m map[K]T, f func(k K, t T) D) []D {
	r := make([]D, 0, len(m))
	for k, v := range m {
		r = append(r, f(k, v))
	}
	return r
}
func SliceToMap[K comparable, T any](list []T, f func(T) K) map[K]T {
	m := make(map[K]T)
	for _, t := range list {
		m[f(t)] = t
	}
	return m
}
func SliceToMapO[K comparable, T, D any](list []T, f func(T) (K, D)) map[K]D {
	m := make(map[K]D)
	for _, t := range list {
		k, v := f(t)
		m[k] = v
	}
	return m
}

func SliceToMapArray[K comparable, T any](list []T, f func(T) K) map[K][]T {
	m := make(map[K][]T)
	for _, t := range list {
		m[f(t)] = append(m[f(t)], t)
	}
	return m
}

func SliceToSliceIds[K comparable, T any](list []T, f func(T) K) []K {
	ids := make([]K, 0)
	for _, t := range list {
		ids = append(ids, f(t))
	}
	return ids
}
func SliceToSlice[S, D any](list []S, f func(S) D) []D {
	ids := make([]D, 0)
	for _, t := range list {
		ids = append(ids, f(t))
	}
	return ids
}
func CopyMaps[K comparable, T any](maps map[K]T) map[K]T {

	temp := make(map[K]T, len(maps))
	for k, t := range maps {
		temp[k] = t
	}

	return temp
}
