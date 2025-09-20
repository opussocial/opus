package elements

import (
	"context"
	"gitlab.com/pedrokoblitz/opus-go/internal/quality"
	"gitlab.com/pedrokoblitz/opus-go/internal/adapters"
	"gitlab.com/pedrokoblitz/opus-go/internal/payloads"
)

/**
 * 
 * PROFILE STORE
 * 
 * 
 * 
 **/

const (
	CreateProfileMySQLQuery = "INSERT INTO profile_relationships (user_id, role_id, profile_element_id) VALUES (?, ?, ?)"
	DeleteProfileByUserIDMySQLQuery = "DELETE FROM profile_relationships WHERE user_id = ?"
)

type ProfileRelationshipStore struct {
	adapter adapters.DatabaseAdapter
}

func NewProfileRelationshipStore(adapter adapters.DatabaseAdapter) *ProfileRelationshipStore {
	return &ProfileRelationshipStore{adapter: adapter}
}

func (s *ProfileRelationshipStore) Create(ctx context.Context, p payloads.Payload) error {
	profile := p.(*ProfileRelationship)
	res, err := s.adapter.ExecContext(ctx,
		CreateProfileMySQLQuery,
		profile.UserID, profile.RoleID, profile.ProfileElementID,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	profile.ID = uint(id)

	return nil
}

func (s *ProfileRelationshipStore) Delete(ctx context.Context, p payloads.Payload) error {
	profile := p.(*payloads.DefaultID)
	_, err := s.adapter.ExecContext(ctx,
		DeleteProfileByUserIDMySQLQuery,
		profile.ID,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}
