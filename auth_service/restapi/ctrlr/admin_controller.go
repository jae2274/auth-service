package ctrlr

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
	"github.com/jae2274/auth-service/auth_service/restapi/service"
)

type AdminController struct {
	userService service.UserService
}

func NewAdminController(userService service.UserService) *AdminController {
	return &AdminController{userService: userService}
}

func (a *AdminController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/admin/authority", a.AddAuthority).Methods("POST")
}

func (a *AdminController) AddAuthority(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.AddAuthorityRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if errorHandler(ctx, w, err) {
		return
	}

	if err := a.userService.AddUserAuthorities(ctx, req.UserId, req.AddedAuthorities); errorHandler(ctx, w, err) {
		return
	}

	w.WriteHeader(http.StatusCreated)
}
