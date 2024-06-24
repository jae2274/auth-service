package ctrlr

import (
	"context"
	"crypto/rand"
	"database/sql"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/jae2274/auth-service/auth_service/common/domain"
	"github.com/jae2274/auth-service/auth_service/common/mysqldb"
	"github.com/jae2274/auth-service/auth_service/models"
	"github.com/jae2274/auth-service/auth_service/restapi/aescryptor"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
	"github.com/jae2274/auth-service/auth_service/restapi/jwtresolver"
	"github.com/jae2274/auth-service/auth_service/restapi/middleware"
	"github.com/jae2274/auth-service/auth_service/restapi/ooauth"
	"github.com/jae2274/auth-service/auth_service/restapi/service"
	"github.com/jae2274/auth-service/auth_service/utils"
	"github.com/jae2274/goutils/llog"
	"github.com/jae2274/goutils/terr"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Controller struct {
	db                *sql.DB
	jwtResolver       *jwtresolver.JwtResolver
	store             *sessions.CookieStore
	afterAuthHtmlTmpl *template.Template
	aesCryptor        *aescryptor.JsonAesCryptor
	googleOauth       ooauth.Ooauth
}

//go:embed after_auth.html
var afterLoginHtml string

func NewController(db *sql.DB, jwtResolver *jwtresolver.JwtResolver, aesCryptor *aescryptor.JsonAesCryptor, googleOauth ooauth.Ooauth) *Controller {

	afterLoginHtmlTmpl, err := template.New("afterLogin").Parse(afterLoginHtml)

	if err != nil {
		panic(err)
	}

	return &Controller{
		db:                db,
		jwtResolver:       jwtResolver,
		store:             sessions.NewCookieStore([]byte("secret")),
		afterAuthHtmlTmpl: afterLoginHtmlTmpl,
		aesCryptor:        aesCryptor,
		googleOauth:       googleOauth,
	}
}

func (c *Controller) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/auth-code-urls", c.AuthCodeUrls).Methods("GET")
	router.HandleFunc("/auth/callback/google", c.Authenticate)
	router.HandleFunc("/auth/sign-in", c.SignIn).Methods("POST")
	router.HandleFunc("/auth/sign-up", c.SignUp).Methods("POST")
	router.HandleFunc("/auth/refresh", c.RefreshJwt).Methods("POST")
	router.HandleFunc("/auth/authority", c.FindUserAuthorities).Methods("GET")
}

func (c *Controller) AuthCodeUrls(w http.ResponseWriter, r *http.Request) {
	state := randToken()
	err := pushSessionState(c.store, w, r, state)
	if errorHandler(r.Context(), w, err) {
		return
	}

	json.NewEncoder(w).Encode(&dto.AuthCodeUrlsResponse{
		AuthCodeUrls: []*dto.AuthCodeUrlRes{
			{AuthServer: string(c.googleOauth.GetAuthServer()), Url: c.googleOauth.GetLoginURL(state)},
		},
	})
}
func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (c *Controller) Authenticate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	state, err := popSessionState(c.store, w, r)
	if errorHandler(ctx, w, err) {
		return
	}

	if state != r.FormValue("state") {
		errorHandler(ctx, w, fmt.Errorf("invalid state. state: %v, r.FormValue('state'): %s", state, r.FormValue("state")))
		return
	}

	token, err := c.googleOauth.GetToken(ctx, r.FormValue("code"))
	if errorHandler(ctx, w, err) {
		return
	}

	userinfo, err := c.googleOauth.GetUserInfo(ctx, &ooauth.OauthToken{Token: token})
	if errorHandler(ctx, w, err) {
		return
	}

	authToken, err := encrypt(c.aesCryptor, &ooauth.OauthToken{
		UserInfo: userinfo,
		Token:    token,
	})
	if errorHandler(ctx, w, err) {
		return
	}

	err = c.afterAuthHtmlTmpl.Execute(w, &dto.AuthenticateResponse{
		AuthToken: authToken,
	})

	if errorHandler(ctx, w, err) {
		return
	}
}

