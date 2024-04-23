package service

import (
	"context"
	"database/sql"
	"strconv"
	"time"
	"userService/usersvc/common/domain"
	"userService/usersvc/models"
	"userService/usersvc/restapi/ctrlr/dto"
	"userService/usersvc/restapi/jwtutils"
	"userService/usersvc/restapi/mapper"
	"userService/usersvc/restapi/ooauth"
	"userService/usersvc/utils"
)

type UserService interface {
	SignIn(ctx context.Context, authorizedBy domain.AuthorizedBy, authorizedId string) (*dto.SignInResponse, error)
	SignUp(ctx context.Context, req *dto.SignUpRequest, userinfo *ooauth.UserInfo) error
}

type UserServiceImpl struct {
	mysqlDB     *sql.DB
	jwtResolver *jwtutils.JwtResolver
}

func NewUserService(mysqlDB *sql.DB, jwtResolver *jwtutils.JwtResolver) UserService {
	return &UserServiceImpl{
		mysqlDB:     mysqlDB,
		jwtResolver: jwtResolver,
	}
}

func (u *UserServiceImpl) SignIn(ctx context.Context, authorizedBy domain.AuthorizedBy, authorizedId string) (*dto.SignInResponse, error) {

	user, isExisted, err := mapper.FindUserByAuthorized(ctx, u.mysqlDB, authorizedBy, authorizedId)
	if err != nil {
		return nil, err
	}

	if isExisted {
		userAgreeds, err := mapper.FindUserAgreements(ctx, u.mysqlDB, user.UserID, true)
		if err != nil {
			return nil, err
		}
		requiredAgreements, err := mapper.FindRequiredAgreements(ctx, u.mysqlDB)
		if err != nil {
			return nil, err
		}

		userAgreedsMap := make(map[int]bool)
		for _, userAgreed := range userAgreeds {
			userAgreedsMap[userAgreed.AgreementID] = true
		}

		needAgreements := []*models.Agreement{}
		for _, requiredAgreement := range requiredAgreements {
			if _, ok := userAgreedsMap[requiredAgreement.AgreementID]; !ok {
				needAgreements = append(needAgreements, requiredAgreement)
			}
		}

		if len(needAgreements) > 0 {
			return u.signInRequireAgreement(ctx, needAgreements)
		}

		return u.signInSuccess(ctx, user)
	} else {
		return u.signInNewUser(ctx)
	}
}

func (u *UserServiceImpl) signInSuccess(ctx context.Context, user *models.User) (*dto.SignInResponse, error) {
	roleNames := make([]string, 0)
	for _, userRole := range user.R.GetUserRoles() {
		if userRole.ExpiryDate.IsZero() || userRole.ExpiryDate.Time.After(time.Now()) { //TODO: 임의로 애플리케이션에서 필터링, 추후 DB에서 필터링으로 변경
			roleNames = append(roleNames, userRole.RoleName)
		}
	}

	jwtToken, err := u.jwtResolver.CreateToken(strconv.Itoa(user.UserID), roleNames)
	if err != nil {
		return nil, err
	}

	return &dto.SignInResponse{
		SignInStatus: dto.SignInSuccess,
		SuccessRes: &dto.SignInSuccessRes{
			Username:     user.Name,
			Roles:        roleNames,
			GrantType:    jwtToken.GrantType,
			AccessToken:  jwtToken.AccessToken,
			RefreshToken: jwtToken.RefreshToken,
		},
	}, nil
}

// TODO: 나머지 구현 및 테스트코드 작성
func (u *UserServiceImpl) signInRequireAgreement(ctx context.Context, requiredAgreements []*models.Agreement) (*dto.SignInResponse, error) {
	agreementRes := make([]*dto.AgreementRes, len(requiredAgreements))
	for i, agreement := range requiredAgreements {
		agreementRes[i] = &dto.AgreementRes{
			AgreementId: agreement.AgreementID,
			Required:    utils.TinyIntToBool(agreement.IsRequired),
			Summary:     agreement.Summary,
		}
	}

	return &dto.SignInResponse{
		SignInStatus: dto.SignInRequireAgreement,
		RequireAgreementRes: &dto.RequireAgreementRes{
			Agreements: agreementRes,
		},
	}, nil
}

func (u *UserServiceImpl) signInNewUser(ctx context.Context) (*dto.SignInResponse, error) {
	agreements, err := mapper.FindAllAgreement(ctx, u.mysqlDB)
	if err != nil {
		return nil, err
	}

	agreementRes := make([]*dto.AgreementRes, len(agreements))
	for i, agreement := range agreements {
		agreementRes[i] = &dto.AgreementRes{
			AgreementId: agreement.AgreementID,
			Required:    utils.TinyIntToBool(agreement.IsRequired),
			Summary:     agreement.Summary,
		}
	}

	return &dto.SignInResponse{
		SignInStatus: dto.SignInNewUser,
		NewUserRes: &dto.SignInNewUserRes{
			Agreements: agreementRes,
		},
	}, nil
}

func saveUser(ctx context.Context, tx *sql.Tx, req *dto.SignUpRequest, userinfo *ooauth.UserInfo) (err error) {
	user, err := mapper.SaveUser(ctx, tx, userinfo.AuthorizedBy, userinfo.AuthorizedID, userinfo.Email, req.Username)
	if err != nil {
		return err
	}

	mAgs := make([]*models.UserAgreement, len(req.Agreements))
	for i, ag := range req.Agreements {
		mAgs[i] = &models.UserAgreement{
			UserID:      user.UserID,
			AgreementID: ag.AgreementId,
			IsAgree:     utils.BoolToTinyInt(ag.IsAgree),
		}
	}

	return user.AddUserAgreements(ctx, tx, true, mAgs...)
}

func (u *UserServiceImpl) SignUp(ctx context.Context, req *dto.SignUpRequest, userinfo *ooauth.UserInfo) error {

	tx, err := u.mysqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := saveUser(ctx, tx, req, userinfo); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
