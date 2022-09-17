package rates

import (
	"gopkg.in/resty.v0"
	"gses2.app/api/services"
	"net/http"
)

type RateServiceFactory interface {
	CreateRateService() ExchangeRateServiceChain
}

type ExchangeRateServiceChain interface {
	SetNext(service *ExchangeRateServiceChain)
	GetExchangeRate(pair services.CurrencyPair) (float64, error)
}

type exchangeRateService struct {
	mediator           *Mediator
	concreteRateClient exchangeRateAPIClient
	next               *ExchangeRateServiceChain
}

func (service *exchangeRateService) SetNext(nextService *ExchangeRateServiceChain) {
	service.next = nextService
}

func (service *exchangeRateService) GetExchangeRate(pair services.CurrencyPair) (float64, error) {
	rate, err := service.getExchangeRate(pair)
	if err != nil && service.next != nil {
		return (*service.next).GetExchangeRate(pair)
	}

	return rate, err
}

func (service *exchangeRateService) getExchangeRate(pair services.CurrencyPair) (float64, error) {
	resp, err := service.makeAPIRequest(pair)
	if err != nil {
		return -1, services.ErrAPIRequestUnsuccessful
	}

	if resp.StatusCode() != http.StatusOK {
		service.notifyMediatorAboutFailureAPIResponseReceived(pair, resp)

		return -1, services.ErrAPIRequestUnsuccessful
	}

	parsedResponse, err := service.concreteRateClient.parseResponseBody(resp.Body)
	if err != nil {
		return -1, err
	}

	service.notifyMediatorAboutSuccessAPIResponseReceived(pair, parsedResponse)
	service.notifyMediatorAboutNewRateReturned(&pair, parsedResponse)

	return parsedResponse.rate, nil
}

func (service *exchangeRateService) makeAPIRequest(pair services.CurrencyPair) (*resty.Response, error) {
	url := service.concreteRateClient.getAPIRequestUrlForGivenCurrencies(pair)
	request := service.concreteRateClient.getAPIRequest()

	return request.Get(url)
}

func (service *exchangeRateService) notifyMediatorAboutFailureAPIResponseReceived(pair services.CurrencyPair, response *resty.Response) {
	(*service.mediator).Notify(FailureAPIResponseReceivedEvent{failureAPIResponseReceivedEventData{
		pair:       &pair,
		provider:   service.concreteRateClient.getName(),
		statusCode: response.StatusCode(),
	}})
}

func (service *exchangeRateService) notifyMediatorAboutSuccessAPIResponseReceived(pair services.CurrencyPair, response *parsedResponse) {
	(*service.mediator).Notify(SuccessAPIResponseReceivedEvent{successAPIResponseReceivedEventData{
		pair:     &pair,
		provider: service.concreteRateClient.getName(),
		response: response,
	}})
}

func (service *exchangeRateService) notifyMediatorAboutNewRateReturned(pair *services.CurrencyPair, response *parsedResponse) {
	(*service.mediator).Notify(NewRateReturnedEvent{newRateReturnedEventData{
		pair:     pair,
		response: response,
	}})
}
