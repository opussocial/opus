package story

import (
	"context"
	"net/http"
	"time"

	"gitlab.com/pedrokoblitz/opus-go/actions"
	"gitlab.com/pedrokoblitz/opus-go/quality"
)

type GraphRelationship struct {
	ID      uint
	StoryID uint

	Name string
	Slug string

	SourceDefinitionID uint
	SourceDefinition   string

	TargetDefinitionID uint
	TargetDefinition   string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *GraphRelationship) FromRequest(r *http.Request) error {
	return nil
}

func (p *GraphRelationship) Validate() error {
	if p.Name == "" {
		return quality.ErrValidation.WithDetail("name is required")
	}
	if p.SourceDefinition == "" && p.SourceDefinitionID == 0 {
		return quality.ErrValidation.WithDetail("source is required")
	}
	if p.TargetDefinition == "" && p.TargetDefinitionID == 0 {
		return quality.ErrValidation.WithDetail("target is required")
	}
	return nil
}

func (p *GraphRelationship) Process() error {
	p.Slug = generateSlug(p.Name)
	return nil
}

func CreateGraphRelationshipAction(p actions.Payload, container actions.Container) error {
	store := NewGraphRelationshipStore(container.DB())
	return store.Create(context.Background(), p)
}

func DeleteGraphRelationshipAction(p actions.Payload, container actions.Container) error {
	store := NewGraphRelationshipStore(container.DB())
	return store.Delete(context.Background(), p)
}
