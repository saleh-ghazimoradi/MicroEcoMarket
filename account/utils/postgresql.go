package utils

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type Postgresql struct {
	Host         string
	Port         string
	User         string
	Password     string
	Name         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
	SSLMode      string
	Timeout      time.Duration
}

type Options func(*Postgresql)

func WithHost(host string) Options {
	return func(p *Postgresql) {
		p.Host = host
	}
}

func WithPort(port string) Options {
	return func(p *Postgresql) {
		p.Port = port
	}
}

func WithUser(user string) Options {
	return func(p *Postgresql) {
		p.User = user
	}
}

func WithPassword(password string) Options {
	return func(p *Postgresql) {
		p.Password = password
	}
}

func WithName(name string) Options {
	return func(p *Postgresql) {
		p.Name = name
	}
}

func WithMaxOpenConns(maxConns int) Options {
	return func(p *Postgresql) {
		p.MaxOpenConns = maxConns
	}
}

func WithMaxIdleConns(maxIdleConns int) Options {
	return func(p *Postgresql) {
		p.MaxIdleConns = maxIdleConns
	}
}

func WithMaxIdleTime(maxIdleTime time.Duration) Options {
	return func(p *Postgresql) {
		p.MaxIdleTime = maxIdleTime
	}
}

func WithSSLMode(mode string) Options {
	return func(p *Postgresql) {
		p.SSLMode = mode
	}
}

func WithTimeout(timeout time.Duration) Options {
	return func(p *Postgresql) {
		p.Timeout = timeout
	}
}

func (p *Postgresql) uri() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", p.Host, p.Port, p.User, p.Password, p.Name, p.SSLMode)
}

func (p *Postgresql) Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", p.uri())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(p.MaxOpenConns)
	db.SetMaxIdleConns(p.MaxIdleConns)
	db.SetConnMaxLifetime(p.Timeout)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func NewPostgresql(opts ...Options) *Postgresql {
	postgresql := &Postgresql{}
	for _, opt := range opts {
		opt(postgresql)
	}
	return postgresql
}
