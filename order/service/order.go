package service

import (
	"context"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/domain"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/dto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/repository"
	"github.com/segmentio/ksuid"
	"time"
)

type OrderService interface {
	CreateOrder(ctx context.Context, input *dto.Order) (*domain.Order, error)
	GetOrdersForAccount(ctx context.Context, accountId string) ([]*domain.Order, error)
}

type orderService struct {
	orderRepository repository.OrderRepository
}

func (o *orderService) CreateOrder(ctx context.Context, input *dto.Order) (*domain.Order, error) {
	order := &domain.Order{
		Id:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		AccountId: input.AccountId,
		Catalogs:  make([]*domain.OrderedCatalog, 0),
	}

	order.TotalPrice = 0.0
	for _, catalog := range input.Catalogs {
		order.TotalPrice += catalog.Price * float64(catalog.Quantity)
	}

	if err := o.orderRepository.CreateOrder(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (o *orderService) GetOrdersForAccount(ctx context.Context, accountId string) ([]*domain.Order, error) {
	return o.orderRepository.GetOrdersForAccount(ctx, accountId)
}

func NewOrderService(orderRepository repository.OrderRepository) OrderService {
	return &orderService{
		orderRepository: orderRepository,
	}
}
