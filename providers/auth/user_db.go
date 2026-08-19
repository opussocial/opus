package auth

import (
	// "fmt"
	"context"
	"database/sql"
	"github.com/opussocialcontent/opus-go/actions"
	"github.com/opussocialcontent/opus-go/quality"
)

const (
	UpdateUserMySQLQuery         = "UPDATE users SET email = ? WHERE id = ?"
	GetUserByIDMySQLQuery   = "SELECT id, email FROM users WHERE id = ?" 
	GetUserByTokenMySQLQuery   = `SELECT 
		u.id, u.email, u.created_at, u.updated_at
		FROM users u
		JOIN tokens on u.id=tokens.user_id 
		WHERE tokens.token = ? 
		AND tokens.expires_at > NOW()`
	RemoveAccountMySQLQuery         = "DELETE FROM users WHERE id = ?"
)

type UserStore struct {
	adapter *sql.DB
}

func NewUserStore(adapter *sql.DB) *UserStore {
	return &UserStore{adapter: adapter}
}

func (s *UserStore) UpdateUser(ctx context.Context, p actions.Payload) error {
	profile := p.(*User)
	_, err := s.adapter.ExecContext(ctx,
		UpdateUserMySQLQuery,
		profile.Email, profile.ID,
	)
	if err != nil {
		return quality.ErrExecution.WithDetail(err.Error())
	}
	return nil
}

func (s *UserStore) RemoveAccount(ctx context.Context, p actions.Payload) error {
	defaultID := p.(*actions.DefaultID)
	_, err := s.adapter.ExecContext(ctx,
		RemoveAccountMySQLQuery,
		defaultID.ID,
	)
	if err != nil {
		return quality.ErrExecution.WithDetail(err.Error())
	}
	return nil
}

func (s *UserStore) GetByToken(ctx context.Context, p actions.Payload) error {
	profile := p.(*User)
	err := s.adapter.QueryRowContext(ctx,
		GetUserByTokenMySQLQuery,
		profile.Token,
	).Scan(&profile.ID, &profile.Email, &profile.CreatedAt, &profile.UpdatedAt)

	if err == sql.ErrNoRows {
		return quality.ErrNotFound
	}
	if err != nil {
		return quality.ErrExecution.WithDetail(err.Error())
	}

	return nil
}

func (s *UserStore) GetByID(ctx context.Context, p actions.Payload) error {
	profile := p.(*User)
	err := s.adapter.QueryRowContext(ctx,
		GetUserByIDMySQLQuery,
		profile.ID,
	).Scan(&profile.ID, &profile.Email)

	if err == sql.ErrNoRows {
		return quality.ErrNotFound
	}
	if err != nil {
		return quality.ErrExecution.WithDetail(err.Error())
	}

	return nil
}

