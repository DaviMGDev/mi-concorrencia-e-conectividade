package utils

// Pacote utils fornece estruturas utilitárias genéricas, incluindo o tipo Map.

import "sync"

// Map representa um mapa genérico com suporte a concorrência.
//
// Campos:
//   - data: armazena os pares chave-valor do mapa.
//   - mutex: garante acesso concorrente seguro ao mapa.
type Map[K comparable, V any] struct {
	data  map[K]V
	mutex sync.Mutex
}

// NewMap cria e retorna um novo mapa (Map) vazio.
//
// Retorno:
//   - *Map[K, V]: ponteiro para o novo mapa.
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{data: make(map[K]V)}
}

// Set define o valor para uma chave no mapa.
//
// Parâmetros:
//   - key: chave a ser definida.
//   - value: valor a ser associado à chave.
func (m *Map[K, V]) Set(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[key] = value
}

// Get retorna o valor associado a uma chave e se ela existe no mapa.
//
// Parâmetros:
//   - key: chave a ser consultada.
//
// Retorno:
//   - V: valor associado à chave.
//   - bool: true se a chave existe, false caso contrário.
func (m *Map[K, V]) Get(key K) (V, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	value, exists := m.data[key]
	return value, exists
}

// Delete remove uma chave e seu valor do mapa.
//
// Parâmetros:
//   - key: chave a ser removida.
func (m *Map[K, V]) Delete(key K) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.data, key)
}

// Keys retorna um slice com todas as chaves do mapa.
//
// Retorno:
//   - []K: slice contendo todas as chaves.
func (m *Map[K, V]) Keys() []K {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	keys := make([]K, 0, len(m.data))
	for key := range m.data {
		keys = append(keys, key)
	}
	return keys
}

// Values retorna um slice com todos os valores do mapa.
//
// Retorno:
//   - []V: slice contendo todos os valores.
func (m *Map[K, V]) Values() []V {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	values := make([]V, 0, len(m.data))
	for _, value := range m.data {
		values = append(values, value)
	}
	return values
}

// Size retorna o número de elementos no mapa.
//
// Retorno:
//   - int: quantidade de elementos.
func (m *Map[K, V]) Size() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return len(m.data)
}

// Clear remove todos os elementos do mapa.
func (m *Map[K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data = make(map[K]V)
}

// ForEach executa a função fornecida para cada par chave-valor do mapa.
//
// Parâmetros:
//   - f: função a ser executada para cada par chave-valor.
func (m *Map[K, V]) ForEach(f func(K, V)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for key, value := range m.data {
		f(key, value)
	}
}
