package rates

import (
	"time"

	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type CacherRateServiceFactory struct {
	MaxTime float64

	Logger services.Logger
}

func (factory CacherRateServiceFactory) CreateRateService() CacherRateService {
	return &inMemoryCacher{
		maximumCacheTimeInMinutes: factory.MaxTime,
		cachedResponses:           make(map[string]parsedResponse),
		logger:                    factory.Logger,
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

	logger services.Logger
}

func (cacher *inMemoryCacher) SetNext(service *ExchangeRateServiceChain) {
	cacher.next = service
}

func (cacher *inMemoryCacher) GetExchangeRate(pair models.CurrencyPair) (*models.ExchangeRate, error) {
	cacher.logger.Debug("Trying to get the exchange rate from an in-memory cache")

	response, ok := cacher.cachedResponses[pair.String()]
	if !ok {
		cacher.logger.Debug("Cache miss, calling GetExchangeRate() on the next client in chain")

		return (*cacher.next).GetExchangeRate(pair)
	}

	if cacher.isCachedResponseOutdated(response) {
		cacher.logger.Debug("Cached rate outdated")

		return (*cacher.next).GetExchangeRate(pair)
	}

	cacher.logger.Debug("Cache hit")

	return models.NewExchangeRate(pair, response.price, response.time), nil
}

func (cacher *inMemoryCacher) Update(pair *models.CurrencyPair, response *parsedResponse) {
	cacher.logger.Debug("Cache update called for " + pair.String())

	cacher.cachedResponses[(*pair).String()] = *response
}

func (cacher *inMemoryCacher) isCachedResponseOutdated(response parsedResponse) bool {
	currentTime := time.Now()

	return currentTime.Sub(response.time).Minutes() >= cacher.maximumCacheTimeInMinutes
}
