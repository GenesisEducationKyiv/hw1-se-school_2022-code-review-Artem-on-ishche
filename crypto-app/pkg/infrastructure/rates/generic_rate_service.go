package rates

import (
	"net/http"

	"gopkg.in/resty.v0"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/models"
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
	mediator           *Mediator
	concreteRateClient exchangeRateAPIClient
	next               *ExchangeRateServiceChain
}

func (service *exchangeRateService) SetNext(nextService *ExchangeRateServiceChain) {
	service.next = nextService
}

func (service *exchangeRateService) GetExchangeRate(pair models.CurrencyPair) (*models.ExchangeRate, error) {
	rate, err := service.getExchangeRate(pair)
	if err != nil && service.next != nil {
		return (*service.next).GetExchangeRate(pair)
	}

	return rate, err
}

func (service *exchangeRateService) getExchangeRate(pair models.CurrencyPair) (*models.ExchangeRate, error) {
	resp, err := service.makeAPIRequest(pair)
	if err != nil {
		return nil, application.ErrAPIRequestUnsuccessful
	}

	if resp.StatusCode() != http.StatusOK {
		service.notifyMediatorAboutFailureAPIResponseReceived(pair, resp)

		return nil, application.ErrAPIRequestUnsuccessful
	}

	parsedResponse, err := service.concreteRateClient.parseResponseBody(resp.Body)
	if err != nil {
		return nil, err
	}

	service.notifyMediatorAboutSuccessAPIResponseReceived(pair, parsedResponse)
	service.notifyMediatorAboutNewRateReturned(&pair, parsedResponse)

	return models.NewExchangeRate(pair, parsedResponse.price, parsedResponse.time), nil
}

func (service *exchangeRateService) makeAPIRequest(pair models.CurrencyPair) (*resty.Response, error) {
	url := service.concreteRateClient.getAPIRequestURLForGivenCurrencies(pair)
	request := service.concreteRateClient.getAPIRequest()

	return request.Get(url)
}

func (service *exchangeRateService) notifyMediatorAboutFailureAPIResponseReceived(pair models.CurrencyPair, response *resty.Response) {
	(*service.mediator).Notify(FailureAPIResponseReceivedEvent{failureAPIResponseReceivedEventData{
		pair:       &pair,
		provider:   service.concreteRateClient.name(),
		statusCode: response.StatusCode(),
	}})
}

func (service *exchangeRateService) notifyMediatorAboutSuccessAPIResponseReceived(pair models.CurrencyPair, response *parsedResponse) {
	(*service.mediator).Notify(SuccessAPIResponseReceivedEvent{successAPIResponseReceivedEventData{
		pair:     &pair,
		provider: service.concreteRateClient.name(),
		response: response,
	}})
}

func (service *exchangeRateService) notifyMediatorAboutNewRateReturned(pair *models.CurrencyPair, response *parsedResponse) {
	(*service.mediator).Notify(NewRateReturnedEvent{newRateReturnedEventData{
		pair:     pair,
		response: response,
	}})
}
