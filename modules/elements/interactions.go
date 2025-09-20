package elements

import (
	"time"
	"context"
	"gitlab.com/pedrokoblitz/opus-go/internal/quality"
	"gitlab.com/pedrokoblitz/opus-go/internal/adapters"
	"gitlab.com/pedrokoblitz/opus-go/internal/payloads"
)

type Interaction struct {
	ID             uint        `json:"id"`
	Scope    string        `json:"scope"`
	ElementID uint      `json:"elementId"`
	CreatedByProfileID uint      `json:"profileId"`
	Rating uint      `json:"rating"`
	Comment string      `json:"comment"`
}

func (p *Interaction) Validate() error {
	if p.Scope == "" {
        return quality.ErrValidation.WithDetail("scope is required")
	} 
	return nil
}

func (p *Interaction) Process() error {
	return nil
}

func CreateInteractionAction(p payloads.Payload, adapter interface{}) error {
	store := NewInteractionStore(adapter.(adapters.DatabaseAdapter))
	return store.Create(context.Background(), p)
}

func DeleteInteractionAction(p payloads.Payload, adapter interface{}) error {
	store := NewInteractionStore(adapter.(adapters.DatabaseAdapter))
	return store.Delete(context.Background(), p)
}

type Keyword struct {
	ID        uint      `json:"id"`
	Term      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (p *Keyword) Validate() error {
	if p.Term == "" {
        return quality.ErrValidation.WithDetail("term is required")
	} 
	return nil
}

func (p *Keyword) Process() error {
	return nil
}

func CreateKeywordAction(p payloads.Payload, adapter interface{}) error {
	store := NewKeywordStore(adapter.(adapters.DatabaseAdapter))
	return store.Create(context.Background(), p)
}

func DeleteKeywordAction(p payloads.Payload, adapter interface{}) error {
	store := NewKeywordStore(adapter.(adapters.DatabaseAdapter))
	return store.Delete(context.Background(), p)
}
