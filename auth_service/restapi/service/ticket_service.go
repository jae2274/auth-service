package service

import (
	"context"
	"database/sql"

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
	count, err := models.TicketSubs(models.TicketSubWhere.TicketID.EQ(ticket.TicketID)).Count(ctx, exec)
	if err != nil {
		return nil, false, terr.Wrap(err)
	}

	return convertToDtoTicket(ticket, count), true, nil
}

// func convert
func convertToDtoTicket(ticket *models.Ticket, usedCount int64) *dto.Ticket {
	ticketAuthorities := make([]*dto.TicketAuthority, 0, len(ticket.R.TicketAuthorities))
	for _, mTicketAuthority := range ticket.R.TicketAuthorities {
		ticketAuthorities = append(ticketAuthorities, convertToDtoTicketAuthority(mTicketAuthority))
	}

	return &dto.Ticket{
		TicketId:          ticket.UUID,
		TicketName:        ticket.TicketName,
		TicketAuthorities: ticketAuthorities,
		CreateUnixMilli:   ticket.CreateDate.UnixMilli(),
		UseableCount:      ticket.UseableCount,
		UsedCount:         int(usedCount),
	}
}

func convertToDtoTicketAuthority(mTicketAuthority *models.TicketAuthority) *dto.TicketAuthority {
	authority := mTicketAuthority.R.Authority
	var expiryDurationMS *int64
	if mTicketAuthority.ExpiryDurationMS.Valid {
		expiryDurationMS = ptr.P(mTicketAuthority.ExpiryDurationMS.Int64)
	}
	return &dto.TicketAuthority{
		AuthorityId:      authority.AuthorityID,
		AuthorityCode:    authority.AuthorityCode,
		AuthorityName:    authority.AuthorityName,
		Summary:          authority.Summary,
		ExpiryDurationMS: expiryDurationMS,
	}

}

func getTicket(ctx context.Context, exec boil.ContextExecutor, ticketId string) (*models.Ticket, bool, error) {
	ticket, err := models.Tickets(
		models.TicketWhere.UUID.EQ(ticketId), qm.Or2(models.TicketWhere.TicketName.EQ(ticketId)),
		qm.Load(models.TicketRels.TicketAuthorities+"."+models.TicketAuthorityRels.Authority)).One(ctx, exec)

	if err != nil && err != sql.ErrNoRows {
		return nil, false, terr.Wrap(err)
	} else if err == sql.ErrNoRows {
		return nil, false, nil
	}

	return ticket, true, nil
}

func CreateTicket(ctx context.Context, tx *sql.Tx, createdByUser int, ticketName string, authorities []*dto.UserAuthorityReq, useableCount int) (string, error) {
	err := attachAuthorityIds(ctx, tx, authorities)
	if err != nil {
		return "", err
	}

	ticket := &models.Ticket{UUID: uuid.New().String(), TicketName: ticketName, CreatedBy: createdByUser, UseableCount: useableCount}

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
	ticket, err := models.Tickets(models.TicketWhere.UUID.EQ(ticketId), qm.Load(models.TicketRels.TicketAuthorities), qm.For("update")).One(ctx, tx)
	if err != nil {
		return terr.Wrap(err)
	}

	usedCount, err := models.TicketSubs(models.TicketSubWhere.TicketID.EQ(ticket.TicketID)).Count(ctx, tx)
	if err != nil {
		return terr.Wrap(err)
	}

	if usedCount >= int64(ticket.UseableCount) {
		return terr.New("no more useable ticket")
	}

	ticketSub := &models.TicketSub{
		TicketID: ticket.TicketID,
		UsedBy:   userId,
	}
	err = ticketSub.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return terr.Wrap(err)
	}

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

func GetAllTickets(ctx context.Context, exec boil.ContextExecutor) ([]*dto.TicketDetail, error) {
	tickets, err := models.Tickets(
		qm.Load(models.TicketRels.TicketAuthorities+"."+models.TicketAuthorityRels.Authority),
		qm.Load(models.TicketRels.TicketSubs+"."+models.TicketSubRels.UsedByUser),
	).All(ctx, exec)
	if err != nil {
		return nil, terr.Wrap(err)
	}

	dtoTickets := make([]*dto.TicketDetail, 0, len(tickets))
	for _, ticket := range tickets {
		usedCount := len(ticket.R.TicketSubs)

		dtoTickets = append(dtoTickets, &dto.TicketDetail{
			Ticket:    *convertToDtoTicket(ticket, int64(usedCount)),
			UsedInfos: convertUsedInfos(ticket.R.TicketSubs),
			CreatedBy: ticket.CreatedBy,
		})
	}

	return dtoTickets, nil
}

func convertUsedInfos(ticketSub []*models.TicketSub) []*dto.UsedInfo {
	usedInfos := make([]*dto.UsedInfo, 0, len(ticketSub))
	for _, ticketSub := range ticketSub {
		usedInfos = append(usedInfos, &dto.UsedInfo{
			UsedBy:        ticketSub.UsedBy,
			UsedUserName:  ticketSub.R.UsedByUser.Name,
			UsedUnixMilli: ticketSub.UsedDate.UnixMilli(),
		})
	}

	return usedInfos
}
