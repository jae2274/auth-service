package ctrlr

import (
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"net/http"
	"text/template"
	"userService/usersvc/common/domain"
	"userService/usersvc/restapi/jwtutils"
	"userService/usersvc/restapi/ooauth"
	"userService/usersvc/restapi/service"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"gopkg.in/validator.v2"
)

type Controller struct {
	googleOauth        ooauth.Ooauth
	jwtResolver        *jwtutils.JwtResolver
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
	session, _ := c.store.Get(r, "session")
	state := session.Values["state"]

	delete(session.Values, "state")
	session.Save(r, w)

	if state != r.FormValue("state") {
		http.Error(w, "Invalid session state", http.StatusUnauthorized)
		return
	}

	token, err := c.googleOauth.GetToken(r.Context(), r.FormValue("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userInfo, err := c.googleOauth.GetUserInfo(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validator.Validate(userInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := c.getUser(userInfo.AuthorizedBy, userInfo.AuthorizedID, userInfo.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jwt, err := c.jwtResolver.CreateToken(user) //TODO: userID 및 roles 설정
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(jwt)

	c.afterLoginHtmlTmpl.Execute(w, jwt)
}

func (c *Controller) getUser(authorizedBy domain.AuthorizedBy, authorizedID string, email string) (domain.User, error) {
	user, err := c.userService.GetUser(authorizedBy, authorizedID)
	if err != nil {
		return domain.User{}, err
	}

	if user == nil {
		err = c.userService.SaveUser(authorizedBy, authorizedID, email)
		if err != nil {
			return domain.User{}, err
		}

		user, err = c.userService.GetUser(authorizedBy, authorizedID)
		if err != nil {
			return domain.User{}, err
		}
	}

	return *user, nil
}
