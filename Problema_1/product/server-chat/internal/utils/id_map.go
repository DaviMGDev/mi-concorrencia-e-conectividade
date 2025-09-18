package utils

import (
	"sync"

	"github.com/google/uuid"
)

type Set[T comparable] struct {
	self  map[uuid.UUID]T
	mutex sync.Mutex
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		self:  make(map[uuid.UUID]T),
		mutex: sync.Mutex{},
	}
}

func (s *Set[T]) Add(value T) uuid.UUID {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	id := uuid.New()
	s.self[id] = value
	return id
}

func (s *Set[T]) Remove(id uuid.UUID) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.self, id)
}

func (s *Set[T]) Get(id uuid.UUID) (T, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	value, exists := s.self[id]
	return value, exists
}

func (s *Set[T]) Values() []T {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	values := make([]T, 0, len(s.self))
	for _, value := range s.self {
		values = append(values, value)
	}
	return values
}

func (s *Set[T]) Size() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return len(s.self)
}

func (s *Set[T]) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.self = make(map[uuid.UUID]T)
}

func (s *Set[T]) Contains(value T) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, v := range s.self {
		if v == value {
			return true
		}
	}
	return false
}

func (s *Set[T]) IDs() []uuid.UUID {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	ids := make([]uuid.UUID, 0, len(s.self))
	for id := range s.self {
		ids = append(ids, id)
	}
	return ids
}
