package test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jae2274/auth-service/auth_service/common/mysqldb"
	"github.com/jae2274/auth-service/auth_service/models"
	"github.com/jae2274/auth-service/test/tinit"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// 아래의 테스트코드는 sqlboiler의 select for update가 의도대로 동작하는지 확인하는 테스트코드입니다.
func TestSelectForUpdate(t *testing.T) {
	db := tinit.DB(t)

	ctx := context.Background()
	authority := models.Authority{
		AuthorityCode: "auth_code",
		AuthorityName: "auth_name",
		Summary:       "auth_summary",
	}

	err := authority.Insert(ctx, db, boil.Infer())
	require.NoError(t, err)

	// select for update로 데이터를 조회하면 다른 트랜잭션에서 해당 데이터를 조회하는 동안 lock이 걸립니다.
	go func() {
		err = mysqldb.WithTransactionVoid(ctx, db, func(tx *sql.Tx) error {
			forUpdateAuth, err := models.Authorities(models.AuthorityWhere.AuthorityID.EQ(authority.AuthorityID), qm.For("update")).One(ctx, tx)
			if err != nil {
				return err
			}

			time.Sleep(500 * time.Millisecond)
			forUpdateAuth.Summary = "updated_summary"
			_, err = forUpdateAuth.Update(ctx, tx, boil.Infer())
			return err
		})
		require.NoError(t, err)
	}()

	time.Sleep(100 * time.Millisecond) // 다른 트랜잭션이 시작되기 전에 대기합니다.
	start := time.Now()
	modified, err := models.Authorities(models.AuthorityWhere.AuthorityID.EQ(authority.AuthorityID), qm.For("update")).One(ctx, db)
	elapsed := time.Since(start)
	require.NoError(t, err)
	require.Greater(t, elapsed, time.Millisecond*350)
	require.Less(t, elapsed, time.Millisecond*450)

	require.Equal(t, "updated_summary", modified.Summary)
}
