package dto

type Catalog struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type CatalogQuery struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
}

type CatalogSearch struct {
	Query  string `json:"query"`
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
}
