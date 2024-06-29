package service

import (
	"context"
	"database/sql"
	"strconv"
	"testing"
	"time"

	"github.com/jae2274/auth-service/auth_service/common/domain"
	"github.com/jae2274/auth-service/auth_service/common/mysqldb"
	"github.com/jae2274/auth-service/auth_service/models"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
	"github.com/jae2274/auth-service/auth_service/restapi/ooauth"
	"github.com/jae2274/auth-service/auth_service/restapi/service"
	"github.com/jae2274/auth-service/test/tinit"
	"github.com/jae2274/goutils/ptr"
	"github.com/stretchr/testify/require"
)

func TestTicketService(t *testing.T) {
	signUpAdminUser := func(t *testing.T, ctx context.Context, db *sql.DB) *models.User {
		user, err := mysqldb.WithTransaction(ctx, db, func(tx *sql.Tx) (*models.User, error) {
			return service.SignUp(ctx, tx, &ooauth.UserInfo{AuthorizedBy: domain.GOOGLE, AuthorizedID: "authId", Email: "targetUser@test.com", Username: "target"}, nil)
		})
		require.NoError(t, err)

		return user
	}

	t.Run("return false if ticket not existed", func(t *testing.T) {
		db := tinit.DB(t)
		ctx := context.Background()

		_, isExisted, err := service.GetTicketInfo(ctx, db, "notExistedTicketId")
		require.NoError(t, err)
		require.False(t, isExisted)
	})

	t.Run("return error if authority not existed", func(t *testing.T) {
		db := tinit.DB(t)
		ctx := context.Background()
		admin := signUpAdminUser(t, ctx, db)

		_, err := createTicketWithTx(ctx, db, admin.UserID, "ticketName", []*dto.UserAuthorityReq{{AuthorityCode: "notExistedAuthority"}})
		require.Error(t, err)
	})

	t.Run("return ticket info if ticket existed", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)
		admin := signUpAdminUser(t, ctx, db)

		ticketName := "ticketName"
		ticketAuthorities := []*dto.UserAuthorityReq{
			{AuthorityCode: authorities[0].AuthorityCode},
			{AuthorityCode: authorities[1].AuthorityCode, ExpiryDurationMS: ptr.P(int64(2 * time.Hour / time.Millisecond))},
		}
		ticketId, err := createTicketWithTx(ctx, db, admin.UserID, ticketName, ticketAuthorities)
		require.NoError(t, err)
		require.NotEmpty(t, ticketId)

		res, isExisted, err := service.GetTicketInfo(ctx, db, ticketId)
		require.NoError(t, err)
		require.True(t, isExisted)
		require.Equal(t, ticketId, res.TicketId)
		require.False(t, res.IsUsed)
		require.Equal(t, ticketName, res.TicketName)
		require.Len(t, res.TicketAuthorities, len(ticketAuthorities))

		for i, ticketAuthority := range res.TicketAuthorities {
			require.Equal(t, ticketAuthorities[i].AuthorityCode, ticketAuthority.AuthorityCode)
			require.Equal(t, authorities[i].AuthorityName, ticketAuthority.AuthorityName) //UserAuthorityReq에서는 존재하지 않는 필드
			require.Equal(t, authorities[i].Summary, ticketAuthority.Summary)             //UserAuthorityReq에서는 존재하지 않는 필드

			if ticketAuthority.ExpiryDurationMS != nil {
				require.Equal(t, int64(2*time.Hour/time.Millisecond), *ticketAuthority.ExpiryDurationMS)
			} else {
				require.Nil(t, ticketAuthority.ExpiryDurationMS)
			}
		}
	})

	t.Run("can get ticket info by ticket name", func(t *testing.T) {

		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)
		admin := signUpAdminUser(t, ctx, db)

		ticketName := "ticketName"
		ticketAuthorities := []*dto.UserAuthorityReq{
			{AuthorityCode: authorities[0].AuthorityCode},
			{AuthorityCode: authorities[1].AuthorityCode, ExpiryDurationMS: ptr.P(int64(2 * time.Hour / time.Millisecond))},
		}
		ticketId, err := createTicketWithTx(ctx, db, admin.UserID, ticketName, ticketAuthorities)
		require.NoError(t, err)
		require.NotEmpty(t, ticketId)

		res, isExisted, err := service.GetTicketInfo(ctx, db, ticketName)
		require.NoError(t, err)
		require.True(t, isExisted)
		require.Equal(t, ticketId, res.TicketId)
		require.False(t, res.IsUsed)
		require.Equal(t, ticketName, res.TicketName)
		require.Len(t, res.TicketAuthorities, len(ticketAuthorities))

		for i, ticketAuthority := range res.TicketAuthorities {
			require.Equal(t, ticketAuthorities[i].AuthorityCode, ticketAuthority.AuthorityCode)
			require.Equal(t, authorities[i].AuthorityName, ticketAuthority.AuthorityName) //UserAuthorityReq에서는 존재하지 않는 필드
			require.Equal(t, authorities[i].Summary, ticketAuthority.Summary)             //UserAuthorityReq에서는 존재하지 않는 필드

			if ticketAuthority.ExpiryDurationMS != nil {
				require.Equal(t, int64(2*time.Hour/time.Millisecond), *ticketAuthority.ExpiryDurationMS)
			} else {
				require.Nil(t, ticketAuthority.ExpiryDurationMS)
			}
		}
	})

	t.Run("return error if create with same ticket name", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)
		admin := signUpAdminUser(t, ctx, db)

		ticketName := "ticketName"
		ticketAuthorities := []*dto.UserAuthorityReq{
			{AuthorityCode: authorities[0].AuthorityCode},
			{AuthorityCode: authorities[1].AuthorityCode, ExpiryDurationMS: ptr.P(int64(2 * time.Hour / time.Millisecond))},
		}
		ticketId, err := createTicketWithTx(ctx, db, admin.UserID, ticketName, ticketAuthorities)
		require.NoError(t, err)
		require.NotEmpty(t, ticketId)

		_, err = createTicketWithTx(ctx, db, admin.UserID, ticketName, ticketAuthorities)
		require.Error(t, err)
	})

	userinfo := &ooauth.UserInfo{
		AuthorizedBy: domain.GOOGLE,
		AuthorizedID: "123456",
		Email:        "email@email.com",
		Username:     "testUsername",
	}

	t.Run("return error before ticket is created", func(t *testing.T) {
		db := tinit.DB(t)

		ctx, _, _, _ := initAgreementFunc(t, db)

		targetUser, err := signUp(ctx, db, userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		err = useTicket(ctx, db, targetUser.UserID, "notExistedTicketId")
		require.Error(t, err)
	})

	t.Run("return authorities after use ticket that has authorities", func(t *testing.T) {
		db := tinit.DB(t)

		ctx, _, _, authorities := initAgreementFunc(t, db)
		admin := signUpAdminUser(t, ctx, db)

		ticketName := "ticketName"
		userAuthorityReqs := []*dto.UserAuthorityReq{
			{AuthorityCode: authorities[0].AuthorityCode},
			{AuthorityCode: authorities[1].AuthorityCode, ExpiryDurationMS: ptr.P(int64(2 * time.Hour / time.Millisecond))},
		}
		ticketId, err := createTicketWithTx(ctx, db, admin.UserID, ticketName, userAuthorityReqs)
		require.NoError(t, err)

		targetUser, err := signUp(ctx, db, userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		err = useTicket(ctx, db, targetUser.UserID, ticketId)
		require.NoError(t, err)

		userAuthorities, err := service.FindValidUserAuthorities(ctx, db, targetUser.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, len(userAuthorityReqs))

		now := time.Now()
		for i, userAuthority := range userAuthorities {
			requireEqualUserRole(t, targetUser.UserID, now, userAuthorityReqs[i], userAuthority)
		}
	})

	t.Run("return error after ticket is used", func(t *testing.T) {
		db := tinit.DB(t)

		ctx, _, _, authorities := initAgreementFunc(t, db)
		admin := signUpAdminUser(t, ctx, db)

		userAuthorityReqs := []*dto.UserAuthorityReq{
			{AuthorityCode: authorities[0].AuthorityCode},
			{AuthorityCode: authorities[1].AuthorityCode, ExpiryDurationMS: ptr.P(int64(2 * time.Hour / time.Millisecond))},
		}

		ticketName := "ticketName"
		ticketId, err := createTicketWithTx(ctx, db, admin.UserID, ticketName, userAuthorityReqs)
		require.NoError(t, err)

		targetUser, err := signUp(ctx, db, userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		err = useTicket(ctx, db, targetUser.UserID, ticketId)
		require.NoError(t, err)

		err = useTicket(ctx, db, targetUser.UserID, ticketId)
		require.Error(t, err)

		userAuthorities, err := service.FindValidUserAuthorities(ctx, db, targetUser.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, len(userAuthorityReqs))

		now := time.Now()
		for i, userAuthority := range userAuthorities {
			requireEqualUserRole(t, targetUser.UserID, now, userAuthorityReqs[i], userAuthority)
		}
	})

	t.Run("return isUsed true after ticket is used", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)
		admin := signUpAdminUser(t, context.Background(), db)
		userAuthorityReqs := []*dto.UserAuthorityReq{
			{AuthorityCode: authorities[0].AuthorityCode},
			{AuthorityCode: authorities[1].AuthorityCode, ExpiryDurationMS: ptr.P(int64(2 * time.Hour / time.Millisecond))},
		}

		ticketName := "ticketName"
		ticketId, err := createTicketWithTx(ctx, db, admin.UserID, ticketName, userAuthorityReqs)
		require.NoError(t, err)

		targetUser, err := signUp(ctx, db, userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		err = useTicket(ctx, db, targetUser.UserID, ticketId)
		require.NoError(t, err)

		res, isExisted, err := service.GetTicketInfo(ctx, db, ticketId)
		require.NoError(t, err)
		require.True(t, isExisted)
		require.True(t, res.IsUsed)
	})

	t.Run("return empty tickets if no ticket existed", func(t *testing.T) {
		db := tinit.DB(t)
		ctx := context.Background()

		tickets, err := service.GetAllTickets(ctx, db)
		require.NoError(t, err)
		require.Empty(t, tickets)
	})

	t.Run("return all tickets if tickets existed", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)
		admin := signUpAdminUser(t, context.Background(), db)
		userAuthorityReqs := [][]*dto.UserAuthorityReq{
			{{AuthorityCode: authorities[0].AuthorityCode}},
			{{AuthorityCode: authorities[1].AuthorityCode, ExpiryDurationMS: ptr.P(int64(2 * time.Hour / time.Millisecond))}},
		}

		ticketIds := make([]string, 0, len(userAuthorityReqs))
		for i, userAuthorityReq := range userAuthorityReqs {
			ticketId, err := createTicketWithTx(ctx, db, admin.UserID, "ticket"+strconv.Itoa(i), userAuthorityReq)
			require.NoError(t, err)

			ticketIds = append(ticketIds, ticketId)
		}

		user, err := signUp(ctx, db, userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)
		err = useTicket(ctx, db, user.UserID, ticketIds[0])
		require.NoError(t, err)

		tickets, err := service.GetAllTickets(ctx, db)
		require.NoError(t, err)
		require.Len(t, tickets, len(userAuthorityReqs))

		require.Equal(t, "ticket0", tickets[0].TicketName)
		require.True(t, tickets[0].IsUsed)
		require.Equal(t, user.UserID, (*tickets[0].UsedInfo).UsedBy)
		require.Equal(t, userinfo.Username, (*tickets[0].UsedInfo).UsedUserName)
		require.WithinDuration(t, time.Now().UTC(), time.UnixMilli(*tickets[0].UsedInfo.UsedUnixMilli).UTC(), time.Second)
		require.Equal(t, userAuthorityReqs[0][0].AuthorityCode, tickets[0].TicketAuthorities[0].AuthorityCode)
		require.Equal(t, authorities[0].AuthorityName, tickets[0].TicketAuthorities[0].AuthorityName)
		require.Equal(t, authorities[0].Summary, tickets[0].TicketAuthorities[0].Summary)
		require.Nil(t, tickets[0].TicketAuthorities[0].ExpiryDurationMS)
		require.Equal(t, tickets[0].CreatedBy, admin.UserID)

		require.Equal(t, "ticket1", tickets[1].TicketName)
		require.False(t, tickets[1].IsUsed)
		require.Nil(t, tickets[1].UsedInfo)
		require.Equal(t, userAuthorityReqs[1][0].AuthorityCode, tickets[1].TicketAuthorities[0].AuthorityCode)
		require.Equal(t, authorities[1].AuthorityName, tickets[1].TicketAuthorities[0].AuthorityName)
		require.Equal(t, authorities[1].Summary, tickets[1].TicketAuthorities[0].Summary)
		require.Equal(t, int64(2*time.Hour/time.Millisecond), *tickets[1].TicketAuthorities[0].ExpiryDurationMS)
	})
}
