package utils

import "sync"

// Set representa um conjunto de elementos únicos, seguro para uso concorrente.
type Set[T comparable] struct {
	data  map[T]struct{}
	mutex sync.Mutex
}

// NewSet cria uma nova instância de Set.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{data: make(map[T]struct{})}
}

// Add adiciona um elemento ao conjunto.
func (s *Set[T]) Add(item T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[item] = struct{}{}
}

// Remove remove um elemento do conjunto.
func (s *Set[T]) Remove(item T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.data, item)
}

// Contains verifica se um elemento está presente no conjunto.
func (s *Set[T]) Contains(item T) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, exists := s.data[item]
	return exists
}

// Size retorna a quantidade de elementos no conjunto.
func (s *Set[T]) Size() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return len(s.data)
}

// Items retorna um slice com todos os elementos do conjunto.
func (s *Set[T]) Items() []T {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	items := make([]T, 0, len(s.data))
	for item := range s.data {
		items = append(items, item)
	}
	return items
}

// Clear remove todos os elementos do conjunto.
func (s *Set[T]) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data = make(map[T]struct{})
}

// ForEach executa uma função para cada elemento do conjunto.
func (s *Set[T]) ForEach(f func(T)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for item := range s.data {
		f(item)
	}
}
