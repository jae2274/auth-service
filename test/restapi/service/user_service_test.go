package service

import (
	"cmp"
	"context"
	"database/sql"
	"slices"
	"testing"
	"time"

	"github.com/jae2274/auth-service/auth_service/common/domain"
	"github.com/jae2274/auth-service/auth_service/models"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
	"github.com/jae2274/auth-service/auth_service/restapi/ooauth"
	"github.com/jae2274/auth-service/auth_service/restapi/service"
	"github.com/jae2274/auth-service/test/tinit"
	"github.com/jae2274/goutils/ptr"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func newNecessaryAgreements() []*models.Agreement {
	return []*models.Agreement{
		{AgreementCode: "code1", Summary: "summary1", IsRequired: 1},
		{AgreementCode: "code2", Summary: "summary2", IsRequired: 1},
	}
}
func newOptionalAgreements() []*models.Agreement {
	return []*models.Agreement{
		{AgreementCode: "code3", Summary: "summary3", IsRequired: 0},
		{AgreementCode: "code4", Summary: "summary4", IsRequired: 0},
	}
}
func initAgreementFunc(t *testing.T, db *sql.DB) (context.Context, []*models.Agreement, []*models.Agreement) { //TODO: 추후 실제 비즈니스 로직을 통해 DB에 저장하는 것으로 변경
	//Given
	ctx := context.Background()
	var requiredAgreements models.AgreementSlice = newNecessaryAgreements()
	for _, agreement := range requiredAgreements {
		err := agreement.Insert(ctx, db, boil.Infer())
		require.NoError(t, err)
	}

	optionalAgreements := newOptionalAgreements()
	for _, agreement := range optionalAgreements {
		err := agreement.Insert(ctx, db, boil.Infer())
		require.NoError(t, err)
	}

	return ctx, requiredAgreements, optionalAgreements
}

func TestUserService(t *testing.T) {
	userinfo := ooauth.UserInfo{
		AuthorizedID: "authorizedID",
		AuthorizedBy: "GOOGLE",
		Email:        "testEmail@testmail.net",
		Username:     "testUsername",
	}

	actionOtherUserSignUP := func(t *testing.T, ctx context.Context, userService service.UserService, agreementReqs ...*dto.UserAgreementReq) *ooauth.UserInfo {
		otherUserInfo := &ooauth.UserInfo{
			AuthorizedBy: "GOOGLE",
			AuthorizedID: "otherAuthorizedID",
			Email:        "other@gmail.com",
			Username:     "other",
		}
		otherUser, err := userService.SignUp(ctx, otherUserInfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		_, err = userService.AddUserRoles(ctx, otherUser.UserID, []*domain.UserRole{
			{RoleName: "ROLE_ADMIN", ExpiryDuration: ptr.P(time.Duration(time.Hour * 24))},
			{RoleName: "ROLE_USER", ExpiryDuration: nil},
		})
		require.NoError(t, err)

		return otherUserInfo
	}

	t.Run("sign up user", func(t *testing.T) {
		ctx := context.Background()
		userService := service.NewUserService(tinit.DB(t))

		signedUpUser, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		require.Equal(t, string(userinfo.AuthorizedBy), signedUpUser.AuthorizedBy)
		require.Equal(t, userinfo.AuthorizedID, signedUpUser.AuthorizedID)
		require.Equal(t, userinfo.Email, signedUpUser.Email)
		require.Equal(t, userinfo.Username, signedUpUser.Name)
	})

	t.Run("return isExistedUser false when user did not sign up", func(t *testing.T) {
		ctx := context.Background()
		userService := service.NewUserService(tinit.DB(t))
		actionOtherUserSignUP(t, ctx, userService)

		_, isExisted, err := userService.FindSignedUpUser(ctx, "authorizedBy", "authorizedID")
		require.NoError(t, err)
		require.False(t, isExisted)
	})

	t.Run("return isExistedUser true when user signed up", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, requiredAgreements, optionalAgreements := initAgreementFunc(t, db)
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService, &dto.UserAgreementReq{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true}, &dto.UserAgreementReq{AgreementId: optionalAgreements[0].AgreementID, IsAgree: true})

		agreementReq := []*dto.UserAgreementReq{{
			AgreementId: requiredAgreements[0].AgreementID,
			IsAgree:     true,
		}, {
			AgreementId: requiredAgreements[1].AgreementID,
			IsAgree:     false,
		}, {
			AgreementId: optionalAgreements[0].AgreementID,
			IsAgree:     true,
		}, {
			AgreementId: optionalAgreements[1].AgreementID,
			IsAgree:     false,
		}}
		_, err := userService.SignUp(ctx, &userinfo, agreementReq)

		require.NoError(t, err)

		user, isExisted, err := userService.FindSignedUpUser(ctx, userinfo.AuthorizedBy, userinfo.AuthorizedID)
		require.NoError(t, err)
		require.True(t, isExisted)
		require.Equal(t, string(userinfo.AuthorizedBy), user.AuthorizedBy)
		require.Equal(t, userinfo.AuthorizedID, user.AuthorizedID)
		require.Equal(t, userinfo.Email, user.Email)
		require.Equal(t, userinfo.Username, user.Name)
	})

	t.Run("return empty agreements when agreements is not saved", func(t *testing.T) {
		ctx := context.Background()
		userService := service.NewUserService(tinit.DB(t))
		actionOtherUserSignUP(t, ctx, userService)

		user, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		agreements, err := userService.FindNecessaryAgreements(ctx, user.UserID)
		require.NoError(t, err)
		require.Empty(t, agreements)
	})

	t.Run("return needed necessary agreements when sign up with not answered", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, requiredAgreements, optionalAgreements := initAgreementFunc(t, db)
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService, &dto.UserAgreementReq{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true}, &dto.UserAgreementReq{AgreementId: optionalAgreements[0].AgreementID, IsAgree: true})

		user, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		agreements, err := userService.FindNecessaryAgreements(ctx, user.UserID)
		require.NoError(t, err)
		require.Equal(t, requiredAgreements, agreements)
	})

	t.Run("return needed necessary agreements when sign up with not agreed", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, requiredAgreements, optionalAgreements := initAgreementFunc(t, db)
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService, &dto.UserAgreementReq{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true}, &dto.UserAgreementReq{AgreementId: optionalAgreements[0].AgreementID, IsAgree: true})

		agreementReq := make([]*dto.UserAgreementReq, 0, len(requiredAgreements))
		for _, agreement := range requiredAgreements {
			agreementReq = append(agreementReq, &dto.UserAgreementReq{
				AgreementId: agreement.AgreementID,
				IsAgree:     false,
			})
		}

		user, err := userService.SignUp(ctx, &userinfo, agreementReq)
		require.NoError(t, err)

		agreements, err := userService.FindNecessaryAgreements(ctx, user.UserID)
		require.NoError(t, err)
		require.Equal(t, requiredAgreements, agreements)
	})

	t.Run("return empty necessary agreements when sign up with all agreed", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, requiredAgreements, optionalAgreements := initAgreementFunc(t, db)
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService, &dto.UserAgreementReq{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true}, &dto.UserAgreementReq{AgreementId: optionalAgreements[0].AgreementID, IsAgree: true})

		agreementReq := make([]*dto.UserAgreementReq, 0, len(requiredAgreements))
		for _, agreement := range requiredAgreements {
			agreementReq = append(agreementReq, &dto.UserAgreementReq{
				AgreementId: agreement.AgreementID,
				IsAgree:     true,
			})
		}

		user, err := userService.SignUp(ctx, &userinfo, agreementReq)
		require.NoError(t, err)

		agreements, err := userService.FindNecessaryAgreements(ctx, user.UserID)
		require.NoError(t, err)
		require.Empty(t, agreements)
	})

	t.Run("return empty necessary agreements when all agreed after sign up", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, requiredAgreements, optionalAgreements := initAgreementFunc(t, db)
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService, &dto.UserAgreementReq{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true}, &dto.UserAgreementReq{AgreementId: optionalAgreements[0].AgreementID, IsAgree: true})

		agreementReq := make([]*dto.UserAgreementReq, 0, len(requiredAgreements))
		for _, agreement := range requiredAgreements {
			agreementReq = append(agreementReq, &dto.UserAgreementReq{
				AgreementId: agreement.AgreementID,
				IsAgree:     false,
			})
		}

		user, err := userService.SignUp(ctx, &userinfo, agreementReq)
		require.NoError(t, err)

		agreementReq = make([]*dto.UserAgreementReq, 0, len(requiredAgreements))
		for _, agreement := range requiredAgreements {
			agreementReq = append(agreementReq, &dto.UserAgreementReq{
				AgreementId: agreement.AgreementID,
				IsAgree:     true,
			})
		}

		err = userService.ApplyUserAgreements(ctx, user.UserID, agreementReq)
		require.NoError(t, err)

		agreements, err := userService.FindNecessaryAgreements(ctx, user.UserID)
		require.NoError(t, err)
		require.Empty(t, agreements)
	})

	t.Run("return empty roles when roles is not saved", func(t *testing.T) {
		ctx := context.Background()
		userService := service.NewUserService(tinit.DB(t))
		actionOtherUserSignUP(t, ctx, userService)

		user, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		roles, err := userService.FindUserRoles(ctx, user.UserID)
		require.NoError(t, err)
		require.Empty(t, roles)
	})

	t.Run("return user's roles when roles is saved", func(t *testing.T) {
		db := tinit.DB(t)
		ctx := context.Background()
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService)

		user, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		userId := user.UserID
		insertedRoles := []*domain.UserRole{
			{RoleName: "ROLE_ADMIN", ExpiryDuration: ptr.P(time.Duration(time.Hour * 24))},
			{RoleName: "ROLE_USER", ExpiryDuration: nil},
		}
		_, err = userService.AddUserRoles(ctx, userId, insertedRoles)
		require.NoError(t, err)

		roles, err := userService.FindUserRoles(ctx, user.UserID)
		require.NoError(t, err)
		require.Len(t, roles, len(insertedRoles))

		slices.SortStableFunc(roles, func(a, b *models.UserRole) int {
			return cmp.Compare(a.RoleName, b.RoleName)
		})
		for i, role := range roles {
			requireEqualUserRole(t, userId, time.Now(), insertedRoles[i], role)
		}
	})

	t.Run("return empty roles when roles is expired", func(t *testing.T) {
		db := tinit.DB(t)
		ctx := context.Background()
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService)

		user, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		_, err = userService.AddUserRoles(ctx, user.UserID, []*domain.UserRole{
			{RoleName: "ROLE_ADMIN", ExpiryDuration: ptr.P(time.Duration(1 * time.Second))}, //2초 후
			{RoleName: "ROLE_USER", ExpiryDuration: ptr.P(time.Duration(1 * time.Second))},  //1초 후
		})
		require.NoError(t, err)
		time.Sleep(time.Second * 2) //2초 대기, 1초 후에 만료되는 ROLE_USER는 만료되었을 것이다.

		roles, err := userService.FindUserRoles(ctx, user.UserID)
		require.NoError(t, err)
		require.Empty(t, roles)
	})

	t.Run("return role with extended expiry date when role was already existed", func(t *testing.T) {

		db := tinit.DB(t)
		ctx := context.Background()
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService)

		user, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		_, err = userService.AddUserRoles(ctx, user.UserID, []*domain.UserRole{
			{RoleName: "ROLE_ADMIN", ExpiryDuration: ptr.P(time.Duration(time.Hour * 24))},
		})
		require.NoError(t, err)

		_, err = userService.AddUserRoles(ctx, user.UserID, []*domain.UserRole{
			{RoleName: "ROLE_ADMIN", ExpiryDuration: ptr.P(time.Duration(time.Hour * 4))},
		})
		require.NoError(t, err)

		roles, err := userService.FindUserRoles(ctx, user.UserID)
		require.NoError(t, err)
		require.Len(t, roles, 1)

		require.Equal(t, user.UserID, roles[0].UserID)
		require.Equal(t, "ROLE_ADMIN", roles[0].RoleName)
		require.WithinDuration(t, time.Now().Add(time.Hour*24).Add(time.Hour*4), roles[0].ExpiryDate.Time, time.Second)
	})

	t.Run("return unexpired role when existed role had not expiry date", func(t *testing.T) {
		db := tinit.DB(t)
		ctx := context.Background()
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService)

		user, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		_, err = userService.AddUserRoles(ctx, user.UserID, []*domain.UserRole{
			{RoleName: "ROLE_ADMIN", ExpiryDuration: nil},
		})
		require.NoError(t, err)

		_, err = userService.AddUserRoles(ctx, user.UserID, []*domain.UserRole{
			{RoleName: "ROLE_ADMIN", ExpiryDuration: ptr.P(time.Duration(time.Hour * 4))},
		})
		require.NoError(t, err)

		roles, err := userService.FindUserRoles(ctx, user.UserID)
		require.NoError(t, err)
		require.Len(t, roles, 1)

		require.Equal(t, user.UserID, roles[0].UserID)
		require.Equal(t, "ROLE_ADMIN", roles[0].RoleName)
		require.Equal(t, false, roles[0].ExpiryDate.Valid)
	})

	//이미 지난 기한에 추가된 역할은 만료된 것으로 간주하고 새로운 만료일을 설정한다.
	t.Run("return role with expiry date from now when existed role was expired", func(t *testing.T) {
		db := tinit.DB(t)
		ctx := context.Background()
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService)

		user, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		_, err = userService.AddUserRoles(ctx, user.UserID, []*domain.UserRole{
			{RoleName: "ROLE_ADMIN", ExpiryDuration: ptr.P(time.Duration(time.Second * 1))},
		})
		require.NoError(t, err)
		time.Sleep(time.Second * 2) //2초 대기, ROLE_ADMIN은 만료되었을 것이다.

		_, err = userService.AddUserRoles(ctx, user.UserID, []*domain.UserRole{
			{RoleName: "ROLE_ADMIN", ExpiryDuration: ptr.P(time.Duration(time.Hour * 4))},
		})
		require.NoError(t, err)

		roles, err := userService.FindUserRoles(ctx, user.UserID)
		require.NoError(t, err)
		require.Len(t, roles, 1)

		require.Equal(t, user.UserID, roles[0].UserID)
		require.Equal(t, "ROLE_ADMIN", roles[0].RoleName)
		require.WithinDuration(t, time.Now().Add(time.Hour*4), roles[0].ExpiryDate.Time, time.Millisecond*500)
	})

	t.Run("return unexpired role when existed role was given with no expiry date", func(t *testing.T) {
		db := tinit.DB(t)
		ctx := context.Background()
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService)

		user, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		_, err = userService.AddUserRoles(ctx, user.UserID, []*domain.UserRole{
			{RoleName: "ROLE_ADMIN", ExpiryDuration: ptr.P(time.Duration(time.Second * 1))},
		})
		require.NoError(t, err)

		_, err = userService.AddUserRoles(ctx, user.UserID, []*domain.UserRole{
			{RoleName: "ROLE_ADMIN", ExpiryDuration: nil},
		})
		require.NoError(t, err)

		roles, err := userService.FindUserRoles(ctx, user.UserID)
		require.NoError(t, err)
		require.Len(t, roles, 1)

		require.Equal(t, user.UserID, roles[0].UserID)
		require.Equal(t, "ROLE_ADMIN", roles[0].RoleName)
		require.Equal(t, false, roles[0].ExpiryDate.Valid)
	})
}

func requireEqualUserRole(t *testing.T, userId int, now time.Time, expected *domain.UserRole, actual *models.UserRole) {
	require.Equal(t, expected.RoleName, actual.RoleName)
	require.Equal(t, userId, actual.UserID)
	if expected.ExpiryDuration != nil {
		require.True(t, actual.ExpiryDate.Valid)
		require.WithinDuration(t, now.Add(*expected.ExpiryDuration), actual.ExpiryDate.Time, time.Second)
	} else {
		require.False(t, actual.ExpiryDate.Valid)
	}
}
