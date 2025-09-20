package data

import (
	"errors"
	"server-of-hope/internal/utils"
)

// RepositoryInterface descreve operações para um repositório genérico de dados.
//
// Métodos:
//   - Create: adiciona um novo item.
//   - Read: retorna um item pelo ID.
//   - Update: atualiza um item existente.
//   - Delete: remove um item pelo ID.
//   - List: retorna todos os itens.
type RepositoryInterface[T any] interface {
	// Create adiciona um novo item ao repositório com o ID especificado.
	//
	// Parâmetros:
	//   - id: identificador do item.
	//   - item: item a ser adicionado.
	//
	// Retorno:
	//   - erro caso não seja possível adicionar.
	Create(id string, item T) error

	// Read retorna o item associado ao ID informado.
	//
	// Parâmetros:
	//   - id: identificador do item.
	//
	// Retorno:
	//   - T: item encontrado.
	//   - erro caso não exista.
	Read(id string) (T, error)

	// Update atualiza o item associado ao ID informado.
	//
	// Parâmetros:
	//   - id: identificador do item.
	//   - item: novo valor do item.
	//
	// Retorno:
	//   - erro caso não exista ou não seja possível atualizar.
	Update(id string, item T) error

	// Delete remove o item associado ao ID informado.
	//
	// Parâmetros:
	//   - id: identificador do item.
	//
	// Retorno:
	//   - erro caso não exista ou não seja possível remover.
	Delete(id string) error

	// List retorna todos os itens do repositório.
	//
	// Retorno:
	//   - slice de itens armazenados.
	//   - erro caso não seja possível listar.
	List() ([]T, error)
}

// InMemoryRepository implementa RepositoryInterface usando um mapa em memória.
//
// Campos:
//   - items: armazena os itens do repositório em memória.
type InMemoryRepository[T any] struct {
	items *utils.Map[string, T]
}

// NewInMemoryRepository cria uma nova instância de InMemoryRepository.
//
// Retorno:
//   - ponteiro para InMemoryRepository.
func NewInMemoryRepository[T any]() *InMemoryRepository[T] {
	return &InMemoryRepository[T]{items: utils.NewMap[string, T]()}
}

// Create adiciona um novo item ao repositório.
//
// Parâmetros:
//   - id: identificador do item.
//   - item: item a ser adicionado.
//
// Retorno:
//   - erro caso não seja possível adicionar.
func (r *InMemoryRepository[T]) Create(id string, item T) error {
	r.items.Set(id, item)
	return nil
}

// Read retorna o item associado ao ID informado, se existir.
//
// Parâmetros:
//   - id: identificador do item.
//
// Retorno:
//   - T: item encontrado.
//   - erro caso não exista.
func (r *InMemoryRepository[T]) Read(id string) (T, error) {
	item, exists := r.items.Get(id)
	if !exists {
		var zero T
		return zero, errors.New("item not found")
	}
	return item, nil
}

// Update atualiza o item associado ao ID informado, se existir.
//
// Parâmetros:
//   - id: identificador do item.
//   - item: novo valor do item.
//
// Retorno:
//   - erro caso não exista ou não seja possível atualizar.
func (r *InMemoryRepository[T]) Update(id string, item T) error {
	_, exists := r.items.Get(id)
	if !exists {
		return errors.New("item not found")
	}
	r.items.Set(id, item)
	return nil
}

// Delete remove o item associado ao ID informado, se existir.
//
// Parâmetros:
//   - id: identificador do item.
//
// Retorno:
//   - erro caso não exista ou não seja possível remover.
func (r *InMemoryRepository[T]) Delete(id string) error {
	_, exists := r.items.Get(id)
	if !exists {
		return errors.New("item not found")
	}
	r.items.Delete(id)
	return nil
}

// List retorna todos os itens armazenados no repositório.
//
// Retorno:
//   - slice de itens armazenados.
//   - erro caso não seja possível listar.
func (r *InMemoryRepository[T]) List() ([]T, error) {
	return r.items.Values(), nil
}
