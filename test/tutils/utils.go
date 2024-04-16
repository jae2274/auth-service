package tutils

import (
	"database/sql"
	"testing"

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
