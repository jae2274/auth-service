package mysqldb

import (
	"database/sql"
	"errors"
)

func CommitOrRollback[T any](tx *sql.Tx, t T, err error) (T, error) {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return t, errors.Join(err, rollbackErr)
		}

		return t, err
	}

	if err := tx.Commit(); err != nil {
		return t, err
	}

	return t, nil
}
