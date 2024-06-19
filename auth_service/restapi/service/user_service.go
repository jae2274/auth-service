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
	"github.com/jae2274/goutils/terr"
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
	FindUserAuthorities(ctx context.Context, userId int) ([]*domain.UserAuthority, error)
	AddUserAuthorities(ctx context.Context, userId int, authorities []*dto.UserAuthorityReq) error
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

func (u *UserServiceImpl) FindUserAuthorities(ctx context.Context, userId int) ([]*domain.UserAuthority, error) {
	userAuthorities, err := models.UserAuthorities(models.UserAuthorityWhere.UserID.EQ(userId),
		qm.Expr(models.UserAuthorityWhere.ExpiryDate.GTE(null.NewTime(time.Now(), true)), qm.Or2(models.UserAuthorityWhere.ExpiryDate.IsNull())),
		qm.OrderBy(models.UserAuthorityColumns.CreatedDate),
		qm.Load(models.UserAuthorityRels.Authority),
	).All(ctx, u.mysqlDB)

	if err != nil {
		return nil, err
	}

	authorities := make([]*domain.UserAuthority, len(userAuthorities))
	for i, userAuthority := range userAuthorities {
		var expiryDate *time.Time = nil
		if userAuthority.ExpiryDate.Valid {
			expiryDate = &userAuthority.ExpiryDate.Time
		}

		authorities[i] = &domain.UserAuthority{
			UserID:        userAuthority.UserID,
			AuthorityID:   userAuthority.R.Authority.AuthorityID,
			AuthorityName: userAuthority.R.Authority.AuthorityName,
			ExpiryDate:    expiryDate,
		}
	}

	return authorities, nil
}

func (u *UserServiceImpl) AddUserAuthorities(ctx context.Context, userId int, dUserAuthorities []*dto.UserAuthorityReq) error {
	err := u.attachAuthorityIds(ctx, userId, dUserAuthorities)
	if err != nil {
		return err
	}

	tx, err := u.mysqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = addUserAuthorities(ctx, tx, userId, dUserAuthorities)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (u *UserServiceImpl) attachAuthorityIds(ctx context.Context, userId int, dUserAuthorities []*dto.UserAuthorityReq) error {
	authorityNames := make([]string, len(dUserAuthorities))
	for i, authority := range dUserAuthorities {
		authorityNames[i] = authority.AuthorityName
	}

	authorities, err := models.Authorities(models.AuthorityWhere.AuthorityName.IN(authorityNames)).All(ctx, u.mysqlDB)
	if err != nil {
		return err
	}
	authorityMap := make(map[string]*models.Authority)
	for _, authority := range authorities {
		authorityMap[authority.AuthorityName] = authority
	}

	for _, userAuthority := range dUserAuthorities {
		authority, ok := authorityMap[userAuthority.AuthorityName]
		if !ok {
			return terr.New("authority not found. authorityName: " + userAuthority.AuthorityName)
		}

		userAuthority.AuthorityID = authority.AuthorityID
	}

	return nil
}

func addUserAuthorities(ctx context.Context, tx *sql.Tx, userId int, dUserAuthorities []*dto.UserAuthorityReq) error {

	now := time.Now()

	// mUserAuthorities := make([]*models.UserAuthority, len(domainUserAuthorities))
	for _, dUserAuthority := range dUserAuthorities {
		mUserAuthority := &models.UserAuthority{
			UserID:      userId,
			AuthorityID: dUserAuthority.AuthorityID,
		}
		err := mUserAuthority.Reload(ctx, tx)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		addExpiryDate := func(date time.Time, duration *time.Duration) null.Time {
			if duration != nil {
				return null.NewTime(date.Add(*duration), true)
			} else {
				return null.NewTime(time.Time{}, false)
			}
		}

		if err == sql.ErrNoRows {
			mUserAuthority.ExpiryDate = addExpiryDate(now, dUserAuthority.ExpiryDuration)

			if err := mUserAuthority.Insert(ctx, tx, boil.Infer()); err != nil {
				return err
			}
		} else {
			if mUserAuthority.ExpiryDate.Valid { //false의 경우, 만료되지 않는 권한으로 간주하여 만료일을 갱신하지 않음
				if mUserAuthority.ExpiryDate.Time.Before(now) {
					mUserAuthority.ExpiryDate = null.NewTime(now, true)
				}

				mUserAuthority.ExpiryDate = addExpiryDate(mUserAuthority.ExpiryDate.Time, dUserAuthority.ExpiryDuration)

				if _, err := mUserAuthority.Update(ctx, tx, boil.Infer()); err != nil {
					return err
				}
			}
		}

		// mUserAuthorities[i] = modelUserAuthority
	}

	return nil
}

func (u *UserServiceImpl) FindAllAgreements(ctx context.Context) ([]*models.Agreement, error) {
	return models.Agreements().All(ctx, u.mysqlDB)
}
