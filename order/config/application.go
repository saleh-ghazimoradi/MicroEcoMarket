package config

type Application struct {
	OrderPort   string `env:"ORDER_PORT"`
	CatalogPort string `env:"CATALOG_PORT"`
	AccountPort string `env:"ACCOUNT_PORT"`
}
