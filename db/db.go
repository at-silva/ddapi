//Package db contains interfaces to abstract the sqlx.DB struct
package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 database/sql.Result

// DB represents a sqlx db
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . DB
type (
	DB interface {
		NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
		NamedQueryContext(ctx context.Context, query string, arg interface{}) (Rows, error)
	}

	// Rows represents a sqlx Rows
	//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Rows
	Rows interface {
		MapScan(map[string]interface{}) error
		Next() bool
	}
)

// New returns a wrapped DB
func New(db *sqlx.DB) DB {
	return &dbOp{db}
}

type dbOp struct {
	*sqlx.DB
}

func (db dbOp) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return db.DB.NamedExecContext(ctx, query, arg)
}

func (db dbOp) NamedQueryContext(ctx context.Context, query string, arg interface{}) (Rows, error) {
	return db.DB.NamedQueryContext(ctx, query, arg)
}
