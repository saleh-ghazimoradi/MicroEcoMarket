package service

import (
	"context"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/domain"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/dto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/repository"
	"github.com/segmentio/ksuid"
)

type AccountService interface {
	CreateAccount(ctx context.Context, input *dto.Account) (*domain.Account, error)
	GetAccountById(ctx context.Context, id string) (*domain.Account, error)
	GetAccounts(ctx context.Context, input *dto.AccountQuery) ([]*domain.Account, error)
}

type accountService struct {
	accountRepository repository.AccountRepository
}

func (a *accountService) CreateAccount(ctx context.Context, input *dto.Account) (*domain.Account, error) {
	account := &domain.Account{
		Id:   ksuid.New().String(),
		Name: input.Name,
	}
	if err := a.accountRepository.CreateAccount(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func (a *accountService) GetAccountById(ctx context.Context, id string) (*domain.Account, error) {
	return a.accountRepository.GetAccountById(ctx, id)
}

func (a *accountService) GetAccounts(ctx context.Context, input *dto.AccountQuery) ([]*domain.Account, error) {
	if input.Limit > 100 || (input.Offset == 0 && input.Limit == 0) {
		input.Limit = 100
	}
	return a.accountRepository.GetAccounts(ctx, input)
}

func NewAccountService(accountRepository repository.AccountRepository) AccountService {
	return &accountService{
		accountRepository: accountRepository,
	}
}
