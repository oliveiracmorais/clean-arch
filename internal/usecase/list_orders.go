package usecase

import (
	"github.com/oliveiracmorais/clean-arch/internal/entity"
	"github.com/oliveiracmorais/clean-arch/internal/event"
	"github.com/oliveiracmorais/clean-arch/pkg/events"
)

type OrderDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrdersListed    event.OrdersListedEventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewListOrdersUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrdersListed event.OrdersListedEventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
		OrdersListed:    OrdersListed,
		EventDispatcher: EventDispatcher,
	}
}

func (l *ListOrdersUseCase) Execute() ([]OrderDTO, error) {
	orders, err := l.OrderRepository.List()
	if err != nil {
		return []OrderDTO{}, err
	}
	ordersDTO := []OrderDTO{}
	for _, order := range orders {
		ordersDTO = append(ordersDTO, OrderDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}

	l.OrdersListed.SetPayload(ordersDTO)
	l.EventDispatcher.Dispatch(l.OrdersListed)

	return ordersDTO, nil
}
