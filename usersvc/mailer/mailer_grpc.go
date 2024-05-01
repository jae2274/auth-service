package mailer

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"userService/usersvc/mailer/mailer_grpc"
	mailerserver "userService/usersvc/mailer/mailer_server"

	"github.com/jae2274/goutils/llog"
	"github.com/jae2274/goutils/mw/grpcmw"
	"github.com/jae2274/goutils/terr"
	"google.golang.org/grpc"
)

func middlewares() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			grpcmw.GetTraceIdUnaryMW(),
			grpcmw.LogErrUnaryMW(),
		),
		grpc.ChainStreamInterceptor(
			grpcmw.GetTraceIdStreamMW(),
			grpcmw.LogErrStreamMW(),
		),
	}
}

func Run(ctx context.Context, grpcPort int, db *sql.DB) error {
	server := mailerserver.NewMailerServer(db)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return terr.Wrap(err)
	}

	grpcServer := grpc.NewServer(middlewares()...)
	mailer_grpc.RegisterUserServer(grpcServer, server)

	llog.Msg("Starting mailer grpc server...").Data("port", grpcPort).Log(ctx)
	err = grpcServer.Serve(listener)
	if err != nil {
		return terr.Wrap(err)
	}

	return nil
}
