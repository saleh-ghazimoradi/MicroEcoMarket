package dto

type Account struct {
	Name string `json:"name"`
}

type AccountQuery struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
}
