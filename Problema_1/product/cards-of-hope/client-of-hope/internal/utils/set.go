// Pacote utils fornece estruturas utilitárias genéricas, incluindo o tipo Set.
package utils

import "sync"

// Set representa um conjunto de elementos únicos do tipo T.
//
// Campos:
//   - data: armazena os elementos do conjunto.
//   - mutex: garante acesso concorrente seguro ao conjunto.
type Set[T comparable] struct {
	data  map[T]struct{}
	mutex sync.Mutex
}

// NewSet cria e retorna um novo conjunto (Set) vazio.
//
// Retorno:
//   - *Set[T]: ponteiro para o novo conjunto.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{data: make(map[T]struct{})}
}

// Add adiciona um item ao conjunto.
//
// Parâmetros:
//   - item: elemento a ser adicionado.
func (s *Set[T]) Add(item T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[item] = struct{}{}
}

// Remove remove um item do conjunto.
//
// Parâmetros:
//   - item: elemento a ser removido.
func (s *Set[T]) Remove(item T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.data, item)
}

// Contains verifica se um item está presente no conjunto.
//
// Parâmetros:
//   - item: elemento a ser verificado.
//
// Retorno:
//   - bool: true se o item estiver presente, false caso contrário.
func (s *Set[T]) Contains(item T) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, exists := s.data[item]
	return exists
}

// Size retorna o número de elementos no conjunto.
//
// Retorno:
//   - int: quantidade de elementos.
func (s *Set[T]) Size() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return len(s.data)
}

// Items retorna um slice com todos os elementos do conjunto.
//
// Retorno:
//   - []T: slice contendo todos os elementos.
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

// ForEach executa a função fornecida para cada elemento do conjunto.
//
// Parâmetros:
//   - f: função a ser executada para cada elemento.
func (s *Set[T]) ForEach(f func(T)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for item := range s.data {
		f(item)
	}
}