func popSessionState(s *sessions.CookieStore, w http.ResponseWriter, r *http.Request) (string, error) {
	session, err := s.Get(r, "session")
	if err != nil {
		return "", err
	}

	state, ok := session.Values["state"]
	if !ok {
		return "", terr.New("state not found")
	}

	delete(session.Values, "state")
	session.Save(r, w)

	return state.(string), nil
}

func pushSessionState(s *sessions.CookieStore, w http.ResponseWriter, r *http.Request, state string) error {
	session, err := s.Get(r, "session")
	if err != nil {
		return err
	}

	session.Options = &sessions.Options{
		Path:   "/auth",
		MaxAge: 300,
	}

	session.Values["state"] = state
	session.Save(r, w)

	return nil
}

func (c *Controller) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.SignInRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if errorHandler(ctx, w, err) {
		return
	}

	ooauthToken, err := decrypt(c.aesCryptor, req.AuthToken)
	if errorHandler(ctx, w, err) {
		return
	}

	res, err := mysqldb.WithTransaction(ctx, c.db, func(tx *sql.Tx) (*dto.SignInResponse, error) {
		return c.signIn(ctx, tx, ooauthToken.UserInfo, req.AdditionalAgreements)
	})

	if errorHandler(ctx, w, err) {
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if errorHandler(ctx, w, err) {
		return
	}
}

func (c *Controller) signIn(ctx context.Context, tx *sql.Tx, userinfo *ooauth.UserInfo, additionalAgreements []*dto.UserAgreementReq) (*dto.SignInResponse, error) {
	user, isExisted, err := service.FindSignedUpUser(ctx, tx, userinfo.AuthorizedBy, userinfo.AuthorizedID)
	if err != nil {
		return nil, err
	}

	if isExisted {
		err = service.ApplyUserAgreements(ctx, tx, user.UserID, additionalAgreements)
		if err != nil {
			return nil, err
		}

		agreements, err := service.FindNecessaryAgreements(ctx, tx, user.UserID)
		if err != nil {
			return nil, err
		}

		if len(agreements) > 0 {
			return &dto.SignInResponse{
				SignInStatus:           dto.SignInNecessaryAgreements,
				NecessaryAgreementsRes: signInNecessaryAgreements(agreements),
			}, nil
		} else {
			successRes, err := signInSuccessRes(ctx, tx, c.jwtResolver, user)
			if err != nil {
				return nil, err
			}

			return &dto.SignInResponse{
				SignInStatus: dto.SignInSuccess,
				SuccessRes:   successRes,
			}, nil
		}
	} else {
		signInNewUserRes, err := signInNewUserRes(ctx, tx, userinfo)
		if err != nil {
			return nil, err
		}
		return &dto.SignInResponse{
			SignInStatus: dto.SignInNewUser,
			NewUserRes:   signInNewUserRes,
		}, nil
	}
}

func signInNecessaryAgreements(necessaryAgreements []*models.Agreement) *dto.SignInNecessaryAgreementsRes {
	agreementRes := make([]*dto.AgreementRes, 0, len(necessaryAgreements))
	for _, agreement := range necessaryAgreements {
		agreementRes = append(agreementRes, &dto.AgreementRes{
			AgreementId: agreement.AgreementID,
			Required:    utils.TinyIntToBool(agreement.IsRequired),
			Summary:     agreement.Summary,
			Priority:    agreement.Priority,
		})
	}

	return &dto.SignInNecessaryAgreementsRes{Agreements: agreementRes}
}

