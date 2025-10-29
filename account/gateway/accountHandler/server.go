package accountHandler

import (
	"context"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/dto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/gateway/proto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/service"
	"google.golang.org/grpc"
	"net"
)

type GRPCAccountServer interface {
	CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error)
	GetAccountById(ctx context.Context, req *proto.GetAccountRequest) (*proto.GetAccountResponse, error)
	GetAccounts(ctx context.Context, req *proto.GetAccountsRequest) (*proto.GetAccountsResponse, error)
	Serve(addr string) error
	Stop() error
}

type gRPCAccountServer struct {
	accountService service.AccountService
	server         *grpc.Server
	proto.UnimplementedAccountServiceServer
}

func (g *gRPCAccountServer) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	account, err := g.accountService.CreateAccount(ctx, &dto.Account{
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}

	return &proto.CreateAccountResponse{
		Account: &proto.Account{
			Id:   account.Id,
			Name: account.Name,
		},
	}, nil
}

func (g *gRPCAccountServer) GetAccountById(ctx context.Context, req *proto.GetAccountRequest) (*proto.GetAccountResponse, error) {
	account, err := g.accountService.GetAccountById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &proto.GetAccountResponse{
		Account: &proto.Account{
			Id:   account.Id,
			Name: account.Name,
		},
	}, nil
}

func (g *gRPCAccountServer) GetAccounts(ctx context.Context, req *proto.GetAccountsRequest) (*proto.GetAccountsResponse, error) {
	accounts, err := g.accountService.GetAccounts(ctx, &dto.AccountQuery{
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return nil, err
	}

	protoAccounts := make([]*proto.Account, len(accounts))
	for i, account := range accounts {
		protoAccounts[i] = &proto.Account{
			Id:   account.Id,
			Name: account.Name,
		}
	}

	return &proto.GetAccountsResponse{
		Accounts: protoAccounts,
	}, nil
}

func (g *gRPCAccountServer) Serve(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	g.server = grpc.NewServer()
	proto.RegisterAccountServiceServer(g.server, g)
	return g.server.Serve(lis)
}

func (g *gRPCAccountServer) Stop() error {
	if g.server != nil {
		g.server.GracefulStop()
	}
	return nil
}

func NewGRPCServer(accountService service.AccountService) GRPCAccountServer {
	return &gRPCAccountServer{
		accountService: accountService,
	}
}
