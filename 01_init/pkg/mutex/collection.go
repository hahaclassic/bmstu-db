package mutex

import "sync"

type Collection[T comparable] struct {
	mu   *sync.RWMutex
	data map[T]struct{}
}

func NewCollection[T comparable]() *Collection[T] {
	return &Collection[T]{
		mu:   &sync.RWMutex{},
		data: make(map[T]struct{}),
	}
}

func (c *Collection[T]) Store(elem T) {
	c.mu.Lock()
	c.data[elem] = struct{}{}
	c.mu.Unlock()
}

func (c *Collection[T]) Contains(elem T) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.data[elem]

	return ok
}

func (c *Collection[T]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.data)
}
