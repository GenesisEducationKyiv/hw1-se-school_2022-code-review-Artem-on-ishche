package publicmocks

import (
	"gses2.app/api/pkg/domain/models"
)

var EmptyCustomersService = emptyCustomersService{}

type emptyCustomersService struct{}

func (emptyCustomersService) CreateCustomer(*models.EmailAddress) error { return nil }
