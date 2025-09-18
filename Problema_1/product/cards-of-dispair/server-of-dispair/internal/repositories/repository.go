package repositories

import (
	"errors"
	"server-of-dispair/internal/utils"
)

type RepositoryInterface[T any] interface {
	Create(id string, entity T) error
	Read(id string) (T, error)
	Update(id string, entity T) error
	Delete(id string) error
	List() ([]T, error)
}

type InMemoryRepository[T any] struct {
	data *utils.Map[string, T]
}

func NewInMemoryRepository[T any]() *InMemoryRepository[T] {
	return &InMemoryRepository[T]{
		data: utils.NewMap[string, T](),
	}
}

func (r *InMemoryRepository[T]) Create(id string, entity T) error {
	if _, exists := r.data.Get(id); exists {
		return errors.New("entity already exists")
	}
	r.data.Set(id, entity)
	return nil
}

func (r *InMemoryRepository[T]) Read(id string) (T, error) {
	if entity, exists := r.data.Get(id); exists {
		return entity, nil
	}
	var zero T
	return zero, errors.New("entity not found")
}

func (r *InMemoryRepository[T]) Update(id string, entity T) error {
	if _, exists := r.data.Get(id); !exists {
		return errors.New("entity not found")
	}
	r.data.Set(id, entity)
	return nil
}

func (r *InMemoryRepository[T]) Delete(id string) error {
	if _, exists := r.data.Get(id); !exists {
		return errors.New("entity not found")
	}
	r.data.Delete(id)
	return nil
}

func (r *InMemoryRepository[T]) List() ([]T, error) {
	entities := make([]T, 0, r.data.Len())
	for _, entity := range r.data.Values() {
		entities = append(entities, entity)
	}
	return entities, nil
}
