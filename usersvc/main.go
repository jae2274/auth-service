package main

import (
	"context"
	"fmt"
	"net/http"
	"userService/usersvc/ctrlr"
	"userService/usersvc/ooauth"
	"userService/usersvc/vars"

	"github.com/gorilla/mux"
	"github.com/jae2274/goutils/llog"
)

func main() {
	mainCtx := context.Background()
	envVars, err := vars.Variables()
	if err != nil {
		llog.LogErr(mainCtx, err)
	}

	googleAuth := ooauth.NewGoogleOauth(envVars.GoogleClientID, envVars.GoogleClientSecret, envVars.RedirectURL)

	router := mux.NewRouter()
	controller := ctrlr.NewController(googleAuth, router)
	controller.RegisterRoutes()

	err = http.ListenAndServe(fmt.Sprintf(":%d", envVars.Port), router)
	if err != nil {
		llog.LogErr(mainCtx, err)
	}
}
