package catalogHandler

import (
	"context"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/domain"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/dto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/gateway/proto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/service"
	"google.golang.org/grpc"
	"net"
)

type GRPCCatalogServer interface {
	CreateCatalog(ctx context.Context, req *proto.CreateCatalogRequest) (*proto.CreateCatalogResponse, error)
	GetCatalogById(ctx context.Context, req *proto.GetCatalogRequest) (*proto.GetCatalogResponse, error)
	GetCatalogs(ctx context.Context, req *proto.GetCatalogsRequest) (*proto.GetCatalogsResponse, error)
	Serve(addr string) error
	Stop() error
}

type gRPCCatalogServer struct {
	catalogService service.CatalogService
	server         *grpc.Server
	proto.UnimplementedCatalogServiceServer
}

func (g *gRPCCatalogServer) CreateCatalog(ctx context.Context, req *proto.CreateCatalogRequest) (*proto.CreateCatalogResponse, error) {
	catalog, err := g.catalogService.CreateCatalog(ctx, &dto.Catalog{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		return nil, err
	}

	return &proto.CreateCatalogResponse{
		Catalog: &proto.Catalog{
			Id:          catalog.Id,
			Name:        catalog.Name,
			Description: catalog.Description,
			Price:       catalog.Price,
		},
	}, nil
}

func (g *gRPCCatalogServer) GetCatalogById(ctx context.Context, req *proto.GetCatalogRequest) (*proto.GetCatalogResponse, error) {
	catalog, err := g.catalogService.GetCatalogById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &proto.GetCatalogResponse{
		Catalog: &proto.Catalog{
			Id:          catalog.Id,
			Name:        catalog.Name,
			Description: catalog.Description,
			Price:       catalog.Price,
		},
	}, nil
}

func (g *gRPCCatalogServer) GetCatalogs(ctx context.Context, req *proto.GetCatalogsRequest) (*proto.GetCatalogsResponse, error) {
	var res []*domain.Catalog
	var err error
	if req.Query != "" {
		res, err = g.catalogService.SearchCatalog(ctx, &dto.SearchCatalog{
			Query:  req.Query,
			Limit:  req.Limit,
			Offset: req.Offset,
		})
	} else if len(req.Ids) != 0 {
		res, err = g.catalogService.GetCatalogsByIds(ctx, req.Ids)
	} else {
		res, err = g.catalogService.GetCatalogs(ctx, &dto.CatalogQuery{
			Limit:  req.Limit,
			Offset: req.Offset,
		})
	}
	if err != nil {
		return nil, err
	}
	catalogs := make([]*proto.Catalog, 0, len(res))
	for _, c := range res {
		catalogs = append(catalogs, &proto.Catalog{
			Id:          c.Id,
			Name:        c.Name,
			Description: c.Description,
			Price:       c.Price,
		})
	}
	return &proto.GetCatalogsResponse{
		Catalogs: catalogs,
	}, nil
}

func (g *gRPCCatalogServer) Serve(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	g.server = grpc.NewServer()
	proto.RegisterCatalogServiceServer(g.server, g)
	return g.server.Serve(lis)
}

func (g *gRPCCatalogServer) Stop() error {
	if g.server != nil {
		g.server.GracefulStop()
	}
	return nil
}

func NewGRPCCatalogServer(catalogService service.CatalogService) GRPCCatalogServer {
	return &gRPCCatalogServer{
		catalogService: catalogService,
	}
}
