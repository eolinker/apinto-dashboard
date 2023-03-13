package common

func SliceToMap[K comparable, T any](list []T, f func(T) K) map[K]T {
	m := make(map[K]T)
	for _, t := range list {
		m[f(t)] = t
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

func CopyMaps[K comparable, T any](maps map[K]T) map[K]T {

	temp := make(map[K]T, len(maps))
	for k, t := range maps {
		temp[k] = t
	}

	return temp
}
