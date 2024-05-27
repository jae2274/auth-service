package service

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"userService/test/tinit"
	"userService/usersvc/common/domain"
	"userService/usersvc/models"
	"userService/usersvc/restapi/ctrlr/dto"
	"userService/usersvc/restapi/jwtutils"
	"userService/usersvc/restapi/ooauth"
	"userService/usersvc/restapi/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jae2274/goutils/terr"
	"github.com/stretchr/testify/require"
)

func TestAdminService(t *testing.T) {
	t.Run("cause error when create ticket with empty roles", func(t *testing.T) {

		ctx := context.Background()
		adminSvc := service.NewAdminService(tinit.DB(t))
		_, err := adminSvc.CreateRoleTicket(ctx, []*models.TicketRole{})

		require.Error(t, err)
	})

	t.Run("success create ticket", func(t *testing.T) {
		ctx := context.Background()
		adminSvc := service.NewAdminService(tinit.DB(t))
		ticketId, err := adminSvc.CreateRoleTicket(ctx, []*models.TicketRole{
			{RoleName: "ROLE_ADMIN"},
		})

		require.NoError(t, err)
		require.NotEmpty(t, ticketId)

	})

	t.Run("cause error when create ticket with duplicated roles", func(t *testing.T) {
		ctx := context.Background()
		adminSvc := service.NewAdminService(tinit.DB(t))
		_, err := adminSvc.CreateRoleTicket(ctx, []*models.TicketRole{
			{RoleName: "ROLE_DUPLICATE"},
			{RoleName: "ROLE_DUPLICATE"},
		})

		require.Error(t, err)
	})

	userinfo := &ooauth.UserInfo{
		AuthorizedBy: domain.GOOGLE,
		AuthorizedID: "testAuthorizedID",
		Email:        "testEmail@test.com",
		Username:     "testUsername",
	}

	t.Run("cause error when use non-existed ticket", func(t *testing.T) {
		ctx := context.Background()
		db := tinit.DB(t)
		byteSecretKey := []byte("secretKey")
		userSvc := service.NewUserService(db, jwtutils.NewJwtUtils(byteSecretKey))
		adminSvc := service.NewAdminService(db)

		userId := signUpAndIn(t, ctx, userSvc, byteSecretKey, userinfo)

		err := adminSvc.UseTicket(ctx, userId, "non-existed-ticket")
		require.Error(t, err)
	})

	t.Run("success use ticket", func(t *testing.T) {
		ctx := context.Background()
		db := tinit.DB(t)
		byteSecretKey := []byte("secretKey")
		userSvc := service.NewUserService(db, jwtutils.NewJwtUtils(byteSecretKey))
		adminSvc := service.NewAdminService(db)

		userId := signUpAndIn(t, ctx, userSvc, byteSecretKey, userinfo)

		ticketId, err := adminSvc.CreateRoleTicket(ctx, []*models.TicketRole{{RoleName: "ROLE_ADMIN"}})
		require.NoError(t, err)

		err = adminSvc.UseTicket(ctx, userId, ticketId)
		require.NoError(t, err)

		claims := signIn(t, ctx, userSvc, byteSecretKey, userinfo)
		require.Contains(t, claims.Roles, "ROLE_ADMIN")
	})

	t.Run("cause error when already used ticket", func(t *testing.T) {
		ctx := context.Background()
		db := tinit.DB(t)
		byteSecretKey := []byte("secretKey")
		userSvc := service.NewUserService(db, jwtutils.NewJwtUtils(byteSecretKey))
		adminSvc := service.NewAdminService(db)

		userId := signUpAndIn(t, ctx, userSvc, byteSecretKey, userinfo)

		ticketId, err := adminSvc.CreateRoleTicket(ctx, []*models.TicketRole{{RoleName: "ROLE_ADMIN"}})
		require.NoError(t, err)

		adminSvc.UseTicket(ctx, userId, ticketId)
		err = adminSvc.UseTicket(ctx, userId, ticketId)
		require.Error(t, err)
	})
}

func signUpAndIn(t *testing.T, ctx context.Context, userSvc service.UserService, secretKey []byte, userinfo *ooauth.UserInfo) int {
	err := userSvc.SignUp(ctx, userinfo, []*dto.UserAgreementReq{})
	require.NoError(t, err)

	res, err := userSvc.SignIn(ctx, userinfo, []*dto.UserAgreementReq{})
	require.NoError(t, err)
	require.Equal(t, res.SignInStatus, dto.SignInSuccess)

	claims, err := parseToken(secretKey, res.SuccessRes.AccessToken)
	require.NoError(t, err)

	userIdInt, err := strconv.Atoi(claims.UserId)
	require.NoError(t, err)

	return userIdInt
}

func signIn(t *testing.T, ctx context.Context, userSvc service.UserService, secretKey []byte, userinfo *ooauth.UserInfo) *jwtutils.CustomClaims {
	res, err := userSvc.SignIn(ctx, userinfo, []*dto.UserAgreementReq{})
	require.NoError(t, err)
	require.Equal(t, res.SignInStatus, dto.SignInSuccess)

	claims, err := parseToken(secretKey, res.SuccessRes.AccessToken)
	require.NoError(t, err)

	return claims
}

func parseToken(secretKey []byte, tokenString string) (*jwtutils.CustomClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(tokenString, &jwtutils.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if jwtToken.Valid {
		if claims, ok := jwtToken.Claims.(*jwtutils.CustomClaims); ok {
			return claims, nil
		} else {
			return &jwtutils.CustomClaims{}, terr.New("invalid token. claims is not CustomClaims type")
		}
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return &jwtutils.CustomClaims{}, nil
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return &jwtutils.CustomClaims{}, terr.New("invalid token. token is malformed")
	} else {
		return &jwtutils.CustomClaims{}, terr.Wrap(err)
	}
}
