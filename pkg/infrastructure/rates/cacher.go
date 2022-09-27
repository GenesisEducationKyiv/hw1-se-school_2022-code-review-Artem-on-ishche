package rates

import (
	"time"

	"gses2.app/api/pkg/domain/models"
)

type CacherRateServiceFactory struct {
	MaxTime float64
}

func (factory CacherRateServiceFactory) CreateRateService() CacherRateService {
	return &inMemoryCacher{
		maximumCacheTimeInMinutes: factory.MaxTime,
		cachedResponses:           make(map[models.CurrencyPair]parsedResponse),
	}
}

type CacherRateService interface {
	SetNext(service *ExchangeRateServiceChain)
	GetExchangeRate(pair models.CurrencyPair) (float64, error)
	Update(pair *models.CurrencyPair, response *parsedResponse)
}

type inMemoryCacher struct {
	next                      *ExchangeRateServiceChain
	maximumCacheTimeInMinutes float64
	cachedResponses           map[models.CurrencyPair]parsedResponse
}

func (cacher *inMemoryCacher) SetNext(service *ExchangeRateServiceChain) {
	cacher.next = service
}

func (cacher *inMemoryCacher) GetExchangeRate(pair models.CurrencyPair) (float64, error) {
	response, ok := cacher.cachedResponses[pair]
	if !ok {
		return (*cacher.next).GetExchangeRate(pair)
	}

	if cacher.isCachedResponseOutdated(response) {
		return (*cacher.next).GetExchangeRate(pair)
	}

	return response.rate, nil
}

func (cacher *inMemoryCacher) Update(pair *models.CurrencyPair, response *parsedResponse) {
	cacher.cachedResponses[*pair] = *response
}

func (cacher *inMemoryCacher) isCachedResponseOutdated(response parsedResponse) bool {
	currentTime := time.Now()

	return currentTime.Sub(response.time).Minutes() >= cacher.maximumCacheTimeInMinutes
}
