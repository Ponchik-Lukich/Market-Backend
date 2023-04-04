package middleware

import "github.com/jmoiron/sqlx"

func RollbackOrCommit(tx *sqlx.Tx, errPtr *error) {
	if p := recover(); p != nil {
		tx.Rollback()
		panic(p)
	} else if *errPtr != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}
