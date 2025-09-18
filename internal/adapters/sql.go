package adapters

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gitlab.com/pedrokoblitz/opus-go/internal/quality"
)

// DatabaseAdapter defines the complete interface for all database operations
type DatabaseAdapter interface {
    // Basic operations
    Close() error
    ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
    QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
    QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}


func NewSQLAdapter(dsn string) (DatabaseAdapter, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, quality.ErrSQLConnection.WithDetail(err.Error())
	}

	if err := db.Ping(); err != nil {
		return nil, quality.ErrSQLConnection.WithDetail(err.Error())
	}

	return db, nil
}

func NewSQLiteAdapter(dsn string) (DatabaseAdapter, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, quality.ErrSQLConnection.WithDetail(err.Error())
	}

	if err := db.Ping(); err != nil {
		return nil, quality.ErrSQLConnection.WithDetail(err.Error())
	}

	return db, nil
}
