package config

import "time"

type ElasticSearch struct {
	Host          string        `env:"ELASTICSEARCH_HOST"`
	Port          string        `env:"ELASTICSEARCH_PORT"`
	ContainerName string        `env:"ELASTICSEARCH_CONTAINER"`
	Username      string        `env:"ELASTICSEARCH_USERNAME"`
	Password      string        `env:"ELASTICSEARCH_PASSWORD"`
	Timeout       time.Duration `env:"ELASTICSEARCH_TIMEOUT"`
}
