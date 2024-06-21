package mysqldb

import (
	"database/sql"
	"errors"
)

func CommitOrRollback[T any](tx *sql.Tx, t T, err error) (T, error) {
	return t, CommitOrRollbackVoid(tx, err)
}

func CommitOrRollbackVoid(tx *sql.Tx, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
