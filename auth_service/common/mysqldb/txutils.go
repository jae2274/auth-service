package mysqldb

import (
	"context"
	"database/sql"
	"errors"
)

func WithTransaction[T any](ctx context.Context, db *sql.DB, f func(tx *sql.Tx) (T, error)) (T, error) {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		var emptyObj T
		return emptyObj, err
	}

	t, err := f(tx)

	return CommitOrRollback(tx, t, err)
}

func WithTransactionVoid(ctx context.Context, db *sql.DB, f func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	err = f(tx)

	return CommitOrRollbackVoid(tx, err)
}

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
