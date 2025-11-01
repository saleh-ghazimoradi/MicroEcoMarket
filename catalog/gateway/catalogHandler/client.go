package catalogHandler

import (
	"context"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/domain"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/dto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/gateway/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCCatalogClient interface {
	CreateCatalog(ctx context.Context, input *dto.Catalog) (*domain.Catalog, error)
	GetCatalogById(ctx context.Context, id string) (*domain.Catalog, error)
	GetCatalogs(ctx context.Context, input *dto.CatalogQuery) ([]*domain.Catalog, error)
	Close() error
}

type gRPCCatalogClient struct {
	conn   *grpc.ClientConn
	client proto.CatalogServiceClient
}

func (g *gRPCCatalogClient) CreateCatalog(ctx context.Context, input *dto.Catalog) (*domain.Catalog, error) {
	req := &proto.CreateCatalogRequest{Name: input.Name, Description: input.Description, Price: input.Price}
	resp, err := g.client.CreateCatalog(ctx, req)
	if err != nil {
		return nil, err
	}
	return &domain.Catalog{
		Id:          resp.Catalog.Id,
		Name:        resp.Catalog.Name,
		Description: resp.Catalog.Description,
		Price:       resp.Catalog.Price,
	}, nil
}

func (g *gRPCCatalogClient) GetCatalogById(ctx context.Context, id string) (*domain.Catalog, error) {
	req := &proto.GetCatalogRequest{Id: id}
	resp, err := g.client.GetCatalogById(ctx, req)
	if err != nil {
		return nil, err
	}

	return &domain.Catalog{
		Id:          resp.Catalog.Id,
		Name:        resp.Catalog.Name,
		Description: resp.Catalog.Description,
		Price:       resp.Catalog.Price,
	}, nil
}

func (g *gRPCCatalogClient) GetCatalogs(ctx context.Context, input *dto.CatalogQuery) ([]*domain.Catalog, error) {
	req := &proto.GetCatalogsRequest{Limit: input.Limit, Offset: input.Offset, Ids: input.Ids, Query: input.Query}
	resp, err := g.client.GetCatalogs(ctx, req)
	if err != nil {
		return nil, err
	}
	catalogs := make([]*domain.Catalog, len(resp.Catalogs))
	for _, catalog := range resp.Catalogs {
		catalogs = append(catalogs, &domain.Catalog{
			Id:          catalog.Id,
			Name:        catalog.Name,
			Description: catalog.Description,
			Price:       catalog.Price,
		})
	}
	return catalogs, nil
}

func (g *gRPCCatalogClient) Close() error {
	return g.conn.Close()
}

func NewGRPCCatalogClient(addr string) (GRPCCatalogClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	catalogServiceClient := proto.NewCatalogServiceClient(conn)
	return &gRPCCatalogClient{
		conn:   conn,
		client: catalogServiceClient,
	}, nil
}
