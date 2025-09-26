package elements

import (
	"context"

	"gitlab.com/pedrokoblitz/opus-go/actions"
	"gitlab.com/pedrokoblitz/opus-go/quality"
)

type ProfileRelationship struct {
	ID               uint   `json:"id"`
	User             string `json:"user"`
	UserID           uint   `json:"userID"`
	Role             string `json:"role"`
	RoleID           uint   `json:"roleID"`
	ProfileElement   string `json:"profileElement"`
	ProfileElementID uint   `json:"profileElementID"`
}

func (p *ProfileRelationship) Validate() error {
	if p.UserID == 0 {
		return quality.ErrValidation.WithDetail("user is required")
	}
	return nil
}

func (p *ProfileRelationship) Process() error {
	return nil
}

func CreateProfileRelationshipAction(p actions.Payload, container actions.Container) error {
	store := NewProfileRelationshipStore(container.DB())
	return store.Create(context.Background(), p)
}

func DeleteProfileRelationshipAction(p actions.Payload, container actions.Container) error {
	store := NewProfileRelationshipStore(container.DB())
	return store.Delete(context.Background(), p)
}
