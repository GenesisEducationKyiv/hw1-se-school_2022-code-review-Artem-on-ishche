package rates

import (
	"fmt"

	"gses2.app/api/pkg/domain/services"
)

type Observer interface {
	Update(event Event)
}

type NewRateReturnedObserver struct {
	Cacher CacherRateService
}

func (observer NewRateReturnedObserver) Update(event Event) {
	eventData := event.GetData().(newRateReturnedEventData)

	observer.Cacher.Update(eventData.pair, eventData.response)
}

type FailureAPIResponseReceivedObserver struct {
	Logger services.Logger
}

func (observer FailureAPIResponseReceivedObserver) Update(event Event) {
	eventData := event.GetData().(failureAPIResponseReceivedEventData)

	text := getFailureAPIResponseReceivedLogText(eventData)

	observer.Logger.Log(text)
}

func getFailureAPIResponseReceivedLogText(eventData failureAPIResponseReceivedEventData) string {
	return fmt.Sprintf(
		`
%s:
  requested price - {%s - %s}
  response - {status code: %v}`,
		eventData.provider, eventData.pair.Base, eventData.pair.Quote, eventData.statusCode)
}

type SuccessAPIResponseReceivedObserver struct {
	Logger services.Logger
}

func (observer SuccessAPIResponseReceivedObserver) Update(event Event) {
	eventData := event.GetData().(successAPIResponseReceivedEventData)

	text := getSuccessAPIResponseReceivedLogText(eventData)

	observer.Logger.Log(text)
}

func getSuccessAPIResponseReceivedLogText(eventData successAPIResponseReceivedEventData) string {
	return fmt.Sprintf(
		`
%s:
  requested price - {%s - %s}
  response - {
    price: %v,
    time: %v
  }`,
		eventData.provider, eventData.pair.Base, eventData.pair.Quote, eventData.response.price, eventData.response.time)
}
