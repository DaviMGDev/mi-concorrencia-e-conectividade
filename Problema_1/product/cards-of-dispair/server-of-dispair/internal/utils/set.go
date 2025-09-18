package utils

type Set[T comparable] Map[T, struct{}]

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		data: make(map[T]struct{}),
	}
}

func (s *Set[T]) Add(value T) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.data[value] = struct{}{}
}

func (s *Set[T]) Remove(value T) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	delete(s.data, value)
}

func (s *Set[T]) Contains(value T) bool {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	_, exists := s.data[value]
	return exists
}

func (s *Set[T]) Values() []T {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	values := make([]T, 0, len(s.data))
	for key := range s.data {
		values = append(values, key)
	}
	return values
}
