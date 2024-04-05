package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"userService/usersvc/ctrlr"
	"userService/usersvc/jwtutils"
	"userService/usersvc/mysqldb"
	"userService/usersvc/ooauth"
	"userService/usersvc/service"
	"userService/usersvc/vars"

	"github.com/gorilla/mux"
	"github.com/jae2274/goutils/llog"
	"github.com/jae2274/goutils/mw"
	"github.com/jae2274/goutils/mw/httpmw"
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

	userService := service.NewUserService(db)

	googleAuth := ooauth.NewGoogleOauth(envVars.GoogleClientID, envVars.GoogleClientSecret, envVars.GoogleRedirectUrl)
	jwtResolver := jwtutils.NewJwtUtils(envVars.SecretKey)
	router := mux.NewRouter()
	router.Use(httpmw.SetTraceIdMW("TODO_REMOVE_ME")) //TODO: 불필요한 파라미터가 잘못 포함되어 있어 이후 라이브러리 수정 필요

	controller := ctrlr.NewController(googleAuth, router, jwtResolver, userService)
	controller.RegisterRoutes()

	llog.Msg("Start api server").Data("port", envVars.ApiPort).Log(mainCtx)
	err = http.ListenAndServe(fmt.Sprintf(":%d", envVars.ApiPort), router)
	check(mainCtx, err)
}

func check(ctx context.Context, err error) {
	if err != nil {
		llog.LogErr(ctx, err)
		panic(err)
	}
}
