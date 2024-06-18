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
	FindAllAgreements(ctx context.Context) ([]*models.Agreement, error) //TODO: 테스트코드 작성
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

	userRoles, err := addUserRoles(ctx, tx, userId, roles)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return userRoles, nil
}

func addUserRoles(ctx context.Context, tx *sql.Tx, userId int, roles []*domain.UserRole) ([]*models.UserRole, error) {
	now := time.Now()
	mRoles := make([]*models.UserRole, len(roles))
	for i, role := range roles {
		userRole := &models.UserRole{
			UserID:   userId,
			RoleName: role.RoleName,
		}
		err := userRole.Reload(ctx, tx)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		addExpiryDate := func(date time.Time, duration *time.Duration) null.Time {
			if duration != nil {
				return null.NewTime(date.Add(*duration), true)
			} else {
				return null.NewTime(time.Time{}, false)
			}
		}

		if err == sql.ErrNoRows {
			userRole.ExpiryDate = addExpiryDate(now, role.ExpiryDuration)

			if err := userRole.Insert(ctx, tx, boil.Infer()); err != nil {
				return nil, err
			}
		} else {
			if userRole.ExpiryDate.Valid { //false의 경우, 만료되지 않는 권한으로 간주하여 만료일을 갱신하지 않음
				if userRole.ExpiryDate.Time.Before(now) {
					userRole.ExpiryDate = null.NewTime(now, true)
				}

				userRole.ExpiryDate = addExpiryDate(userRole.ExpiryDate.Time, role.ExpiryDuration)

				if _, err := userRole.Update(ctx, tx, boil.Infer()); err != nil {
					return nil, err
				}
			}
		}

		mRoles[i] = userRole
	}

	return mRoles, nil
}

func (u *UserServiceImpl) FindAllAgreements(ctx context.Context) ([]*models.Agreement, error) {
	return models.Agreements().All(ctx, u.mysqlDB)
}
