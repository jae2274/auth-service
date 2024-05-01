package tutils

import (
	"database/sql"
	"fmt"
	"testing"
	"userService/usersvc/common/domain"
	"userService/usersvc/models"

	"github.com/stretchr/testify/require"
)

func TxCommit(t *testing.T, db *sql.DB, action func(*sql.Tx)) {
	tx, err := db.Begin()
	if err != nil {
		require.NoError(t, err)
	}
	action(tx)
	require.NoError(t, tx.Commit())
}

func TxRollback(t *testing.T, db *sql.DB, action func(*sql.Tx)) {
	tx, err := db.Begin()
	if err != nil {
		require.NoError(t, err)
	}
	action(tx)
	require.NoError(t, tx.Rollback())
}

func NewUser(n int) *models.User {
	return &models.User{
		AuthorizedBy: string(domain.GOOGLE),
		AuthorizedID: fmt.Sprintf("test%d", n),
		Name:         fmt.Sprintf("testName%d", n),
		Email:        fmt.Sprintf("test%d@mail.com", n),
	}
}
