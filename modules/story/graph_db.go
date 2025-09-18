package story

import (
	"context"
	"gitlab.com/pedrokoblitz/opus-go/internal/quality"
	"gitlab.com/pedrokoblitz/opus-go/internal/adapters"
	"gitlab.com/pedrokoblitz/opus-go/internal/payloads"
)


const (
	CreateGraphRelationshipMySQLQuery = "INSERT INTO graph_relationships (name, slug, source_id, target_id) VALUES (?, ?, ?, ?)"
	DeleteGraphRelationshipMySQLQuery = "DELETE FROM graph_relationships WHERE id = ?"
)

/**
 * 
 * GRAPH RELATIONSHIP STORE
 * 
 * 
 **/

type GraphRelationshipStore struct {
	adapter adapters.DatabaseAdapter
}

func NewGraphRelationshipStore(adapter adapters.DatabaseAdapter) *GraphRelationshipStore {
	return &GraphRelationshipStore{adapter: adapter}
}

func (s *GraphRelationshipStore) Create(ctx context.Context, p payloads.Payload) error {
	relation := p.(*GraphRelationship)
	res, err := s.adapter.ExecContext(ctx,
		CreateGraphRelationshipMySQLQuery,
		relation.Name, relation.Slug, relation.SourceDefinitionID, relation.TargetDefinitionID, 
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	relation.ID = uint(id)
	return nil
}

func (s *GraphRelationshipStore) Delete(ctx context.Context, p payloads.Payload) error {
	relation := p.(*payloads.DefaultID)
	_, err := s.adapter.ExecContext(ctx,
		DeleteGraphRelationshipMySQLQuery,
		relation.ID,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}

