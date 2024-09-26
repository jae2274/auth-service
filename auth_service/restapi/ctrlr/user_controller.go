package ctrlr

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"net/http"
	"strconv"

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
		return useTicket(ctx, tx, userId, claims.AuthorizedBy, claims.AuthorizedID, ticketUUID)
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

func useTicket(ctx context.Context, tx *sql.Tx, userId int, authBy domain.AuthorizedBy, authId string, ticketId string) (*dto.UseTicketResponse, error) {
	_, err := service.UseTicket(ctx, tx, userId, ticketId)

	var ticketStatus dto.TicketStatus
	switch err {
	case service.ErrTicketNotFound:
		ticketStatus = dto.NOT_EXISTED
	case service.ErrNoMoreUseableTicket:
		ticketStatus = dto.NO_MORE_USEABLE
	case service.ErrAlreadyUsedTicket:
		ticketStatus = dto.ALREADY_USED
	case nil:
		ticketStatus = dto.SUCCESSFULLY_USED
	default:
		return nil, err
	}

	return &dto.UseTicketResponse{
		TicketStatus: ticketStatus,
	}, nil
}
