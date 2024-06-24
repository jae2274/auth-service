package ctrlr

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jae2274/auth-service/auth_service/common/mysqldb"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
	"github.com/jae2274/auth-service/auth_service/restapi/service"
	"github.com/jae2274/goutils/terr"
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

	err = mysqldb.WithTransactionVoid(ctx, a.db, func(tx *sql.Tx) error {
		return service.AddUserAuthorities(ctx, tx, req.UserId, req.Authorities)
	})

	if errorHandler(ctx, w, err) {
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

	err = mysqldb.WithTransactionVoid(ctx, a.db, func(tx *sql.Tx) error {
		return service.RemoveAuthority(ctx, tx, req.UserId, req.AuthorityCode)
	})

	if errorHandler(ctx, w, err) {
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

	ticketId, err := mysqldb.WithTransaction(ctx, a.db, func(tx *sql.Tx) (string, error) {
		return service.CreateTicket(ctx, tx, req.TicketAuthorities)
	})
	if errorHandler(ctx, w, err) {
		return
	}

	ticket, isExisted, err := service.GetTicketInfo(ctx, a.db, ticketId)
	if errorHandler(ctx, w, err) {
		return
	}

	if !isExisted {
		errorHandler(ctx, w, terr.New("ticket not found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(ticket)

	if errorHandler(ctx, w, err) {
		return
	}
}
