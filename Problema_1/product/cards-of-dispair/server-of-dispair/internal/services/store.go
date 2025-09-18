package services

import "server-of-dispair/internal/domain"

type StoreService struct {
	Package chan domain.CardPackage
}

func NewStoreService() *StoreService {
	return &StoreService{
		Package: make(chan domain.CardPackage, 100),
	}
}

func (service *StoreService) AddPackage(cardPackage domain.CardPackage) {
	service.Package <- cardPackage
}

func (service *StoreService) GetPackage() domain.CardPackage {
	return <-service.Package
}
