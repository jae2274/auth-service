package ctrlr

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
	"github.com/jae2274/auth-service/auth_service/restapi/service"
)

type AdminController struct {
	db *sql.DB
}

func NewAdminController(db *sql.DB) *AdminController {
	return &AdminController{db: db}
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

	if err := service.AddUserAuthorities(ctx, a.db, req.UserId, req.Authorities); errorHandler(ctx, w, err) {
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

	if err := service.RemoveAuthority(ctx, a.db, req.UserId, req.AuthorityCode); errorHandler(ctx, w, err) {
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

	ticketId, err := service.CreateTicket(ctx, a.db, req.TicketAuthorities)
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
