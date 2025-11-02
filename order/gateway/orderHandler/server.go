package orderHandler

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/gateway/accountHandler"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/dto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/gateway/catalogHandler"
	orderDTO "github.com/saleh-ghazimoradi/MircoEcoMarket/order/dto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/gateway/proto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

type GRPCOrderServer interface {
	CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error)
	GetOrdersForAccount(ctx context.Context, req *proto.GetOrdersForAccountRequest) (*proto.GetOrdersForAccountResponse, error)
	Serve(addr string) error
	Stop() error
}

type gRPCOrderServer struct {
	orderService  service.OrderService
	accountClient accountHandler.GRPCAccountClient
	catalogClient catalogHandler.GRPCCatalogClient
	server        *grpc.Server
	proto.UnimplementedOrderServiceServer
}

func (g *gRPCOrderServer) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	_, err := g.accountClient.GetAccountById(ctx, req.AccountId)
	if err != nil {
		log.Println("Error getting account by id:", err)
		return nil, errors.New("account not found")
	}

	var catalogIds []string
	for _, p := range req.Catalogs {
		catalogIds = append(catalogIds, p.CatalogId)
	}

	orderedCatalog, err := g.catalogClient.GetCatalogs(ctx, &dto.CatalogQuery{
		Limit:  0,
		Offset: 0,
		Query:  "",
		Ids:    catalogIds,
	})
	if err != nil {
		log.Println("Error getting catalogs:", err)
	}

	var catalogs []*orderDTO.OrderedCatalog
	for _, o := range orderedCatalog {
		catalog := &orderDTO.OrderedCatalog{
			Id:          o.Id,
			Quantity:    0,
			Price:       o.Price,
			Name:        o.Name,
			Description: o.Description,
		}
		for _, rc := range req.Catalogs {
			if rc.CatalogId == o.Id {
				catalog.Quantity = rc.Quantity
				break
			}
		}
		if catalog.Quantity != 0 {
			catalogs = append(catalogs, catalog)
		}
	}

	order, err := g.orderService.CreateOrder(ctx, &orderDTO.Order{
		AccountId: req.AccountId,
		Catalogs:  catalogs,
	})

	if err != nil {
		return nil, errors.New("could not create order")
	}

	orderProto := &proto.Order{
		Id:         order.Id,
		AccountId:  order.AccountId,
		TotalPrice: order.TotalPrice,
		Catalogs:   []*proto.Order_OrderCatalog{},
	}

	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()
	for _, o := range order.Catalogs {
		orderProto.Catalogs = append(orderProto.Catalogs, &proto.Order_OrderCatalog{
			Id:          o.Id,
			Name:        o.Name,
			Description: o.Description,
			Price:       o.Price,
			Quantity:    o.Quantity,
		})
	}

	return &proto.CreateOrderResponse{
		Order: orderProto,
	}, nil

}

func (g *gRPCOrderServer) GetOrdersForAccount(ctx context.Context, req *proto.GetOrdersForAccountRequest) (*proto.GetOrdersForAccountResponse, error) {
	accountOrders, err := g.orderService.GetOrdersForAccount(ctx, req.AccountId)
	if err != nil {
		return nil, errors.New("could not get orders")
	}

	catalogIdMap := make(map[string]bool)
	for _, o := range accountOrders {
		for _, rc := range o.Catalogs {
			catalogIdMap[rc.Id] = true
		}
	}

	var catalogsIDs []string
	for id := range catalogIdMap {
		catalogsIDs = append(catalogsIDs, id)
	}

	catalogs, err := g.catalogClient.GetCatalogs(ctx, &dto.CatalogQuery{
		Limit:  0,
		Offset: 0,
		Query:  "",
		Ids:    catalogsIDs,
	})

	if err != nil {
		return nil, errors.New("could not get catalogs")
	}

	var orders []*proto.Order
	for _, o := range accountOrders {
		op := &proto.Order{
			AccountId:  o.AccountId,
			Id:         o.Id,
			TotalPrice: o.TotalPrice,
			Catalogs:   []*proto.Order_OrderCatalog{},
		}
		op.CreatedAt, _ = o.CreatedAt.MarshalBinary()

		for _, catalog := range o.Catalogs {
			for _, c := range catalogs {
				if c.Id == catalog.Id {
					catalog.Name = c.Name
					catalog.Description = c.Description
					catalog.Price = c.Price
					break
				}
			}
			op.Catalogs = append(op.Catalogs, &proto.Order_OrderCatalog{
				Id:          catalog.Id,
				Name:        catalog.Name,
				Description: catalog.Description,
				Price:       catalog.Price,
				Quantity:    catalog.Quantity,
			})
		}
		orders = append(orders, op)
	}
	return &proto.GetOrdersForAccountResponse{
		Orders: orders,
	}, nil
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

func NewGRPCOrderServer(orderService service.OrderService, accountClient accountHandler.GRPCAccountClient, catalogClient catalogHandler.GRPCCatalogClient) GRPCOrderServer {
	return &gRPCOrderServer{
		orderService:  orderService,
		accountClient: accountClient,
		catalogClient: catalogClient,
	}
}
