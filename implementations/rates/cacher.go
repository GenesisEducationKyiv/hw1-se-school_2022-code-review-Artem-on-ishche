package rates

import (
	"time"

	"gses2.app/api/services"
)

type CacherRateServiceFactory struct {
	MaxTime float64
}

func (factory CacherRateServiceFactory) CreateRateService() CacherRateService {
	return &inMemoryCacher{
		maximumCacheTimeInMinutes: factory.MaxTime,
		cachedResponses:           make(map[services.CurrencyPair]parsedResponse),
	}
}

type CacherRateService interface {
	SetNext(service *ExchangeRateServiceChain)
	GetExchangeRate(pair services.CurrencyPair) (float64, error)
	Update(pair *services.CurrencyPair, response *parsedResponse)
}

type inMemoryCacher struct {
	next                      *ExchangeRateServiceChain
	maximumCacheTimeInMinutes float64
	cachedResponses           map[services.CurrencyPair]parsedResponse
}

func (cacher *inMemoryCacher) SetNext(service *ExchangeRateServiceChain) {
	cacher.next = service
}

func (cacher *inMemoryCacher) GetExchangeRate(pair services.CurrencyPair) (float64, error) {
	response, ok := cacher.cachedResponses[pair]
	if !ok {
		return (*cacher.next).GetExchangeRate(pair)
	}

	if cacher.isCachedResponseOutdated(response) {
		return (*cacher.next).GetExchangeRate(pair)
	}

	return response.rate, nil
}

func (cacher *inMemoryCacher) Update(pair *services.CurrencyPair, response *parsedResponse) {
	cacher.cachedResponses[*pair] = *response
}

func (cacher *inMemoryCacher) isCachedResponseOutdated(response parsedResponse) bool {
	currentTime := time.Now()

	return currentTime.Sub(response.time).Minutes() >= cacher.maximumCacheTimeInMinutes
}
