package services

import "gses2.app/api/pkg/domain/models"

type CustomersService interface {
	CreateCustomer(emailAddress *models.EmailAddress) error
}
