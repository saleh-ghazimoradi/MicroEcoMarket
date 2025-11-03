package config

type Application struct {
	AccountPort string `env:"ACCOUNT_PORT"`
	CatalogPort string `env:"CATALOG_PORT"`
	OrderPort   string `env:"ORDER_PORT"`
}
