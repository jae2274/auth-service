package ctrlr

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
	"github.com/jae2274/auth-service/auth_service/restapi/service"
)

type AdminController struct {
	userService   service.UserService
	ticketService *service.TicketService
}

func NewAdminController(userService service.UserService, ticketService *service.TicketService) *AdminController {
	return &AdminController{userService: userService, ticketService: ticketService}
}

func (a *AdminController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/admin/authority", a.AddAuthority).Methods("POST")
	router.HandleFunc("/auth/admin/authority", a.RemoveAuthority).Methods("DELETE")
	router.HandleFunc("/auth/admin/ticket", a.CreateTicket).Methods("POST")
}

func (a *AdminController) AddAuthority(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.AddAuthorityRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if errorHandler(ctx, w, err) {
		return
	}

	if err := a.userService.AddUserAuthorities(ctx, req.UserId, req.Authorities); errorHandler(ctx, w, err) {
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *AdminController) RemoveAuthority(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.RemoveAuthorityRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if errorHandler(ctx, w, err) {
		return
	}

	if err := a.userService.RemoveAuthority(ctx, req.UserId, req.AuthorityCode); errorHandler(ctx, w, err) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *AdminController) CreateTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req *dto.CreateTicketRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if errorHandler(ctx, w, err) {
		return
	}

	ticketId, err := a.ticketService.CreateTicket(ctx, req.TicketAuthorities)
	if errorHandler(ctx, w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(struct {
		TicketUUID string `json:"ticketUUID"`
	}{TicketUUID: ticketId})

	if errorHandler(ctx, w, err) {
		return
	}
}
