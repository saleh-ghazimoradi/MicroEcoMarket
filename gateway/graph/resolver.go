package graph

import (
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/gateway/accountHandler"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/gateway/catalogHandler"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/order/gateway/orderHandler"
)

type Resolver struct {
	AccountClient accountHandler.GRPCAccountClient
	CatalogClient catalogHandler.GRPCCatalogClient
	OrderClient   orderHandler.GRPCOrderClient
}
