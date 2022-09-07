//go:generate mockgen -source server.go -destination mock/server_mock.go -package mock

package grpcserver

import (
	"context"
	"net"

	"github.com/Calmantara/go-common/logger"
	"google.golang.org/grpc"
)

type GRPCServer interface {
	// function to return GRPC server
	GetServer() *grpc.Server
	// function to return listener
	GetListener() net.Listener
	// serving GRPC server
	SERVE()
}

type GRPCServerImpl struct {
	sugar      logger.CustomLogger
	Port       string
	grpcServer *grpc.Server
	listener   net.Listener
}

// new constructor method for grpc server
func NewGRPCServer(sugar logger.CustomLogger, port string) GRPCServer {
	// define new struct
	grpcServer := &GRPCServerImpl{
		sugar: sugar,
	}

	//setup grpc listener
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		sugar.Logger().Fatalf("failed to listen: %v", err)
	}
	//setup grpc server
	server := grpc.NewServer()
	//populate all to one struct
	grpcServer.Port = port
	grpcServer.listener = lis
	grpcServer.grpcServer = server

	return grpcServer
}

// function to return GRPC server
func (g *GRPCServerImpl) GetServer() *grpc.Server {
	// get grpc server
	return g.grpcServer
}

// function to return listener
func (g *GRPCServerImpl) GetListener() net.Listener {
	return g.listener
}

// serving GRPC server
func (g *GRPCServerImpl) SERVE() {
	g.sugar.WithContext(context.Background()).Infof("serving grpc server in port:%v", g.Port)
	if err := g.grpcServer.Serve(g.listener); err != nil {
		panic(err)
	}
}
