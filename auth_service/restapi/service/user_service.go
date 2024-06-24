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

func SignUp(ctx context.Context, tx *sql.Tx, userinfo *ooauth.UserInfo, agreements []*dto.UserAgreementReq) (*models.User, error) {
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

func FindSignedUpUser(ctx context.Context, exec boil.ContextExecutor, authorizedBy domain.AuthorizedBy, authorizedID string) (*models.User, bool, error) {
	return mapper.FindUserByAuthorized(ctx, exec, authorizedBy, authorizedID)
}

func ApplyUserAgreements(ctx context.Context, tx *sql.Tx, userId int, agreements []*dto.UserAgreementReq) error {
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

func FindNecessaryAgreements(ctx context.Context, exec boil.ContextExecutor, userId int) ([]*models.Agreement, error) {
	userAgreeds, err := mapper.FindUserAgreements(ctx, exec, userId, true)
	if err != nil {
		return nil, err
	}

	requiredAgreements, err := mapper.FindRequiredAgreements(ctx, exec)
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

func FindUserAuthorities(ctx context.Context, exec boil.ContextExecutor, userId int) ([]*domain.UserAuthority, error) {
	userAuthorities, err := models.UserAuthorities(models.UserAuthorityWhere.UserID.EQ(userId),
		qm.Expr(models.UserAuthorityWhere.ExpiryDate.GTE(null.NewTime(time.Now(), true)), qm.Or2(models.UserAuthorityWhere.ExpiryDate.IsNull())),
		qm.OrderBy(models.UserAuthorityColumns.CreatedDate),
		qm.Load(models.UserAuthorityRels.Authority),
	).All(ctx, exec)

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
			AuthorityCode: userAuthority.R.Authority.AuthorityCode,
			ExpiryDate:    expiryDate,
		}
	}

	return authorities, nil
}

func FindUserAuthoritiesByAuthorityIds(ctx context.Context, exec boil.ContextExecutor, userId int, authorityId []int) ([]*domain.UserAuthority, error) {
	userAuthorities, err := models.UserAuthorities(
		models.UserAuthorityWhere.UserID.EQ(userId), models.UserAuthorityWhere.AuthorityID.IN(authorityId),
		qm.OrderBy(models.UserAuthorityColumns.CreatedDate),
		qm.Load(models.UserAuthorityRels.Authority),
	).All(ctx, exec)
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
			AuthorityCode: userAuthority.R.Authority.AuthorityCode,
			ExpiryDate:    expiryDate,
		}
	}

	return authorities, nil
}

func checkHasAuthorityAdmin(userAuthorities []*dto.UserAuthorityReq) bool {
	for _, userAuthority := range userAuthorities {
		if userAuthority.AuthorityCode == domain.AuthorityAdmin {
			return true
		}
	}

	return false
}

func AddUserAuthorities(ctx context.Context, tx *sql.Tx, userId int, dUserAuthorities []*dto.UserAuthorityReq) error {
	if checkHasAuthorityAdmin(dUserAuthorities) {
		return terr.New("cannot add authority admin")
	}

	err := attachAuthorityIds(ctx, tx, dUserAuthorities)
	if err != nil {
		return err
	}

	return addUserAuthorities(ctx, tx, userId, dUserAuthorities)
}

func FindAllAgreements(ctx context.Context, exec boil.ContextExecutor) ([]*models.Agreement, error) {
	return models.Agreements().All(ctx, exec)
}

func RemoveAuthority(ctx context.Context, tx *sql.Tx, userId int, authorityCode string) error {
	mAuthority, err := models.Authorities(models.AuthorityWhere.AuthorityCode.EQ(authorityCode)).One(ctx, tx)
	if err != nil && err != sql.ErrNoRows {
		return terr.Wrap(err)
	}

	if err == sql.ErrNoRows {
		return terr.New("authority not found")
	}

	_, err = models.UserAuthorities(models.UserAuthorityWhere.UserID.EQ(userId), models.UserAuthorityWhere.AuthorityID.EQ(mAuthority.AuthorityID)).DeleteAll(ctx, tx)

	return err
}
