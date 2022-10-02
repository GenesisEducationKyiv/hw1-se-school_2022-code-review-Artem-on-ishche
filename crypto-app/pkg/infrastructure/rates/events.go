package rates

import (
	"gses2.app/api/pkg/domain/models"
)

type Event interface {
	GetName() string
	GetData() interface{}
}

type newRateReturnedEventData struct {
	pair     *models.CurrencyPair
	response *parsedResponse
}

type NewRateReturnedEvent struct {
	data newRateReturnedEventData
}

func (e NewRateReturnedEvent) GetName() string {
	return "NewRateReturnedEvent"
}

func (e NewRateReturnedEvent) GetData() interface{} {
	return e.data
}

type failureAPIResponseReceivedEventData struct {
	pair       *models.CurrencyPair
	provider   string
	statusCode int
}

type FailureAPIResponseReceivedEvent struct {
	data failureAPIResponseReceivedEventData
}

func (e FailureAPIResponseReceivedEvent) GetName() string {
	return "FailureAPIResponseReceivedEvent"
}

func (e FailureAPIResponseReceivedEvent) GetData() interface{} {
	return e.data
}

type successAPIResponseReceivedEventData struct {
	pair     *models.CurrencyPair
	provider string
	response *parsedResponse
}

type SuccessAPIResponseReceivedEvent struct {
	data successAPIResponseReceivedEventData
}

func (e SuccessAPIResponseReceivedEvent) GetName() string {
	return "SuccessAPIResponseReceivedEvent"
}

func (e SuccessAPIResponseReceivedEvent) GetData() interface{} {
	return e.data
}
