package mutex

import "sync"

type Slice[T any] struct {
	mu   *sync.RWMutex
	data []T
}

func NewSlice[T any]() *Slice[T] {
	return &Slice[T]{
		mu:   &sync.RWMutex{},
		data: make([]T, 0),
	}
}

func (s *Slice[T]) Add(elem T) {
	s.mu.Lock()
	s.data = append(s.data, elem)
	s.mu.Unlock()
}

func (s *Slice[T]) Get(idx int) T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if idx < s.Len() {
		return s.data[idx]
	}

	var val T

	return val
}

func (s *Slice[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}
