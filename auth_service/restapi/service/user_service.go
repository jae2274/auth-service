package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/jae2274/auth-service/auth_service/common/domain"
	"github.com/jae2274/auth-service/auth_service/models"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
	"github.com/jae2274/auth-service/auth_service/restapi/mapper"
	"github.com/jae2274/auth-service/auth_service/restapi/ooauth"
	"github.com/jae2274/auth-service/auth_service/utils"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type UserService interface {
	// SignIn(ctx context.Context, userinfo *ooauth.UserInfo, addAgreements []*dto.UserAgreementReq) (*dto.SignInResponse, error)
	SignUp(ctx context.Context, userinfo *ooauth.UserInfo, additionalAgreements []*dto.UserAgreementReq) (*models.User, error)
	FindSignedUpUser(ctx context.Context, authorizedBy domain.AuthorizedBy, authorizedID string) (*models.User, bool, error)
	ApplyUserAgreements(ctx context.Context, userId int, agreements []*dto.UserAgreementReq) error
	FindNecessaryAgreements(ctx context.Context, userId int) ([]*models.Agreement, error)
	FindUserRoles(ctx context.Context, userId int) ([]*models.UserRole, error)
	AddUserRoles(ctx context.Context, userId int, roles []*domain.UserRole) ([]*models.UserRole, error)
	// RefreshJwt(ctx context.Context, refreshToken string) (dto.RefreshJwtResponse, bool, error)
}

type UserServiceImpl struct {
	mysqlDB *sql.DB
}

func NewUserService(mysqlDB *sql.DB) UserService {
	return &UserServiceImpl{
		mysqlDB: mysqlDB,
	}
}

