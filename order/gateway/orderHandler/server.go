package orderHandler

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/gateway/accountHandler"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/dto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/gateway/catalogHandler"
	orderDTO "github.com/saleh-ghazimoradi/MircoEcoMarket/order/dto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/gateway/proto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/service"
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
	orderService  service.OrderService
	accountClient accountHandler.GRPCAccountClient
	catalogClient catalogHandler.GRPCCatalogClient
	server        *grpc.Server
	proto.UnimplementedOrderServiceServer
}

func (g *gRPCOrderServer) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	// 1. Verify account exists
	_, err := g.accountClient.GetAccountById(ctx, req.AccountId)
	if err != nil {
		return nil, fmt.Errorf("account not found")
	}

	// 2. Collect catalog IDs
	var catalogIds []string
	for _, c := range req.Catalogs {
		catalogIds = append(catalogIds, c.CatalogId)
	}

	// 3. Fetch catalog details from Catalog service
	catalogsFromService, err := g.catalogClient.GetCatalogs(ctx, &dto.CatalogQuery{
		Ids: catalogIds,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch catalog details: %v", err)
	}

	// 4. Build ordered catalogs with quantity
	var orderedCatalogs []*orderDTO.OrderedCatalog
	for _, catalog := range catalogsFromService {
		var quantity uint32
		for _, c := range req.Catalogs {
			if c.CatalogId == catalog.Id {
				quantity = c.Quantity
				break
			}
		}
		if quantity == 0 {
			continue
		}
		orderedCatalogs = append(orderedCatalogs, &orderDTO.OrderedCatalog{
			Id:          catalog.Id,
			Name:        catalog.Name,
			Description: catalog.Description,
			Price:       catalog.Price,
			Quantity:    quantity,
		})
	}

	// 5. Create order
	order, err := g.orderService.CreateOrder(ctx, &orderDTO.Order{
		AccountId: req.AccountId,
		Catalogs:  orderedCatalogs,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create order: %v", err)
	}

	// 6. Convert to gRPC Order response
	orderProto := &proto.Order{
		Id:         order.Id,
		AccountId:  order.AccountId,
		TotalPrice: order.TotalPrice,
		Catalogs:   []*proto.Order_OrderCatalog{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()

	for _, c := range order.Catalogs {
		orderProto.Catalogs = append(orderProto.Catalogs, &proto.Order_OrderCatalog{
			Id:          c.Id,
			Name:        c.Name,
			Description: c.Description,
			Price:       c.Price,
			Quantity:    c.Quantity,
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
