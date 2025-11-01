package domain

import "time"

type Order struct {
	Id         string            `json:"id"`
	CreatedAt  time.Time         `json:"created_at"`
	TotalPrice float64           `json:"total_price"`
	AccountId  string            `json:"account_id"`
	Catalogs   []*OrderedCatalog `json:"catalogs"`
}

type OrderedCatalog struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    uint32  `json:"quantity"`
}
