package ctrlr

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jae2274/auth-service/auth_service/common/mysqldb"
	"github.com/jae2274/auth-service/auth_service/models"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
	"github.com/jae2274/auth-service/auth_service/restapi/middleware"
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
	router.HandleFunc("/auth/admin/authority", a.GetAllAuthorities).Methods("GET")
	router.HandleFunc("/auth/admin/authority", a.AddAuthority).Methods("POST")
	router.HandleFunc("/auth/admin/authority", a.RemoveAuthority).Methods("DELETE")
	router.HandleFunc("/auth/admin/ticket", a.GetAllTickets).Methods("GET")
	router.HandleFunc("/auth/admin/ticket", a.CreateTicket).Methods("POST")
	router.HandleFunc("/auth/admin/user", a.GetAllUsers).Methods("GET")
}

func (a *AdminController) GetAllAuthorities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authorities, err := service.GetAllAuthorities(ctx, a.db)
	if errorHandler(ctx, w, err) {
		return
	}

	authoritiesRes := make([]*dto.AuthorityRes, 0, len(authorities))
	for _, authority := range authorities {
		authoritiesRes = append(authoritiesRes, &dto.AuthorityRes{
			AuthorityCode: authority.AuthorityCode,
			AuthorityName: authority.AuthorityName,
			Summary:       authority.Summary,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(&dto.GetAllAuthoritiesResponse{Authorities: authoritiesRes})
	if errorHandler(ctx, w, err) {
		return
	}
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

	if err == service.ErrCannotControlAuthorityAdmin {
		http.Error(w, "cannot control authority admin", http.StatusBadRequest)
	}

	if errorHandler(ctx, w, err) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getAllTicket(ctx context.Context, db *sql.DB, byMe bool) ([]*dto.TicketDetail, error) {
	if byMe {
		claims := GetClaimsOrPatal(ctx)

		adminUserId, err := strconv.Atoi(claims.UserId)
		if err != nil {
			return nil, err
		}

		return service.GetAllTicketsByUserId(ctx, db, adminUserId)
	}

	return service.GetAllTickets(ctx, db)
}

func (a *AdminController) GetAllTickets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	byMeStr := r.URL.Query().Get("by_me")

	var tickets []*dto.TicketDetail
	var err error
	if byMeStr == "true" {
		tickets, err = getAllTicket(ctx, a.db, true)
	} else {
		tickets, err = getAllTicket(ctx, a.db, false)
	}

	if errorHandler(ctx, w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(&dto.GetAllTicketsResponse{Tickets: tickets})
	if errorHandler(ctx, w, err) {
		return
	}
}

func (a *AdminController) CreateTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, isExisted := middleware.GetClaims(ctx)
	if !isExisted {
		errorHandler(ctx, w, terr.New("no claims"))
		return
	}
	adminUserId, err := strconv.Atoi(claims.UserId)
	if errorHandler(ctx, w, err) {
		return
	}

	var req *dto.CreateTicketRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if errorHandler(ctx, w, err) {
		return
	}

	createdTicket, err := mysqldb.WithTransaction(ctx, a.db, func(tx *sql.Tx) (*models.Ticket, error) {
		return service.CreateTicket(ctx, tx, adminUserId, req.TicketName, req.TicketAuthorities, req.UseableCount)
	})
	if errorHandler(ctx, w, err) {
		return
	}

	ticket, isExisted, err := service.GetTicketInfo(ctx, a.db, createdTicket.UUID)
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

func (a *AdminController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := service.GetAllUsers(ctx, a.db)
	if errorHandler(ctx, w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(&dto.GetAllUsersResponse{Users: users})
	if errorHandler(ctx, w, err) {
		return
	}
}
