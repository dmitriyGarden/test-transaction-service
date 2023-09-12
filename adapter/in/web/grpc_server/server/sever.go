package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/dmitriyGarden/test-transaction-service/adapter/in/web/grpc_server/api"
	"github.com/dmitriyGarden/test-transaction-service/model"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type IConfig interface {
	GetListen() string
}

type GrpcServer struct {
	api.UnimplementedTransactionServer
	transactionService model.ITransaction
	cfg                IConfig
	l                  model.ILogger
	v                  *validator.Validate
}

func (g *GrpcServer) BalanceUp(ctx context.Context, request *api.BalanceRequest) (*api.BalanceResponse, error) {
	uid, err := g.getUserId(ctx)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, model.ErrInvalidToken) || errors.Is(err, model.ErrAuthRequired) {
			code = codes.Unauthenticated
		}
		g.l.Errorf("getUserId: %v", err)
		return nil, status.Errorf(code, err.Error())
	}
	err = g.transactionService.AddBalance(ctx, uid, int64(request.GetAmount()))
	if err != nil {
		g.l.Error("AddBalance: ", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.BalanceResponse{
		Success: true,
		Message: "The account has been credited",
	}, nil
}

func (g *GrpcServer) BalanceDown(ctx context.Context, request *api.BalanceRequest) (*api.BalanceResponse, error) {
	uid, err := g.getUserId(ctx)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, model.ErrInvalidToken) || errors.Is(err, model.ErrAuthRequired) {
			code = codes.Unauthenticated
		}
		g.l.Errorf("getUserId: %v", err)
		return nil, status.Errorf(code, err.Error())
	}
	err = g.transactionService.AddBalance(ctx, uid, int64(request.GetAmount())*-1)
	if err != nil {
		g.l.Error("AddBalance: ", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.BalanceResponse{
		Success: true,
		Message: "The account has been debited",
	}, nil
}

func New(cnf IConfig, tr model.ITransaction, l model.ILogger) (*GrpcServer, error) {
	return &GrpcServer{
		cfg:                cnf,
		transactionService: tr,
		l:                  l,
		v:                  validator.New(),
	}, nil
}

func (g *GrpcServer) setRequestID(ctx context.Context) context.Context {
	reqID := ""
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		headers := md.Get("x-request-id")
		if len(headers) > 0 {
			reqID = headers[0]
		}
	}
	if reqID == "" {
		reqID = uuid.NewString()
	}
	return context.WithValue(ctx, "requestID", reqID)
}

func (g *GrpcServer) getUserId(ctx context.Context) (uuid.UUID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil, model.ErrAuthRequired
	}
	headers := md.Get("authorization")
	if len(headers) == 0 {
		return uuid.Nil, model.ErrAuthRequired
	}
	token := strings.TrimPrefix(headers[0], "Bearer ")
	uid, err := g.transactionService.GetUserFromToken(token)
	if err != nil {
		return uuid.Nil, fmt.Errorf("userService.GetUserFromToken: %w", err)
	}
	return uid, nil
}

func (g *GrpcServer) Run(ctx context.Context) error {
	s := grpc.NewServer()
	api.RegisterTransactionServer(s, g)
	addr := g.cfg.GetListen()
	l, err := new(net.ListenConfig).Listen(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}
	g.l.Infoln("Listen: ", addr)
	err = s.Serve(l)
	if err != nil {
		return fmt.Errorf("s.Serve: %w", err)
	}
	return nil
}
