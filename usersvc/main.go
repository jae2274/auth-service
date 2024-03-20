package main

import (
	"context"
	"fmt"
	"net/http"
	"userService/usersvc/ctrlr"
	"userService/usersvc/jwtutils"
	"userService/usersvc/mysqldb"
	"userService/usersvc/ooauth"
	"userService/usersvc/service"
	"userService/usersvc/vars"

	"github.com/gorilla/mux"
	"github.com/jae2274/goutils/llog"
)

func main() {
	mainCtx := context.Background()
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

func check(ctx context.Context, err error) {
	if err != nil {
		llog.LogErr(ctx, err)
		panic(err)
	}
}
