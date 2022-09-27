package rates

import (
	"fmt"
	"reflect"
)

type Mediator interface {
	Attach(observer Observer, eventName string) error
	Detach(observer Observer, eventName string) error
	Notify(event Event)
}

func NewMediator() Mediator {
	return &mediator{make(map[string]map[string]Observer)}
}

type mediator struct {
	observers map[string]map[string]Observer
}

func (m *mediator) Attach(observer Observer, eventName string) error {
	observerType := reflect.TypeOf(observer).String()

	if _, ok := m.observers[eventName]; !ok {
		m.observers[eventName] = make(map[string]Observer)
	}

	if _, ok := m.observers[eventName][observerType]; ok {
		return fmt.Errorf("the '%s' event is already attached to '%s'", eventName, observerType)
	}

	m.observers[eventName][observerType] = observer

	return nil
}

func (m *mediator) Detach(observer Observer, eventName string) error {
	observerType := reflect.TypeOf(observer).String()

	if _, ok := m.observers[eventName][observerType]; !ok {
		return fmt.Errorf("the '%s' event is not attached to '%s'", eventName, observerType)
	}

	delete(m.observers[eventName], observerType)

	return nil
}

func (m *mediator) Notify(event Event) {
	eventName := event.GetName()

	observersForEvent := m.observers[eventName]

	for _, observer := range observersForEvent {
		observer.Update(event)
	}
}
