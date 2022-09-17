package rates

import "gses2.app/api/services"

type Event interface {
	GetName() string
	GetData() interface{}
}

type newRateReturnedEventData struct {
	pair     *services.CurrencyPair
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
	pair       *services.CurrencyPair
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
	pair     *services.CurrencyPair
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