func signInSuccessRes(ctx context.Context, db boil.ContextExecutor, jwtResolver *jwtresolver.JwtResolver, user *models.User) (*dto.SignInSuccessRes, error) {
	userAuthorities, err := service.FindUserAuthorities(ctx, db, user.UserID)
	if err != nil {
		return nil, err
	}
	authorityCodes := make([]string, 0, len(userAuthorities))
	for _, authority := range userAuthorities {
		authorityCodes = append(authorityCodes, authority.AuthorityCode)
	}
	token, err := jwtResolver.CreateToken(strconv.Itoa(user.UserID), authorityCodes, time.Now())
	if err != nil {
		return nil, err
	}

	return &dto.SignInSuccessRes{
		Username:     user.Name,
		Authorities:  authorityCodes,
		GrantType:    "Bearer",
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func signInNewUserRes(ctx context.Context, db boil.ContextExecutor, userinfo *ooauth.UserInfo) (*dto.SignInNewUserRes, error) {
	agreements, err := service.FindAllAgreements(ctx, db)
	if err != nil {
		return nil, err
	}

	agreementRes := make([]*dto.AgreementRes, len(agreements))
	for i, agreement := range agreements {
		agreementRes[i] = &dto.AgreementRes{
			AgreementId: agreement.AgreementID,
			Required:    utils.TinyIntToBool(agreement.IsRequired),
			Summary:     agreement.Summary,
		}
	}

	return &dto.SignInNewUserRes{
		Email:      userinfo.Email,
		Username:   userinfo.Username,
		Agreements: agreementRes,
	}, nil
}

func (c *Controller) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if errorHandler(ctx, w, err) {
		return
	}

	ooauthToken, err := decrypt(c.aesCryptor, req.AuthToken)
	if errorHandler(ctx, w, err) {
		return
	}

	err = mysqldb.WithTransactionVoid(ctx, c.db, func(tx *sql.Tx) error {
		_, err = service.SignUp(ctx, tx, ooauthToken.UserInfo, req.Agreements)
		return err
	})

	if errorHandler(ctx, w, err) {
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func errorHandler(ctx context.Context, w http.ResponseWriter, err error) bool {
	if err != nil {
		llog.LogErr(ctx, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return true
	}
	return false
}

func encrypt(aesCryptor *aescryptor.JsonAesCryptor, ooauthToken *ooauth.OauthToken) (string, error) {
	authToken, err := aesCryptor.Encrypt(ooauthToken)
	if err != nil {
		return "", err
	}
	return authToken, nil
}

func decrypt(aesCryptor *aescryptor.JsonAesCryptor, authToken string) (*ooauth.OauthToken, error) {
	ooauthToken := &ooauth.OauthToken{}
	err := aesCryptor.Decrypt(authToken, ooauthToken)
	if err != nil {
		return nil, err
	}
	return ooauthToken, nil
}

func (c *Controller) RefreshJwt(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.RefreshJwtRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if errorHandler(ctx, w, err) {
		return
	}

	claims, isValid, err := c.jwtResolver.ParseToken(req.RefreshToken)
	if errorHandler(ctx, w, err) {
		return
	}

	if !isValid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userId, err := strconv.Atoi(claims.UserId)
	if errorHandler(ctx, w, err) {
		return
	}

	res, err := mysqldb.WithTransaction(ctx, c.db, func(tx *sql.Tx) (*dto.RefreshJwtResponse, error) {
		dUserAuthorities, err := service.FindUserAuthorities(ctx, tx, userId)
		if err != nil {
			return nil, err
		}

		authorityCodes := make([]string, 0, len(dUserAuthorities))
		for _, authority := range dUserAuthorities {
			authorityCodes = append(authorityCodes, authority.AuthorityCode)
		}

		tokens, err := c.jwtResolver.CreateToken(claims.UserId, authorityCodes, time.Now())
		if err != nil {
			return nil, err
		}

		return &dto.RefreshJwtResponse{
			AccessToken: tokens.AccessToken,
			Authorities: authorityCodes,
		}, nil
	})

	if errorHandler(ctx, w, err) {
		return
	}

	err = json.NewEncoder(w).Encode(res)

	if errorHandler(ctx, w, err) {
		return
	}
}

func (c *Controller) FindUserAuthorities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, isExisted := middleware.GetClaims(ctx)
	if !isExisted {
		http.Error(w, "no claims in context", http.StatusUnauthorized)
		return
	}

	userId, err := strconv.Atoi(claims.UserId)
	if errorHandler(ctx, w, err) {
		return
	}

	userAuthorities, err := service.FindUserAuthorities(ctx, c.db, userId)
	if errorHandler(ctx, w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Authorities []*domain.UserAuthority `json:"authorities"`
	}{Authorities: userAuthorities})

	if errorHandler(ctx, w, err) {
		return
	}
}
