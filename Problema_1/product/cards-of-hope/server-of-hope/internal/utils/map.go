package utils

import "sync"

// Map representa um mapa seguro para uso concorrente.
type Map[K comparable, V any] struct {
	data  map[K]V
	mutex sync.Mutex
}

// NewMap cria uma nova instância de Map.
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{data: make(map[K]V)}
}

// Set define o valor de uma chave no mapa.
func (m *Map[K, V]) Set(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[key] = value
}

// Get retorna o valor associado a uma chave e se ela existe no mapa.
func (m *Map[K, V]) Get(key K) (V, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	value, exists := m.data[key]
	return value, exists
}

// Delete remove uma chave e seu valor do mapa.
func (m *Map[K, V]) Delete(key K) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.data, key)
}

// Keys retorna um slice com todas as chaves do mapa.
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
func (m *Map[K, V]) Values() []V {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	values := make([]V, 0, len(m.data))
	for _, value := range m.data {
		values = append(values, value)
	}
	return values
}

// Size retorna a quantidade de pares chave-valor no mapa.
func (m *Map[K, V]) Size() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return len(m.data)
}

// Clear remove todos os pares chave-valor do mapa.
func (m *Map[K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data = make(map[K]V)
}

// ForEach executa uma função para cada par chave-valor do mapa.
func (m *Map[K, V]) ForEach(f func(K, V)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for key, value := range m.data {
		f(key, value)
	}
}
