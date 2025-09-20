package elements

import (
	"time"
	"context"
	"gitlab.com/pedrokoblitz/opus-go/internal/quality"
	"gitlab.com/pedrokoblitz/opus-go/internal/adapters"
	"gitlab.com/pedrokoblitz/opus-go/internal/payloads"
	"gitlab.com/pedrokoblitz/opus-go/modules/story"
)

type Graph struct {
	ID                  uint
	CreatedAt           time.Time
	UpdatedAt           time.Time
	GraphRelationship   story.GraphRelationship
	GraphRelationshipID uint
	Profile          Element
	ProfileID        uint
	Source          Element
	SourceID        uint
	Target          Element
	TargetID        uint
}

func (p *Graph) Validate() error {
	if p.TargetID == 0 {
        return quality.ErrValidation.WithDetail("target is required")
	} 
	return nil
}

func (p *Graph) Process() error {
	return nil
}

func CreateGraphAction(p payloads.Payload, adapter interface{}) error {
	store := NewGraphStore(adapter.(adapters.DatabaseAdapter))
	return store.Create(context.Background(), p)
}

func DeleteGraphAction(p payloads.Payload, adapter interface{}) error {
	store := NewGraphStore(adapter.(adapters.DatabaseAdapter))
	return store.Delete(context.Background(), p)
}

