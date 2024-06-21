package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/jae2274/auth-service/auth_service/common/domain"
	"github.com/jae2274/auth-service/auth_service/common/mysqldb"
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
	SignUp(ctx context.Context, userinfo *ooauth.UserInfo, additionalAgreements []*dto.UserAgreementReq) (*models.User, error)
	FindSignedUpUser(ctx context.Context, authorizedBy domain.AuthorizedBy, authorizedID string) (*models.User, bool, error)
	ApplyUserAgreements(ctx context.Context, userId int, agreements []*dto.UserAgreementReq) error
	FindNecessaryAgreements(ctx context.Context, userId int) ([]*models.Agreement, error)
	FindUserAuthorities(ctx context.Context, userId int) ([]*domain.UserAuthority, error)
	AddUserAuthorities(ctx context.Context, userId int, authorities []*dto.UserAuthorityReq) error
	FindAllAgreements(ctx context.Context) ([]*models.Agreement, error) //TODO: 테스트코드 작성
	RemoveAuthority(ctx context.Context, userId int, authorityName string) error
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
	return mysqldb.CommitOrRollback(tx, user, err)
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

	err = applyUserAgreements(ctx, tx, userId, agreements)
	return mysqldb.CommitOrRollbackVoid(tx, err)
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

func checkHasAuthorityAdmin(userAuthorities []*dto.UserAuthorityReq) bool {
	for _, userAuthority := range userAuthorities {
		if userAuthority.AuthorityName == domain.AuthorityAdmin {
			return true
		}
	}

	return false
}

func (u *UserServiceImpl) AddUserAuthorities(ctx context.Context, userId int, dUserAuthorities []*dto.UserAuthorityReq) error {
	if checkHasAuthorityAdmin(dUserAuthorities) {
		return terr.New("cannot add authority admin")
	}

	err := attachAuthorityIds(ctx, u.mysqlDB, dUserAuthorities)
	if err != nil {
		return err
	}

	tx, err := u.mysqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = addUserAuthorities(ctx, tx, userId, dUserAuthorities)
	return mysqldb.CommitOrRollbackVoid(tx, err)
}

func (u *UserServiceImpl) FindAllAgreements(ctx context.Context) ([]*models.Agreement, error) {
	return models.Agreements().All(ctx, u.mysqlDB)
}

func (u *UserServiceImpl) RemoveAuthority(ctx context.Context, userId int, authorityName string) error {
	mAuthority, err := models.Authorities(models.AuthorityWhere.AuthorityName.EQ(authorityName)).One(ctx, u.mysqlDB)
	if err != nil && err != sql.ErrNoRows {
		return terr.Wrap(err)
	}

	if err == sql.ErrNoRows {
		return terr.New("authority not found")
	}

	_, err = models.UserAuthorities(models.UserAuthorityWhere.UserID.EQ(userId), models.UserAuthorityWhere.AuthorityID.EQ(mAuthority.AuthorityID)).DeleteAll(ctx, u.mysqlDB)

	return err

}
