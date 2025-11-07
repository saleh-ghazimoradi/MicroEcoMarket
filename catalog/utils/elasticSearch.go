package utils

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticSearch struct {
	Host     string
	Port     string
	Username string
	Password string
	Timeout  time.Duration
}

type Option func(*ElasticSearch)

func WithHost(host string) Option {
	return func(e *ElasticSearch) {
		e.Host = host
	}
}

func WithPort(port string) Option {
	return func(e *ElasticSearch) {
		e.Port = port
	}
}

func WithUsername(username string) Option {
	return func(e *ElasticSearch) {
		e.Username = username
	}
}

func WithPassword(password string) Option {
	return func(e *ElasticSearch) {
		e.Password = password
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(e *ElasticSearch) {
		e.Timeout = timeout
	}
}

func (e *ElasticSearch) uri() string {
	return fmt.Sprintf("http://%s:%s", e.Host, e.Port)
}

func (e *ElasticSearch) Connect() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses:  []string{e.uri()},
		Username:   e.Username,
		Password:   e.Password,
		MaxRetries: 3,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.Timeout)
	defer cancel()

	res, err := client.Info(client.Info.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to ping elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("elasticsearch healthcheck failed with status: %s", res.Status())
	}

	return client, nil
}

func NewElasticSearch(opts ...Option) *ElasticSearch {
	es := &ElasticSearch{}
	for _, opt := range opts {
		opt(es)
	}
	return es
}
