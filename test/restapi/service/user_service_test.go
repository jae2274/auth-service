package service

import (
	"cmp"
	"context"
	"database/sql"
	"slices"
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

func TestUserService(t *testing.T) {
	userinfo := ooauth.UserInfo{
		AuthorizedID: "authorizedID",
		AuthorizedBy: "GOOGLE",
		Email:        "testEmail@testmail.net",
		Username:     "testUsername",
	}

	actionOtherUserSignUP := func(t *testing.T, ctx context.Context, db *sql.DB, agreementReqs ...*dto.UserAgreementReq) *ooauth.UserInfo {
		otherUserInfo := &ooauth.UserInfo{
			AuthorizedBy: "GOOGLE",
			AuthorizedID: "otherAuthorizedID",
			Email:        "other@gmail.com",
			Username:     "other",
		}
		otherUser, err := service.SignUp(ctx, db, otherUserInfo, agreementReqs)
		require.NoError(t, err)

		err = service.AddUserAuthorities(ctx, db, otherUser.UserID, []*dto.UserAuthorityReq{
			{AuthorityCode: "AUTHORITY_USER", ExpiryDuration: ptr.P(dto.Duration(time.Hour * 24))},
			{AuthorityCode: "AUTHORITY_GUEST", ExpiryDuration: nil},
		})
		require.NoError(t, err)

		return otherUserInfo
	}

	t.Run("sign up user", func(t *testing.T) {
		db := tinit.DB(t)
		ctx := context.Background()

		signedUpUser, err := service.SignUp(ctx, db, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		require.Equal(t, string(userinfo.AuthorizedBy), signedUpUser.AuthorizedBy)
		require.Equal(t, userinfo.AuthorizedID, signedUpUser.AuthorizedID)
		require.Equal(t, userinfo.Email, signedUpUser.Email)
		require.Equal(t, userinfo.Username, signedUpUser.Name)
	})

	t.Run("return isExistedUser false when user did not sign up", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, _ := initAgreementFunc(t, db)
		actionOtherUserSignUP(t, ctx, db)

		_, isExisted, err := service.FindSignedUpUser(ctx, db, "authorizedBy", "authorizedID")
		require.NoError(t, err)
		require.False(t, isExisted)
	})

	t.Run("return isExistedUser true when user signed up", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, requiredAgreements, optionalAgreements, _ := initAgreementFunc(t, db)

		actionOtherUserSignUP(t, ctx, db, &dto.UserAgreementReq{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true}, &dto.UserAgreementReq{AgreementId: optionalAgreements[0].AgreementID, IsAgree: true})

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
		_, err := service.SignUp(ctx, db, &userinfo, agreementReq)

		require.NoError(t, err)

		user, isExisted, err := service.FindSignedUpUser(ctx, db, userinfo.AuthorizedBy, userinfo.AuthorizedID)
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

		actionOtherUserSignUP(t, ctx, db, &dto.UserAgreementReq{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true}, &dto.UserAgreementReq{AgreementId: optionalAgreements[0].AgreementID, IsAgree: true})

		user, err := service.SignUp(ctx, db, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		agreements, err := service.FindNecessaryAgreements(ctx, db, user.UserID)
		require.NoError(t, err)
		require.Equal(t, requiredAgreements, agreements)
	})

	t.Run("return needed necessary agreements when sign up with not agreed", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, requiredAgreements, optionalAgreements, _ := initAgreementFunc(t, db)

		actionOtherUserSignUP(t, ctx, db, &dto.UserAgreementReq{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true}, &dto.UserAgreementReq{AgreementId: optionalAgreements[0].AgreementID, IsAgree: true})

		agreementReq := make([]*dto.UserAgreementReq, 0, len(requiredAgreements))
		for _, agreement := range requiredAgreements {
			agreementReq = append(agreementReq, &dto.UserAgreementReq{
				AgreementId: agreement.AgreementID,
				IsAgree:     false,
			})
		}

		user, err := service.SignUp(ctx, db, &userinfo, agreementReq)
		require.NoError(t, err)

		agreements, err := service.FindNecessaryAgreements(ctx, db, user.UserID)
		require.NoError(t, err)
		require.Equal(t, requiredAgreements, agreements)
	})

	t.Run("return empty necessary agreements when sign up with all agreed", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, requiredAgreements, optionalAgreements, _ := initAgreementFunc(t, db)

		actionOtherUserSignUP(t, ctx, db, &dto.UserAgreementReq{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true}, &dto.UserAgreementReq{AgreementId: optionalAgreements[0].AgreementID, IsAgree: true})

		agreementReq := make([]*dto.UserAgreementReq, 0, len(requiredAgreements))
		for _, agreement := range requiredAgreements {
			agreementReq = append(agreementReq, &dto.UserAgreementReq{
				AgreementId: agreement.AgreementID,
				IsAgree:     true,
			})
		}

		user, err := service.SignUp(ctx, db, &userinfo, agreementReq)
		require.NoError(t, err)

		agreements, err := service.FindNecessaryAgreements(ctx, db, user.UserID)
		require.NoError(t, err)
		require.Empty(t, agreements)
	})

	t.Run("return empty necessary agreements when all agreed after sign up", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, requiredAgreements, optionalAgreements, _ := initAgreementFunc(t, db)

		actionOtherUserSignUP(t, ctx, db, &dto.UserAgreementReq{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true}, &dto.UserAgreementReq{AgreementId: optionalAgreements[0].AgreementID, IsAgree: true})

		agreementReq := make([]*dto.UserAgreementReq, 0, len(requiredAgreements))
		for _, agreement := range requiredAgreements {
			agreementReq = append(agreementReq, &dto.UserAgreementReq{
				AgreementId: agreement.AgreementID,
				IsAgree:     false,
			})
		}

		user, err := service.SignUp(ctx, db, &userinfo, agreementReq)
		require.NoError(t, err)

		agreementReq = make([]*dto.UserAgreementReq, 0, len(requiredAgreements))
		for _, agreement := range requiredAgreements {
			agreementReq = append(agreementReq, &dto.UserAgreementReq{
				AgreementId: agreement.AgreementID,
				IsAgree:     true,
			})
		}

		err = service.ApplyUserAgreements(ctx, db, user.UserID, agreementReq)
		require.NoError(t, err)

		agreements, err := service.FindNecessaryAgreements(ctx, db, user.UserID)
		require.NoError(t, err)
		require.Empty(t, agreements)
	})

	t.Run("return empty authorities when authorities is not saved", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, _ := initAgreementFunc(t, db)

		actionOtherUserSignUP(t, ctx, db)

		user, err := service.SignUp(ctx, db, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		authorities, err := service.FindUserAuthorities(ctx, db, user.UserID)
		require.NoError(t, err)
		require.Empty(t, authorities)
	})

	t.Run("return error when try to add not existed authority", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, _ := initAgreementFunc(t, db)

		actionOtherUserSignUP(t, ctx, db)

		user, err := service.SignUp(ctx, db, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		err = service.AddUserAuthorities(ctx, db, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityCode: "notExistedAuthority", ExpiryDuration: nil},
		})
		require.Error(t, err)
	})

	t.Run("return user's authorities when authorities is saved", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)

		actionOtherUserSignUP(t, ctx, db)

		user, err := service.SignUp(ctx, db, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		userId := user.UserID
		insertedAuthorities := []*dto.UserAuthorityReq{
			{AuthorityCode: authorities[0].AuthorityCode, ExpiryDuration: ptr.P(dto.Duration(time.Hour * 24))},
			{AuthorityCode: authorities[1].AuthorityCode, ExpiryDuration: nil},
		}
		err = service.AddUserAuthorities(ctx, db, userId, insertedAuthorities)
		require.NoError(t, err)

		userAuthorities, err := service.FindUserAuthorities(ctx, db, user.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, len(insertedAuthorities))

		slices.SortStableFunc(userAuthorities, func(a, b *domain.UserAuthority) int {
			return -cmp.Compare(a.AuthorityCode, b.AuthorityCode)
		})
		for i, authority := range userAuthorities {
			requireEqualUserRole(t, userId, time.Now(), insertedAuthorities[i], authority)
		}
	})

	t.Run("return empty authorities when authorities is expired", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)

		actionOtherUserSignUP(t, ctx, db)

		user, err := service.SignUp(ctx, db, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		err = service.AddUserAuthorities(ctx, db, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityCode: authorities[0].AuthorityCode, ExpiryDuration: ptr.P(dto.Duration(1 * time.Second))}, //2초 후
			{AuthorityCode: authorities[1].AuthorityCode, ExpiryDuration: ptr.P(dto.Duration(1 * time.Second))}, //1초 후
		})
		require.NoError(t, err)
		time.Sleep(time.Second * 2) //2초 대기, 1초 후에 만료되는 AUTHORITY_USER는 만료되었을 것이다.

		userAuthorities, err := service.FindUserAuthorities(ctx, db, user.UserID)
		require.NoError(t, err)
		require.Empty(t, userAuthorities)
	})

	t.Run("return authority with extended expiry date when authority was already existed", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)

		actionOtherUserSignUP(t, ctx, db)

		user, err := service.SignUp(ctx, db, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		sameAuthority := authorities[0].AuthorityCode
		err = service.AddUserAuthorities(ctx, db, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityCode: sameAuthority, ExpiryDuration: ptr.P(dto.Duration(time.Hour * 24))},
		})
		require.NoError(t, err)

		err = service.AddUserAuthorities(ctx, db, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityCode: sameAuthority, ExpiryDuration: ptr.P(dto.Duration(time.Hour * 4))},
		})
		require.NoError(t, err)

		userAuthorities, err := service.FindUserAuthorities(ctx, db, user.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, 1)

		require.Equal(t, sameAuthority, userAuthorities[0].AuthorityCode)
		require.WithinDuration(t, time.Now().Add(time.Hour*24).Add(time.Hour*4), *userAuthorities[0].ExpiryDate, time.Second)
	})

	t.Run("return unexpired authority when existed authority had not expiry date", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)

		actionOtherUserSignUP(t, ctx, db)

		user, err := service.SignUp(ctx, db, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		sameAuthority := authorities[0].AuthorityCode
		err = service.AddUserAuthorities(ctx, db, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityCode: sameAuthority, ExpiryDuration: nil},
		})
		require.NoError(t, err)

		err = service.AddUserAuthorities(ctx, db, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityCode: sameAuthority, ExpiryDuration: ptr.P(dto.Duration(time.Hour * 4))},
		})
		require.NoError(t, err)

		userAuthorities, err := service.FindUserAuthorities(ctx, db, user.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, 1)

		require.Equal(t, user.UserID, userAuthorities[0].UserID)
		require.Equal(t, sameAuthority, userAuthorities[0].AuthorityCode)
		require.Nil(t, userAuthorities[0].ExpiryDate)
	})

	//이미 지난 기한에 추가된 역할은 만료된 것으로 간주하고 새로운 만료일을 설정한다.
	t.Run("return authority with expiry date from now when existed authority was expired", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)

		actionOtherUserSignUP(t, ctx, db)

		user, err := service.SignUp(ctx, db, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		sameAuthority := authorities[0].AuthorityCode
		err = service.AddUserAuthorities(ctx, db, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityCode: sameAuthority, ExpiryDuration: ptr.P(dto.Duration(time.Second * 1))},
		})
		require.NoError(t, err)
		time.Sleep(time.Second * 2) //2초 대기, AUTHORITY_ADMIN은 만료되었을 것이다.

		err = service.AddUserAuthorities(ctx, db, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityCode: sameAuthority, ExpiryDuration: ptr.P(dto.Duration(time.Hour * 4))},
		})
		require.NoError(t, err)

		userAuthorities, err := service.FindUserAuthorities(ctx, db, user.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, 1)

		require.Equal(t, sameAuthority, userAuthorities[0].AuthorityCode)
		require.WithinDuration(t, time.Now().Add(time.Hour*4), *userAuthorities[0].ExpiryDate, time.Millisecond*600)
	})

	t.Run("return unexpired authority when existed authority was given with no expiry date", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)

		actionOtherUserSignUP(t, ctx, db)

		user, err := service.SignUp(ctx, db, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		sameAuthority := authorities[0].AuthorityCode
		err = service.AddUserAuthorities(ctx, db, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityCode: sameAuthority, ExpiryDuration: ptr.P(dto.Duration(time.Second * 1))},
		})
		require.NoError(t, err)

		err = service.AddUserAuthorities(ctx, db, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityCode: sameAuthority, ExpiryDuration: nil},
		})
		require.NoError(t, err)

		dUserAuthorities, err := service.FindUserAuthorities(ctx, db, user.UserID)
		require.NoError(t, err)
		require.Len(t, dUserAuthorities, 1)

		require.Equal(t, sameAuthority, dUserAuthorities[0].AuthorityCode)
		require.Nil(t, dUserAuthorities[0].ExpiryDate)
	})

	t.Run("can't add AUTHORITY_ADMIN to user", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, _ := initAgreementFunc(t, db)

		actionOtherUserSignUP(t, ctx, db)

		user, err := service.SignUp(ctx, db, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		err = service.AddUserAuthorities(ctx, db, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityCode: domain.AuthorityAdmin, ExpiryDuration: nil},
		})
		require.Error(t, err)
	})

	t.Run("return authorities without removed authority", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)

		actionOtherUserSignUP(t, ctx, db)

		user, err := service.SignUp(ctx, db, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		err = service.AddUserAuthorities(ctx, db, user.UserID, []*dto.UserAuthorityReq{
			{AuthorityCode: authorities[0].AuthorityCode, ExpiryDuration: nil},
			{AuthorityCode: authorities[1].AuthorityCode, ExpiryDuration: nil},
		})
		require.NoError(t, err)

		err = service.RemoveAuthority(ctx, db, user.UserID, authorities[1].AuthorityCode)
		require.NoError(t, err)

		userAuthorities, err := service.FindUserAuthorities(ctx, db, user.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, 1)
		require.Equal(t, authorities[0].AuthorityCode, userAuthorities[0].AuthorityCode)
	})

	t.Run("return error when try to remove not existed authority", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, _ := initAgreementFunc(t, db)

		actionOtherUserSignUP(t, ctx, db)

		user, err := service.SignUp(ctx, db, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		err = service.RemoveAuthority(ctx, db, user.UserID, "notExistedAuthority")
		require.Error(t, err)
	})

	t.Run("can remove even if no authority to remove", func(t *testing.T) {
		db := tinit.DB(t)
		ctx, _, _, authorities := initAgreementFunc(t, db)

		actionOtherUserSignUP(t, ctx, db)

		user, err := service.SignUp(ctx, db, &userinfo, []*dto.UserAgreementReq{})
		require.NoError(t, err)

		err = service.RemoveAuthority(ctx, db, user.UserID, authorities[0].AuthorityCode)
		require.NoError(t, err)
	})
}
