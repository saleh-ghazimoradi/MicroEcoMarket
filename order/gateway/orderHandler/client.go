package orderHandler

import (
	"context"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/domain"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/dto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/gateway/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type GRPCOrderClient interface {
	CreateOrder(ctx context.Context, input *dto.Order) (*domain.Order, error)
	GetOrdersForAccount(ctx context.Context, accountId string) ([]*domain.Order, error)
	Close() error
}

type gRPCOrderClient struct {
	conn   *grpc.ClientConn
	client proto.OrderServiceClient
}

func (g *gRPCOrderClient) CreateOrder(ctx context.Context, input *dto.Order) (*domain.Order, error) {
	var protoCatalogs []*proto.CreateOrderRequest_OrderCatalog

	for _, c := range input.Catalogs {
		protoCatalogs = append(protoCatalogs, &proto.CreateOrderRequest_OrderCatalog{
			CatalogId: c.Id,
			Quantity:  c.Quantity,
		})
	}

	req, err := g.client.CreateOrder(ctx, &proto.CreateOrderRequest{
		AccountId: input.AccountId,
		Catalogs:  protoCatalogs,
	})
	if err != nil {
		log.Printf("error creating order: %v", err)
		return nil, err
	}

	var createdAt time.Time

	if len(req.Order.CreatedAt) > 0 {
		if err := createdAt.UnmarshalBinary(req.Order.CreatedAt); err != nil {
			log.Printf("Error unmarshaling CreatedAt: %v", err)
			return nil, err
		}
	}

	catalogs := make([]*domain.OrderedCatalog, len(req.Order.Catalogs))
	for i, c := range req.Order.Catalogs {
		catalogs[i] = &domain.OrderedCatalog{
			Id:          c.Id,
			Name:        c.Name,
			Description: c.Description,
			Price:       c.Price,
			Quantity:    c.Quantity,
		}
	}

	return &domain.Order{
		Id:         req.Order.Id,
		CreatedAt:  createdAt,
		TotalPrice: req.Order.TotalPrice,
		AccountId:  req.Order.AccountId,
		Catalogs:   catalogs,
	}, nil
}

func (g *gRPCOrderClient) GetOrdersForAccount(ctx context.Context, accountId string) ([]*domain.Order, error) {
	resp, err := g.client.GetOrdersForAccount(ctx, &proto.GetOrdersForAccountRequest{
		AccountId: accountId,
	})
	if err != nil {
		log.Printf("Error getting orders for account %s: %v", accountId, err)
		return nil, err
	}

	orders := make([]*domain.Order, len(resp.Orders))
	for i, o := range resp.Orders {
		var createdAt time.Time
		if len(o.CreatedAt) > 0 {
			if err := createdAt.UnmarshalBinary(o.CreatedAt); err != nil {
				log.Printf("Error unmarshaling CreatedAt: %v", err)
				return nil, err
			}
		}
		catalogs := make([]*domain.OrderedCatalog, len(o.Catalogs))
		for i, c := range o.Catalogs {
			catalogs[i] = &domain.OrderedCatalog{
				Id:          c.Id,
				Name:        c.Name,
				Description: c.Description,
				Price:       c.Price,
				Quantity:    c.Quantity,
			}
		}
		orders[i] = &domain.Order{
			Id:         o.Id,
			CreatedAt:  createdAt,
			TotalPrice: o.TotalPrice,
			AccountId:  o.AccountId,
			Catalogs:   catalogs,
		}
	}
	return orders, nil
}

func (g *gRPCOrderClient) Close() error {
	return g.conn.Close()
}

func NewGRPCOrderHandler(addr string) (GRPCOrderClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	orderServiceClient := proto.NewOrderServiceClient(conn)
	return &gRPCOrderClient{
		conn:   conn,
		client: orderServiceClient,
	}, nil
}
