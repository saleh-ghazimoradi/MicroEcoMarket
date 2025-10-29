package config

import (
	"github.com/caarlos0/env/v11"
	"sync"
)

var (
	instance *Config
	once     sync.Once
	initErr  error
)

type Config struct {
	Application   Application
	ElasticSearch ElasticSearch
}

func NewConfig() (*Config, error) {
	once.Do(func() {
		instance = &Config{}
		initErr = env.Parse(instance)
		if initErr != nil {
			instance = nil
		}
	})
	return instance, initErr
}
