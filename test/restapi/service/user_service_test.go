package service

import (
	"context"
	"testing"
	"userService/test/tinit"
	"userService/usersvc/common/domain"
	"userService/usersvc/models"
	"userService/usersvc/restapi/ctrlr/dto"
	"userService/usersvc/restapi/jwtutils"
	"userService/usersvc/restapi/ooauth"
	"userService/usersvc/restapi/service"

	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestUsers(t *testing.T) {

	t.Run("return new_user", func(t *testing.T) {
		ctx := context.Background()
		userSvc := service.NewUserService(tinit.DB(t), jwtutils.NewJwtUtils([]byte("secretKey")))

		res, err := userSvc.SignIn(ctx, domain.GOOGLE, "authId")
		require.NoError(t, err)

		require.Equal(t, dto.SignInNewUser, res.SignInStatus)
	})

	t.Run("return success", func(t *testing.T) {
		t.Run("if required agreements not existed", func(t *testing.T) {
			ctx := context.Background()
			userSvc := service.NewUserService(tinit.DB(t), jwtutils.NewJwtUtils([]byte("secretKey")))

			userSvc.SignUp(ctx, &dto.SignUpRequest{
				Username:   "testUsername",
				Agreements: []*dto.UserAgreementReq{},
			}, &ooauth.UserInfo{
				AuthorizedBy: domain.GOOGLE,
				AuthorizedID: "authId",
				Email:        "",
			})

			res, err := userSvc.SignIn(ctx, domain.GOOGLE, "authId")
			require.NoError(t, err)

			require.Equal(t, dto.SignInSuccess, res.SignInStatus)
		})

		t.Run("if all required agreements are agreed", func(t *testing.T) {
			ctx := context.Background()
			db := tinit.DB(t)
			userSvc := service.NewUserService(db, jwtutils.NewJwtUtils([]byte("secretKey")))
			var requiredAgreements models.AgreementSlice = newRequiredAgreements()
			for _, agreement := range requiredAgreements {
				err := agreement.Insert(ctx, db, boil.Infer())
				require.NoError(t, err)
			}

			userSvc.SignUp(ctx, &dto.SignUpRequest{
				Username: "testUsername",
				Agreements: []*dto.UserAgreementReq{
					{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true},
					{AgreementId: requiredAgreements[1].AgreementID, IsAgree: true},
					{AgreementId: requiredAgreements[2].AgreementID, IsAgree: false},
				},
			}, &ooauth.UserInfo{
				AuthorizedBy: domain.GOOGLE,
				AuthorizedID: "authId",
				Email:        "",
			})

			res, err := userSvc.SignIn(ctx, domain.GOOGLE, "authId")
			require.NoError(t, err)

			require.Equal(t, dto.SignInSuccess, res.SignInStatus)
		})
	})

	t.Run("return require_agreement ", func(t *testing.T) {
		t.Run("if not all required agreements are now agreed", func(t *testing.T) {
			ctx := context.Background()
			db := tinit.DB(t)
			userSvc := service.NewUserService(db, jwtutils.NewJwtUtils([]byte("secretKey")))
			var requiredAgreements models.AgreementSlice = newRequiredAgreements()
			for _, agreement := range requiredAgreements {
				err := agreement.Insert(ctx, db, boil.Infer())
				require.NoError(t, err)
			}

			userSvc.SignUp(ctx, &dto.SignUpRequest{
				Username: "testUsername",
				Agreements: []*dto.UserAgreementReq{
					{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true},
					{AgreementId: requiredAgreements[1].AgreementID, IsAgree: false},
					{AgreementId: requiredAgreements[2].AgreementID, IsAgree: false},
				},
			}, &ooauth.UserInfo{
				AuthorizedBy: domain.GOOGLE,
				AuthorizedID: "authId",
				Email:        "",
			})

			res, err := userSvc.SignIn(ctx, domain.GOOGLE, "authId")
			require.NoError(t, err)

			require.Equal(t, dto.SignInRequireAgreement, res.SignInStatus)
		})
	})
}

func TestAgreements(t *testing.T) {
	initFunc := func(t *testing.T) (context.Context, service.UserService, models.AgreementSlice) {
		//Given
		ctx := context.Background()
		db := tinit.DB(t)
		userSvc := service.NewUserService(db, jwtutils.NewJwtUtils([]byte("secretKey")))
		var requiredAgreements models.AgreementSlice = newRequiredAgreements()
		for _, agreement := range requiredAgreements {
			err := agreement.Insert(ctx, db, boil.Infer())
			require.NoError(t, err)
		}
		return ctx, userSvc, requiredAgreements
	}

	t.Run("return all required agreements if all required agreements are not checked", func(t *testing.T) {
		//Given
		ctx, userSvc, requiredAgreements := initFunc(t)

		userSvc.SignUp(ctx, &dto.SignUpRequest{
			Username:   "testUsername",
			Agreements: []*dto.UserAgreementReq{}, //empty
		}, &ooauth.UserInfo{
			AuthorizedBy: domain.GOOGLE,
			AuthorizedID: "authId",
		})

		//When
		agreementsRes, err := userSvc.RequiredAgreements(ctx, domain.GOOGLE, "authId")
		require.NoError(t, err)

		//Then
		require.Len(t, agreementsRes.Agreements, 2)
		require.Equal(t, requiredAgreements[0].AgreementID, agreementsRes.Agreements[0].AgreementId)
		require.Equal(t, requiredAgreements[1].AgreementID, agreementsRes.Agreements[1].AgreementId)
	})

	t.Run("return all required agreements if all required agreements are not agreed", func(t *testing.T) {
		//Given
		ctx, userSvc, requiredAgreements := initFunc(t)

		userSvc.SignUp(ctx, &dto.SignUpRequest{
			Username: "testUsername",
			Agreements: []*dto.UserAgreementReq{
				{AgreementId: requiredAgreements[0].AgreementID, IsAgree: false},
				{AgreementId: requiredAgreements[1].AgreementID, IsAgree: false},
			},
		}, &ooauth.UserInfo{
			AuthorizedBy: domain.GOOGLE,
			AuthorizedID: "authId",
		})

		//When
		agreementsRes, err := userSvc.RequiredAgreements(ctx, domain.GOOGLE, "authId")
		require.NoError(t, err)

		//Then
		require.Len(t, agreementsRes.Agreements, 2)
		require.Equal(t, requiredAgreements[0].AgreementID, agreementsRes.Agreements[0].AgreementId)
		require.Equal(t, requiredAgreements[1].AgreementID, agreementsRes.Agreements[1].AgreementId)
	})

	t.Run("return one agreement that is not agreed", func(t *testing.T) {
		//Given
		ctx, userSvc, requiredAgreements := initFunc(t)

		userSvc.SignUp(ctx, &dto.SignUpRequest{
			Username: "testUsername",
			Agreements: []*dto.UserAgreementReq{
				{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true},
				{AgreementId: requiredAgreements[1].AgreementID, IsAgree: false},
			},
		}, &ooauth.UserInfo{
			AuthorizedBy: domain.GOOGLE,
			AuthorizedID: "authId",
		})

		//When
		agreementsRes, err := userSvc.RequiredAgreements(ctx, domain.GOOGLE, "authId")
		require.NoError(t, err)

		//Then
		require.Len(t, agreementsRes.Agreements, 1)
		require.Equal(t, requiredAgreements[1].AgreementID, agreementsRes.Agreements[0].AgreementId)
	})

	t.Run("return empty agreements if all required agreements are agreed", func(t *testing.T) {
		//Given
		ctx := context.Background()
		db := tinit.DB(t)
		userSvc := service.NewUserService(db, jwtutils.NewJwtUtils([]byte("secretKey")))
		var requiredAgreements models.AgreementSlice = newRequiredAgreements()
		for _, agreement := range requiredAgreements {
			err := agreement.Insert(ctx, db, boil.Infer())
			require.NoError(t, err)
		}

		userSvc.SignUp(ctx, &dto.SignUpRequest{
			Username: "testUsername",
			Agreements: []*dto.UserAgreementReq{
				{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true},
				{AgreementId: requiredAgreements[1].AgreementID, IsAgree: true},
			},
		}, &ooauth.UserInfo{
			AuthorizedBy: domain.GOOGLE,
			AuthorizedID: "authId",
		})

		//When
		agreementsRes, err := userSvc.RequiredAgreements(ctx, domain.GOOGLE, "authId")
		require.NoError(t, err)

		//Then
		require.Len(t, agreementsRes.Agreements, 0)
	})

}

func newRequiredAgreements() []*models.Agreement {
	return []*models.Agreement{
		{AgreementCode: "code1", Summary: "summary1", IsRequired: 1},
		{AgreementCode: "code2", Summary: "summary2", IsRequired: 1},
		{AgreementCode: "code3", Summary: "summary3", IsRequired: 0},
	}
}
