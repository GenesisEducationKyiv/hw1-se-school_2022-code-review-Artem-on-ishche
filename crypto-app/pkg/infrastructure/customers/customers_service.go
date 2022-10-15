package customers

import (
	"fmt"
	"net/http"

	"gopkg.in/resty.v0"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type ServiceImpl struct {
	customersServiceUrl string
	createCustomerRoute string

	logger services.Logger
}

func NewServiceImpl(customersServiceUrl string, createCustomerRoute string, logger services.Logger) *ServiceImpl {
	return &ServiceImpl{customersServiceUrl: customersServiceUrl, createCustomerRoute: createCustomerRoute, logger: logger}
}

func (service ServiceImpl) CreateCustomer(emailAddress *models.EmailAddress) error {
	service.logger.Debug(fmt.Sprintf("CreateCustomer() called with emailAddress = {%s}", emailAddress.String()))

	resp, err := resty.R().Post(service.url(emailAddress))
	if resp == nil {
		service.logger.Debug(fmt.Sprintf("API request to customers service returned no response, err = {%v}", err))
	} else {
		service.logger.Debug(fmt.Sprintf("API request to customers service returned status code = {%v}, err = {%v}",
			resp.StatusCode(), err))
	}

	if err != nil || resp == nil || resp.StatusCode() != http.StatusOK {
		err = application.ErrCustomersServiceRequestUnsuccessful
	}

	return err
}

func (service ServiceImpl) url(emailAddress *models.EmailAddress) string {
	url := fmt.Sprintf("%s%s?emailAddress=%s",
		service.customersServiceUrl, service.createCustomerRoute, emailAddress.String())
	service.logger.Debug("url to send request to is " + url)

	return url
}
