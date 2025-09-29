package elements

import (
	"context"
	"database/sql"
	"gitlab.com/pedrokoblitz/opus-go/quality"
	"gitlab.com/pedrokoblitz/opus-go/actions"
)

/**
 * 
 * GRAPH STORE
 * 
 * 
 **/

const (
	CreateGraphMySQLQuery = "INSERT INTO graph (element_profile_id, element_source_id, element_target_id) VALUES (?, ?, ?)"
	DeleteGraphMySQLQuery = "DELETE FROM graph WHERE id = ?"
)

type GraphStore struct {
	adapter *sql.DB
}

func NewGraphStore(adapter *sql.DB) *GraphStore {
	return &GraphStore{adapter: adapter}
}

func (s *GraphStore) Create(ctx context.Context, p actions.Payload) error {
	graph := p.(*Graph)
	res, err := s.adapter.ExecContext(ctx,
		CreateGraphMySQLQuery,
		graph.SourceID, graph.SourceID, graph.SourceID,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	graph.ID = uint(id)

	return nil
}

func (s *GraphStore) Delete(ctx context.Context, p actions.Payload) error {
	graph := p.(*actions.DefaultID)
	_, err := s.adapter.ExecContext(ctx,
		DeleteGraphMySQLQuery,
		graph.ID,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}
