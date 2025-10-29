package dto

type Catalog struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"omitempty"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}

type CatalogQuery struct {
	Limit  uint64   `json:"limit"`
	Offset uint64   `json:"offset"`
	Ids    []string `json:"ids,omitempty"`
	Query  string   `json:"query,omitempty"`
}