func (u *UserServiceImpl) SignUp(ctx context.Context, userinfo *ooauth.UserInfo, agreements []*dto.UserAgreementReq) (*models.User, error) {
	tx, err := u.mysqlDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	user, err := signUp(ctx, tx, userinfo, agreements)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func signUp(ctx context.Context, tx *sql.Tx, userinfo *ooauth.UserInfo, agreements []*dto.UserAgreementReq) (*models.User, error) {
	user, err := mapper.SaveUser(ctx, tx, userinfo.AuthorizedBy, userinfo.AuthorizedID, userinfo.Email, userinfo.Username)
	if err != nil {
		return nil, err
	}

	mAgs := make([]*models.UserAgreement, len(agreements))
	for i, ag := range agreements {
		mAgs[i] = &models.UserAgreement{
			UserID:      user.UserID,
			AgreementID: ag.AgreementId,
			IsAgree:     utils.BoolToTinyInt(ag.IsAgree),
		}
	}

	if err := user.AddUserAgreements(ctx, tx, true, mAgs...); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserServiceImpl) FindSignedUpUser(ctx context.Context, authorizedBy domain.AuthorizedBy, authorizedID string) (*models.User, bool, error) {
	return mapper.FindUserByAuthorized(ctx, u.mysqlDB, authorizedBy, authorizedID)
}

func (u *UserServiceImpl) ApplyUserAgreements(ctx context.Context, userId int, agreements []*dto.UserAgreementReq) error {

	tx, err := u.mysqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := applyUserAgreements(ctx, tx, userId, agreements); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func applyUserAgreements(ctx context.Context, tx *sql.Tx, userId int, agreements []*dto.UserAgreementReq) error {
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

func (u *UserServiceImpl) FindNecessaryAgreements(ctx context.Context, userId int) ([]*models.Agreement, error) {
	return u.necessaryAgreements(ctx, userId)
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

func (u *UserServiceImpl) FindUserRoles(ctx context.Context, userId int) ([]*models.UserRole, error) {
	return models.UserRoles(models.UserRoleWhere.UserID.EQ(userId),
		qm.Expr(models.UserRoleWhere.ExpiryDate.GTE(null.NewTime(time.Now(), true)), qm.Or2(models.UserRoleWhere.ExpiryDate.IsNull())),
		qm.OrderBy(models.UserRoleColumns.CreatedDate)).All(ctx, u.mysqlDB)
	// return nil, nil
}

func (u *UserServiceImpl) AddUserRoles(ctx context.Context, userId int, roles []*domain.UserRole) ([]*models.UserRole, error) {
	tx, err := u.mysqlDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	mRoles := make([]*models.UserRole, len(roles))
	for i, role := range roles {
		expiryDate := null.NewTime(time.Time{}, false)

		if role.ExpiryDuration != nil {
			expiryDate = null.NewTime(now.Add(*role.ExpiryDuration), true)
		}

		mRoles[i] = &models.UserRole{
			UserID:     userId,
			RoleName:   role.RoleName,
			ExpiryDate: expiryDate,
		}

		if err := mRoles[i].Insert(ctx, tx, boil.Infer()); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return mRoles, nil
}

// func (u *UserServiceImpl) RefreshJwt(ctx context.Context, refreshToken string) (dto.RefreshJwtResponse, bool, error) {
// 	return dto.RefreshJwtResponse{}, false, nil
// }

// func (u *UserServiceImpl) applyUserAgreements(ctx context.Context, tx *sql.Tx, userId int, agreements []*dto.UserAgreementReq) error {
// 	for _, addAgreement := range agreements {
// 		userAgreement := models.UserAgreement{
// 			UserID:      userId,
// 			AgreementID: addAgreement.AgreementId,
// 			IsAgree:     utils.BoolToTinyInt(addAgreement.IsAgree),
// 		}
// 		isExisted, err := userAgreement.Exists(ctx, tx)
// 		if err != nil {
// 			return err
// 		}
// 		if isExisted {
// 			_, err := userAgreement.Update(ctx, tx, boil.Infer())
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			err := userAgreement.Insert(ctx, tx, boil.Infer())
// 			if err != nil {
// 				return err
// 			}

// 		}
// 	}

// 	return nil
// }

// func (u *UserServiceImpl) SignIn(ctx context.Context, userinfo *ooauth.UserInfo, additionalAgreements []*dto.UserAgreementReq) (*dto.SignInResponse, error) {

// 	user, isExisted, err := mapper.FindUserByAuthorized(ctx, u.mysqlDB, userinfo.AuthorizedBy, userinfo.AuthorizedID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if isExisted {
// 		tx, err := u.mysqlDB.BeginTx(ctx, nil)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if err := u.applyUserAgreements(ctx, tx, user.UserID, additionalAgreements); err != nil {
// 			tx.Rollback()
// 			return nil, err
// 		}
// 		if err := tx.Commit(); err != nil {
// 			return nil, err
// 		}

// 		nAgreements, err := u.necessaryAgreements(ctx, user.UserID)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if len(nAgreements) > 0 {
// 			return u.signInRequireAgreement(ctx, nAgreements)
// 		}

// 		return u.signInSuccess(ctx, user)
// 	} else {
// 		return u.signInNewUser(ctx, userinfo)
// 	}
// }

// func (u *UserServiceImpl) signInSuccess(ctx context.Context, user *models.User) (*dto.SignInResponse, error) {
// 	roleNames := make([]string, 0)
// 	for _, userRole := range user.R.GetUserRoles() {
// 		if userRole.ExpiryDate.IsZero() || userRole.ExpiryDate.Time.After(time.Now()) { //TODO: 임의로 애플리케이션에서 필터링, 추후 DB에서 필터링으로 변경
// 			roleNames = append(roleNames, userRole.RoleName)
// 		}
// 	}

// 	jwtToken, err := u.jwtResolver.CreateToken(strconv.Itoa(user.UserID), roleNames, time.Now())
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &dto.SignInResponse{
// 		SignInStatus: dto.SignInSuccess,
// 		SuccessRes: &dto.SignInSuccessRes{
// 			Username:     user.Name,
// 			Roles:        roleNames,
// 			GrantType:    jwtToken.GrantType,
// 			AccessToken:  jwtToken.AccessToken,
// 			RefreshToken: jwtToken.RefreshToken,
// 		},
// 	}, nil
// }

// // TODO: 나머지 구현 및 테스트코드 작성
// func (u *UserServiceImpl) signInRequireAgreement(ctx context.Context, requiredAgreements []*models.Agreement) (*dto.SignInResponse, error) {
// 	agreementRes := make([]*dto.AgreementRes, len(requiredAgreements))
// 	for i, agreement := range requiredAgreements {
// 		agreementRes[i] = &dto.AgreementRes{
// 			AgreementId: agreement.AgreementID,
// 			Required:    utils.TinyIntToBool(agreement.IsRequired),
// 			Summary:     agreement.Summary,
// 		}
// 	}

// 	return &dto.SignInResponse{
// 		SignInStatus: dto.SignInNecessaryAgreements,
// 		NecessaryAgreementsRes: &dto.SignInNecessaryAgreementsRes{
// 			Agreements: agreementRes,
// 		},
// 	}, nil
// }

// func (u *UserServiceImpl) signInNewUser(ctx context.Context, userinfo *ooauth.UserInfo) (*dto.SignInResponse, error) {
// 	agreements, err := mapper.FindAllAgreement(ctx, u.mysqlDB)
// 	if err != nil {
// 		return nil, err
// 	}

// 	agreementRes := make([]*dto.AgreementRes, len(agreements))
// 	for i, agreement := range agreements {
// 		agreementRes[i] = &dto.AgreementRes{
// 			AgreementId: agreement.AgreementID,
// 			Required:    utils.TinyIntToBool(agreement.IsRequired),
// 			Summary:     agreement.Summary,
// 		}
// 	}

// 	return &dto.SignInResponse{
// 		SignInStatus: dto.SignInNewUser,
// 		NewUserRes: &dto.SignInNewUserRes{
// 			Email:      userinfo.Email,
// 			Username:   userinfo.Username,
// 			Agreements: agreementRes,
// 		},
// 	}, nil
// }

// func (u *UserServiceImpl) necessaryAgreements(ctx context.Context, userId int) ([]*models.Agreement, error) {

// 	userAgreeds, err := mapper.FindUserAgreements(ctx, u.mysqlDB, userId, true)
// 	if err != nil {
// 		return nil, err
// 	}

// 	requiredAgreements, err := mapper.FindRequiredAgreements(ctx, u.mysqlDB)
// 	if err != nil {
// 		return nil, err
// 	}

// 	userAgreedsMap := make(map[int]bool)
// 	for _, userAgreed := range userAgreeds {
// 		userAgreedsMap[userAgreed.AgreementID] = true
// 	}

// 	var necessaryAgreements []*models.Agreement
// 	for _, requiredAgreement := range requiredAgreements {
// 		if _, ok := userAgreedsMap[requiredAgreement.AgreementID]; !ok {
// 			necessaryAgreements = append(necessaryAgreements, requiredAgreement)
// 		}
// 	}

// 	return necessaryAgreements, nil
// }

// func (u *UserServiceImpl) SignUp(ctx context.Context, userinfo *ooauth.UserInfo, agreements []*dto.UserAgreementReq) error {

// 	tx, err := u.mysqlDB.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}

// 	if err := signUp(ctx, tx, userinfo, agreements); err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	if err := tx.Commit(); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func signUp(ctx context.Context, tx *sql.Tx, userinfo *ooauth.UserInfo, agreements []*dto.UserAgreementReq) (err error) {
// 	user, err := mapper.SaveUser(ctx, tx, userinfo.AuthorizedBy, userinfo.AuthorizedID, userinfo.Email, userinfo.Username)
// 	if err != nil {
// 		return err
// 	}

// 	mAgs := make([]*models.UserAgreement, len(agreements))
// 	for i, ag := range agreements {
// 		mAgs[i] = &models.UserAgreement{
// 			UserID:      user.UserID,
// 			AgreementID: ag.AgreementId,
// 			IsAgree:     utils.BoolToTinyInt(ag.IsAgree),
// 		}
// 	}

// 	return user.AddUserAgreements(ctx, tx, true, mAgs...)
// }

// func (u *UserServiceImpl) RefreshJwt(ctx context.Context, refreshToken string) (dto.RefreshJwtResponse, bool, error) {
// 	// claims, isValid,err := u.jwtResolver.ParseToken(refreshToken)
// 	// if !isValid {
// 	// }

// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// userId := claims.UserId
// 	// roles := claims.Roles

// 	// jwtToken, err := u.jwtResolver.CreateToken(userId, roles)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// return &dto.TokenInfo{
// 	// 	GrantType:    "Bearer",
// 	// 	AccessToken:  jwtToken.AccessToken,
// 	// 	RefreshToken: jwtToken.RefreshToken,
// 	// }, nil

// 	return dto.RefreshJwtResponse{}, false, nil
// }
