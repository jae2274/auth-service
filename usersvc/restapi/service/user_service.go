package service

import (
	"context"
	"database/sql"
	"strconv"
	"time"
	"userService/usersvc/common/domain"
	"userService/usersvc/models"
	"userService/usersvc/restapi/aescryptor"
	"userService/usersvc/restapi/ctrlr/dto"
	"userService/usersvc/restapi/jwtutils"
	"userService/usersvc/restapi/mapper"
	"userService/usersvc/restapi/ooauth"
	"userService/usersvc/utils"
)

type UserService interface {
	AuthCodeUrls(ctx context.Context, state string) *dto.AuthCodeUrlsResponse
	Authenticate(ctx context.Context, code string) (*dto.AuthenticateResponse, error)
	SignIn(ctx context.Context, authToken string) (*dto.SignInResponse, error)
	SignUp(ctx context.Context, req *dto.SignUpRequest) error
}

type UserServiceImpl struct {
	mysqlDB     *sql.DB
	aesCryptor  *aescryptor.JsonAesCryptor
	googleOauth ooauth.Ooauth
	jwtResolver *jwtutils.JwtResolver
}

func NewUserService(mysqlDB *sql.DB, aesCryptor *aescryptor.JsonAesCryptor, googleOauth ooauth.Ooauth, jwtResolver *jwtutils.JwtResolver) UserService {
	return &UserServiceImpl{
		mysqlDB:     mysqlDB,
		aesCryptor:  aesCryptor,
		googleOauth: googleOauth,
		jwtResolver: jwtResolver,
	}
}

func (u *UserServiceImpl) AuthCodeUrls(ctx context.Context, state string) *dto.AuthCodeUrlsResponse {
	return &dto.AuthCodeUrlsResponse{
		AuthCodeUrls: []*dto.AuthCodeUrlRes{
			{AuthServer: string(u.googleOauth.GetAuthServer()), Url: u.googleOauth.GetLoginURL(state)},
		},
	}
}

func (u *UserServiceImpl) Authenticate(ctx context.Context, code string) (*dto.AuthenticateResponse, error) {
	token, err := u.googleOauth.GetToken(ctx, code)
	if err != nil {
		return nil, err
	}

	authToken, err := u.aesCryptor.Encrypt(token)
	if err != nil {
		return nil, err
	}

	return &dto.AuthenticateResponse{
		AuthToken: authToken,
	}, nil
}

func (u *UserServiceImpl) SignIn(ctx context.Context, authToken string) (*dto.SignInResponse, error) {
	token := &ooauth.OauthToken{}
	err := u.aesCryptor.Decrypt(authToken, token)
	if err != nil {
		return nil, err
	}

	userinfo, err := u.googleOauth.GetUserInfo(ctx, token)
	if err != nil {
		return nil, err
	}

	user, isExisted, err := mapper.FindUserByAuthorized(ctx, u.mysqlDB, userinfo.AuthorizedBy, userinfo.AuthorizedID)
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
		requiredAgreements = []*models.Agreement{}
		for _, userAgreed := range userAgreeds {
			userAgreedsMap[userAgreed.AgreementID] = true
		}

		for _, requiredAgreement := range requiredAgreements {
			if _, ok := userAgreedsMap[requiredAgreement.AgreementID]; !ok {
				requiredAgreements = append(requiredAgreements, requiredAgreement)
			}
		}

		if len(requiredAgreements) > 0 {
			return u.signInRequireAgreement(ctx, requiredAgreements)
		}

		return u.signInSuccess(ctx, user)
	} else {
		return u.signInNewUser(ctx, userinfo.Email)
	}
}

func (u *UserServiceImpl) signInSuccess(ctx context.Context, user *models.User) (*dto.SignInResponse, error) {
	roleNames := make([]string, 0)
	for _, userRole := range user.R.GetUserRoles() {
		if userRole.ExpiryDate.IsZero() || userRole.ExpiryDate.Time.After(time.Now()) { //TODO: 임의로 애플리케이션에서 필터링, 추후 DB에서 필터링으로 변경
			roleNames = append(roleNames, userRole.RoleName)
		}
	}

	jwtToken, err := u.jwtResolver.CreateToken(strconv.Itoa(user.UserID), user.Email, roleNames)
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
		SignInStatus: dto.SignInNewUser,
		RequireAgreementRes: &dto.RequireAgreementRes{
			Agreements: agreementRes,
		},
	}, nil
}

func (u *UserServiceImpl) signInNewUser(ctx context.Context, email string) (*dto.SignInResponse, error) {
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
			Email:      email,
			Agreements: agreementRes,
		},
	}, nil
}

func saveUser(ctx context.Context, tx *sql.Tx, authorizedBy domain.AuthorizedBy, authorizedID, email string, req *dto.SignUpRequest) (err error) {
	user, err := mapper.SaveUser(ctx, tx, authorizedBy, authorizedID, email, req.Username)
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

func (u *UserServiceImpl) SignUp(ctx context.Context, req *dto.SignUpRequest) error {
	token := &ooauth.OauthToken{}
	err := u.aesCryptor.Decrypt(req.AuthToken, token)
	if err != nil {
		return err
	}

	userinfo, err := u.googleOauth.GetUserInfo(ctx, token)
	if err != nil {
		return err
	}

	tx, err := u.mysqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := saveUser(ctx, tx, userinfo.AuthorizedBy, userinfo.AuthorizedID, userinfo.Email, req); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
