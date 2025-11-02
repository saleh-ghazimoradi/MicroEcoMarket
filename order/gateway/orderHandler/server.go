package orderHandler

import (
	"context"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/gateway/proto"
	"google.golang.org/grpc"
	"net"
)

type GRPCOrderServer interface {
	CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error)
	GetOrdersForAccount(ctx context.Context, req *proto.GetOrdersForAccountRequest) (*proto.GetOrdersForAccountResponse, error)
	Serve(addr string) error
	Stop() error
}

type gRPCOrderServer struct {
	server *grpc.Server
	proto.UnimplementedOrderServiceServer
}

func (g *gRPCOrderServer) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	return nil, nil
}

func (g *gRPCOrderServer) GetOrdersForAccount(ctx context.Context, req *proto.GetOrdersForAccountRequest) (*proto.GetOrdersForAccountResponse, error) {
	return nil, nil
}

func (g *gRPCOrderServer) Serve(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	g.server = grpc.NewServer()
	proto.RegisterOrderServiceServer(g.server, g)
	return g.server.Serve(lis)
}

func (g *gRPCOrderServer) Stop() error {
	if g.server != nil {
		g.server.GracefulStop()
	}
	return nil
}

func NewGRPCOrderServer() GRPCOrderServer {
	return &gRPCOrderServer{}
}
