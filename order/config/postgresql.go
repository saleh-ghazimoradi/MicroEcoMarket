package config

import "time"

type Postgresql struct {
	Host         string        `env:"POSTGRESQL_HOST"`
	Port         string        `env:"POSTGRESQL_PORT"`
	User         string        `env:"POSTGRESQL_USER"`
	Password     string        `env:"POSTGRESQL_PASSWORD"`
	Name         string        `env:"POSTGRESQL_NAME"`
	MaxOpenConns int           `env:"POSTGRESQL_MAX_OPEN_CONNS"`
	MaxIdleConns int           `env:"POSTGRESQL_MAX_IDLE_CONNS"`
	MaxIdleTime  time.Duration `env:"POSTGRESQL_MAX_IDLE_TIME"`
	SSLMode      string        `env:"POSTGRESQL_SSL_MODE"`
	Timeout      time.Duration `env:"POSTGRESQL_TIMEOUT"`
}
