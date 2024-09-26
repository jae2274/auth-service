package ctrlr

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jae2274/auth-service/auth_service/common/domain"
	"github.com/jae2274/auth-service/auth_service/common/mysqldb"
	"github.com/jae2274/auth-service/auth_service/restapi/middleware"
	"github.com/jae2274/auth-service/auth_service/restapi/service"

	"github.com/gorilla/mux"
)

type UserController struct {
	db *sql.DB
}

func NewUserController(db *sql.DB) *UserController {
	return &UserController{
		db: db,
	}
}

func (c *UserController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/authority", c.FindAllUserAuthorities).Methods("GET")
	router.HandleFunc("/auth/withdrawal", c.Withdrawal).Methods("DELETE")
}

func (c *UserController) FindAllUserAuthorities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, isExisted := middleware.GetClaims(ctx)
	if !isExisted {
		http.Error(w, "no claims in context", http.StatusUnauthorized)
		return
	}

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

	claims, isExisted := middleware.GetClaims(ctx)
	if !isExisted {
		http.Error(w, "no claims in context", http.StatusUnauthorized)
		return
	}

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
