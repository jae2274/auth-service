package ctrlr

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jae2274/auth-service/auth_service/common/mysqldb"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
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
	router.HandleFunc("/auth/ticket/{ticketUUID}", c.UseTicket).Methods("PATCH")
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

	res, err := mysqldb.WithTransaction(ctx, c.db, func(tx *sql.Tx) (*dto.UseTicketResponse, error) {
		return c.useTicket(ctx, tx, userId, ticketUUID)
	})

	if errorHandler(ctx, w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if errorHandler(ctx, w, err) {
		return
	}
}

func (c *TicketController) useTicket(ctx context.Context, tx *sql.Tx, userId int, ticketId string) (*dto.UseTicketResponse, error) {
	res := &dto.UseTicketResponse{}
	ticket, isExisted, err := service.GetTicketInfo(ctx, tx, ticketId)
	if err != nil {
		return nil, err
	}

	if !isExisted {
		res.TicketStatus = dto.NOT_EXISTED
		return res, nil
	} else if ticket.IsUsed {
		res.TicketStatus = dto.ALREADY_USED
		return res, nil
	}

	err = service.UseTicket(ctx, tx, userId, ticketId)
	if err != nil {
		return nil, err
	}
	authorityIds := make([]int, len(ticket.TicketAuthorities))
	for i, authority := range ticket.TicketAuthorities {
		authorityIds[i] = authority.AuthorityId
	}

	/*
		추후 해당 기능이 사용될 여지가 있을 것으로 판단되어 주석처리하였습니다.
	*/
	// userAuthorities, err := service.FindUserAuthoritiesByAuthorityIds(ctx, tx, userId, authorityIds)
	// if err != nil {
	// 	return nil, err
	// }
	allAuthorities, err := service.FindValidUserAuthorities(ctx, tx, userId)
	if err != nil {
		return nil, err
	}
	allAuthorityCodes := make([]string, 0, len(allAuthorities))
	for _, authority := range allAuthorities {
		allAuthorityCodes = append(allAuthorityCodes, authority.AuthorityCode)
	}
	tokens, err := c.jwtResolver.CreateToken(strconv.Itoa(userId), allAuthorityCodes, time.Now())
	if err != nil {
		return nil, err
	}
	res.TicketStatus = dto.SUCCESSFULLY_USED
	res.AccessToken = &tokens.AccessToken
	res.Authorities = allAuthorityCodes
	// res.AppliedAuthorities = userAuthorities
	return res, nil
}
