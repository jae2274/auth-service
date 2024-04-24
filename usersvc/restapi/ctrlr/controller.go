package ctrlr

import (
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"userService/usersvc/restapi/aescryptor"
	"userService/usersvc/restapi/ctrlr/dto"
	"userService/usersvc/restapi/ooauth"
	"userService/usersvc/restapi/service"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jae2274/goutils/llog"
	"github.com/jae2274/goutils/terr"
)

type Controller struct {
	userService       service.UserService
	router            *mux.Router
	store             *sessions.CookieStore
	afterAuthHtmlTmpl *template.Template
	aesCryptor        *aescryptor.JsonAesCryptor
	googleOauth       ooauth.Ooauth
}

//go:embed after_auth.html
var afterLoginHtml string

func NewController(router *mux.Router, userService service.UserService, aesCryptor *aescryptor.JsonAesCryptor, googleOauth ooauth.Ooauth) *Controller {

	afterLoginHtmlTmpl, err := template.New("afterLogin").Parse(afterLoginHtml)

	if err != nil {
		panic(err)
	}

	return &Controller{
		router:            router,
		userService:       userService,
		store:             sessions.NewCookieStore([]byte("secret")),
		afterAuthHtmlTmpl: afterLoginHtmlTmpl,
		aesCryptor:        aesCryptor,
		googleOauth:       googleOauth,
	}
}

func (c *Controller) RegisterRoutes() {
	c.router.HandleFunc("/auth/auth-code-urls", c.AuthCodeUrls)
	c.router.HandleFunc("/auth/callback/google", c.Authenticate)
	c.router.HandleFunc("/auth/user-info", c.UserInfo)
	c.router.HandleFunc("/auth/sign-in", c.SignIn)
	c.router.HandleFunc("/auth/sign-up", c.SignUp)
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

func (c *Controller) UserInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.UserInfoRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if errorHandler(ctx, w, err) {
		return
	}

	ooauthToken, err := decrypt(c.aesCryptor, req.AuthToken)
	if errorHandler(ctx, w, err) {
		return
	}

	json.NewEncoder(w).Encode(&dto.UserInfoResponse{
		Email:    ooauthToken.UserInfo.Email,
		Username: ooauthToken.UserInfo.Username,
	})
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

	res, err := c.userService.SignIn(ctx, ooauthToken.UserInfo.AuthorizedBy, ooauthToken.UserInfo.AuthorizedID)
	if errorHandler(ctx, w, err) {
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if errorHandler(ctx, w, err) {
		return
	}
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

	err = c.userService.SignUp(ctx, req.Username, req.Agreements, ooauthToken.UserInfo.AuthorizedBy, ooauthToken.UserInfo.AuthorizedID, ooauthToken.UserInfo.Email)
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
