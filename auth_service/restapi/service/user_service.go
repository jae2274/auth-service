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
	"github.com/jae2274/goutils/ptr"
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

func FindValidUserAuthorities(ctx context.Context, exec boil.ContextExecutor, userId int) ([]*domain.UserAuthority, error) {
	userAuthorities, err := models.UserAuthorities(models.UserAuthorityWhere.UserID.EQ(userId),
		qm.Expr(models.UserAuthorityWhere.ExpiryDate.GTE(null.NewTime(time.Now(), true)), qm.Or2(models.UserAuthorityWhere.ExpiryDate.IsNull())),
		qm.OrderBy(models.UserAuthorityColumns.CreatedDate),
		qm.Load(models.UserAuthorityRels.Authority),
	).All(ctx, exec)

	if err != nil {
		return nil, err
	}

	return convertToUserAuthority(userAuthorities), nil
}

func FindAllUserAuthorities(ctx context.Context, exec boil.ContextExecutor, userId int) ([]*domain.UserAuthority, error) {
	userAuthorities, err := models.UserAuthorities(models.UserAuthorityWhere.UserID.EQ(userId),
		qm.OrderBy(models.UserAuthorityColumns.CreatedDate),
		qm.Load(models.UserAuthorityRels.Authority),
	).All(ctx, exec)

	if err != nil {
		return nil, err
	}

	return convertToUserAuthority(userAuthorities), nil
}

func convertToUserAuthority(mUserAuthorities models.UserAuthoritySlice) []*domain.UserAuthority {
	authorities := make([]*domain.UserAuthority, len(mUserAuthorities))
	for i, userAuthority := range mUserAuthorities {
		var expiryDate *int64 = nil
		if userAuthority.ExpiryDate.Valid {
			expiryDate = ptr.P(userAuthority.ExpiryDate.Time.UnixMilli())
		}

		authorities[i] = &domain.UserAuthority{
			UserID:          userAuthority.UserID,
			AuthorityID:     userAuthority.R.Authority.AuthorityID,
			AuthorityCode:   userAuthority.R.Authority.AuthorityCode,
			AuthorityName:   userAuthority.R.Authority.AuthorityName,
			Summary:         userAuthority.R.Authority.Summary,
			ExpiryUnixMilli: expiryDate,
		}
	}

	return authorities
}

/*
	추후 해당 기능이 사용될 여지가 있을 것으로 판단되어 주석처리하였습니다.
*/
// func FindUserAuthoritiesByAuthorityIds(ctx context.Context, exec boil.ContextExecutor, userId int, authorityId []int) ([]*domain.UserAuthority, error) {
// 	userAuthorities, err := models.UserAuthorities(
// 		models.UserAuthorityWhere.UserID.EQ(userId), models.UserAuthorityWhere.AuthorityID.IN(authorityId),
// 		qm.OrderBy(models.UserAuthorityColumns.CreatedDate),
// 		qm.Load(models.UserAuthorityRels.Authority),
// 	).All(ctx, exec)
// 	if err != nil {
// 		return nil, err
// 	}

// 	authorities := make([]*domain.UserAuthority, len(userAuthorities))
// 	for i, userAuthority := range userAuthorities {
// 		var expiryDate *int64 = nil
// 		if userAuthority.ExpiryDate.Valid {
// 			expiryDate = ptr.P(userAuthority.ExpiryDate.Time.UnixMilli())
// 		}

// 		authorities[i] = &domain.UserAuthority{
// 			UserID:          userAuthority.UserID,
// 			AuthorityID:     userAuthority.R.Authority.AuthorityID,
// 			AuthorityCode:   userAuthority.R.Authority.AuthorityCode,
// 			AuthorityName:   userAuthority.R.Authority.AuthorityName,
// 			Summary:         userAuthority.R.Authority.Summary,
// 			ExpiryUnixMilli: expiryDate,
// 		}
// 	}

// 	return authorities, nil
// }

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

var ErrCannotControlAuthorityAdmin = terr.New("cannot control authority admin")

func RemoveAuthority(ctx context.Context, tx *sql.Tx, userId int, authorityCode string) error {
	if authorityCode == domain.AuthorityAdmin {
		return ErrCannotControlAuthorityAdmin
	}

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

func GetAllUsers(ctx context.Context, exec boil.ContextExecutor) ([]*domain.User, error) {
	mUsers, err := models.Users(qm.Load(models.UserRels.UserAuthorities+"."+models.UserAuthorityRels.Authority)).All(ctx, exec)
	if err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(mUsers))
	for i, mUser := range mUsers {
		users[i] = convertToUser(mUser)
	}

	return users, nil
}

func convertToUser(mUser *models.User) *domain.User {
	userAuthorities := make([]*domain.UserAuthority, len(mUser.R.UserAuthorities))
	for i, mUserAuthority := range mUser.R.UserAuthorities {
		var expiryUnixMilli *int64 = nil
		if mUserAuthority.ExpiryDate.Valid {
			expiryUnixMilli = ptr.P(mUserAuthority.ExpiryDate.Time.UnixMilli())
		}
		userAuthorities[i] = &domain.UserAuthority{
			UserID:          mUserAuthority.UserID,
			AuthorityID:     mUserAuthority.AuthorityID,
			AuthorityCode:   mUserAuthority.R.Authority.AuthorityCode,
			AuthorityName:   mUserAuthority.R.Authority.AuthorityName,
			Summary:         mUserAuthority.R.Authority.Summary,
			ExpiryUnixMilli: expiryUnixMilli,
		}
	}

	return &domain.User{
		UserID:           mUser.UserID,
		AuthorizedBy:     domain.AuthorizedBy(mUser.AuthorizedBy),
		AuthorizedID:     mUser.AuthorizedID,
		UserName:         mUser.Name,
		Email:            mUser.Email,
		Authorities:      userAuthorities,
		CreatedUnixMilli: mUser.CreateDate.UnixMilli(),
	}
}
