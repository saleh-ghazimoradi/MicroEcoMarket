package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/domain"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/dto"
	"time"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account *domain.Account) error
	GetAccountById(ctx context.Context, id string) (*domain.Account, error)
	GetAccounts(ctx context.Context, accountQuery *dto.AccountQuery) ([]*domain.Account, error)
}

type accountRepository struct {
	dbWrite *sql.DB
	dbRead  *sql.DB
}

func (a *accountRepository) CreateAccount(ctx context.Context, account *domain.Account) error {
	query := `INSERT INTO account(id,name) VALUES ($1,$2)`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := a.dbWrite.ExecContext(ctx, query, account.Id, account.Name)
	return err
}

func (a *accountRepository) GetAccountById(ctx context.Context, id string) (*domain.Account, error) {
	query := `SELECT id, name FROM account WHERE id=$1`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var account domain.Account
	if err := a.dbRead.QueryRowContext(ctx, query, id).Scan(&account.Id, &account.Name); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}
	return &account, nil
}

func (a *accountRepository) GetAccounts(ctx context.Context, accountQuery *dto.AccountQuery) ([]*domain.Account, error) {
	query := `SELECT id, name FROM account ORDER BY id DESC LIMIT $1 OFFSET $2`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	rows, err := a.dbRead.QueryContext(ctx, query, accountQuery.Limit, accountQuery.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var accounts []*domain.Account
	for rows.Next() {
		var account domain.Account
		if err := rows.Scan(&account.Id, &account.Name); err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
}

func NewAccountRepository(dbWrite, dbRead *sql.DB) AccountRepository {
	return &accountRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
