package repository

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/domain"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	GetOrdersForAccount(ctx context.Context, accountId string) ([]*domain.Order, error)
}

type orderRepository struct {
	dbWrite *sql.DB
	dbRead  *sql.DB
}

func (o *orderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	return nil
}

func (o *orderRepository) GetOrdersForAccount(ctx context.Context, accountId string) ([]*domain.Order, error) {
	return nil, nil
}

func NewOrderRepository(dbWrite, dbRead *sql.DB) OrderRepository {
	return &orderRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
