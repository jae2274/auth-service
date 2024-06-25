package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jae2274/auth-service/auth_service/common/domain"
	"github.com/jae2274/auth-service/auth_service/common/mysqldb"
	"github.com/jae2274/auth-service/auth_service/models"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
	"github.com/jae2274/auth-service/auth_service/restapi/ooauth"
	"github.com/jae2274/auth-service/auth_service/restapi/service"
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
	adminAuthority := &models.Authority{AuthorityCode: domain.AuthorityAdmin, AuthorityName: "관리자", Summary: "관리자 권한"}
	err := adminAuthority.Insert(ctx, db, boil.Infer())
	require.NoError(t, err)

	return ctx, requiredAgreements, optionalAgreements, authorities
}

func newAuthorities() []*models.Authority {
	return []*models.Authority{
		{AuthorityCode: "AUTHORITY_USER", AuthorityName: "사용자", Summary: "사용자 권한"},
		{AuthorityCode: "AUTHORITY_GUEST", AuthorityName: "손님", Summary: "게스트 권한"},
	}
}

func requireEqualUserRole(t *testing.T, userId int, now time.Time, expected *dto.UserAuthorityReq, actual *domain.UserAuthority) {
	require.Equal(t, expected.AuthorityCode, actual.AuthorityCode)
	require.Equal(t, userId, actual.UserID)
	if expected.ExpiryDurationMS != nil {
		require.WithinDuration(t, now.Add(time.Duration(*expected.ExpiryDurationMS)*time.Millisecond).UTC(), time.UnixMilli(*actual.ExpiryUnixMilli).UTC(), time.Second)
	} else {
		require.Nil(t, actual.ExpiryUnixMilli)
	}
}

func createTicketWithTx(ctx context.Context, db *sql.DB, authorities []*dto.UserAuthorityReq) (string, error) {
	return mysqldb.WithTransaction(ctx, db, func(tx *sql.Tx) (string, error) {
		return service.CreateTicket(ctx, tx, authorities)
	})
}

func signUp(ctx context.Context, db *sql.DB, userinfo *ooauth.UserInfo, agreements []*dto.UserAgreementReq) (*models.User, error) {
	return mysqldb.WithTransaction(ctx, db, func(tx *sql.Tx) (*models.User, error) {
		return service.SignUp(ctx, tx, userinfo, agreements)
	})
}

func useTicket(ctx context.Context, db *sql.DB, userId int, ticketId string) error {
	return mysqldb.WithTransactionVoid(ctx, db, func(tx *sql.Tx) error {
		return service.UseTicket(ctx, tx, userId, ticketId)
	})
}

func addUserAuthorities(ctx context.Context, db *sql.DB, userId int, dUserAuthorities []*dto.UserAuthorityReq) error {
	return mysqldb.WithTransactionVoid(ctx, db, func(tx *sql.Tx) error {
		return service.AddUserAuthorities(ctx, tx, userId, dUserAuthorities)
	})
}

func applyUserAgreements(ctx context.Context, db *sql.DB, userId int, agreements []*dto.UserAgreementReq) error {
	return mysqldb.WithTransactionVoid(ctx, db, func(tx *sql.Tx) error {
		return service.ApplyUserAgreements(ctx, tx, userId, agreements)
	})
}

func removeAuthority(ctx context.Context, db *sql.DB, userId int, authorityCode string) error {
	return mysqldb.WithTransactionVoid(ctx, db, func(tx *sql.Tx) error {
		return service.RemoveAuthority(ctx, tx, userId, authorityCode)
	})
}
