package story

import (
	"time"
	"context"
    "net/http"
	"gitlab.com/pedrokoblitz/opus-go/internal/quality"
	"gitlab.com/pedrokoblitz/opus-go/internal/adapters"
	"gitlab.com/pedrokoblitz/opus-go/internal/payloads"
)

type GraphRelationship struct {
	ID           uint
	StoryID		 uint

	Name         string
	Slug         string
	
	SourceDefinitionID uint
	SourceDefinition       string
	
	TargetDefinitionID uint
	TargetDefinition       string
	
	CreatedAt    time.Time
	UpdatedAt    time.Time
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

func CreateGraphRelationshipAction(p payloads.Payload, adapter interface{}) error {
	store := NewGraphRelationshipStore(adapter.(adapters.DatabaseAdapter))
	return store.Create(context.Background(), p)
}

func DeleteGraphRelationshipAction(p payloads.Payload, adapter interface{}) error {
	store := NewGraphRelationshipStore(adapter.(adapters.DatabaseAdapter))
	return store.Delete(context.Background(), p)
}


