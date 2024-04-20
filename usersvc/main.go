package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"userService/usersvc/common/mysqldb"
	"userService/usersvc/common/vars"
	"userService/usersvc/restapi/aescryptor"
	"userService/usersvc/restapi/ctrlr"
	"userService/usersvc/restapi/jwtutils"
	"userService/usersvc/restapi/ooauth"
	"userService/usersvc/restapi/service"
	"userService/usersvc/utils"

	"github.com/gorilla/handlers"
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

	secretKey := utils.CreateHash(envVars.SecretKey)
	jwtResolver := jwtutils.NewJwtUtils(secretKey)
	aesCryptor, err := aescryptor.NewJsonAesCryptor(secretKey)
	check(mainCtx, err)

	router := mux.NewRouter()
	router.Use(httpmw.SetTraceIdMW()) //TODO: 불필요한 파라미터가 잘못 포함되어 있어 이후 라이브러리 수정 필요

	controller := ctrlr.NewController(googleAuth, router, jwtResolver, aesCryptor, userService)
	controller.RegisterRoutes()

	var allowOrigins []string
	if envVars.AccessControlAllowOrigin != nil {
		allowOrigins = append(allowOrigins, *envVars.AccessControlAllowOrigin)
	}
	originsOk := handlers.AllowedOrigins(allowOrigins)
	credentialsOk := handlers.AllowCredentials()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	llog.Msg("Start api server").Data("port", envVars.ApiPort).Log(mainCtx)
	err = http.ListenAndServe(fmt.Sprintf(":%d", envVars.ApiPort), handlers.CORS(originsOk, credentialsOk, headersOk, methodsOk)(router))
	check(mainCtx, err)
}

func check(ctx context.Context, err error) {
	if err != nil {
		llog.LogErr(ctx, err)
		panic(err)
	}
}
