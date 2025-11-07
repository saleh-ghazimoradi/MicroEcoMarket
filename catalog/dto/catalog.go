package dto

type Catalog struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"omitempty"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}

type CatalogQuery struct {
	Limit  uint64   `json:"limit" validate:"omitempty,gte=0,lte=100"`
	Offset uint64   `json:"offset" validate:"omitempty,gte=0"`
	Query  string   `json:"query,omitempty"` // Fixed: omitempty (was "query,omitezero")
	Ids    []string `json:"ids,omitempty"`   // Fixed: omitempty
}

type SearchCatalog struct { // Consistent naming
	Query  string `json:"query" validate:"required"`
	Limit  uint64 `json:"limit" validate:"omitempty,gte=0,lte=100"`
	Offset uint64 `json:"offset" validate:"omitempty,gte=0"`
}
