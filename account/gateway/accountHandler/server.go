package accountHandler

import (
	"context"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/dto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/gateway/proto"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/service"
	"google.golang.org/grpc"
	"net"
)

type GRPCAccountServer struct {
	accountService service.AccountService
	server         *grpc.Server
	proto.UnimplementedAccountServiceServer
}

func (g *GRPCAccountServer) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
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

func (g *GRPCAccountServer) GetAccountById(ctx context.Context, req *proto.GetAccountRequest) (*proto.GetAccountResponse, error) {
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

func (g *GRPCAccountServer) GetAccounts(ctx context.Context, req *proto.GetAccountsRequest) (*proto.GetAccountsResponse, error) {
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

func (g *GRPCAccountServer) Serve(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	proto.RegisterAccountServiceServer(grpcServer, g)
	return grpcServer.Serve(lis)
}

func (g *GRPCAccountServer) Stop() error {
	if g.server != nil {
		g.server.GracefulStop()
	}
	return nil
}

func NewGRPCServer(accountService service.AccountService) *GRPCAccountServer {
	return &GRPCAccountServer{
		accountService: accountService,
	}
}
