package repository

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/domain"
	"time"
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
	tx, err := o.dbWrite.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO "order"(id, created_at, account_id, total_price) VALUES($1, $2, $3, $4)`,
		order.Id,
		order.CreatedAt,
		order.AccountId,
		order.TotalPrice,
	)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, pq.CopyIn("order_catalog", "order_id", "catalog_id", "quantity"))
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, c := range order.Catalogs {
		_, err = stmt.ExecContext(ctx, order.Id, c.Id, c.Quantity)
		if err != nil {
			return err
		}
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		_ = tx.Rollback()
		return commitErr
	}

	return nil
}

func (o *orderRepository) GetOrdersForAccount(ctx context.Context, accountId string) ([]*domain.Order, error) {
	rows, err := o.dbRead.QueryContext(ctx, `
SELECT
  o.id,
  o.created_at,
  o.account_id,
  o.total_price::money::numeric::float8,
  oc.catalog_id,
  oc.quantity
FROM "order" o
LEFT JOIN order_catalog oc ON o.id = oc.order_id
WHERE o.account_id = $1
ORDER BY o.id;
`, accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	var currOrder *domain.Order
	var currOrderID string
	var catalogs []*domain.OrderedCatalog

	for rows.Next() {
		var (
			orderID    string
			createdAt  time.Time
			accountID  string
			totalPrice float64
			catalogID  sql.NullString
			quantity   sql.NullInt64
		)

		if err := rows.Scan(&orderID, &createdAt, &accountID, &totalPrice, &catalogID, &quantity); err != nil {
			return nil, err
		}

		// detect a new order boundary
		if currOrderID == "" || currOrderID != orderID {
			// push previous order if exists
			if currOrder != nil {
				currOrder.Catalogs = catalogs
				orders = append(orders, currOrder)
			}

			currOrderID = orderID
			catalogs = make([]*domain.OrderedCatalog, 0)
			currOrder = &domain.Order{
				Id:         orderID,
				CreatedAt:  createdAt,
				AccountId:  accountID,
				TotalPrice: totalPrice,
				Catalogs:   nil,
			}
		}

		// append ordered product if exists
		if catalogID.Valid {
			q := uint32(0)
			if quantity.Valid {
				q = uint32(quantity.Int64)
			}
			catalogs = append(catalogs, &domain.OrderedCatalog{
				Id:       catalogID.String,
				Quantity: q,
			})
		}
	}

	// push the last order
	if currOrder != nil {
		currOrder.Catalogs = catalogs
		orders = append(orders, currOrder)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func NewOrderRepository(dbWrite, dbRead *sql.DB) OrderRepository {
	return &orderRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
