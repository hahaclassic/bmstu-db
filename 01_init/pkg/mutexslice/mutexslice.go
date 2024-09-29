package mutexslice

import "sync"

type Slice[T any] struct {
	mu   *sync.RWMutex
	data []T
}

func (s *Slice[T]) Add(elem T) {
	s.mu.Lock()
	s.data = append(s.data, elem)
	s.mu.Unlock()
}

func (s *Slice[T]) Get(idx int) T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data[idx]
}

func (s *Slice[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}
