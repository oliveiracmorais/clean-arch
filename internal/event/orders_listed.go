package event

import (
	"time"

	"github.com/oliveiracmorais/clean-arch/pkg/events"
)

type OrdersListedEventInterface interface {
	events.EventInterface
}

type OrdersListed struct {
	Name    string
	Payload interface{}
}

func NewOrdersListed() *OrdersListed {
	return &OrdersListed{
		Name: "OrdersListed",
	}
}

func (e *OrdersListed) GetName() string {
	return e.Name
}

func (e *OrdersListed) GetPayload() interface{} {
	return e.Payload
}

func (e *OrdersListed) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *OrdersListed) GetDateTime() time.Time {
	return time.Now()
}
