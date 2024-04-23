package service

import (
	"context"
	"testing"
	"userService/test/tinit"
	"userService/usersvc/common/domain"
	"userService/usersvc/restapi/ctrlr/dto"
	"userService/usersvc/restapi/jwtutils"
	"userService/usersvc/restapi/ooauth"
	"userService/usersvc/restapi/service"

	"github.com/stretchr/testify/require"
)

func TestUserService(t *testing.T) {

	t.Run("return new_user", func(t *testing.T) {
		ctx := context.Background()
		userSvc := initUserService(t)

		res, err := userSvc.SignIn(ctx, domain.GOOGLE, "authId", "test@gmail.com")
		require.NoError(t, err)

		require.Equal(t, dto.SignInNewUser, res.SignInStatus)
	})

	t.Run("return require_agreement", func(t *testing.T) {
		ctx := context.Background()
		userSvc := initUserService(t)

		userSvc.SignUp(ctx, &dto.SignUpRequest{
			Username:   "testUsername",
			Agreements: []*dto.UserAgreementReq{},
		}, &ooauth.UserInfo{
			AuthorizedBy: domain.GOOGLE,
			AuthorizedID: "authId",
			Email:        "",
		})

		res, err := userSvc.SignIn(ctx, domain.GOOGLE, "authId", "test@gmail.com")
		require.NoError(t, err)

		require.Equal(t, dto.SignInSuccess, res.SignInStatus)
	})

	t.Run("return success", func(t *testing.T) {
	})

}

func initUserService(t *testing.T) service.UserService {
	return service.NewUserService(tinit.DB(t), jwtutils.NewJwtUtils([]byte("secretKey")))
}
