package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/domain"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/dto"
	"net/http"
)

type CatalogRepository interface {
	CreateCatalog(ctx context.Context, catalog *domain.Catalog) error
	GetCatalogById(ctx context.Context, id string) (*domain.Catalog, error)
	GetCatalogs(ctx context.Context, input *dto.CatalogQuery) ([]*domain.Catalog, error)
	GetCatalogsByIds(ctx context.Context, ids []string) ([]*domain.Catalog, error)
	SearchCatalog(ctx context.Context, input *dto.SearchCatalog) ([]*domain.Catalog, error)
}

type catalogRepository struct {
	client *elasticsearch.Client
	index  string
}

func (c *catalogRepository) CreateCatalog(ctx context.Context, catalog *domain.Catalog) error {
	data, err := json.Marshal(catalog)
	if err != nil {
		return fmt.Errorf("failed to marshal catalog: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      c.index,
		DocumentID: catalog.Id,
		Body:       bytes.NewReader(data),
		Refresh:    "true", // Immediate visibility
	}
	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to index catalog: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("elasticsearch index failed with status %d", res.StatusCode)
	}
	return nil
}

func (c *catalogRepository) GetCatalogById(ctx context.Context, id string) (*domain.Catalog, error) {
	req := esapi.GetRequest{
		Index:      c.index,
		DocumentID: id,
	}
	res, err := req.Do(ctx, c.client)
	if err != nil {
		return nil, fmt.Errorf("failed to get catalog: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("elasticsearch get failed with status %d", res.StatusCode)
	}

	var catalog domain.Catalog
	if err := json.NewDecoder(res.Body).Decode(&catalog); err != nil {
		return nil, fmt.Errorf("failed to decode catalog: %w", err)
	}
	return &catalog, nil
}

func (c *catalogRepository) GetCatalogs(ctx context.Context, input *dto.CatalogQuery) ([]*domain.Catalog, error) {
	// Build dynamic search body based on input (query or ids or all)
	body, err := c.buildSearchBody(input)
	if err != nil {
		return nil, err
	}

	var searchReq *esapi.SearchRequest
	if body != nil {
		searchReq = &esapi.SearchRequest{
			Index: []string{c.index},
			Body:  body,
		}
	} else {
		from := int(input.Offset) // Temporary int for *int
		size := int(input.Limit)  // Temporary int for *int
		searchReq = &esapi.SearchRequest{
			Index: []string{c.index},
			From:  &from,
			Size:  &size,
		}
	}

	res, err := searchReq.Do(ctx, c.client)
	if err != nil {
		return nil, fmt.Errorf("failed to get catalogs: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("elasticsearch search failed with status %d", res.StatusCode)
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source domain.Catalog `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search results: %w", err)
	}

	catalogs := make([]*domain.Catalog, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		catalogs = append(catalogs, &hit.Source)
	}
	return catalogs, nil
}

func (c *catalogRepository) GetCatalogsByIds(ctx context.Context, ids []string) ([]*domain.Catalog, error) {
	var mgetReq esapi.MgetRequest
	mgetReq = esapi.MgetRequest{
		Index: c.index,
		Body:  c.buildMgetBody(ids),
	}

	res, err := mgetReq.Do(ctx, c.client)
	if err != nil {
		return nil, fmt.Errorf("failed to multi-get catalogs: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("elasticsearch mget failed with status %d", res.StatusCode)
	}

	var result struct {
		Docs []struct {
			Found  bool           `json:"found"`
			Source domain.Catalog `json:"_source"`
		} `json:"docs"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode mget results: %w", err)
	}

	catalogs := make([]*domain.Catalog, 0)
	for _, doc := range result.Docs {
		if doc.Found {
			catalogs = append(catalogs, &doc.Source)
		}
	}
	return catalogs, nil
}

func (c *catalogRepository) SearchCatalog(ctx context.Context, input *dto.SearchCatalog) ([]*domain.Catalog, error) {
	body, err := c.buildSearchBody(&dto.CatalogQuery{Query: input.Query, Limit: input.Limit, Offset: input.Offset})
	if err != nil {
		return nil, err
	}

	searchReq := esapi.SearchRequest{
		Index: []string{c.index},
		Body:  body,
	}

	res, err := searchReq.Do(ctx, c.client)
	if err != nil {
		return nil, fmt.Errorf("failed to search catalogs: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("elasticsearch search failed with status %d", res.StatusCode)
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source domain.Catalog `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search results: %w", err)
	}

	catalogs := make([]*domain.Catalog, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		catalogs = append(catalogs, &hit.Source)
	}
	return catalogs, nil
}

func (c *catalogRepository) buildSearchBody(input *dto.CatalogQuery) (*bytes.Reader, error) {
	if input.Query != "" {
		queryBody := map[string]interface{}{
			"query": map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":  input.Query,
					"fields": []string{"name", "description"},
				},
			},
			"from": input.Offset,
			"size": input.Limit,
		}
		data, err := json.Marshal(queryBody)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(data), nil
	}
	if len(input.Ids) > 0 {
		idsBody := map[string]interface{}{
			"query": map[string]interface{}{
				"terms": map[string]interface{}{
					"id": input.Ids,
				},
			},
			"from": input.Offset,
			"size": input.Limit,
		}
		data, err := json.Marshal(idsBody)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(data), nil
	}
	// For paginate without query/ids, return nil (use From/Size in SearchRequest)
	return nil, nil
}

func (c *catalogRepository) buildMgetBody(ids []string) *bytes.Reader {
	body := map[string]interface{}{
		"docs": []map[string]interface{}{},
	}
	for _, id := range ids {
		body["docs"] = append(body["docs"].([]map[string]interface{}), map[string]interface{}{
			"_index": c.index,
			"_id":    id,
		})
	}
	data, _ := json.Marshal(body)
	return bytes.NewReader(data)
}

func NewCatalogRepository(client *elasticsearch.Client, index string) CatalogRepository {
	return &catalogRepository{
		client: client,
		index:  index,
	}
}
