package adapters

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gitlab.com/pedrokoblitz/opus-go/quality"
)

func NewSQLAdapter(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, quality.ErrSQLConnection.WithDetail(err.Error())
	}

	if err := db.Ping(); err != nil {
		return nil, quality.ErrSQLConnection.WithDetail(err.Error())
	}

	return db, nil
}

func NewSQLiteAdapter(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, quality.ErrSQLConnection.WithDetail(err.Error())
	}

	if err := db.Ping(); err != nil {
		return nil, quality.ErrSQLConnection.WithDetail(err.Error())
	}

	return db, nil
}
