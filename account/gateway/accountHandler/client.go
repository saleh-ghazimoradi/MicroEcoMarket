package accountHandler

import (
	"context"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/domain"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/dto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/gateway/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCAccountClient struct {
	conn   *grpc.ClientConn
	client proto.AccountServiceClient
}

func (g *GRPCAccountClient) CreateAccount(ctx context.Context, input *dto.Account) (*domain.Account, error) {
	req := &proto.CreateAccountRequest{Name: input.Name}
	resp, err := g.client.CreateAccount(ctx, req)
	if err != nil {
		return nil, err
	}

	return &domain.Account{
		Id:   resp.Account.Id,
		Name: resp.Account.Name,
	}, nil
}

func (g *GRPCAccountClient) GetAccountById(ctx context.Context, id string) (*domain.Account, error) {
	req := &proto.GetAccountRequest{Id: id}
	resp, err := g.client.GetAccountById(ctx, req)
	if err != nil {
		return nil, err
	}
	return &domain.Account{
		Id:   resp.Account.Id,
		Name: resp.Account.Name,
	}, nil
}

func (g *GRPCAccountClient) GetAccounts(ctx context.Context, input *dto.AccountQuery) ([]*domain.Account, error) {
	req := &proto.GetAccountsRequest{Limit: input.Limit, Offset: input.Offset}
	resp, err := g.client.GetAccounts(ctx, req)
	if err != nil {
		return nil, err
	}

	accounts := make([]*domain.Account, len(resp.Accounts))
	for i, account := range resp.Accounts {
		accounts[i] = &domain.Account{
			Id:   account.Id,
			Name: account.Name,
		}
	}
	return accounts, nil
}

func (g *GRPCAccountClient) Close() error {
	return g.conn.Close()
}

func NewGRPCClient(addr string) (*GRPCAccountClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	accountServiceClient := proto.NewAccountServiceClient(conn)
	return &GRPCAccountClient{
		conn:   conn,
		client: accountServiceClient,
	}, nil
}
