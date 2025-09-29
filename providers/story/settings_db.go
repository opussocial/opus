package story

import (
	"context"
	"database/sql"

	"gitlab.com/pedrokoblitz/opus-go/actions"
	"gitlab.com/pedrokoblitz/opus-go/quality"
)

const (
	CreateSettingMySQLQuery = "INSERT INTO settings (scope, name, slug, value) VALUES (?, ?, ?, ?)"
	CreateStorySettingMySQLQuery = "INSERT INTO settings (scope, name, slug, value, story_id) VALUES (?, ?, ?, ?, ?)"
	CreateUserSettingMySQLQuery = "INSERT INTO settings (scope, name, slug, value, user_id) VALUES (?, ?, ?, ?, ?)"
	CreateProfileSettingMySQLQuery = "INSERT INTO settings (scope, name, slug, value, user_id, story_id) VALUES (?, ?, ?, ?, ?, ?)"
	DeleteSettingMySQLQuery = "DELETE FROM settings WHERE id = ?"
	DeleteSettingByStoryIDMySQLQuery = "DELETE FROM settings WHERE story_id = ?"
	DeleteSettingByUserIDMySQLQuery = "DELETE FROM settings WHERE user_id = ?"
)


/**
 * 
 * ROLE STORE
 * 
 * 
 **/

type SettingStore struct {
	adapter *sql.DB
}

func NewSettingStore(adapter *sql.DB) *SettingStore {
	return &SettingStore{adapter: adapter}
}

func (s *SettingStore) Create(ctx context.Context, p actions.Payload) error {
	setting := p.(*Setting)
	res, err := s.adapter.ExecContext(ctx,
		CreateSettingMySQLQuery,
		setting.Name, setting.Slug,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	setting.ID = uint(id)

	return nil
}

func (us *SettingStore) Delete(ctx context.Context, p actions.Payload) error {
	setting := p.(*actions.DefaultID)
	_, err := us.adapter.ExecContext(ctx,
		DeleteSettingMySQLQuery,
		setting.ID,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}

func (us *SettingStore) DeleteByStoryID(ctx context.Context, p actions.Payload) error {
	setting := p.(*actions.DefaultID)
	_, err := us.adapter.ExecContext(ctx,
		DeleteSettingByStoryIDMySQLQuery,
		setting.ID,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}


