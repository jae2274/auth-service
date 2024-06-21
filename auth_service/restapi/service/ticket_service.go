package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jae2274/auth-service/auth_service/common/mysqldb"
	"github.com/jae2274/auth-service/auth_service/models"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
	"github.com/jae2274/goutils/ptr"
	"github.com/jae2274/goutils/terr"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TicketService struct {
	mysqlDB *sql.DB
}

func NewTicketService(mysqlDB *sql.DB) TicketService {
	return TicketService{
		mysqlDB: mysqlDB,
	}
}

func (t *TicketService) CreateTicket(ctx context.Context, authorities []*dto.UserAuthorityReq) (string, error) {
	err := attachAuthorityIds(ctx, t.mysqlDB, authorities)
	if err != nil {
		return "", err
	}

	tx, err := t.mysqlDB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}

	ticketId, err := t.createTicket(ctx, tx, authorities)

	return mysqldb.CommitOrRollback(tx, ticketId, err)
}

func (t *TicketService) createTicket(ctx context.Context, tx *sql.Tx, authorities []*dto.UserAuthorityReq) (string, error) {
	ticket := &models.Ticket{
		UUID: uuid.New().String(),
	}

	if err := ticket.Insert(ctx, tx, boil.Infer()); err != nil {
		return "", terr.Wrap(err)
	}

	ticketAuthorities := make([]*models.TicketAuthority, len(authorities))
	for i, authority := range authorities {
		expiryDurationMS := null.NewInt64(0, false)
		if authority.ExpiryDuration != nil {
			expiryDurationMS = null.NewInt64(int64(time.Duration((*authority.ExpiryDuration))/time.Millisecond), true)
		}
		ticketAuthorities[i] = &models.TicketAuthority{
			TicketID:         ticket.TicketID,
			AuthorityID:      authority.AuthorityID,
			ExpiryDurationMS: expiryDurationMS,
		}
	}

	if err := ticket.AddTicketAuthorities(ctx, tx, true, ticketAuthorities...); err != nil {
		return "", terr.Wrap(err)
	}

	return ticket.UUID, nil
}

func (t *TicketService) UseTicket(ctx context.Context, userId int, ticketId string) (bool, error) {
	tx, err := t.mysqlDB.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}

	isExisted, err := useTicket(ctx, tx, userId, ticketId)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return isExisted, errors.Join(err, rollbackErr)
		}

		return isExisted, err
	}

	if err := tx.Commit(); err != nil {
		return isExisted, err
	}

	return isExisted, nil
}

func useTicket(ctx context.Context, tx *sql.Tx, userId int, ticketId string) (bool, error) {
	ticket, err := models.Tickets(models.TicketWhere.UUID.EQ(ticketId), qm.Load(models.TicketRels.TicketAuthorities)).One(ctx, tx)
	if err != nil && err != sql.ErrNoRows {
		return false, terr.Wrap(err)
	} else if err == sql.ErrNoRows {
		return false, nil
	}

	dUserAuthorities := make([]*dto.UserAuthorityReq, 0, len(ticket.R.TicketAuthorities))
	for _, ticketAuthority := range ticket.R.TicketAuthorities {
		var expiryDuration *dto.Duration
		if ticketAuthority.ExpiryDurationMS.Valid {
			expiryDuration = ptr.P(dto.Duration(time.Duration(ticketAuthority.ExpiryDurationMS.Int64) * time.Millisecond))
		}

		dUserAuthorities = append(dUserAuthorities, &dto.UserAuthorityReq{
			AuthorityID:    ticketAuthority.AuthorityID,
			ExpiryDuration: expiryDuration,
		})
	}

	err = addUserAuthorities(ctx, tx, userId, dUserAuthorities)

	if err != nil {
		return false, terr.Wrap(err)
	}

	return true, nil
}
