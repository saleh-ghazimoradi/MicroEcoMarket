package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/domain"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/dto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/repository"
	"github.com/segmentio/ksuid"
)

type CatalogService interface {
	CreateCatalog(ctx context.Context, input *dto.Catalog) (*domain.Catalog, error)
	GetCatalogById(ctx context.Context, id string) (*domain.Catalog, error)
	GetCatalogs(ctx context.Context, input *dto.CatalogQuery) ([]*domain.Catalog, error)
}

type catalogService struct {
	catalogRepository repository.CatalogRepository
}

func (c *catalogService) CreateCatalog(ctx context.Context, input *dto.Catalog) (*domain.Catalog, error) {
	if input.Price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}
	catalog := &domain.Catalog{
		Id:          ksuid.New().String(),
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
	}
	if err := c.catalogRepository.CreateCatalog(ctx, catalog); err != nil {
		return nil, fmt.Errorf("create catalog failed: %w", err)
	}

	return catalog, nil
}

func (c *catalogService) GetCatalogById(ctx context.Context, id string) (*domain.Catalog, error) {
	return c.catalogRepository.GetCatalogById(ctx, id)
}

func (c *catalogService) GetCatalogs(ctx context.Context, input *dto.CatalogQuery) ([]*domain.Catalog, error) {
	return c.catalogRepository.GetCatalogs(ctx, input)
}

func NewCatalogService(catalogRepository repository.CatalogRepository) CatalogService {
	return &catalogService{
		catalogRepository: catalogRepository,
	}
}
