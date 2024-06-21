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
	db := tinit.DB(t)
	ticketService := service.NewTicketService(db)
	userService := service.NewUserService(db)

	t.Run("return error if authority not existed", func(t *testing.T) {
		tinit.DB(t)
		ctx := context.Background()

		_, err := ticketService.CreateTicket(ctx, []*dto.UserAuthorityReq{{AuthorityName: "notExistedAuthority"}})
		require.Error(t, err)
	})

	t.Run("return ticket id if authority existed", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)

		ticketId, err := ticketService.CreateTicket(ctx, []*dto.UserAuthorityReq{{AuthorityName: authorities[0].AuthorityName}})
		require.NoError(t, err)
		require.NotEmpty(t, ticketId)
	})

	userinfo := &ooauth.UserInfo{
		AuthorizedBy: domain.GOOGLE,
		AuthorizedID: "123456",
		Email:        "email@email.com",
		Username:     "testUsername",
	}

	t.Run("return authorities after use ticket that has authorities", func(t *testing.T) {
		db := tinit.DB(t)

		ctx, _, _, authorities := initAgreementFunc(t, db)

		userAuthorityReqs := []*dto.UserAuthorityReq{
			{AuthorityName: authorities[0].AuthorityName},
			{AuthorityName: authorities[1].AuthorityName, ExpiryDuration: ptr.P(dto.Duration(2 * time.Hour))},
		}
		ticketId, err := ticketService.CreateTicket(ctx, userAuthorityReqs)
		require.NoError(t, err)

		targetUser, err := userService.SignUp(ctx, userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		isExisted, err := ticketService.UseTicket(ctx, targetUser.UserID, ticketId)
		require.NoError(t, err)
		require.True(t, isExisted)

		userAuthorities, err := userService.FindUserAuthorities(ctx, targetUser.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, 1)

		require.Equal(t, authorities[0].AuthorityName, userAuthorities[0].AuthorityName)
		require.Nil(t, userAuthorities[0].ExpiryDate)
		require.Equal(t, authorities[1].AuthorityName, userAuthorities[1].AuthorityName)
		require.WithinDuration(t, time.Now().Add(2*time.Hour), *userAuthorities[1].ExpiryDate, 1*time.Second)
	})
}
