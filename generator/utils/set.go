package utils

import (
	"maps"
	"slices"
)

type Set[T comparable] map[T]struct{}

func (s Set[T]) Contains(t T) bool {
	_, ok := s[t]
	return ok
}

func (s Set[T]) Add(t ...T) {
	for _, val := range t {
		s[val] = struct{}{}
	}
}

func (s Set[T]) AsSlice() []T {
	return slices.Collect(maps.Keys(s))
}

func SetFrom[T comparable](vals ...T) Set[T] {
	out := make(Set[T], len(vals))
	out.Add(vals...)
	return out
}
