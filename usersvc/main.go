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
)

const (
	app = "user-service"
	svc = "careerhub"

	ctxKeyTraceID = "trace_id"
)

func main() {
	mainCtx := context.Background()
	err := initLogger(mainCtx)
	check(mainCtx, err)

	envVars, err := vars.Variables()
	check(mainCtx, err)

	db, err := mysqldb.DB(envVars.DbVars)
	check(mainCtx, err)
	defer db.Close()

	userService := service.NewUserService(db)

	googleAuth := ooauth.NewGoogleOauth(envVars.GoogleClientID, envVars.GoogleClientSecret, envVars.GoogleRedirectUrl)
	jwtResolver := jwtutils.NewJwtUtils(envVars.SecretKey)
	router := mux.NewRouter()
	controller := ctrlr.NewController(googleAuth, router, jwtResolver, userService)
	controller.RegisterRoutes()

	err = http.ListenAndServe(fmt.Sprintf(":%d", envVars.ApiPort), router)
	check(mainCtx, err)
}

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

func check(ctx context.Context, err error) {
	if err != nil {
		llog.LogErr(ctx, err)
		panic(err)
	}
}
