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
	ticket, isExisted, err := getTicket(ctx, exec, ticketId)
	if err != nil {
		return nil, false, terr.Wrap(err)
	} else if !isExisted {
		return nil, false, nil
	}

	// mTicketAuthorities, err := models.TicketAuthorities(models.TicketAuthorityWhere.TicketID.EQ(ticket.TicketID), qm.Load(models.TicketAuthorityRels.Authority)).All(ctx, exec)
	// if err != nil {
	// 	return nil, false, terr.Wrap(err)
	// }

	return convertToDtoTicket(ticket), true, nil
}

// func convert
func convertToDtoTicket(ticket *models.Ticket) *dto.Ticket {
	ticketAuthorities := make([]*dto.TicketAuthority, 0, len(ticket.R.TicketAuthorities))
	for _, mTicketAuthority := range ticket.R.TicketAuthorities {
		ticketAuthorities = append(ticketAuthorities, convertToDtoTicketAuthority(mTicketAuthority))
	}

	return &dto.Ticket{
		TicketId:          ticket.UUID,
		IsUsed:            ticket.UsedBy.Valid,
		TicketAuthorities: ticketAuthorities,
	}
}

func convertToDtoTicketAuthority(mTicketAuthority *models.TicketAuthority) *dto.TicketAuthority {
	authority := mTicketAuthority.R.Authority
	var expiryDurationMS *int64
	var expiryDuration *dto.Duration
	if mTicketAuthority.ExpiryDurationMS.Valid {
		expiryDurationMS = ptr.P(mTicketAuthority.ExpiryDurationMS.Int64)
		expiryDuration = ptr.P(dto.Duration(time.Duration(mTicketAuthority.ExpiryDurationMS.Int64) * time.Millisecond))
	}
	return &dto.TicketAuthority{
		AuthorityId:      authority.AuthorityID,
		AuthorityCode:    authority.AuthorityCode,
		AuthorityName:    authority.AuthorityName,
		Summary:          authority.Summary,
		ExpiryDuration:   expiryDuration,
		ExpiryDurationMS: expiryDurationMS,
	}

}

func getTicket(ctx context.Context, exec boil.ContextExecutor, ticketId string) (*models.Ticket, bool, error) {
	ticket, err := models.Tickets(models.TicketWhere.UUID.EQ(ticketId), qm.Load(models.TicketRels.TicketAuthorities+"."+models.TicketAuthorityRels.Authority)).One(ctx, exec)

	if err != nil && err != sql.ErrNoRows {
		return nil, false, terr.Wrap(err)
	} else if err == sql.ErrNoRows {
		return nil, false, nil
	}

	return ticket, true, nil
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
		if authority.ExpiryDurationMS != nil {
			expiryDurationMS = null.NewInt64(*authority.ExpiryDurationMS, true)
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

func UseTicket(ctx context.Context, tx *sql.Tx, userId int, ticketId string) error {
	ticket, err := models.Tickets(models.TicketWhere.UUID.EQ(ticketId), models.TicketWhere.UsedBy.IsNull(), qm.Load(models.TicketRels.TicketAuthorities)).One(ctx, tx)
	if err != nil {
		return terr.Wrap(err)
	}

	ticket.UsedBy = null.IntFrom(userId)
	ticket.Update(ctx, tx, boil.Infer())

	dUserAuthorities := make([]*dto.UserAuthorityReq, 0, len(ticket.R.TicketAuthorities))
	for _, ticketAuthority := range ticket.R.TicketAuthorities {
		var expiryDurationMS *int64
		if ticketAuthority.ExpiryDurationMS.Valid {
			expiryDurationMS = ptr.P(ticketAuthority.ExpiryDurationMS.Int64)
		}

		dUserAuthorities = append(dUserAuthorities, &dto.UserAuthorityReq{
			AuthorityID:      ticketAuthority.AuthorityID,
			ExpiryDurationMS: expiryDurationMS,
		})
	}

	err = addUserAuthorities(ctx, tx, userId, dUserAuthorities)

	if err != nil {
		return terr.Wrap(err)
	}

	return nil
}

func GetAllTickets(ctx context.Context, exec boil.ContextExecutor) ([]*dto.Ticket, error) {
	tickets, err := models.Tickets(qm.Load(models.TicketRels.TicketAuthorities+"."+models.TicketAuthorityRels.Authority)).All(ctx, exec)
	if err != nil {
		return nil, terr.Wrap(err)
	}

	dtoTickets := make([]*dto.Ticket, 0, len(tickets))
	for _, ticket := range tickets {
		dtoTickets = append(dtoTickets, convertToDtoTicket(ticket))
	}

	return dtoTickets, nil
}
