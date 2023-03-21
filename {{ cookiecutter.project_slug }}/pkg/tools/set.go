package tools

import (
	"sync"
)

type Set[T comparable] struct {
	mu  sync.RWMutex
	set map[T]struct{}
}

func NewSetWithCapacity[T comparable](c int) *Set[T] {
	set := make(map[T]struct{}, c)
	return &Set[T]{set: set}
}

func NewSet[T comparable](elems ...T) *Set[T] {
	set := make(map[T]struct{}, len(elems))
	for _, elem := range elems {
		set[elem] = struct{}{}
	}
	return &Set[T]{set: set}
}

func (s *Set[T]) Add(ele ...T) *Set[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range ele {
		s.set[e] = struct{}{}
	}
	return s
}

func (s *Set[T]) Del(ele ...T) *Set[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range ele {
		delete(s.set, e)
	}
	return s
}

func (s *Set[T]) Has(ele T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.set[ele]
	return exists
}

// 取交集.
func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ret := NewSet[T]()
	for ele := range s.set {
		if other.Has(ele) {
			ret.Add(ele)
		}
	}
	return ret
}

func (s *Set[T]) Len() int {
	return len(s.set)
}

func (s *Set[T]) Set() map[T]struct{} {
	return s.set
}

func (s *Set[T]) Values() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ret := make([]T, 0, len(s.set))
	for ele := range s.set {
		ret = append(ret, ele)
	}
	return ret
}

func (s *Set[T]) Compare(other *Set[T]) (less, more *Set[T]) {
	less = NewSet[T]()
	more = NewSet[T]()
	// less: s - other
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()
	for ele := range other.set {
		if !s.Has(ele) {
			less.Add(ele)
		}
	}
	// more: other - s
	for ele := range s.set {
		if !other.Has(ele) {
			more.Add(ele)
		}
	}
	return less, more
}
