package restapi

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
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
	"github.com/jae2274/goutils/mw/httpmw"
)

func Run(ctx context.Context, grpcPort int, envVars *vars.Vars, db *sql.DB) error {

	jwtResolver := jwtutils.NewJwtUtils([]byte(envVars.SecretKey))
	userService := service.NewUserService(db, jwtResolver)

	router := mux.NewRouter()
	router.Use(httpmw.SetTraceIdMW()) //TODO: 불필요한 파라미터가 잘못 포함되어 있어 이후 라이브러리 수정 필요

	aesCryptor, err := aescryptor.NewJsonAesCryptor(utils.CreateHash(envVars.SecretKey))
	if err != nil {
		return err
	}

	googleAuth := ooauth.NewGoogleOauth(envVars.GoogleClientID, envVars.GoogleClientSecret, envVars.GoogleRedirectUrl)
	controller := ctrlr.NewController(router, userService, aesCryptor, googleAuth)
	controller.RegisterRoutes()

	var allowOrigins []string
	if envVars.AccessControlAllowOrigin != nil {
		allowOrigins = append(allowOrigins, *envVars.AccessControlAllowOrigin)
	}
	originsOk := handlers.AllowedOrigins(allowOrigins)
	credentialsOk := handlers.AllowCredentials()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	llog.Msg("Starting restapi grpc server...").Data("port", envVars.ApiPort).Log(ctx)
	err = http.ListenAndServe(fmt.Sprintf(":%d", envVars.ApiPort), handlers.CORS(originsOk, credentialsOk, headersOk, methodsOk)(router))
	if err != nil {
		return err
	}

	return nil
}
