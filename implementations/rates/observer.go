package rates

import (
	"fmt"
	"gses2.app/api/services"
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
  requested rate - {%s - %s}
  response - {status code: %v}`,
		eventData.provider, eventData.pair.From, eventData.pair.To, eventData.statusCode)
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
  requested rate - {%s - %s}
  response - {
    rate: %v,
    time: %v
  }`,
		eventData.provider, eventData.pair.From, eventData.pair.To, eventData.response.rate, eventData.response.time)
}
