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

var (
	ErrInvalidInput = errors.New("invalid input: name required, price > 0")
)

type CatalogService interface {
	CreateCatalog(ctx context.Context, input *dto.Catalog) (*domain.Catalog, error)
	GetCatalogById(ctx context.Context, id string) (*domain.Catalog, error)
	GetCatalogs(ctx context.Context, input *dto.CatalogQuery) ([]*domain.Catalog, error)
	GetCatalogsByIds(ctx context.Context, ids []string) ([]*domain.Catalog, error)
	SearchCatalog(ctx context.Context, input *dto.SearchCatalog) ([]*domain.Catalog, error)
}

type catalogService struct {
	catalogRepository repository.CatalogRepository
}

func (c *catalogService) CreateCatalog(ctx context.Context, input *dto.Catalog) (*domain.Catalog, error) {
	if input.Name == "" || input.Price <= 0 {
		return nil, ErrInvalidInput
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
	return catalog, nil // Returns with ID
}

func (c *catalogService) GetCatalogById(ctx context.Context, id string) (*domain.Catalog, error) {
	cat, err := c.catalogRepository.GetCatalogById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get catalog by id failed: %w", err)
	}
	return cat, nil
}

func (c *catalogService) GetCatalogs(ctx context.Context, input *dto.CatalogQuery) ([]*domain.Catalog, error) {
	if input.Limit == 0 || input.Limit > 100 {
		input.Limit = 10 // Default
	}
	if input.Offset == 0 {
		input.Offset = 0
	}
	return c.catalogRepository.GetCatalogs(ctx, input)
}

func (c *catalogService) GetCatalogsByIds(ctx context.Context, ids []string) ([]*domain.Catalog, error) {
	if len(ids) == 0 {
		return []*domain.Catalog{}, nil
	}
	if len(ids) > 50 {
		return nil, errors.New("too many IDs")
	}
	return c.catalogRepository.GetCatalogsByIds(ctx, ids)
}

func (c *catalogService) SearchCatalog(ctx context.Context, input *dto.SearchCatalog) ([]*domain.Catalog, error) {
	if input.Query == "" {
		return nil, errors.New("search query required")
	}
	if input.Limit == 0 || input.Limit > 100 {
		input.Limit = 10
	}
	if input.Offset == 0 {
		input.Offset = 0
	}
	return c.catalogRepository.SearchCatalog(ctx, input)
}

func NewCatalogService(catalogRepository repository.CatalogRepository) CatalogService {
	return &catalogService{
		catalogRepository: catalogRepository,
	}
}
