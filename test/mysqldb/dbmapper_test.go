package mysqldb

import (
	"database/sql"
	"testing"
	"userService/test/tinit"
	"userService/usersvc/domain"
	"userService/usersvc/entity"
	"userService/usersvc/mysqldb"

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

	t.Run("Save and Find User in same tx", func(t *testing.T) {
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
	t.Run("Save and Find User in two tx", func(t *testing.T) {
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
