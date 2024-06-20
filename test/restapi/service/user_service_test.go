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
func initAgreementFunc(t *testing.T, db *sql.DB) (context.Context, []*models.Agreement, []*models.Agreement, []*models.Authority) { //TODO: 추후 실제 비즈니스 로직을 통해 DB에 저장하는 것으로 변경
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

	authorities := newAuthorities()
	for _, authority := range authorities {
		err := authority.Insert(ctx, db, boil.Infer())
		require.NoError(t, err)
	}
	adminAuthority := &models.Authority{AuthorityName: domain.AuthorityAdmin, Summary: "관리자 권한"}
	err := adminAuthority.Insert(ctx, db, boil.Infer())
	require.NoError(t, err)

	return ctx, requiredAgreements, optionalAgreements, authorities
}

func newAuthorities() []*models.Authority {
	return []*models.Authority{
		{AuthorityName: "AUTHORITY_USER", Summary: "사용자 권한"},
		{AuthorityName: "AUTHORITY_GUEST", Summary: "게스트 권한"},
	}
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
		otherUser, err := userService.SignUp(ctx, otherUserInfo, agreementReqs)
		require.NoError(t, err)

		err = userService.AddUserAuthorities(ctx, otherUser.UserID, []*dto.UserAuthorityReq{
			{AuthorityName: "AUTHORITY_USER", ExpiryDuration: ptr.P(dto.Duration(time.Hour * 24))},
			{AuthorityName: "AUTHORITY_GUEST", ExpiryDuration: nil},
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
		db := tinit.DB(t)
		ctx, _, _, _ := initAgreementFunc(t, db)
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService)

		_, isExisted, err := userService.FindSignedUpUser(ctx, "authorizedBy", "authorizedID")
		require.NoError(t, err)
		require.False(t, isExisted)
	})

	t.Run("return isExistedUser true when user signed up", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, requiredAgreements, optionalAgreements, _ := initAgreementFunc(t, db)
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

	t.Run("return needed necessary agreements when sign up with not answered", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, requiredAgreements, optionalAgreements, _ := initAgreementFunc(t, db)
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
		ctx, requiredAgreements, optionalAgreements, _ := initAgreementFunc(t, db)
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
		ctx, requiredAgreements, optionalAgreements, _ := initAgreementFunc(t, db)
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
		ctx, requiredAgreements, optionalAgreements, _ := initAgreementFunc(t, db)
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

	t.Run("return empty authorities when authorities is not saved", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, _ := initAgreementFunc(t, db)
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService)

		user, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		authorities, err := userService.FindUserAuthorities(ctx, user.UserID)
		require.NoError(t, err)
		require.Empty(t, authorities)
	})

	t.Run("return user's authorities when authorities is saved", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService)

		user, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		userId := user.UserID
		insertedAuthorities := []*dto.UserAuthorityReq{
			{AuthorityName: authorities[0].AuthorityName, ExpiryDuration: ptr.P(dto.Duration(time.Hour * 24))},
			{AuthorityName: authorities[1].AuthorityName, ExpiryDuration: nil},
		}
		err = userService.AddUserAuthorities(ctx, userId, insertedAuthorities)
		require.NoError(t, err)

		userAuthorities, err := userService.FindUserAuthorities(ctx, user.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, len(insertedAuthorities))

		slices.SortStableFunc(userAuthorities, func(a, b *domain.UserAuthority) int {
			return -cmp.Compare(a.AuthorityName, b.AuthorityName)
		})
		for i, authority := range userAuthorities {
			requireEqualUserRole(t, userId, time.Now(), insertedAuthorities[i], authority)
		}
	})

	t.Run("return empty authorities when authorities is expired", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService)

		user, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		err = userService.AddUserAuthorities(ctx, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityName: authorities[0].AuthorityName, ExpiryDuration: ptr.P(dto.Duration(1 * time.Second))}, //2초 후
			{AuthorityName: authorities[1].AuthorityName, ExpiryDuration: ptr.P(dto.Duration(1 * time.Second))}, //1초 후
		})
		require.NoError(t, err)
		time.Sleep(time.Second * 2) //2초 대기, 1초 후에 만료되는 AUTHORITY_USER는 만료되었을 것이다.

		userAuthorities, err := userService.FindUserAuthorities(ctx, user.UserID)
		require.NoError(t, err)
		require.Empty(t, userAuthorities)
	})

	t.Run("return authority with extended expiry date when authority was already existed", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService)

		user, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		sameAuthority := authorities[0].AuthorityName
		err = userService.AddUserAuthorities(ctx, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityName: sameAuthority, ExpiryDuration: ptr.P(dto.Duration(time.Hour * 24))},
		})
		require.NoError(t, err)

		err = userService.AddUserAuthorities(ctx, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityName: sameAuthority, ExpiryDuration: ptr.P(dto.Duration(time.Hour * 4))},
		})
		require.NoError(t, err)

		userAuthorities, err := userService.FindUserAuthorities(ctx, user.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, 1)

		require.Equal(t, sameAuthority, userAuthorities[0].AuthorityName)
		require.WithinDuration(t, time.Now().Add(time.Hour*24).Add(time.Hour*4), *userAuthorities[0].ExpiryDate, time.Second)
	})

	t.Run("return unexpired authority when existed authority had not expiry date", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService)

		user, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		sameAuthority := authorities[0].AuthorityName
		err = userService.AddUserAuthorities(ctx, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityName: sameAuthority, ExpiryDuration: nil},
		})
		require.NoError(t, err)

		err = userService.AddUserAuthorities(ctx, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityName: sameAuthority, ExpiryDuration: ptr.P(dto.Duration(time.Hour * 4))},
		})
		require.NoError(t, err)

		userAuthorities, err := userService.FindUserAuthorities(ctx, user.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, 1)

		require.Equal(t, user.UserID, userAuthorities[0].UserID)
		require.Equal(t, sameAuthority, userAuthorities[0].AuthorityName)
		require.Nil(t, userAuthorities[0].ExpiryDate)
	})

	//이미 지난 기한에 추가된 역할은 만료된 것으로 간주하고 새로운 만료일을 설정한다.
	t.Run("return authority with expiry date from now when existed authority was expired", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService)

		user, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		sameAuthority := authorities[0].AuthorityName
		err = userService.AddUserAuthorities(ctx, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityName: sameAuthority, ExpiryDuration: ptr.P(dto.Duration(time.Second * 1))},
		})
		require.NoError(t, err)
		time.Sleep(time.Second * 2) //2초 대기, AUTHORITY_ADMIN은 만료되었을 것이다.

		err = userService.AddUserAuthorities(ctx, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityName: sameAuthority, ExpiryDuration: ptr.P(dto.Duration(time.Hour * 4))},
		})
		require.NoError(t, err)

		userAuthorities, err := userService.FindUserAuthorities(ctx, user.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, 1)

		require.Equal(t, sameAuthority, userAuthorities[0].AuthorityName)
		require.WithinDuration(t, time.Now().Add(time.Hour*4), *userAuthorities[0].ExpiryDate, time.Millisecond*600)
	})

	t.Run("return unexpired authority when existed authority was given with no expiry date", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService)

		user, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		sameAuthority := authorities[0].AuthorityName
		err = userService.AddUserAuthorities(ctx, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityName: sameAuthority, ExpiryDuration: ptr.P(dto.Duration(time.Second * 1))},
		})
		require.NoError(t, err)

		err = userService.AddUserAuthorities(ctx, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityName: sameAuthority, ExpiryDuration: nil},
		})
		require.NoError(t, err)

		dUserAuthorities, err := userService.FindUserAuthorities(ctx, user.UserID)
		require.NoError(t, err)
		require.Len(t, dUserAuthorities, 1)

		require.Equal(t, sameAuthority, dUserAuthorities[0].AuthorityName)
		require.Nil(t, dUserAuthorities[0].ExpiryDate)
	})

	t.Run("can't add AUTHORITY_ADMIN to user", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, _ := initAgreementFunc(t, db)
		userService := service.NewUserService(db)
		actionOtherUserSignUP(t, ctx, userService)

		user, err := userService.SignUp(ctx, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		err = userService.AddUserAuthorities(ctx, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityName: domain.AuthorityAdmin, ExpiryDuration: nil},
		})
		require.Error(t, err)
	})
}

func requireEqualUserRole(t *testing.T, userId int, now time.Time, expected *dto.UserAuthorityReq, actual *domain.UserAuthority) {
	require.Equal(t, expected.AuthorityName, actual.AuthorityName)
	require.Equal(t, userId, actual.UserID)
	if expected.ExpiryDuration != nil {
		require.WithinDuration(t, now.Add(time.Duration(*expected.ExpiryDuration)), *actual.ExpiryDate, time.Second)
	} else {
		require.Nil(t, actual.ExpiryDate)
	}
}
