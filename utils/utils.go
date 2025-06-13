package utils

import "time"

func MapSlice[A any, B any](in []A, mapper func(A) B) []B {
	out := make([]B, len(in))
	for i, v := range in {
		out[i] = mapper(v)
	}
	return out
}

func SliceToMap[A any, K comparable](in []A, keySelector func(A) K) map[K]A {
	m := make(map[K]A, len(in))
	for _, v := range in {
		k := keySelector(v)
		m[k] = v
	}
	return m
}

func ToPointerTime(t time.Time) *time.Time {
	return &t
}
