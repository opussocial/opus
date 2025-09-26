package auth

import (
	"time"
	"context"
	"database/sql"
	"gitlab.com/pedrokoblitz/opus-go/quality"
	"gitlab.com/pedrokoblitz/opus-go/actions"
)

const (
	GetUserByEmailMySQLQuery   = "SELECT id, email, password FROM users WHERE email = ? AND confirmed_at IS NOT NULL"
	InsertAuthTokenMySQLQuery  = "INSERT INTO tokens (token, user_id, scope, expires_at) VALUES (?, ?, 'auth', ?)"
	DeleteAuthTokenMySQLQuery  = "DELETE FROM tokens WHERE token = ?"
)

type LoginStore struct {
	adapter actions.DatabaseAdapter
}

func NewLoginStore(adapter actions.DatabaseAdapter) *LoginStore {
	return &LoginStore{adapter: adapter}
}

func (s *LoginStore) GetByEmail(ctx context.Context, p payloads.Payload) error {
	login := p.(*Login)
	err := s.adapter.QueryRowContext(ctx,
		GetUserByEmailMySQLQuery,
		login.Email,
	).Scan(&login.ID, &login.Email, &login.Password)

	if err == sql.ErrNoRows {
		return quality.ErrNotFound
	}
	if err != nil {
		return quality.ErrExecution.WithDetail(err.Error())
	}
	return nil
}

func (s *LoginStore) CreateAuthToken(ctx context.Context, p payloads.Payload) error {
	login := p.(*Login)
	_, err := s.adapter.ExecContext(ctx,
		InsertAuthTokenMySQLQuery,
		login.Token, login.ID, time.Now().Add(1 * time.Hour),
	)
	if err != nil {
		return quality.ErrExecution.WithDetail(err.Error())
	}
	return nil
}

func (s *LoginStore) DeleteAuthToken(ctx context.Context, p payloads.Payload) error {
	token := p.(*Token)
	_, err := s.adapter.ExecContext(ctx,
		DeleteAuthTokenMySQLQuery,
		token.Token,
	)
	if err != nil {
		return quality.ErrExecution.WithDetail(err.Error())
	}
	return nil
}
