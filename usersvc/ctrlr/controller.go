package ctrlr

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"text/template"
	"userService/usersvc/jwtutils"
	"userService/usersvc/ooauth"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Controller struct {
	googleOauth ooauth.Ooauth
	jwtResolver *jwtutils.JwtResolver
	router      *mux.Router
	store       *sessions.CookieStore
}

func NewController(googleOauth ooauth.Ooauth, router *mux.Router, jwtResolver *jwtutils.JwtResolver) *Controller {
	return &Controller{
		googleOauth: googleOauth,
		router:      router,
		jwtResolver: jwtResolver,
		store:       sessions.NewCookieStore([]byte("secret")),
	}
}

func (c *Controller) RegisterRoutes() {
	c.router.HandleFunc("/auth", c.RenderAuthView)
	c.router.HandleFunc("/auth/callback", c.Authenticate)
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, _ := template.ParseFiles(name)
	tmpl.Execute(w, data)
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
	renderTemplate(w, "usersvc/ctrlr/auth.html", c.googleOauth.GetLoginURL(state))
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

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

	jwt, err := c.jwtResolver.CreateToken(userInfo.Email, userInfo.Email, []string{"user"}) //TODO: userID 및 roles 설정
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(jwt)) //TODO: 임시로 jwt를 브라우저에 노출
	http.Redirect(w, r, "/", http.StatusFound)
}
