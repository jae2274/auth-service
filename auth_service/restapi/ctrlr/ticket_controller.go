package ctrlr

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jae2274/auth-service/auth_service/restapi/middleware"
	"github.com/jae2274/auth-service/auth_service/restapi/service"
)

type TicketController struct {
	ticketService *service.TicketService
}

func NewTicketController(ticketService *service.TicketService) *TicketController {
	return &TicketController{
		ticketService: ticketService,
	}
}

func (c *TicketController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/ticket/{ticketUUID}", c.GetTicketInfo).Methods("GET")
	router.HandleFunc("auth/ticket/{ticketUUID}", c.UseTicket).Methods("PATCH")
}

func (c *TicketController) GetTicketInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	ticketUUID := vars["ticketUUID"]

	ticket, isExisted, err := c.ticketService.GetTicketInfo(ctx, ticketUUID)
	if errorHandler(ctx, w, err) {
		return
	}

	if !isExisted {
		http.Error(w, "ticket not existed", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ticket)
	if errorHandler(ctx, w, err) {
		return
	}
}

func (c *TicketController) UseTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	ticketUUID := vars["ticketUUID"]

	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		http.Error(w, "no claims in context", http.StatusUnauthorized)
		return
	}
	userId, err := strconv.Atoi(claims.UserId)
	if errorHandler(ctx, w, err) {
		return
	}

	isExisted, err := c.ticketService.UseTicket(ctx, userId, ticketUUID)
	if errorHandler(ctx, w, err) {
		return
	}

	if !isExisted {
		http.Error(w, "ticket not existed", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
