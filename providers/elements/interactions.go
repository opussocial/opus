package elements

import (
	"context"
	"time"

	"gitlab.com/pedrokoblitz/opus-go/actions"
	"gitlab.com/pedrokoblitz/opus-go/quality"
)

type Interaction struct {
	ID                 uint   `json:"id"`
	Scope              string `json:"scope"`
	ElementID          uint   `json:"elementId"`
	CreatedByProfileID uint   `json:"profileId"`
	Rating             uint   `json:"rating"`
	Comment            string `json:"comment"`
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

func CreateInteractionAction(p actions.Payload, container actions.Container) error {
	store := NewInteractionStore(container.DB())
	return store.Create(context.Background(), p)
}

func DeleteInteractionAction(p actions.Payload, container actions.Container) error {
	store := NewInteractionStore(container.DB())
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

func CreateKeywordAction(p actions.Payload, container actions.Container) error {
	store := NewKeywordStore(container.DB())
	return store.Create(context.Background(), p)
}

func DeleteKeywordAction(p actions.Payload, container actions.Container) error {
	store := NewKeywordStore(container.DB())
	return store.Delete(context.Background(), p)
}
