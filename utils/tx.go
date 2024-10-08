package utils

import "gorm.io/gorm"

type TxFn func(tx *gorm.DB) error

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`
// nolint: gocritic // no need to lint
func WithTransaction(db *gorm.DB, fn TxFn) (err error) {
	tx := db.Begin()

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and re-panic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			_ = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
