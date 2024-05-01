package main

import (
	"context"
	"os"
	"userService/usersvc/common/mysqldb"
	"userService/usersvc/common/vars"
	"userService/usersvc/restapi"

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

func main() {
	mainCtx := context.Background()

	err := initLogger(mainCtx)
	check(mainCtx, err)
	llog.Info(mainCtx, "Start Application")

	envVars, err := vars.Variables()
	check(mainCtx, err)

	db, err := mysqldb.DB(envVars.DbVars)
	check(mainCtx, err)
	defer db.Close()

	err = restapi.Run(mainCtx, envVars.ApiPort, envVars, db)
	check(mainCtx, err)
}

func check(ctx context.Context, err error) {
	if err != nil {
		llog.LogErr(ctx, err)
		panic(err)
	}
}
