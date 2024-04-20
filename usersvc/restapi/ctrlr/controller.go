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
	"userService/usersvc/common/domain"
	"userService/usersvc/restapi/aescryptor"
	"userService/usersvc/restapi/ctrlr/dto"
	"userService/usersvc/restapi/jwtutils"
	"userService/usersvc/restapi/ooauth"
	"userService/usersvc/restapi/service"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jae2274/goutils/llog"
)

type Controller struct {
	googleOauth       ooauth.Ooauth
	jwtResolver       *jwtutils.JwtResolver
	aesCryptor        *aescryptor.JsonAesCryptor
	userService       service.UserService
	router            *mux.Router
	store             *sessions.CookieStore
	afterAuthHtmlTmpl *template.Template
}

//go:embed after_auth.html
var afterLoginHtml string

func NewController(googleOauth ooauth.Ooauth, router *mux.Router, jwtResolver *jwtutils.JwtResolver, aesCryptor *aescryptor.JsonAesCryptor, userService service.UserService) *Controller {

	afterLoginHtmlTmpl, err := template.New("afterLogin").Parse(afterLoginHtml)

	if err != nil {
		panic(err)
	}

	return &Controller{
		googleOauth:       googleOauth,
		router:            router,
		jwtResolver:       jwtResolver,
		aesCryptor:        aesCryptor,
		userService:       userService,
		store:             sessions.NewCookieStore([]byte("secret")),
		afterAuthHtmlTmpl: afterLoginHtmlTmpl,
	}
}

func (c *Controller) RegisterRoutes() {
	c.router.HandleFunc("/auth/auth-code-urls", c.AuthCodeUrls)
	c.router.HandleFunc("/auth/callback/google", c.Authenticate)
	c.router.HandleFunc("/auth/sign-in", c.SignIn)
	c.router.HandleFunc("/auth/sign-up", c.SignUp)
}

func (c *Controller) AuthCodeUrls(w http.ResponseWriter, r *http.Request) {
	session, _ := c.store.Get(r, "session")
	session.Options = &sessions.Options{
		Path:   "/auth",
		MaxAge: 300,
	}
	state := randToken()
	session.Values["state"] = state
	session.Save(r, w)

	gottenState := session.Values["state"]
	fmt.Printf("state: %v\n", gottenState)

	res := &dto.AuthCodeUrlsResponse{
		AuthCodeUrls: []*dto.AuthCodeUrlRes{
			{AuthServer: string(c.googleOauth.GetAuthServer()), Url: c.googleOauth.GetLoginURL(state)},
		},
	}

	json.NewEncoder(w).Encode(res)
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (c *Controller) Authenticate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		if err != nil {
			c.authFailed(ctx, w, err)
		}
	}()

	session, _ := c.store.Get(r, "session")
	state := session.Values["state"]

	delete(session.Values, "state")
	session.Save(r, w)

	if state != r.FormValue("state") {
		err = fmt.Errorf("invalid state. state: %v, r.FormValue('state'): %s", state, r.FormValue("state"))
		return
	}

	token, err := c.googleOauth.GetToken(r.Context(), r.FormValue("code"))
	if err != nil {
		return
	}

	authToken, err := c.aesCryptor.Encrypt(token)
	if err != nil {
		return
	}

	c.authSuccess(ctx, w, authToken)
}

func (c *Controller) authFailed(ctx context.Context, w http.ResponseWriter, err error) {
	llog.LogErr(ctx, err)
	c.afterAuthHtmlTmpl.Execute(w, &dto.AfterAuthViewVars{
		AuthStatus: dto.AuthFailed,
	})
}

func (c *Controller) authSuccess(ctx context.Context, w http.ResponseWriter, authToken string) {
	c.afterAuthHtmlTmpl.Execute(w, &dto.AfterAuthViewVars{
		AuthStatus: dto.AuthSuccess,
		AuthToken:  authToken,
	})
}

func (c *Controller) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		if err != nil {
			c.signInError(ctx, w, err)
		}
	}()

	var req dto.SignInRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return
	}

	token := &ooauth.OauthToken{}
	err = c.aesCryptor.Decrypt(req.AuthToken, token)
	if err != nil {
		return
	}

	user, err := c.googleOauth.GetUserInfo(ctx, token)
	if err != nil {
		return
	}

	userinfo, isExisted, err := c.userService.GetUser(ctx, user.AuthorizedBy, user.AuthorizedID)
	if err != nil {
		return
	}

	if isExisted {
		c.signInSuccess(ctx, w, userinfo)
		return
	}

	agreements, err := c.userService.GetAgreements(ctx)
	if err != nil {
		return
	}

	c.signInNewUser(ctx, w, user.Email, agreements)
}

func (c *Controller) signInSuccess(ctx context.Context, w http.ResponseWriter, user *domain.User) {
	jwt, err := c.jwtResolver.CreateToken(user)
	if err != nil {
		c.signInError(ctx, w, err)
		return
	}

	res := &dto.SignInResponse{
		SignInStatus: dto.SignInSuccess,
		SuccessRes: &dto.SignInSuccessRes{
			GrantType:    jwt.GrantType,
			AccessToken:  jwt.AccessToken,
			RefreshToken: jwt.RefreshToken,
		},
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		c.signInError(ctx, w, err)
		return
	}
}

func (c *Controller) signInNewUser(ctx context.Context, w http.ResponseWriter, email string, ags []*domain.Agreement) {
	agreements := make([]*dto.AgreementRes, len(ags))
	for i, ag := range ags {
		agreements[i] = &dto.AgreementRes{
			AgreementCode: ag.AgreementCode,
			IsRequired:    ag.IsRequired,
			Summary:       ag.Summary,
		}
	}

	res := &dto.SignInResponse{
		SignInStatus: dto.SignInNewUser,
		NewUserRes: &dto.SignInNewUserRes{
			Email:      email,
			Agreements: agreements,
		},
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		c.signInError(ctx, w, err)
		return
	}
}

func (c *Controller) signInError(ctx context.Context, w http.ResponseWriter, err error) {
	llog.LogErr(ctx, err)

	res := &dto.SignInResponse{
		SignInStatus: dto.SignInFailed,
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (c *Controller) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	var req dto.SignUpRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return
	}

	token := &ooauth.OauthToken{}
	err = c.aesCryptor.Decrypt(req.AuthToken, token)
	if err != nil {
		return
	}

	user, err := c.googleOauth.GetUserInfo(ctx, token)
	if err != nil {
		return
	}

	ags := make([]*domain.UserAgreement, len(req.Agreements))
	for i, ag := range req.Agreements {
		ags[i] = &domain.UserAgreement{
			AgreementID: ag.AgreementID,
			IsAgree:     ag.IsAgree,
		}
	}

	err = c.userService.SaveUser(ctx, user.AuthorizedBy, user.AuthorizedID, user.Email, ags)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusCreated)
}
