package auth

import (
	// "fmt"
	"time"
	"context"
	// "database/sql"
	"gitlab.com/pedrokoblitz/opus-go/internal/quality"
	"gitlab.com/pedrokoblitz/opus-go/internal/adapters"
	"gitlab.com/pedrokoblitz/opus-go/internal/payloads"
)

const (
	InsertConfirmationTokenMySQLQuery  = "INSERT INTO tokens (token, user_id, scope, expires_at) VALUES (?, ?, 'confirmation', ?)"	
	SignUpUserMySQLQuery          = "INSERT INTO users (email, password) VALUES (?, ?)"
	SignUpConfirmedUserMySQLQuery = "INSERT INTO users (email, password, confirmed_at) VALUES (?, ?, NOW())"
	ConfirmUserMySQLQuery           = `
		UPDATE users JOIN tokens ON users.id=tokens.user_id SET confirmed_at = NOW() WHERE tokens.token = ?
	`
	DeleteConfirmationTokenMySQLQuery  = "DELETE FROM tokens WHERE token = ?"
)

type SignUpStore struct {
	adapter adapters.DatabaseAdapter
}

func NewSignUpStore(adapter adapters.DatabaseAdapter) *SignUpStore {
	return &SignUpStore{adapter: adapter}
}

func (s *SignUpStore) SignUp(ctx context.Context, p payloads.Payload) error {
	signUp := p.(*SignUp)
	res, err := s.adapter.ExecContext(ctx,
		SignUpUserMySQLQuery,
		signUp.Email, signUp.Password,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return quality.ErrExecution.WithDetail(err.Error())
	}
	signUp.ID = uint(id)

	return nil
}

func (s *SignUpStore) CreateConfirmationToken(ctx context.Context, p payloads.Payload) error {
	signUp := p.(*SignUp)
	_, err := s.adapter.ExecContext(ctx,
		InsertConfirmationTokenMySQLQuery,
		signUp.Token, signUp.ID, time.Now().Add(1 * time.Hour),
	)
	if err != nil {
		return quality.ErrExecution.WithDetail(err.Error())
	}
	return nil
}

func (s *SignUpStore) Confirm(ctx context.Context, p payloads.Payload) error {
	confirm := p.(*Token)
	res, err := s.adapter.ExecContext(ctx,
		ConfirmUserMySQLQuery,
		confirm.Token,
	)
	if err != nil {
		return quality.ErrExecution.WithDetail(err.Error())
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return quality.ErrExecution.WithDetail(err.Error())
	}
	if affected == 0 {
		return quality.ErrExecution
	}
	return nil
}

func (s *SignUpStore) DeleteConfirmationToken(ctx context.Context, p payloads.Payload) error {
	token := p.(*Token)
	_, err := s.adapter.ExecContext(ctx,
		DeleteConfirmationTokenMySQLQuery,
		token.Token,
	)
	if err != nil {
		return quality.ErrExecution.WithDetail(err.Error())
	}
	return nil
}
