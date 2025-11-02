package orderHandler

import (
	"context"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/domain"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/dto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/gateway/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	return nil, nil
}

func (g *gRPCOrderClient) GetOrdersForAccount(ctx context.Context, accountId string) ([]*domain.Order, error) {
	return nil, nil
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
