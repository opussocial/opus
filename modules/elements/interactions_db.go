package elements

import (
	"context"
	"gitlab.com/pedrokoblitz/opus-go/internal/quality"
	"gitlab.com/pedrokoblitz/opus-go/internal/adapters"
	"gitlab.com/pedrokoblitz/opus-go/internal/payloads"
)

/**
 * 
 * SCOPE STORE
 * 
 * 

-- For element_taxonomy (linking elements and taxonomy terms)
InsertElementTaxonomyMySQLQuery = 
  "INSERT INTO element_taxonomy (element_id, taxonomy_id) VALUES (?, ?)"
DeleteElementTaxononyMySQLQuery = 
  "DELETE FROM element_taxonomy WHERE element_id = ? AND taxonomy_id = ?"

-- For taxonomy (term definitions)
InsertTaxonomyMySQLQuery = 
  "INSERT INTO taxonomy (scope, term, slug) VALUES (?, ?, ?)"
UpdateTaxonomyMySQLQuery = 
  "UPDATE taxonomy SET scope = ?, term = ?, slug = ? WHERE id = ?"
DeleteTaxonomyMySQLQuery = 
  "DELETE FROM taxonomy WHERE id = ?"

 * 
 **/

const (
	CreateInteractionMySQLQuery = "INSERT INTO interactions (scope, element_id, created_by_profile_id, rating, comment) VALUES (?, ?, ?, ?, ?)"
	DeleteInteractionMySQLQuery = "DELETE FROM interactions WHERE id = ?"
)

type InteractionStore struct {
	adapter adapters.DatabaseAdapter
}

func NewInteractionStore(adapter adapters.DatabaseAdapter) *InteractionStore {
	return &InteractionStore{adapter: adapter}
}

func (s *InteractionStore) Create(ctx context.Context, p payloads.Payload) error {
	interaction := p.(*Interaction)
	res, err := s.adapter.ExecContext(ctx,
		CreateInteractionMySQLQuery,
		interaction.Scope, 
		interaction.ElementID,
		interaction.CreatedByProfileID,
		interaction.Rating,
		interaction.Comment,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	interaction.ID = uint(id)

	return nil
}

func (s *InteractionStore) Delete(ctx context.Context, p payloads.Payload) error {
	interaction := p.(*payloads.DefaultID)
	_, err := s.adapter.ExecContext(ctx,
		DeleteInteractionMySQLQuery,
		interaction.ID,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}


/**
 * 
 * KEYWORD STORE
 * 
 * 
 **/

const (
	CreateKeywordMySQLQuery = "INSERT INTO taxonomy (term, slug) VALUES (?, ?)"
	DeleteKeywordMySQLQuery = "DELETE FROM taxonomy WHERE slug = ?"
)

type KeywordStore struct {
	adapter adapters.DatabaseAdapter
}

func NewKeywordStore(adapter adapters.DatabaseAdapter) *KeywordStore {
	return &KeywordStore{adapter: adapter}
}

func (s *KeywordStore) Create(ctx context.Context, p payloads.Payload) error {
	keyword := p.(*Keyword)
	res, err := s.adapter.ExecContext(ctx,
		CreateKeywordMySQLQuery,
		keyword.Term, 
		keyword.Slug, 
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	keyword.ID = uint(id)

	return nil
}

func (s *KeywordStore) Delete(ctx context.Context, p payloads.Payload) error {
	keyword := p.(*payloads.DefaultID)
	_, err := s.adapter.ExecContext(ctx,
		DeleteKeywordMySQLQuery,
		keyword.ID,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}
