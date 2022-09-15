package rates

import (
	"gopkg.in/resty.v0"
	"gses2.app/api/services"
	"net/http"
	"time"
)

type parsedResponse struct {
	rate float64
	time time.Time
}

type exchangeRateAPIClient interface {
	getName() string
	getAPIRequestUrlForGivenCurrencies(from, to services.Currency) string
	getAPIRequest() *resty.Request
	parseResponseBody(responseBody []byte) (*parsedResponse, error)
}

type exchangeRateService struct {
	concreteRateClient exchangeRateAPIClient
	next               *services.ExchangeRateService
}

func (service *exchangeRateService) SetNext(nextService *services.ExchangeRateService) {
	service.next = nextService
}

func (service *exchangeRateService) GetExchangeRate(from, to services.Currency) (float64, error) {
	rate, err := service.getExchangeRate(from, to)
	if err != nil && service.next != nil {
		return (*service.next).GetExchangeRate(from, to)
	}

	return rate, err
}

func (service *exchangeRateService) getExchangeRate(from, to services.Currency) (float64, error) {
	resp, err := service.makeAPIRequest(from, to)
	if !isAPIRequestSuccessful(resp, err) {
		return -1, services.ErrAPIRequestUnsuccessful
	}

	parsedResponse, err := service.concreteRateClient.parseResponseBody(resp.Body)
	if err != nil {
		return -1, err
	}

	return parsedResponse.rate, nil
}

func (service *exchangeRateService) makeAPIRequest(from, to services.Currency) (*resty.Response, error) {
	url := service.concreteRateClient.getAPIRequestUrlForGivenCurrencies(from, to)
	request := service.concreteRateClient.getAPIRequest()

	return request.Get(url)
}

func isAPIRequestSuccessful(resp *resty.Response, err error) bool {
	return err == nil && resp.StatusCode() == http.StatusOK
}
