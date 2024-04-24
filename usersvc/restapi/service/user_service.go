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
	"userService/usersvc/utils"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UserService interface {
	SignIn(ctx context.Context, authorizedBy domain.AuthorizedBy, authorizedId string, addAgreements []*dto.UserAgreementReq) (*dto.SignInResponse, error)
	SignUp(ctx context.Context, username string, additionalAgreements []*dto.UserAgreementReq, authBy domain.AuthorizedBy, authId string, email string) error
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

func (u *UserServiceImpl) applyUserAgreements(ctx context.Context, tx *sql.Tx, userId int, agreements []*dto.UserAgreementReq) error {
	for _, addAgreement := range agreements {
		userAgreement := models.UserAgreement{
			UserID:      userId,
			AgreementID: addAgreement.AgreementId,
			IsAgree:     utils.BoolToTinyInt(addAgreement.IsAgree),
		}
		isExisted, err := userAgreement.Exists(ctx, tx)
		if err != nil {
			return err
		}
		if isExisted {
			_, err := userAgreement.Update(ctx, tx, boil.Infer())
			if err != nil {
				return err
			}
		} else {
			err := userAgreement.Insert(ctx, tx, boil.Infer())
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func (u *UserServiceImpl) SignIn(ctx context.Context, authorizedBy domain.AuthorizedBy, authorizedId string, additionalAgreements []*dto.UserAgreementReq) (*dto.SignInResponse, error) {

	user, isExisted, err := mapper.FindUserByAuthorized(ctx, u.mysqlDB, authorizedBy, authorizedId)
	if err != nil {
		return nil, err
	}

	if isExisted {
		tx, err := u.mysqlDB.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		if err := u.applyUserAgreements(ctx, tx, user.UserID, additionalAgreements); err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Commit(); err != nil {
			return nil, err
		}

		nAgreements, err := u.necessaryAgreements(ctx, user.UserID)
		if err != nil {
			return nil, err
		}

		if len(nAgreements) > 0 {
			return u.signInRequireAgreement(ctx, nAgreements)
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
		SignInStatus: dto.SignInNecessaryAgreements,
		NecessaryAgreementsRes: &dto.SignInNecessaryAgreementsRes{
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

func (u *UserServiceImpl) necessaryAgreements(ctx context.Context, userId int) ([]*models.Agreement, error) {

	userAgreeds, err := mapper.FindUserAgreements(ctx, u.mysqlDB, userId, true)
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

	var necessaryAgreements []*models.Agreement
	for _, requiredAgreement := range requiredAgreements {
		if _, ok := userAgreedsMap[requiredAgreement.AgreementID]; !ok {
			necessaryAgreements = append(necessaryAgreements, requiredAgreement)
		}
	}

	return necessaryAgreements, nil
}

func (u *UserServiceImpl) SignUp(ctx context.Context, username string, agreements []*dto.UserAgreementReq, authBy domain.AuthorizedBy, authId string, email string) error {

	tx, err := u.mysqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := signUp(ctx, tx, username, agreements, authBy, authId, email); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func signUp(ctx context.Context, tx *sql.Tx, username string, agreements []*dto.UserAgreementReq, authBy domain.AuthorizedBy, authId string, email string) (err error) {
	user, err := mapper.SaveUser(ctx, tx, authBy, authId, email, username)
	if err != nil {
		return err
	}

	mAgs := make([]*models.UserAgreement, len(agreements))
	for i, ag := range agreements {
		mAgs[i] = &models.UserAgreement{
			UserID:      user.UserID,
			AgreementID: ag.AgreementId,
			IsAgree:     utils.BoolToTinyInt(ag.IsAgree),
		}
	}

	return user.AddUserAgreements(ctx, tx, true, mAgs...)
}
