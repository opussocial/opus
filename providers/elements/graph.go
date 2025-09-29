package elements

import (
	"context"
	"time"

	"gitlab.com/pedrokoblitz/opus-go/actions"
	"gitlab.com/pedrokoblitz/opus-go/quality"
	"gitlab.com/pedrokoblitz/opus-go/modules/story"
)

type Graph struct {
	ID                  uint
	CreatedAt           time.Time
	UpdatedAt           time.Time
	GraphRelationship   story.GraphRelationship
	GraphRelationshipID uint
	Profile             Element
	ProfileID           uint
	Source              Element
	SourceID            uint
	Target              Element
	TargetID            uint
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

func CreateGraphAction(p actions.Payload, container actions.Container) error {
	store := NewGraphStore(container.DB())
	return store.Create(context.Background(), p)
}

func DeleteGraphAction(p actions.Payload, container actions.Container) error {
	store := NewGraphStore(container.DB())
	return store.Delete(context.Background(), p)
}
