package ctrlr

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jae2274/auth-service/auth_service/restapi/jwtresolver"
	"github.com/jae2274/auth-service/auth_service/restapi/middleware"
	"github.com/jae2274/auth-service/auth_service/restapi/service"
)

type TicketController struct {
	db          *sql.DB
	jwtResolver *jwtresolver.JwtResolver
}

func NewTicketController(db *sql.DB, jwtResolver *jwtresolver.JwtResolver) *TicketController {
	return &TicketController{db: db, jwtResolver: jwtResolver}
}

func (c *TicketController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/ticket", c.GetTicketInfo).Methods("GET")
}

func (c *TicketController) GetTicketInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ticketCode := r.URL.Query().Get("ticket_code")
	if ticketCode == "" {
		http.Error(w, "ticket_code is required", http.StatusBadRequest)
		return
	}

	ticket, isExisted, err := service.GetTicketInfo(ctx, c.db, ticketCode)
	if errorHandler(ctx, w, err) {
		return
	}

	if !isExisted {
		http.Error(w, "ticket not existed", http.StatusNotFound)
		return
	}

	claims, isExisted := middleware.GetClaims(ctx)
	if isExisted {
		userId, err := strconv.Atoi(claims.UserId)
		if errorHandler(ctx, w, err) {
			return
		}
		isUsed, err := service.CheckUseTicket(ctx, c.db, userId, ticket.TicketIndexId)
		if errorHandler(ctx, w, err) {
			return
		}

		ticket.AlreadyUsed = isUsed
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ticket)
	if errorHandler(ctx, w, err) {
		return
	}
}
