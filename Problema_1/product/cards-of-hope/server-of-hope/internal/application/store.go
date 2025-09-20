package application

import "server-of-hope/internal/domain"

// StoreServiceInterface descreve as operações para manipulação de pacotes de cartas na loja.
//
// Métodos:
//   - AddPackage: adiciona um novo pacote de cartas.
//   - GetPackage: recupera um pacote de cartas.
type StoreServiceInterface interface {
	// AddPackage adiciona um novo pacote de cartas à loja.
	//
	// Parâmetros:
	//   - cardPackage: pacote de cartas a ser adicionado.
	AddPackage(cardPackage domain.CardPackage)

	// GetPackage recupera um pacote de cartas da loja.
	//
	// Retorno:
	//   - CardPackage: pacote de cartas recuperado.
	GetPackage() domain.CardPackage
}

// StoreService implementa a lógica de armazenamento de pacotes de cartas usando um canal.
//
// Campos:
//   - packages: canal responsável pelo armazenamento dos pacotes de cartas.
type StoreService struct {
	// packages é o canal responsável pelo armazenamento dos pacotes de cartas.
	packages chan domain.CardPackage
}

// NewStoreService cria uma nova instância de StoreService.
//
// Retorno:
//   - ponteiro para StoreService.
func NewStoreService() *StoreService {
	return &StoreService{
		packages: make(chan domain.CardPackage, 10),
	}
}

// AddPackage adiciona um novo pacote de cartas ao canal de armazenamento.
//
// Parâmetros:
//   - cardPackage: pacote de cartas a ser adicionado.
func (s *StoreService) AddPackage(cardPackage domain.CardPackage) {
	s.packages <- cardPackage
}

// GetPackage recupera um pacote de cartas do canal de armazenamento.
//
// Retorno:
//   - CardPackage: pacote de cartas recuperado.
func (s *StoreService) GetPackage() domain.CardPackage {
	return <-s.packages
}
