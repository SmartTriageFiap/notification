package repository

import "github.com/SmartTriageFiap/notification/entities"

type Repository interface {
	GetPartnersAddresses() ([]entities.Partner, error)
}

func New() Repository {
	return newMongoRepository()
}
