package ctrlr

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jae2274/auth-service/auth_service/common/domain"
	"github.com/jae2274/auth-service/auth_service/common/mysqldb"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
	"github.com/jae2274/auth-service/auth_service/restapi/jwtresolver"
	"github.com/jae2274/auth-service/auth_service/restapi/service"

	"github.com/gorilla/mux"
)

type UserController struct {
	db          *sql.DB
	jwtResolver *jwtresolver.JwtResolver
}

func NewUserController(db *sql.DB, jwtResolver *jwtresolver.JwtResolver) *UserController {
	return &UserController{
		db:          db,
		jwtResolver: jwtResolver,
	}
}

func (c *UserController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/authority", c.FindAllUserAuthorities).Methods("GET")
	router.HandleFunc("/auth/withdrawal", c.Withdrawal).Methods("DELETE")
	router.HandleFunc("/auth/ticket/{ticketUUID}", c.UseTicket).Methods("PATCH")
}

func (c *UserController) FindAllUserAuthorities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims := GetClaimsOrPatal(ctx)

	userId, err := strconv.Atoi(claims.UserId)
	if errorHandler(ctx, w, err) {
		return
	}

	userAuthorities, err := service.FindAllUserAuthorities(ctx, c.db, userId)
	if errorHandler(ctx, w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Authorities []*domain.UserAuthority `json:"authorities"`
	}{Authorities: userAuthorities})

	if errorHandler(ctx, w, err) {
		return
	}
}

func (c *UserController) Withdrawal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims := GetClaimsOrPatal(ctx)

	userId, err := strconv.Atoi(claims.UserId)
	if errorHandler(ctx, w, err) {
		return
	}

	err = mysqldb.WithTransactionVoid(ctx, c.db, func(tx *sql.Tx) error {
		return service.Withdrawal(ctx, tx, userId)
	})

	if errorHandler(ctx, w, err) {
		return
	}
}

func (c *UserController) UseTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	ticketUUID := vars["ticketUUID"]

	claims := GetClaimsOrPatal(ctx)
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

func (c *UserController) useTicket(ctx context.Context, tx *sql.Tx, userId int, ticketId string) (*dto.UseTicketResponse, error) {
	res := &dto.UseTicketResponse{}
	ticket, err := service.UseTicket(ctx, tx, userId, ticketId)
	switch err {
	case service.ErrTicketNotFound:
		res.TicketStatus = dto.NOT_EXISTED
		return res, nil
	case service.ErrNoMoreUseableTicket:
		res.TicketStatus = dto.NO_MORE_USEABLE
		return res, nil
	case service.ErrAlreadyUsedTicket:
		res.TicketStatus = dto.ALREADY_USED
		return res, nil
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
