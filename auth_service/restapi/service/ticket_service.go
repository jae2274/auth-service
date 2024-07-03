package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
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

	count, err := models.TicketUseds(models.TicketUsedWhere.TicketID.EQ(ticket.TicketID)).Count(ctx, exec)
	if err != nil {
		return nil, false, terr.Wrap(err)
	}

	return convertToDtoTicket(ticket, count), true, nil
}

func CheckUseTicket(ctx context.Context, exec boil.ContextExecutor, userId, ticketId int) (bool, error) {
	return models.TicketUseds(models.TicketUsedWhere.TicketID.EQ(ticketId), models.TicketUsedWhere.UsedBy.EQ(userId)).Exists(ctx, exec)
}

// func convert
func convertToDtoTicket(ticket *models.Ticket, usedCount int64) *dto.Ticket {
	ticketAuthorities := make([]*dto.TicketAuthority, 0, len(ticket.R.TicketAuthorities))
	for _, mTicketAuthority := range ticket.R.TicketAuthorities {
		ticketAuthorities = append(ticketAuthorities, convertToDtoTicketAuthority(mTicketAuthority))
	}

	return &dto.Ticket{
		TicketIndexId:     ticket.TicketID,
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

func CreateTicket(ctx context.Context, tx *sql.Tx, createdByUser int, ticketName string, authorities []*dto.UserAuthorityReq, useableCount int) (*models.Ticket, error) {
	if useableCount <= 0 {
		return nil, terr.New("useableCount must be greater than 0")
	}
	err := attachAuthorityIds(ctx, tx, authorities)
	if err != nil {
		return nil, err
	}

	ticket := &models.Ticket{UUID: uuid.New().String(), TicketName: ticketName, CreatedBy: createdByUser, UseableCount: useableCount}

	if err := ticket.Insert(ctx, tx, boil.Infer()); err != nil {
		return nil, terr.Wrap(err)
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
		return nil, terr.Wrap(err)
	}

	return ticket, nil
}

var ErrTicketNotFound = errors.New("ticket not found")
var ErrNoMoreUseableTicket = errors.New("no more useable ticket")
var ErrAlreadyUsedTicket = errors.New("already used ticket")

func UseTicket(ctx context.Context, tx *sql.Tx, userId int, ticketId string) (*dto.Ticket, error) {
	ticket, err := models.Tickets(models.TicketWhere.UUID.EQ(ticketId),
		qm.Load(models.TicketRels.TicketAuthorities+"."+models.TicketAuthorityRels.Authority),
		qm.For("update")).One(ctx, tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, terr.Wrap(err)
	} else if err == sql.ErrNoRows {
		return nil, ErrTicketNotFound
	}

	usedCount, err := models.TicketUseds(models.TicketUsedWhere.TicketID.EQ(ticket.TicketID)).Count(ctx, tx)
	if err != nil {
		return nil, terr.Wrap(err)
	}

	if usedCount >= int64(ticket.UseableCount) {
		return nil, ErrNoMoreUseableTicket
	}

	ticketSub := &models.TicketUsed{
		TicketID: ticket.TicketID,
		UsedBy:   userId,
	}
	err = ticketSub.Insert(ctx, tx, boil.Infer())
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 { // duplicate key error
			return nil, ErrAlreadyUsedTicket
		}

		return nil, terr.Wrap(err)
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
		return nil, terr.Wrap(err)
	}

	return convertToDtoTicket(ticket, usedCount), nil
}

func GetAllTickets(ctx context.Context, exec boil.ContextExecutor) ([]*dto.TicketDetail, error) {
	tickets, err := models.Tickets(
		qm.Load(models.TicketRels.TicketAuthorities+"."+models.TicketAuthorityRels.Authority),
		qm.Load(models.TicketRels.TicketUseds+"."+models.TicketUsedRels.UsedByUser),
	).All(ctx, exec)
	if err != nil {
		return nil, terr.Wrap(err)
	}

	dtoTickets := make([]*dto.TicketDetail, 0, len(tickets))
	for _, ticket := range tickets {
		usedCount := len(ticket.R.TicketUseds)

		dtoTickets = append(dtoTickets, &dto.TicketDetail{
			Ticket:    *convertToDtoTicket(ticket, int64(usedCount)),
			UsedInfos: convertUsedInfos(ticket.R.TicketUseds),
			CreatedBy: ticket.CreatedBy,
		})
	}

	return dtoTickets, nil
}

func convertUsedInfos(ticketUseds []*models.TicketUsed) []*dto.UsedInfo {
	usedInfos := make([]*dto.UsedInfo, 0, len(ticketUseds))
	for _, ticketUsed := range ticketUseds {
		usedInfos = append(usedInfos, &dto.UsedInfo{
			UsedBy:        ticketUsed.UsedBy,
			UsedUserName:  ticketUsed.R.UsedByUser.Name,
			UsedUnixMilli: ticketUsed.UsedDate.UnixMilli(),
		})
	}

	return usedInfos
}
