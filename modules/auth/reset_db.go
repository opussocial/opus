package auth

import (
	"time"
	"context"
	"database/sql"
	"gitlab.com/pedrokoblitz/opus-go/actions"
	"gitlab.com/pedrokoblitz/opus-go/quality"
)

const (
	ResetPasswordMySQLQuery = "UPDATE users SET password = ? WHERE id = (SELECT user_id FROM tokens WHERE token = ? AND scope = 'reset')"
	InsertResetTokenMySQLQuery = "INSERT INTO tokens (token, user_id, scope, expires_at) VALUES (?, ?, 'reset', ?)"
	DeleteResetTokenMySQLQuery = "DELETE FROM tokens WHERE token = ? AND scope = 'reset'"
)

type ResetStore struct {
	adapter *sql.DB
}

func NewResetStore(adapter *sql.DB) *ResetStore {
	return &ResetStore{adapter: adapter}
}

func (s *ResetStore) CreateResetToken(ctx context.Context, p payloads.Payload) error {
	reset := p.(*ResetRequest)
	_, err := s.adapter.ExecContext(ctx,
		InsertResetTokenMySQLQuery,
		reset.Token, reset.ID, time.Now().Add(1 * time.Hour),
	)
	if err != nil {
		return quality.ErrExecution.WithDetail(err.Error())
	}
	return nil
}

func (s *ResetStore) DeleteResetToken(ctx context.Context, p payloads.Payload) error {
	token := p.(*Reset)
	_, err := s.adapter.ExecContext(ctx,
		DeleteResetTokenMySQLQuery,
		token.Token,
	)
	if err != nil {
		return quality.ErrExecution.WithDetail(err.Error())
	}
	return nil
}

func (s *ResetStore) ResetPassword(ctx context.Context, p payloads.Payload) error {
	password := p.(*Reset)
	_, err := s.adapter.ExecContext(ctx,
		ResetPasswordMySQLQuery,
		password.Password, password.Token,
	)
	if err != nil {
		return quality.ErrExecution.WithDetail(err.Error())
	}
	return nil
}
