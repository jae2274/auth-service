package service

import (
	"context"
	"testing"
	"userService/test/tinit"
	"userService/usersvc/common/domain"
	"userService/usersvc/models"
	"userService/usersvc/restapi/ctrlr/dto"
	"userService/usersvc/restapi/jwtutils"
	"userService/usersvc/restapi/service"

	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func initAgreementFunc(t *testing.T) (context.Context, service.UserService, models.AgreementSlice, models.AgreementSlice) {
	//Given
	ctx := context.Background()
	db := tinit.DB(t)
	userSvc := service.NewUserService(db, jwtutils.NewJwtUtils([]byte("secretKey")))
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

	return ctx, userSvc, requiredAgreements, optionalAgreements
}

func TestUsers(t *testing.T) {
	t.Run("return new_user", func(t *testing.T) {
		t.Run("when agreements empty", func(t *testing.T) {
			ctx := context.Background()
			userSvc := service.NewUserService(tinit.DB(t), jwtutils.NewJwtUtils([]byte("secretKey")))

			res, err := userSvc.SignIn(ctx, domain.GOOGLE, "authId", []*dto.UserAgreementReq{})
			require.NoError(t, err)

			require.Equal(t, dto.SignInNewUser, res.SignInStatus)
			require.Len(t, res.NewUserRes.Agreements, 0)
		})

		t.Run("when agreements existed", func(t *testing.T) {
			ctx, userSvc, requireAgreement, optionalAgreement := initAgreementFunc(t)
			savedAgreements := append(requireAgreement, optionalAgreement...)

			res, err := userSvc.SignIn(ctx, domain.GOOGLE, "authId", []*dto.UserAgreementReq{})
			require.NoError(t, err)

			require.Equal(t, dto.SignInNewUser, res.SignInStatus)
			require.Len(t, res.NewUserRes.Agreements, len(savedAgreements))

			for i, agreement := range res.NewUserRes.Agreements {
				require.Equal(t, savedAgreements[i].AgreementID, agreement.AgreementId)
				require.Equal(t, savedAgreements[i].Summary, agreement.Summary)
			}
		})
	})

	t.Run("return success", func(t *testing.T) {
		t.Run("if required agreements empty", func(t *testing.T) {
			ctx := context.Background()
			userSvc := service.NewUserService(tinit.DB(t), jwtutils.NewJwtUtils([]byte("secretKey")))

			userSvc.SignUp(ctx,
				"testUsername",
				[]*dto.UserAgreementReq{},
				domain.GOOGLE,
				"authId",
				"email")

			res, err := userSvc.SignIn(ctx, domain.GOOGLE, "authId", []*dto.UserAgreementReq{})
			require.NoError(t, err)

			require.Equal(t, dto.SignInSuccess, res.SignInStatus)
		})

		t.Run("if all required agreements are agreed", func(t *testing.T) {
			ctx, userSvc, requiredAgreements, _ := initAgreementFunc(t)

			userSvc.SignUp(ctx,
				"testUsername",
				[]*dto.UserAgreementReq{
					{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true},
					{AgreementId: requiredAgreements[1].AgreementID, IsAgree: true},
				},
				domain.GOOGLE,
				"authId",
				"email")

			res, err := userSvc.SignIn(ctx, domain.GOOGLE, "authId", []*dto.UserAgreementReq{})
			require.NoError(t, err)

			require.Equal(t, dto.SignInSuccess, res.SignInStatus)
		})

		t.Run("if one required agreement is not checked, but agreed when sign in ", func(t *testing.T) {
			ctx, userSvc, requiredAgreements, _ := initAgreementFunc(t)

			userSvc.SignUp(ctx,
				"testUsername",
				[]*dto.UserAgreementReq{
					{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true},
				},
				domain.GOOGLE,
				"authId",
				"email")

			res, err := userSvc.SignIn(ctx, domain.GOOGLE, "authId", []*dto.UserAgreementReq{
				{AgreementId: requiredAgreements[1].AgreementID, IsAgree: true},
			})
			require.NoError(t, err)

			require.Equal(t, dto.SignInSuccess, res.SignInStatus)
		})
		t.Run("if one required agreement is not agreed, but agreed when sign in", func(t *testing.T) {
			ctx, userSvc, requiredAgreements, _ := initAgreementFunc(t)

			userSvc.SignUp(ctx,
				"testUsername",
				[]*dto.UserAgreementReq{
					{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true},
					{AgreementId: requiredAgreements[1].AgreementID, IsAgree: false},
				},
				domain.GOOGLE,
				"authId",
				"email")

			res, err := userSvc.SignIn(ctx, domain.GOOGLE, "authId", []*dto.UserAgreementReq{
				{AgreementId: requiredAgreements[1].AgreementID, IsAgree: true},
			})
			require.NoError(t, err)

			require.Equal(t, dto.SignInSuccess, res.SignInStatus)
		})
	})

	t.Run("return necessary_agreements ", func(t *testing.T) {
		t.Run("if all required agreements are not checked", func(t *testing.T) {
			ctx, userSvc, requiredAgreements, _ := initAgreementFunc(t)

			userSvc.SignUp(ctx,
				"testUsername",
				[]*dto.UserAgreementReq{},
				domain.GOOGLE,
				"authId",
				"email")

			res, err := userSvc.SignIn(ctx, domain.GOOGLE, "authId", []*dto.UserAgreementReq{})
			require.NoError(t, err)

			require.Equal(t, dto.SignInNecessaryAgreements, res.SignInStatus)
			require.Len(t, res.NecessaryAgreementsRes.Agreements, 2)
			require.Equal(t, requiredAgreements[0].AgreementID, res.NecessaryAgreementsRes.Agreements[0].AgreementId)
			require.Equal(t, requiredAgreements[1].AgreementID, res.NecessaryAgreementsRes.Agreements[1].AgreementId)
		})
		t.Run("if all required agreements are not agreed", func(t *testing.T) {
			ctx, userSvc, requiredAgreements, _ := initAgreementFunc(t)

			userSvc.SignUp(ctx,
				"testUsername",
				[]*dto.UserAgreementReq{
					{AgreementId: requiredAgreements[0].AgreementID, IsAgree: false},
					{AgreementId: requiredAgreements[1].AgreementID, IsAgree: false},
				},
				domain.GOOGLE,
				"authId",
				"email")

			res, err := userSvc.SignIn(ctx, domain.GOOGLE, "authId", []*dto.UserAgreementReq{})
			require.NoError(t, err)

			require.Equal(t, dto.SignInNecessaryAgreements, res.SignInStatus)
			require.Len(t, res.NecessaryAgreementsRes.Agreements, 2)
			require.Equal(t, requiredAgreements[0].AgreementID, res.NecessaryAgreementsRes.Agreements[0].AgreementId)
			require.Equal(t, requiredAgreements[1].AgreementID, res.NecessaryAgreementsRes.Agreements[1].AgreementId)
		})

		t.Run("if one required agreement is not agreed", func(t *testing.T) {
			ctx, userSvc, requiredAgreements, _ := initAgreementFunc(t)

			userSvc.SignUp(ctx,
				"testUsername",
				[]*dto.UserAgreementReq{
					{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true},
					{AgreementId: requiredAgreements[1].AgreementID, IsAgree: false},
				},
				domain.GOOGLE,
				"authId",
				"email")

			res, err := userSvc.SignIn(ctx, domain.GOOGLE, "authId", []*dto.UserAgreementReq{})
			require.NoError(t, err)

			require.Equal(t, dto.SignInNecessaryAgreements, res.SignInStatus)
			require.Len(t, res.NecessaryAgreementsRes.Agreements, 1)
			require.Equal(t, requiredAgreements[1].AgreementID, res.NecessaryAgreementsRes.Agreements[0].AgreementId)
		})

		t.Run("if all required agreements are not checked, and one agreement agreed when sign in", func(t *testing.T) {
			ctx, userSvc, requiredAgreements, _ := initAgreementFunc(t)

			userSvc.SignUp(ctx,
				"testUsername",
				[]*dto.UserAgreementReq{},
				domain.GOOGLE,
				"authId",
				"email")

			res, err := userSvc.SignIn(ctx, domain.GOOGLE, "authId", []*dto.UserAgreementReq{
				{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true},
			})
			require.NoError(t, err)

			require.Equal(t, dto.SignInNecessaryAgreements, res.SignInStatus)
			require.Len(t, res.NecessaryAgreementsRes.Agreements, 1)
			require.Equal(t, requiredAgreements[1].AgreementID, res.NecessaryAgreementsRes.Agreements[0].AgreementId)
		})
		t.Run("if all required agreements are not agreed, and one agreement agreed when sign in", func(t *testing.T) {
			ctx, userSvc, requiredAgreements, _ := initAgreementFunc(t)

			userSvc.SignUp(ctx,
				"testUsername",
				[]*dto.UserAgreementReq{
					{AgreementId: requiredAgreements[0].AgreementID, IsAgree: false},
					{AgreementId: requiredAgreements[1].AgreementID, IsAgree: false},
				},
				domain.GOOGLE,
				"authId",
				"email")

			res, err := userSvc.SignIn(ctx, domain.GOOGLE, "authId", []*dto.UserAgreementReq{
				{AgreementId: requiredAgreements[0].AgreementID, IsAgree: true},
			})
			require.NoError(t, err)

			require.Equal(t, dto.SignInNecessaryAgreements, res.SignInStatus)
			require.Len(t, res.NecessaryAgreementsRes.Agreements, 1)
			require.Equal(t, requiredAgreements[1].AgreementID, res.NecessaryAgreementsRes.Agreements[0].AgreementId)
		})
	})
}

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
