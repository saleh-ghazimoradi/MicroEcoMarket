package repository

import (
	"context"
	"encoding/json"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/domain"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/dto"
	"gopkg.in/olivere/elastic.v5"
)

type CatalogRepository interface {
	CreateCatalog(ctx context.Context, catalog *domain.Catalog) error
	GetCatalogById(ctx context.Context, id string) (*domain.Catalog, error)
	GetCatalogs(ctx context.Context, input *dto.CatalogQuery) ([]*domain.Catalog, error)
	GetCatalogsByIds(ctx context.Context, ids []string) ([]*domain.Catalog, error)
	SearchCatalog(ctx context.Context, input *dto.SearchCatalog) ([]*domain.Catalog, error)
}

type catalogRepository struct {
	client *elastic.Client
	index  string
}

func (c *catalogRepository) CreateCatalog(ctx context.Context, catalog *domain.Catalog) error {
	_, err := c.client.Index().Index(c.index).Type("catalog").Id(catalog.Id).BodyJson(catalog).Do(ctx)
	return err
}

func (c *catalogRepository) GetCatalogById(ctx context.Context, id string) (*domain.Catalog, error) {
	res, err := c.client.Get().Index(c.index).Type("catalog").Id(id).Do(ctx)
	if err != nil {
		return nil, err
	}

	if !res.Found {
		return nil, ErrNotFound
	}

	var catalog domain.Catalog
	if err := json.Unmarshal(*res.Source, &catalog); err != nil {
		return nil, err
	}

	return &catalog, nil
}

func (c *catalogRepository) GetCatalogs(ctx context.Context, input *dto.CatalogQuery) ([]*domain.Catalog, error) {
	res, err := c.client.Search().Index(c.index).Type("catalog").Query(elastic.NewMatchAllQuery()).From(int(input.Offset)).Size(int(input.Limit)).Do(ctx)
	if err != nil {
		return nil, err
	}

	var catalogs []*domain.Catalog
	for _, hit := range res.Hits.Hits {
		var catalog domain.Catalog
		if err := json.Unmarshal(*hit.Source, &catalog); err != nil {
			return nil, err
		}
		catalogs = append(catalogs, &catalog)
	}
	return catalogs, nil
}

func (c *catalogRepository) GetCatalogsByIds(ctx context.Context, ids []string) ([]*domain.Catalog, error) {
	var items []*elastic.MultiGetItem
	for _, id := range ids {
		items = append(items, elastic.NewMultiGetItem().Index(c.index).Type("catalog").Id(id))
	}

	res, err := c.client.MultiGet().Add(items...).Do(ctx)
	if err != nil {
		return nil, err
	}
	var catalogs []*domain.Catalog
	for _, doc := range res.Docs {
		var catalog domain.Catalog
		if err := json.Unmarshal(*doc.Source, &catalog); err != nil {
			return nil, err
		}
		catalogs = append(catalogs, &catalog)
	}
	return catalogs, nil
}

func (c *catalogRepository) SearchCatalog(ctx context.Context, input *dto.SearchCatalog) ([]*domain.Catalog, error) {
	res, err := c.client.Search().Index(c.index).Type("catalog").Query(elastic.NewMultiMatchQuery(input.Query, "name", "description")).From(int(input.Offset)).Size(int(input.Limit)).Do(ctx)
	if err != nil {
		return nil, err
	}
	var catalogs []*domain.Catalog
	for _, hit := range res.Hits.Hits {
		var catalog domain.Catalog
		if err := json.Unmarshal(*hit.Source, &catalog); err != nil {
			return nil, err
		}
		catalogs = append(catalogs, &catalog)
	}
	return catalogs, nil
}

func NewCatalogRepository(client *elastic.Client, index string) CatalogRepository {
	return &catalogRepository{
		client: client,
		index:  index,
	}
}
