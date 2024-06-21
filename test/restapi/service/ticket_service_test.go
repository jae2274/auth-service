package service

import (
	"context"
	"testing"
	"time"

	"github.com/jae2274/auth-service/auth_service/common/domain"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
	"github.com/jae2274/auth-service/auth_service/restapi/ooauth"
	"github.com/jae2274/auth-service/auth_service/restapi/service"
	"github.com/jae2274/auth-service/test/tinit"
	"github.com/jae2274/goutils/ptr"
	"github.com/stretchr/testify/require"
)

func TestTicketService(t *testing.T) {

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

		_, err := service.CreateTicket(ctx, db, []*dto.UserAuthorityReq{{AuthorityCode: "notExistedAuthority"}})
		require.Error(t, err)
	})

	t.Run("return ticket info if ticket existed", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)

		ticketAuthorities := []*dto.UserAuthorityReq{
			{AuthorityCode: authorities[0].AuthorityCode},
			{AuthorityCode: authorities[1].AuthorityCode, ExpiryDuration: ptr.P(dto.Duration(2 * time.Hour))},
		}
		ticketId, err := service.CreateTicket(ctx, db, ticketAuthorities)
		require.NoError(t, err)
		require.NotEmpty(t, ticketId)

		res, isExisted, err := service.GetTicketInfo(ctx, db, ticketId)
		require.NoError(t, err)
		require.True(t, isExisted)
		require.Equal(t, ticketId, res.TicketId)
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

	userinfo := &ooauth.UserInfo{
		AuthorizedBy: domain.GOOGLE,
		AuthorizedID: "123456",
		Email:        "email@email.com",
		Username:     "testUsername",
	}

	t.Run("return isExisted false before ticket is created", func(t *testing.T) {
		db := tinit.DB(t)

		ctx, _, _, _ := initAgreementFunc(t, db)

		targetUser, err := service.SignUp(ctx, db, userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		isExisted, err := service.UseTicket(ctx, db, targetUser.UserID, "notExistedTicketId")
		require.NoError(t, err)
		require.False(t, isExisted)
	})

	t.Run("return authorities after use ticket that has authorities", func(t *testing.T) {
		db := tinit.DB(t)

		ctx, _, _, authorities := initAgreementFunc(t, db)

		userAuthorityReqs := []*dto.UserAuthorityReq{
			{AuthorityCode: authorities[0].AuthorityCode},
			{AuthorityCode: authorities[1].AuthorityCode, ExpiryDuration: ptr.P(dto.Duration(2 * time.Hour))},
		}
		ticketId, err := service.CreateTicket(ctx, db, userAuthorityReqs)
		require.NoError(t, err)

		targetUser, err := service.SignUp(ctx, db, userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		isExisted, err := service.UseTicket(ctx, db, targetUser.UserID, ticketId)
		require.NoError(t, err)
		require.True(t, isExisted)

		userAuthorities, err := service.FindUserAuthorities(ctx, db, targetUser.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, len(userAuthorityReqs))

		now := time.Now()
		for i, userAuthority := range userAuthorities {
			requireEqualUserRole(t, targetUser.UserID, now, userAuthorityReqs[i], userAuthority)
		}
	})

	t.Run("return isExisted false after ticket is used", func(t *testing.T) {
		db := tinit.DB(t)

		ctx, _, _, authorities := initAgreementFunc(t, db)
		userAuthorityReqs := []*dto.UserAuthorityReq{
			{AuthorityCode: authorities[0].AuthorityCode},
			{AuthorityCode: authorities[1].AuthorityCode, ExpiryDuration: ptr.P(dto.Duration(2 * time.Hour))},
		}

		ticketId, err := service.CreateTicket(ctx, db, userAuthorityReqs)
		require.NoError(t, err)

		targetUser, err := service.SignUp(ctx, db, userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		isExisted, err := service.UseTicket(ctx, db, targetUser.UserID, ticketId)
		require.NoError(t, err)
		require.True(t, isExisted)

		isExisted, err = service.UseTicket(ctx, db, targetUser.UserID, ticketId)
		require.NoError(t, err)
		require.False(t, isExisted)

		userAuthorities, err := service.FindUserAuthorities(ctx, db, targetUser.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, len(userAuthorityReqs))

		now := time.Now()
		for i, userAuthority := range userAuthorities {
			requireEqualUserRole(t, targetUser.UserID, now, userAuthorityReqs[i], userAuthority)
		}
	})

}
