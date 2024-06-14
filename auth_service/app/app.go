package app

import (
	"context"
	"os"

	"github.com/jae2274/auth-service/auth_service/common/mysqldb"
	"github.com/jae2274/auth-service/auth_service/common/vars"
	"github.com/jae2274/auth-service/auth_service/mailer"
	"github.com/jae2274/auth-service/auth_service/restapi"

	"github.com/jae2274/goutils/llog"
	"github.com/jae2274/goutils/mw"
)

const (
	app = "user-service"
	svc = "careerhub"

	ctxKeyTraceID = string(mw.CtxKeyTraceID)
)

func initLogger(ctx context.Context) error {
	llog.SetMetadata("service", svc)
	llog.SetMetadata("app", app)
	llog.SetDefaultContextData(ctxKeyTraceID)

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	llog.SetMetadata("hostname", hostname)

	return nil
}

func Run(mainCtx context.Context) {

	err := initLogger(mainCtx)
	check(mainCtx, err)
	llog.Info(mainCtx, "Start Application")

	envVars, err := vars.Variables()
	check(mainCtx, err)

	db, err := mysqldb.DB(envVars.DbVars)
	check(mainCtx, err)
	defer db.Close()

	errChan := make(chan error)
	go func() {
		err := restapi.Run(mainCtx, envVars, db)
		errChan <- err
	}()

	go func() {
		err := mailer.Run(mainCtx, envVars.MailerGrpcPort, db)
		errChan <- err
	}()

	select {
	case <-mainCtx.Done():
		llog.Info(mainCtx, "Finished Application")
	case err := <-errChan:
		check(mainCtx, err)
	}
}

func check(ctx context.Context, err error) {
	if err != nil {
		llog.LogErr(ctx, err)
		panic(err)
	}
}
