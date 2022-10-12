package rates

import (
	"fmt"
	"net/http"

	"gopkg.in/resty.v0"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

const timeLayout = "2006-01-02T15:04:05.999Z"

type RateServiceFactory interface {
	CreateRateService() ExchangeRateServiceChain
}

type ExchangeRateServiceChain interface {
	SetNext(service *ExchangeRateServiceChain)
	GetExchangeRate(pair models.CurrencyPair) (*models.ExchangeRate, error)
}

type exchangeRateService struct {
	concreteRateClient exchangeRateAPIClient
	next               *ExchangeRateServiceChain
	cacher             CacherRateService

	logger services.Logger
}

func (service *exchangeRateService) SetNext(nextService *ExchangeRateServiceChain) {
	service.next = nextService
}

func (service *exchangeRateService) GetExchangeRate(pair models.CurrencyPair) (*models.ExchangeRate, error) {
	service.logger.Debug(fmt.Sprintf("GetExchangeRate() called on %s rate service with pair = %s",
		service.concreteRateClient.name(), pair.String()))

	rate, err := service.getExchangeRate(pair)
	if err != nil && service.next != nil {
		service.logger.Debug("Calling GetExchangeRate() on the next client in chain")

		return (*service.next).GetExchangeRate(pair)
	}

	return rate, err
}

func (service *exchangeRateService) getExchangeRate(pair models.CurrencyPair) (*models.ExchangeRate, error) {
	service.logger.Debug("Trying to get the exchange rate from " + service.concreteRateClient.name())

	resp, err := service.makeAPIRequest(pair)
	if err != nil {
		service.logger.Error("API request failed")

		return nil, application.ErrAPIRequestUnsuccessful
	}

	if resp.StatusCode() != http.StatusOK {
		service.logUnsuccessfulAPIResponse(&pair, resp.StatusCode())

		return nil, application.ErrAPIRequestUnsuccessful
	}

	parsedResponse, err := service.concreteRateClient.parseResponseBody(resp.Body)
	if err != nil {
		service.logger.Error("API returned a success code but I failed to parse the response")

		return nil, err
	}

	service.logSuccessfulAPIResponse(&pair, parsedResponse)
	service.cacher.Update(&pair, parsedResponse)

	return models.NewExchangeRate(pair, parsedResponse.price, parsedResponse.time), nil
}

func (service *exchangeRateService) makeAPIRequest(pair models.CurrencyPair) (*resty.Response, error) {
	url := service.concreteRateClient.getAPIRequestURLForGivenCurrencies(pair)
	request := service.concreteRateClient.getAPIRequest()

	return request.Get(url)
}

func (service *exchangeRateService) logUnsuccessfulAPIResponse(pair *models.CurrencyPair, statusCode int) {
	service.logger.Error(fmt.Sprintf(
		`
%s:
  requested rate - %s
  response - {status code: %v}`,
		service.concreteRateClient.name(), pair.String(), statusCode))
}

func (service *exchangeRateService) logSuccessfulAPIResponse(pair *models.CurrencyPair, resp *parsedResponse) {
	service.logger.Debug(fmt.Sprintf(
		`
%s:
  requested rate - %s
  response - {
    price: %v,
    time: %v
  }`,
		service.concreteRateClient.name(), pair.String(), resp.price, resp.time))
}
