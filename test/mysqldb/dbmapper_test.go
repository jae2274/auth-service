package mysqldb

import (
	"database/sql"
	"errors"
	"testing"
	"time"
	"userService/test/tinit"
	"userService/usersvc/domain"
	"userService/usersvc/entity"
	"userService/usersvc/mysqldb"

	"github.com/go-sql-driver/mysql"
	"github.com/jae2274/goutils/terr"
	"github.com/stretchr/testify/require"
)

func TestDBMapper(t *testing.T) {
	willSavedUserVO := entity.UserVO{
		AuthorizedBy: domain.GOOGLE,
		AuthorizedID: "test",
		Email:        "test@mail.com",
	}

	t.Run("Find Non-Existent User", func(t *testing.T) {
		sqlDB := tinit.DB(t)
		defer sqlDB.Close()

		txCommit(t, sqlDB, func(tx *sql.Tx) {
			user, err := mysqldb.FindByAuthorized(tx, domain.GOOGLE, "test")
			require.NoError(t, err)
			require.Nil(t, user)
		})
	})

	t.Run("Save and Find User", func(t *testing.T) {
		t.Run("In same tx", func(t *testing.T) {
			sqlDB := tinit.DB(t)
			defer sqlDB.Close()
			txCommit(t, sqlDB, func(tx *sql.Tx) {

				require.NoError(t, mysqldb.SaveUser(tx, willSavedUserVO))

				user, err := mysqldb.FindByAuthorized(tx, domain.GOOGLE, "test")
				require.NoError(t, err)
				require.NotNil(t, user)

				require.Equal(t, int64(1), user.UserID)
				require.Equal(t, willSavedUserVO.Email, user.Email)
				require.Equal(t, willSavedUserVO.AuthorizedBy, user.AuthorizedBy)
				require.Equal(t, willSavedUserVO.AuthorizedID, user.AuthorizedID)
			})

		})
		t.Run("In two tx", func(t *testing.T) {
			sqlDB := tinit.DB(t)
			defer sqlDB.Close()
			txCommit(t, sqlDB, func(tx *sql.Tx) {

				require.NoError(t, mysqldb.SaveUser(tx, willSavedUserVO))
			})

			txCommit(t, sqlDB, func(tx *sql.Tx) {
				user, err := mysqldb.FindByAuthorized(tx, domain.GOOGLE, "test")
				require.NoError(t, err)
				require.NotNil(t, user)

				require.Equal(t, int64(1), user.UserID)
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
		txRollback(t, sqlDB, func(tx *sql.Tx) {
			require.NoError(t, mysqldb.SaveUser(tx, willSavedUserVO))
		})

		txCommit(t, sqlDB, func(tx *sql.Tx) {
			user, err := mysqldb.FindByAuthorized(tx, domain.GOOGLE, "test")
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

		sameUser := entity.UserVO{
			AuthorizedBy: willSavedUserVO.AuthorizedBy,
			AuthorizedID: willSavedUserVO.AuthorizedID,
			Email:        "test2@naver.com",
		}

		t.Run("In same tx", func(t *testing.T) {
			sqlDB := tinit.DB(t)
			defer sqlDB.Close()

			txCommit(t, sqlDB, func(tx *sql.Tx) {
				require.NoError(t, mysqldb.SaveUser(tx, willSavedUserVO))
				err := mysqldb.SaveUser(tx, sameUser)
				require.Error(t, err)
				require.True(t, isDuplicate(err))
			})
		})

		t.Run("In two tx", func(t *testing.T) {
			sqlDB := tinit.DB(t)
			defer sqlDB.Close()

			txCommit(t, sqlDB, func(tx *sql.Tx) {
				require.NoError(t, mysqldb.SaveUser(tx, willSavedUserVO))
			})

			txCommit(t, sqlDB, func(tx *sql.Tx) {
				err := mysqldb.SaveUser(tx, sameUser)
				require.Error(t, err)
				require.True(t, isDuplicate(err))
			})
		})
	})
}

func txCommit(t *testing.T, db *sql.DB, action func(*sql.Tx)) {
	tx, err := db.Begin()
	if err != nil {
		require.NoError(t, err)
	}
	action(tx)
	require.NoError(t, tx.Commit())
}

func txRollback(t *testing.T, db *sql.DB, action func(*sql.Tx)) {
	tx, err := db.Begin()
	if err != nil {
		require.NoError(t, err)
	}
	action(tx)
	require.NoError(t, tx.Rollback())
}
