package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jae2274/auth-service/auth_service/common/domain"
	"github.com/jae2274/auth-service/auth_service/models"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func newNecessaryAgreements() []*models.Agreement {
	return []*models.Agreement{
		{AgreementCode: "code1", Summary: "summary1", IsRequired: 1},
		{AgreementCode: "code2", Summary: "summary2", IsRequired: 1},
	}
}
func newOptionalAgreements() []*models.Agreement {
	return []*models.Agreement{
		{AgreementCode: "code3", Summary: "summary3", IsRequired: 0},
		{AgreementCode: "code4", Summary: "summary4", IsRequired: 0},
	}
}
func initAgreementFunc(t *testing.T, db *sql.DB) (context.Context, []*models.Agreement, []*models.Agreement, []*models.Authority) { //TODO: 추후 실제 비즈니스 로직을 통해 DB에 저장하는 것으로 변경
	//Given
	ctx := context.Background()
	var requiredAgreements models.AgreementSlice = newNecessaryAgreements()
	for _, agreement := range requiredAgreements {
		err := agreement.Insert(ctx, db, boil.Infer())
		require.NoError(t, err)
	}

	optionalAgreements := newOptionalAgreements()
	for _, agreement := range optionalAgreements {
		err := agreement.Insert(ctx, db, boil.Infer())
		require.NoError(t, err)
	}

	authorities := newAuthorities()
	for _, authority := range authorities {
		err := authority.Insert(ctx, db, boil.Infer())
		require.NoError(t, err)
	}
	adminAuthority := &models.Authority{AuthorityName: domain.AuthorityAdmin, Summary: "관리자 권한"}
	err := adminAuthority.Insert(ctx, db, boil.Infer())
	require.NoError(t, err)

	return ctx, requiredAgreements, optionalAgreements, authorities
}

func newAuthorities() []*models.Authority {
	return []*models.Authority{
		{AuthorityName: "AUTHORITY_USER", Summary: "사용자 권한"},
		{AuthorityName: "AUTHORITY_GUEST", Summary: "게스트 권한"},
	}
}

func requireEqualUserRole(t *testing.T, userId int, now time.Time, expected *dto.UserAuthorityReq, actual *domain.UserAuthority) {
	require.Equal(t, expected.AuthorityName, actual.AuthorityName)
	require.Equal(t, userId, actual.UserID)
	if expected.ExpiryDuration != nil {
		require.WithinDuration(t, now.Add(time.Duration(*expected.ExpiryDuration)), *actual.ExpiryDate, time.Second)
	} else {
		require.Nil(t, actual.ExpiryDate)
	}
}
