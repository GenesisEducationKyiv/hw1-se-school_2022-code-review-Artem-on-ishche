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
		cachedResponses:           make(map[string]parsedResponse),
	}
}

type CacherRateService interface {
	SetNext(service *ExchangeRateServiceChain)
	GetExchangeRate(pair models.CurrencyPair) (*models.ExchangeRate, error)
	Update(pair *models.CurrencyPair, response *parsedResponse)
}

type inMemoryCacher struct {
	next                      *ExchangeRateServiceChain
	maximumCacheTimeInMinutes float64
	cachedResponses           map[string]parsedResponse
}

func (cacher *inMemoryCacher) SetNext(service *ExchangeRateServiceChain) {
	cacher.next = service
}

func (cacher *inMemoryCacher) GetExchangeRate(pair models.CurrencyPair) (*models.ExchangeRate, error) {
	response, ok := cacher.cachedResponses[pair.String()]
	if !ok {
		return (*cacher.next).GetExchangeRate(pair)
	}

	if cacher.isCachedResponseOutdated(response) {
		return (*cacher.next).GetExchangeRate(pair)
	}

	return models.NewExchangeRate(pair, response.price, response.time), nil
}

func (cacher *inMemoryCacher) Update(pair *models.CurrencyPair, response *parsedResponse) {
	cacher.cachedResponses[(*pair).String()] = *response
}

func (cacher *inMemoryCacher) isCachedResponseOutdated(response parsedResponse) bool {
	currentTime := time.Now()

	return currentTime.Sub(response.time).Minutes() >= cacher.maximumCacheTimeInMinutes
}
