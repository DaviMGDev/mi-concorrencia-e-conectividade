// COMPLETED
package utils

import "sync"

type Map[K comparable, V any] struct {
	data  map[K]V
	Mutex sync.Mutex
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		data: make(map[K]V),
	}
}

func (m *Map[K, V]) Set(key K, value V) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.data[key] = value
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	value, exists := m.data[key]
	return value, exists
}

func (m *Map[K, V]) Delete(key K) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	delete(m.data, key)
}

func (m *Map[K, V]) Keys() []K {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	keys := make([]K, 0, len(m.data))
	for key := range m.data {
		keys = append(keys, key)
	}
	return keys
}

func (m *Map[K, V]) Values() []V {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	values := make([]V, 0, len(m.data))
	for _, value := range m.data {
		values = append(values, value)
	}
	return values
}

func (m *Map[K, V]) Len() int {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	return len(m.data)
}

func (m *Map[K, V]) Clear() {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.data = make(map[K]V)
}

func (m *Map[K, V]) Has(key K) bool {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	_, exists := m.data[key]
	return exists
}

func (m *Map[K, V]) ForEach(f func(key K, value V)) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	for key, value := range m.data {
		f(key, value)
	}
}
