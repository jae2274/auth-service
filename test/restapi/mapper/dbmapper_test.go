package mysqldb

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
	"userService/models"
	"userService/test/tinit"
	"userService/test/tutils"
	"userService/usersvc/common/domain"
	"userService/usersvc/restapi/mapper"

	"github.com/go-sql-driver/mysql"
	"github.com/jae2274/goutils/terr"
	"github.com/stretchr/testify/require"
)

func TestDBMapper(t *testing.T) {
	willSavedUserVO := models.User{
		AuthorizedBy: string(domain.GOOGLE),
		AuthorizedID: "test",
		Email:        "test@mail.com",
	}

	t.Run("Find Non-Existent User", func(t *testing.T) {
		sqlDB := tinit.DB(t)
		defer sqlDB.Close()

		ctx := context.Background()
		tutils.TxCommit(t, sqlDB, func(tx *sql.Tx) {
			user, err := mapper.FindByAuthorized(ctx, tx, domain.GOOGLE, "test")
			require.NoError(t, err)
			require.Nil(t, user)
		})
	})

	t.Run("Save and Find User", func(t *testing.T) {
		t.Run("In same tx", func(t *testing.T) {
			sqlDB := tinit.DB(t)
			defer sqlDB.Close()

			ctx := context.Background()
			tutils.TxCommit(t, sqlDB, func(tx *sql.Tx) {

				require.NoError(t, mapper.SaveUser(ctx, tx, domain.AuthorizedBy(willSavedUserVO.AuthorizedBy), willSavedUserVO.AuthorizedID, willSavedUserVO.Email))

				user, err := mapper.FindByAuthorized(ctx, tx, domain.GOOGLE, "test")
				require.NoError(t, err)
				require.NotNil(t, user)

				require.Equal(t, 1, user.UserID)
				require.Equal(t, willSavedUserVO.Email, user.Email)
				require.Equal(t, willSavedUserVO.AuthorizedBy, user.AuthorizedBy)
				require.Equal(t, willSavedUserVO.AuthorizedID, user.AuthorizedID)
			})

		})
		t.Run("In two tx", func(t *testing.T) {
			sqlDB := tinit.DB(t)
			defer sqlDB.Close()

			ctx := context.Background()
			tutils.TxCommit(t, sqlDB, func(tx *sql.Tx) {
				require.NoError(t, mapper.SaveUser(ctx, tx, domain.AuthorizedBy(willSavedUserVO.AuthorizedBy), willSavedUserVO.AuthorizedID, willSavedUserVO.Email))
			})

			tutils.TxCommit(t, sqlDB, func(tx *sql.Tx) {
				user, err := mapper.FindByAuthorized(ctx, tx, domain.GOOGLE, "test")
				require.NoError(t, err)
				require.NotNil(t, user)

				require.Equal(t, 1, user.UserID)
				require.Equal(t, willSavedUserVO.Email, user.Email)
				require.Equal(t, willSavedUserVO.AuthorizedBy, user.AuthorizedBy)
				require.Equal(t, willSavedUserVO.AuthorizedID, user.AuthorizedID)
				require.GreaterOrEqual(t, time.Now().UTC(), user.CreateDate.UTC())
				require.LessOrEqual(t, time.Now().UTC().Add(-time.Second), user.CreateDate.UTC())
			})
		})
	})

	t.Run("Save but rollback", func(t *testing.T) {
		sqlDB := tinit.DB(t)
		defer sqlDB.Close()

		ctx := context.Background()
		tutils.TxRollback(t, sqlDB, func(tx *sql.Tx) {
			require.NoError(t, mapper.SaveUser(ctx, tx, domain.AuthorizedBy(willSavedUserVO.AuthorizedBy), willSavedUserVO.AuthorizedID, willSavedUserVO.Email))
		})

		tutils.TxCommit(t, sqlDB, func(tx *sql.Tx) {
			user, err := mapper.FindByAuthorized(ctx, tx, domain.GOOGLE, "test")
			require.NoError(t, err)
			require.Nil(t, user)
		})
	})

	t.Run("Insert Duplicate User", func(t *testing.T) {
		isDuplicate := func(err error) bool {
			var mysqlErr *mysql.MySQLError
			err = terr.UnWrap(err)
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				return true
			}
			return false
		}

		sameUser := models.User{
			AuthorizedBy: willSavedUserVO.AuthorizedBy,
			AuthorizedID: willSavedUserVO.AuthorizedID,
			Email:        "test2@naver.com",
		}

		t.Run("In same tx", func(t *testing.T) {
			sqlDB := tinit.DB(t)
			defer sqlDB.Close()

			ctx := context.Background()
			tutils.TxCommit(t, sqlDB, func(tx *sql.Tx) {
				require.NoError(t, mapper.SaveUser(ctx, tx, domain.AuthorizedBy(willSavedUserVO.AuthorizedBy), willSavedUserVO.AuthorizedID, willSavedUserVO.Email))
				err := mapper.SaveUser(ctx, tx, domain.AuthorizedBy(sameUser.AuthorizedBy), sameUser.AuthorizedID, sameUser.Email)
				require.Error(t, err)
				require.True(t, isDuplicate(err))
			})
		})

		t.Run("In two tx", func(t *testing.T) {
			sqlDB := tinit.DB(t)
			defer sqlDB.Close()

			ctx := context.Background()
			tutils.TxCommit(t, sqlDB, func(tx *sql.Tx) {
				require.NoError(t, mapper.SaveUser(ctx, tx, domain.AuthorizedBy(willSavedUserVO.AuthorizedBy), willSavedUserVO.AuthorizedID, willSavedUserVO.Email))
			})

			tutils.TxCommit(t, sqlDB, func(tx *sql.Tx) {
				err := mapper.SaveUser(ctx, tx, domain.AuthorizedBy(sameUser.AuthorizedBy), sameUser.AuthorizedID, sameUser.Email)
				require.Error(t, err)
				require.True(t, isDuplicate(err))
			})
		})
	})
}
