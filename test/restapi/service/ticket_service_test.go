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

		_, err := createTicketWithTx(ctx, db, admin.UserID, "ticketName", []*dto.UserAuthorityReq{{AuthorityCode: "notExistedAuthority"}}, 1)
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
		ticket, err := createTicketWithTx(ctx, db, admin.UserID, ticketName, ticketAuthorities, 1)
		require.NoError(t, err)

		res, isExisted, err := service.GetTicketInfo(ctx, db, ticket.UUID)
		require.NoError(t, err)
		require.True(t, isExisted)
		require.Equal(t, ticket.UUID, res.TicketId)
		require.Equal(t, res.UseableCount, 1)
		require.Equal(t, res.UsedCount, 0)
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
		ticket, err := createTicketWithTx(ctx, db, admin.UserID, ticketName, ticketAuthorities, 1)
		require.NoError(t, err)

		res, isExisted, err := service.GetTicketInfo(ctx, db, ticketName)
		require.NoError(t, err)
		require.True(t, isExisted)
		require.Equal(t, ticket.UUID, res.TicketId)
		require.Equal(t, res.UseableCount, 1)
		require.Equal(t, res.UsedCount, 0)
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
		ticketId, err := createTicketWithTx(ctx, db, admin.UserID, ticketName, ticketAuthorities, 1)
		require.NoError(t, err)
		require.NotEmpty(t, ticketId)

		_, err = createTicketWithTx(ctx, db, admin.UserID, ticketName, ticketAuthorities, 1)
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

		_, err = useTicket(ctx, db, targetUser.UserID, "notExistedTicketId")
		require.ErrorIs(t, err, service.ErrTicketNotFound)
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
		ticket, err := createTicketWithTx(ctx, db, admin.UserID, ticketName, userAuthorityReqs, 1)
		require.NoError(t, err)

		targetUser, err := signUp(ctx, db, userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		mTicket, err := useTicket(ctx, db, targetUser.UserID, ticket.UUID)
		require.NoError(t, err)
		for i, authority := range mTicket.TicketAuthorities {
			require.Equal(t, userAuthorityReqs[i].AuthorityCode, authority.AuthorityCode)
		}

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
		ticket, err := createTicketWithTx(ctx, db, admin.UserID, ticketName, userAuthorityReqs, 1)
		require.NoError(t, err)

		targetUser, err := signUp(ctx, db, userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		_, err = useTicket(ctx, db, targetUser.UserID, ticket.UUID)
		require.NoError(t, err)

		_, err = useTicket(ctx, db, targetUser.UserID, ticket.UUID)
		require.Error(t, err)

		userAuthorities, err := service.FindValidUserAuthorities(ctx, db, targetUser.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, len(userAuthorityReqs))

		now := time.Now()
		for i, userAuthority := range userAuthorities {
			requireEqualUserRole(t, targetUser.UserID, now, userAuthorityReqs[i], userAuthority)
		}
	})

	anotherUser1 := &ooauth.UserInfo{
		AuthorizedBy: domain.GOOGLE,
		AuthorizedID: "anotherUser1",
		Email:        "anotherUser1@email.com",
		Username:     "anotherUser1",
	}

	anotherUser2 := &ooauth.UserInfo{
		AuthorizedBy: domain.GOOGLE,
		AuthorizedID: "anotherUser2",
		Email:        "anotherUser2@email.com",
		Username:     "anotherUser2",
	}

	t.Run("can use ticket many times as useableCount", func(t *testing.T) {
		db := tinit.DB(t)

		ctx, _, _, authorities := initAgreementFunc(t, db)
		admin := signUpAdminUser(t, ctx, db)

		userAuthorityReqs := []*dto.UserAuthorityReq{
			{AuthorityCode: authorities[0].AuthorityCode},
			{AuthorityCode: authorities[1].AuthorityCode, ExpiryDurationMS: ptr.P(int64(2 * time.Hour / time.Millisecond))},
		}

		ticketName := "ticketName"
		ticket, err := createTicketWithTx(ctx, db, admin.UserID, ticketName, userAuthorityReqs, 2) // 사용가능 횟수 2회
		require.NoError(t, err)

		users := make([]*models.User, 0, 3)
		for _, targetUser := range []*ooauth.UserInfo{userinfo, anotherUser1, anotherUser2} {
			user, err := signUp(ctx, db, targetUser, []*dto.UserAgreementReq{})
			require.NoError(t, err)
			users = append(users, user)
		}

		for _, targetUser := range users[:2] {
			now := time.Now()
			_, err := useTicket(ctx, db, targetUser.UserID, ticket.UUID)
			require.NoError(t, err)

			userAuthorities, err := service.FindValidUserAuthorities(ctx, db, targetUser.UserID)
			require.NoError(t, err)
			require.Len(t, userAuthorities, len(userAuthorityReqs))

			for i, userAuthority := range userAuthorities {
				requireEqualUserRole(t, targetUser.UserID, now, userAuthorityReqs[i], userAuthority)
			}
		}

		_, err = useTicket(ctx, db, users[2].UserID, ticket.UUID) //사용가능 횟수 초과
		require.ErrorIs(t, err, service.ErrNoMoreUseableTicket)

		userAuthorities, err := service.FindValidUserAuthorities(ctx, db, users[2].UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, 0)

		ticketInfo, isExisted, err := service.GetTicketInfo(ctx, db, ticket.UUID)
		require.NoError(t, err)
		require.True(t, isExisted)
		require.Equal(t, ticketInfo.UseableCount, 2)
		require.Equal(t, ticketInfo.UsedCount, 2)
	})

	t.Run("same ticket cannot be used by same user", func(t *testing.T) {
		db := tinit.DB(t)

		ctx, _, _, authorities := initAgreementFunc(t, db)
		admin := signUpAdminUser(t, ctx, db)

		userAuthorityReqs := []*dto.UserAuthorityReq{
			{AuthorityCode: authorities[0].AuthorityCode},
			{AuthorityCode: authorities[1].AuthorityCode, ExpiryDurationMS: ptr.P(int64(2 * time.Hour / time.Millisecond))},
		}

		ticketName := "ticketName"
		ticket, err := createTicketWithTx(ctx, db, admin.UserID, ticketName, userAuthorityReqs, 2)
		require.NoError(t, err)

		user, err := signUp(ctx, db, userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		now := time.Now()
		_, err = useTicket(ctx, db, user.UserID, ticket.UUID)
		require.NoError(t, err)

		userAuthorities, err := service.FindValidUserAuthorities(ctx, db, user.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, len(userAuthorityReqs))

		for i, userAuthority := range userAuthorities {
			requireEqualUserRole(t, user.UserID, now, userAuthorityReqs[i], userAuthority)
		}

		_, err = useTicket(ctx, db, user.UserID, ticket.UUID) // 이미 사용한 티켓
		require.ErrorIs(t, err, service.ErrAlreadyUsedTicket)

		userAuthorities, err = service.FindValidUserAuthorities(ctx, db, user.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, len(userAuthorityReqs))

		for i, userAuthority := range userAuthorities {
			requireEqualUserRole(t, user.UserID, now, userAuthorityReqs[i], userAuthority)
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
		ticket, err := createTicketWithTx(ctx, db, admin.UserID, ticketName, userAuthorityReqs, 1)
		require.NoError(t, err)

		targetUser, err := signUp(ctx, db, userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		_, err = useTicket(ctx, db, targetUser.UserID, ticket.UUID)
		require.NoError(t, err)

		res, isExisted, err := service.GetTicketInfo(ctx, db, ticket.UUID)
		require.NoError(t, err)
		require.True(t, isExisted)

		require.Equal(t, res.UseableCount, 1)
		require.Equal(t, res.UsedCount, 1)
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
			ticket, err := createTicketWithTx(ctx, db, admin.UserID, "ticket"+strconv.Itoa(i), userAuthorityReq, 1)
			require.NoError(t, err)

			ticketIds = append(ticketIds, ticket.UUID)
		}

		user, err := signUp(ctx, db, userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)
		_, err = useTicket(ctx, db, user.UserID, ticketIds[0])
		require.NoError(t, err)

		tickets, err := service.GetAllTickets(ctx, db)
		require.NoError(t, err)
		require.Len(t, tickets, len(userAuthorityReqs))

		ticket0 := tickets[0]
		require.Equal(t, "ticket0", ticket0.TicketName)
		require.Equal(t, ticket0.UseableCount, 1)
		require.Equal(t, ticket0.UsedCount, 1)
		require.Len(t, ticket0.UsedInfos, 1)
		require.Equal(t, user.UserID, ticket0.UsedInfos[0].UsedBy)
		require.Equal(t, userinfo.Username, ticket0.UsedInfos[0].UsedUserName)
		require.WithinDuration(t, time.Now().UTC(), time.UnixMilli(ticket0.UsedInfos[0].UsedUnixMilli).UTC(), time.Second)
		require.Equal(t, userAuthorityReqs[0][0].AuthorityCode, ticket0.TicketAuthorities[0].AuthorityCode)
		require.Equal(t, authorities[0].AuthorityName, ticket0.TicketAuthorities[0].AuthorityName)
		require.Equal(t, authorities[0].Summary, ticket0.TicketAuthorities[0].Summary)
		require.Nil(t, ticket0.TicketAuthorities[0].ExpiryDurationMS)
		require.Equal(t, ticket0.CreatedBy, admin.UserID)

		ticket1 := tickets[1]
		require.Equal(t, "ticket1", ticket1.TicketName)
		require.Equal(t, ticket1.UseableCount, 1)
		require.Equal(t, ticket1.UsedCount, 0)
		require.Len(t, ticket1.UsedInfos, 0)
		require.Equal(t, userAuthorityReqs[1][0].AuthorityCode, ticket1.TicketAuthorities[0].AuthorityCode)
		require.Equal(t, authorities[1].AuthorityName, ticket1.TicketAuthorities[0].AuthorityName)
		require.Equal(t, authorities[1].Summary, ticket1.TicketAuthorities[0].Summary)
		require.Equal(t, int64(2*time.Hour/time.Millisecond), *ticket1.TicketAuthorities[0].ExpiryDurationMS)
	})

	t.Run("return isUsed false if ticket is not used", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)
		user, err := signUp(ctx, db, userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		ticket, err := createTicketWithTx(ctx, db, user.UserID, "ticketName", []*dto.UserAuthorityReq{{AuthorityCode: authorities[0].AuthorityCode}}, 1)
		require.NoError(t, err)

		isUsed, err := service.CheckUseTicket(ctx, db, user.UserID, ticket.TicketID)
		require.NoError(t, err)
		require.False(t, isUsed)
	})

	t.Run("return isUsed true if ticket is used", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)
		user, err := signUp(ctx, db, userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		ticket, err := createTicketWithTx(ctx, db, user.UserID, "ticketName", []*dto.UserAuthorityReq{{AuthorityCode: authorities[0].AuthorityCode}}, 1)
		require.NoError(t, err)

		_, err = useTicket(ctx, db, user.UserID, ticket.UUID)
		require.NoError(t, err)

		isUsed, err := service.CheckUseTicket(ctx, db, user.UserID, ticket.TicketID)
		require.NoError(t, err)
		require.True(t, isUsed)
	})

	t.Run("return all created tickets by specific user", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)
		admin1 := signUpAdminUser(t, ctx, db)

		admin2, err := signUp(ctx, db, &ooauth.UserInfo{
			AuthorizedBy: domain.GOOGLE,
			AuthorizedID: "anotherAdmin",
			Email:        "anotherAdmin@google.com",
			Username:     "anotherAdmin",
		}, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		_, err = createTicketWithTx(ctx, db, admin1.UserID, "ticket1", []*dto.UserAuthorityReq{
			{AuthorityCode: authorities[0].AuthorityCode, ExpiryDurationMS: ptr.P(int64(2 * time.Hour / time.Millisecond))},
			{AuthorityCode: authorities[1].AuthorityCode},
		}, 1)
		require.NoError(t, err)

		_, err = createTicketWithTx(ctx, db, admin1.UserID, "ticket2", []*dto.UserAuthorityReq{
			{AuthorityCode: authorities[0].AuthorityCode},
		}, 2)
		require.NoError(t, err)

		_, err = createTicketWithTx(ctx, db, admin2.UserID, "ticket3", []*dto.UserAuthorityReq{
			{AuthorityCode: authorities[1].AuthorityCode},
		}, 3)
		require.NoError(t, err)

		tickets, err := service.GetAllTicketsByUserId(ctx, db, admin1.UserID)
		require.NoError(t, err)
		require.Len(t, tickets, 2)

		ticket1 := tickets[0]
		require.Equal(t, "ticket1", ticket1.TicketName)
		require.Equal(t, 1, ticket1.UseableCount)
		require.Equal(t, 0, ticket1.UsedCount)
		require.Len(t, ticket1.TicketAuthorities, 2)
		require.Equal(t, authorities[0].AuthorityName, ticket1.TicketAuthorities[0].AuthorityName)
		require.Equal(t, int64(2*time.Hour/time.Millisecond), *ticket1.TicketAuthorities[0].ExpiryDurationMS)
		require.Equal(t, authorities[1].AuthorityName, ticket1.TicketAuthorities[1].AuthorityName)
		require.Nil(t, ticket1.TicketAuthorities[1].ExpiryDurationMS)

		ticket2 := tickets[1]
		require.Equal(t, "ticket2", ticket2.TicketName)
		require.Equal(t, 2, ticket2.UseableCount)
		require.Equal(t, 0, ticket2.UsedCount)
		require.Len(t, ticket2.TicketAuthorities, 1)
		require.Equal(t, authorities[0].AuthorityName, ticket2.TicketAuthorities[0].AuthorityName)
		require.Nil(t, ticket2.TicketAuthorities[0].ExpiryDurationMS)

		tickets, err = service.GetAllTicketsByUserId(ctx, db, admin2.UserID)
		require.NoError(t, err)
		require.Len(t, tickets, 1)

		ticket3 := tickets[0]
		require.Equal(t, "ticket3", ticket3.TicketName)
		require.Equal(t, 3, ticket3.UseableCount)
		require.Equal(t, 0, ticket3.UsedCount)
		require.Len(t, ticket3.TicketAuthorities, 1)
		require.Equal(t, authorities[1].AuthorityName, ticket3.TicketAuthorities[0].AuthorityName)
		require.Nil(t, ticket3.TicketAuthorities[0].ExpiryDurationMS)
	})
}
