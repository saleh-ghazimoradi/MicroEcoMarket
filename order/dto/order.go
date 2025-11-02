package dto

type Order struct {
	AccountId string `json:"account_id"`
	Catalogs  []*OrderedCatalog
}

type OrderedCatalog struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    uint32  `json:"quantity"`
}
