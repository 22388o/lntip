package models

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

// checkErr checks for each sql transaction if an error was returned
// the transaction is rolledback if any is found, otherwise the transaction is commited
func checkErr(tx *sqlx.Tx, err error) {
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			zap.S().Error(err)
		}
	} else {
		err = tx.Commit()
		if err != nil {
			zap.S().Error(err)
		}
	}
}
