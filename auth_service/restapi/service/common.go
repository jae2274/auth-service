package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/jae2274/auth-service/auth_service/models"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
	"github.com/jae2274/goutils/terr"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func attachAuthorityIds(ctx context.Context, exec boil.ContextExecutor, dUserAuthorities []*dto.UserAuthorityReq) error {
	authorityCodes := make([]string, len(dUserAuthorities))
	for i, authority := range dUserAuthorities {
		authorityCodes[i] = authority.AuthorityCode
	}

	authorities, err := models.Authorities(models.AuthorityWhere.AuthorityCode.IN(authorityCodes)).All(ctx, exec)
	if err != nil {
		return err
	}
	authorityMap := make(map[string]*models.Authority)
	for _, authority := range authorities {
		authorityMap[authority.AuthorityCode] = authority
	}

	for _, userAuthority := range dUserAuthorities {
		authority, ok := authorityMap[userAuthority.AuthorityCode]
		if !ok {
			return terr.New("authority not found. authorityCode: " + userAuthority.AuthorityCode)
		}

		userAuthority.AuthorityID = authority.AuthorityID
	}

	return nil
}

func addExpiryDate(date time.Time, durationMS *int64) null.Time {
	if durationMS != nil {
		return null.NewTime(date.Add(time.Duration(*durationMS)*time.Millisecond), true)
	} else {
		return null.NewTime(time.Time{}, false)
	}
}

func addUserAuthorities(ctx context.Context, exec boil.ContextExecutor, userId int, dUserAuthorities []*dto.UserAuthorityReq) error {

	now := time.Now()

	// mUserAuthorities := make([]*models.UserAuthority, len(domainUserAuthorities))
	for _, dUserAuthority := range dUserAuthorities {
		mUserAuthority := &models.UserAuthority{ //로직 최적화 필요
			UserID:      userId,
			AuthorityID: dUserAuthority.AuthorityID,
		}
		err := mUserAuthority.Reload(ctx, exec)
		if err != nil && err != sql.ErrNoRows {
			return terr.Wrap(err)
		}

		if err == sql.ErrNoRows {
			mUserAuthority.ExpiryDate = addExpiryDate(now, dUserAuthority.ExpiryDurationMS)

			if err := mUserAuthority.Insert(ctx, exec, boil.Infer()); err != nil {
				return terr.Wrap(err)
			}
		} else {
			if mUserAuthority.ExpiryDate.Valid { //false의 경우, 만료되지 않는 권한으로 간주하여 만료일을 갱신하지 않음
				if mUserAuthority.ExpiryDate.Time.Before(now) {
					mUserAuthority.ExpiryDate = null.NewTime(now, true)
				}

				mUserAuthority.ExpiryDate = addExpiryDate(mUserAuthority.ExpiryDate.Time, dUserAuthority.ExpiryDurationMS)

				if _, err := mUserAuthority.Update(ctx, exec, boil.Infer()); err != nil {
					return terr.Wrap(err)
				}
			}
		}

		// mUserAuthorities[i] = modelUserAuthority
	}

	return nil
}
