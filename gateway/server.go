package main

import (
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/gateway/accountHandler"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/gateway/catalogHandler"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/gateway/config"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/gateway/orderHandler"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/gateway/graph"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	accountClient, err := accountHandler.NewGRPCAccountClient(cfg.Application.AccountPort)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := accountClient.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	catalogClient, err := catalogHandler.NewGRPCCatalogClient(cfg.Application.CatalogPort)
	if err != nil {
		defer func() {
			if err := accountClient.Close(); err != nil {
				log.Fatal(err)
			}
		}()
		log.Fatal(err)
	}

	defer func() {
		if err := catalogClient.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	orderClient, err := orderHandler.NewGRPCOrderClient(cfg.Application.OrderPort)
	if err != nil {
		defer func() {
			if err := accountClient.Close(); err != nil {
				log.Fatal(err)
			}
		}()
		defer func() {
			if err := catalogClient.Close(); err != nil {
				log.Fatal(err)
			}
		}()
		log.Fatal(err)
	}

	defer func() {
		if err := orderClient.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{AccountClient: accountClient, CatalogClient: catalogClient, OrderClient: orderClient}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
