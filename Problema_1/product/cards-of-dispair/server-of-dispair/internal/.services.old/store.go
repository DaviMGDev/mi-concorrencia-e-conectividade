package services

import (
	"math/rand"
	"server-of-dispair/internal/entities"
)

type StoreServiceInterface interface {
	BuyPackage() (entities.PackageInterface, error)
	RestockPackage()
}

type StoreService struct {
	Store *entities.Store
}

func NewStoreService(store *entities.Store) *StoreService {
	return &StoreService{
		Store: store,
	}
}

func (s *StoreService) RestockPackage() {
	pkg := entities.Package{
		Cards: [3]entities.CardInterface{
			entities.NewCard("rock", uint(rand.Intn(5)+1)),
			entities.NewCard("paper", uint(rand.Intn(5)+1)),
			entities.NewCard("scissors", uint(rand.Intn(5)+1)),
		},
	}
	s.Store.AddPackage(pkg)
}

func (s *StoreService) BuyPackage() (entities.PackageInterface, error) {
	return s.Store.GetPackage()
}