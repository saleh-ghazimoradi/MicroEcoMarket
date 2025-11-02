package service

import (
	"context"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/domain"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/dto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/repository"
)

type OrderService interface {
	CreateOrder(ctx context.Context, input *dto.Order) (*domain.Order, error)
	GetOrdersForAccount(ctx context.Context, accountId string) ([]*domain.Order, error)
}

type orderService struct {
	orderRepository repository.OrderRepository
}

func (o *orderService) CreateOrder(ctx context.Context, input *dto.Order) (*domain.Order, error) {
	return nil, nil
}

func (o *orderService) GetOrdersForAccount(ctx context.Context, accountId string) ([]*domain.Order, error) {
	return nil, nil
}

func NewOrderService(orderRepository repository.OrderRepository) OrderService {
	return &orderService{
		orderRepository: orderRepository,
	}
}
