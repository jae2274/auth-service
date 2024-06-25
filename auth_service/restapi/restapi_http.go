package restapi

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/jae2274/auth-service/auth_service/common/domain"
	"github.com/jae2274/auth-service/auth_service/common/vars"
	"github.com/jae2274/auth-service/auth_service/restapi/aescryptor"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr"
	"github.com/jae2274/auth-service/auth_service/restapi/jwtresolver"
	"github.com/jae2274/auth-service/auth_service/restapi/middleware"
	"github.com/jae2274/auth-service/auth_service/restapi/ooauth"
	"github.com/jae2274/auth-service/auth_service/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jae2274/goutils/llog"
	"github.com/jae2274/goutils/mw/httpmw"
)

func Run(ctx context.Context, envVars *vars.Vars, db *sql.DB) error {
	router := mux.NewRouter()
	router.Use(httpmw.SetTraceIdMW()) //TODO: 불필요한 파라미터가 잘못 포함되어 있어 이후 라이브러리 수정 필요
	jwtResolver := jwtresolver.NewJwtResolver(envVars.SecretKey)
	router.Use(middleware.SetClaimsMW(jwtResolver))

	aesCryptor, err := aescryptor.NewJsonAesCryptor(utils.CreateHash(envVars.SecretKey))
	if err != nil {
		return err
	}
	googleAuth := ooauth.NewGoogleOauth(envVars.GoogleClientID, envVars.GoogleClientSecret, envVars.GoogleRedirectUrl)
	controller := ctrlr.NewController(db, jwtResolver, aesCryptor, googleAuth)
	controller.RegisterRoutes(router)

	ticketController := ctrlr.NewTicketController(db, jwtResolver)
	ticketController.RegisterRoutes(router)

	adminRouter := router.NewRoute().Subrouter()
	adminController := ctrlr.NewAdminController(db)
	adminController.RegisterRoutes(adminRouter)
	adminRouter.Use(middleware.CheckHasAuthority(domain.AuthorityAdmin))

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
