package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"userService/usersvc/common/domain"
	"userService/usersvc/models"

	"github.com/google/uuid"
	"github.com/jae2274/goutils/terr"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type AdminService struct {
	mysqlDB *sql.DB
}

func NewAdminService(mysqlDB *sql.DB) *AdminService {
	return &AdminService{
		mysqlDB: mysqlDB,
	}
}

var ErrEmptyRole = fmt.Errorf("empty role")

func (a *AdminService) CreateRoleTicket(ctx context.Context, ticketRoles []*models.TicketRole) (string, error) {
	if len(ticketRoles) == 0 {
		return "", terr.Wrap(ErrEmptyRole)
	}
	tx, err := a.mysqlDB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}

	ticketId, err := createRoleTicket(ctx, tx, ticketRoles)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = errors.Join(rollbackErr)
		}

		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return ticketId, nil
}

func createRoleTicket(ctx context.Context, tx *sql.Tx, ticketRoles []*models.TicketRole) (string, error) {

	ticketUuid := uuid.New().String()
	ticket := models.Ticket{
		UUID: ticketUuid,
	}

	if err := ticket.Insert(ctx, tx, boil.Infer()); err != nil {
		return "", err
	}

	//TODO: bulk insert
	for _, role := range ticketRoles {
		err := ticket.SetTicketRole(ctx, tx, true, role)
		if err != nil {
			return "", err
		}
	}

	return ticketUuid, nil
}

func (a *AdminService) UseTicket(ctx context.Context, userId int, ticketUuid string) error {

	tx, err := a.mysqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = useTicket(ctx, tx, userId, ticketUuid)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = errors.Join(rollbackErr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func useTicket(ctx context.Context, tx *sql.Tx, userId int, ticketUuid string) error {

	ticket, err := models.Tickets(models.TicketWhere.UUID.EQ(ticketUuid), models.TicketWhere.UsedBy.IsNull()).One(ctx, tx)
	if err != nil {
		return err
	}

	ticket.UsedBy = null.NewInt(userId, true)
	_, err = ticket.Update(ctx, tx, boil.Infer())
	if err != nil {
		return err
	}

	ticketRoles, err := models.TicketRoles(models.TicketRoleWhere.TicketID.EQ(ticket.TicketID)).All(ctx, tx)
	if err != nil {
		return err
	}

	for _, role := range ticketRoles {
		var expiryDate null.Time

		if role.ExpiryTerm.Valid {
			now := time.Now()
			expiryDuration := role.ExpiryTerm.Time.Sub(now)
			expiryDate = null.NewTime(now.Add(expiryDuration), true)
		}

		userRole := models.UserRole{
			UserID:      userId,
			RoleName:    role.RoleName,
			GrantedType: string(domain.TICKET),
			GrantedBy:   ticket.TicketID,
			ExpiryDate:  expiryDate,
		}
		err := userRole.Insert(ctx, tx, boil.Infer()) //TODO: bulk insert
		if err != nil {
			return err
		}
	}

	return nil
}
