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
	"userService/usersvc/restapi/ctrlr/dto"
	"userService/usersvc/restapi/service"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jae2274/goutils/llog"
)

type Controller struct {
	userService       service.UserService
	router            *mux.Router
	store             *sessions.CookieStore
	afterAuthHtmlTmpl *template.Template
}

//go:embed after_auth.html
var afterLoginHtml string

func NewController(router *mux.Router, userService service.UserService) *Controller {

	afterLoginHtmlTmpl, err := template.New("afterLogin").Parse(afterLoginHtml)

	if err != nil {
		panic(err)
	}

	return &Controller{
		router:            router,
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

	res := c.userService.AuthCodeUrls(r.Context(), state)

	json.NewEncoder(w).Encode(res)
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (c *Controller) Authenticate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, err := c.store.Get(r, "session")
	if errorHandler(ctx, w, err) {
		return
	}

	state, ok := session.Values["state"]
	if !ok {
		errorHandler(ctx, w, fmt.Errorf("state not found"))
		return
	}

	delete(session.Values, "state")
	session.Save(r, w)

	if state != r.FormValue("state") {
		errorHandler(ctx, w, fmt.Errorf("invalid state. state: %v, r.FormValue('state'): %s", state, r.FormValue("state")))
		return
	}

	res, err := c.userService.Authenticate(ctx, r.FormValue("code"))
	if errorHandler(ctx, w, err) {
		return
	}

	err = c.afterAuthHtmlTmpl.Execute(w, res)
	if errorHandler(ctx, w, err) {
		return
	}
}

func (c *Controller) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.SignInRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if errorHandler(ctx, w, err) {
		return
	}

	res, err := c.userService.SignIn(ctx, req.AuthToken)
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

	err = c.userService.SignUp(ctx, &req)
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
