package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jae2274/auth-service/auth_service/models"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
	"github.com/jae2274/goutils/ptr"
	"github.com/jae2274/goutils/terr"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func GetTicketInfo(ctx context.Context, exec boil.ContextExecutor, ticketId string) (*dto.Ticket, bool, error) {
	ticket, err := models.Tickets(models.TicketWhere.UUID.EQ(ticketId)).One(ctx, exec)

	if err != nil && err != sql.ErrNoRows {
		return nil, false, terr.Wrap(err)
	} else if err == sql.ErrNoRows {
		return nil, false, nil
	}

	mTicketAuthorities, err := models.TicketAuthorities(models.TicketAuthorityWhere.TicketID.EQ(ticket.TicketID), qm.Load(models.TicketAuthorityRels.Authority)).All(ctx, exec)
	if err != nil {
		return nil, false, terr.Wrap(err)
	}

	ticketAuthorities := make([]*dto.TicketAuthority, 0, len(mTicketAuthorities))
	for _, mTicketAuthority := range mTicketAuthorities {
		authority := mTicketAuthority.R.Authority
		var expiryDurationMS *int64
		if mTicketAuthority.ExpiryDurationMS.Valid {
			expiryDurationMS = ptr.P(mTicketAuthority.ExpiryDurationMS.Int64)
		}
		ticketAuthorities = append(ticketAuthorities, &dto.TicketAuthority{
			AuthorityCode:    authority.AuthorityCode,
			AuthorityName:    authority.AuthorityName,
			Summary:          authority.Summary,
			ExpiryDurationMS: expiryDurationMS,
		})
	}

	return &dto.Ticket{
		TicketId:          ticket.UUID,
		TicketAuthorities: ticketAuthorities,
	}, true, nil
}

func CreateTicket(ctx context.Context, tx *sql.Tx, authorities []*dto.UserAuthorityReq) (string, error) {
	err := attachAuthorityIds(ctx, tx, authorities)
	if err != nil {
		return "", err
	}

	ticket := &models.Ticket{UUID: uuid.New().String()}

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

func UseTicket(ctx context.Context, tx *sql.Tx, userId int, ticketId string) (bool, error) {
	ticket, err := models.Tickets(models.TicketWhere.UUID.EQ(ticketId), models.TicketWhere.UsedBy.IsNull(), qm.Load(models.TicketRels.TicketAuthorities)).One(ctx, tx)
	if err != nil && err != sql.ErrNoRows {
		return false, terr.Wrap(err)
	} else if err == sql.ErrNoRows {
		return false, nil
	}

	ticket.UsedBy = null.IntFrom(userId)
	ticket.Update(ctx, tx, boil.Infer())

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
