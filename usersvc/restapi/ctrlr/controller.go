package ctrlr

import (
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/base64"
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
	"github.com/jae2274/goutils/terr"
	"golang.org/x/oauth2"
	"gopkg.in/validator.v2"
)

type Controller struct {
	googleOauth        ooauth.Ooauth
	jwtResolver        *jwtutils.JwtResolver
	aesCryptor         *aescryptor.JsonAesCryptor
	userService        service.UserService
	router             *mux.Router
	store              *sessions.CookieStore
	authHtmlTmpl       *template.Template
	afterLoginHtmlTmpl *template.Template
}

//go:embed auth.html
var authHtml string

//go:embed after_login.html
var afterLoginHtml string

func NewController(googleOauth ooauth.Ooauth, router *mux.Router, jwtResolver *jwtutils.JwtResolver, userService service.UserService) *Controller {

	authHtmlTmpl, err := template.New("auth").Parse(authHtml)

	if err != nil {
		panic(err)
	}

	afterLoginHtmlTmpl, err := template.New("afterLogin").Parse(afterLoginHtml)

	if err != nil {
		panic(err)
	}

	return &Controller{
		googleOauth:        googleOauth,
		router:             router,
		jwtResolver:        jwtResolver,
		userService:        userService,
		store:              sessions.NewCookieStore([]byte("secret")),
		authHtmlTmpl:       authHtmlTmpl,
		afterLoginHtmlTmpl: afterLoginHtmlTmpl,
	}
}

func (c *Controller) RegisterRoutes() {
	c.router.HandleFunc("/auth/login", c.RenderAuthView)
	c.router.HandleFunc("/auth/callback/google", c.Authenticate)
}

func (c *Controller) RenderAuthView(w http.ResponseWriter, r *http.Request) {
	session, _ := c.store.Get(r, "session")
	session.Options = &sessions.Options{
		Path:   "/auth",
		MaxAge: 300,
	}
	state := randToken()
	session.Values["state"] = state
	session.Save(r, w)

	c.authHtmlTmpl.Execute(w, c.googleOauth.GetLoginURL(state))
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// TODO: error handling, 현재는 임의로 외부에 에러 내용을 노출하고 있다.
func (c *Controller) Authenticate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session, _ := c.store.Get(r, "session")
	state := session.Values["state"]

	delete(session.Values, "state")
	session.Save(r, w)

	if state != r.FormValue("state") {
		c.loginFailed(ctx, w, terr.New("invalid status"))
		return
	}

	token, err := c.googleOauth.GetToken(r.Context(), r.FormValue("code"))
	if err != nil {
		c.loginFailed(ctx, w, err)
		return
	}

	if err != nil {
		c.loginFailed(ctx, w, err)
		return
	}

	userInfo, err := c.googleOauth.GetUserInfo(r.Context(), token)
	if err != nil {
		c.loginFailed(ctx, w, err)
		return
	}
	if err := validator.Validate(userInfo); err != nil {
		c.loginFailed(ctx, w, err)
		return
	}

	user, isExisted, err := c.userService.GetUser(r.Context(), userInfo.AuthorizedBy, userInfo.AuthorizedID)
	if err != nil {
		c.loginFailed(ctx, w, err)
		return
	}

	if isExisted {
		c.loginSuccess(ctx, w, user)
		return
	} else {
		c.loginNewUser(ctx, w, token, userInfo)
		return
	}
}

func (c *Controller) loginSuccess(ctx context.Context, w http.ResponseWriter, user *domain.User) {
	jwt, err := c.jwtResolver.CreateToken(user)
	if err != nil {
		c.loginFailed(ctx, w, err)
		return
	}

	viewVars := &dto.AfterLoginViewVars{}

	viewVars.LoginStatus = dto.LoginSuccess

	viewVars.GrantType = jwt.GrantType
	viewVars.AccessToken = jwt.AccessToken
	viewVars.RefreshToken = jwt.RefreshToken

	llog.Level(llog.DEBUG).Msg("login success").Data("user", user)
	c.afterLoginHtmlTmpl.Execute(w, viewVars)
}

func (c *Controller) loginNewUser(ctx context.Context, w http.ResponseWriter, token *oauth2.Token, userInfo *ooauth.UserInfo) {
	authToken, err := c.aesCryptor.Encrypt(token)
	if err != nil {
		c.loginFailed(ctx, w, err)
		return
	}

	viewVars := &dto.AfterLoginViewVars{}

	viewVars.LoginStatus = dto.LoginNewUser

	viewVars.AuthToken = authToken
	viewVars.Email = userInfo.Email

	llog.Level(llog.DEBUG).Msg("new user").Data("email", userInfo.Email)
	c.afterLoginHtmlTmpl.Execute(w, viewVars)
}

func (c *Controller) loginFailed(ctx context.Context, w http.ResponseWriter, err error) {
	viewVars := &dto.AfterLoginViewVars{}

	viewVars.LoginStatus = dto.LoginFailed

	llog.LogErr(ctx, err)
	c.afterLoginHtmlTmpl.Execute(w, viewVars)
}
