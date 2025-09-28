package story

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"database/sql"

	"gitlab.com/pedrokoblitz/opus-go/actions"
	"gitlab.com/pedrokoblitz/opus-go/quality"
)

type RoleResults []Role

func (p *RoleResults) FromRequest(r *http.Request) error {
	return nil
}

func (p *RoleResults) Validate() error {
	return nil
}

func (p *RoleResults) Process() error {
	return nil
}

type Role struct {
	ID uint `json:"id"`

	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`

	StoryID    uint   `json:"storyId"`
	Story      string `json:"story"`
	StoryModel Story  `json:"story"`

	PermissionModels []Permission `json:"permissionModels"`
	Permissions      []string     `json:"permissions"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// MarshalJSON implements json.Marshaler - controls what gets output to JSON
func (p Role) MarshalJSON() ([]byte, error) {
	type roleJson struct {
		ID          uint         `json:"id"`
		Name        string       `json:"name"`
		Slug        string       `json:"slug"`
		Description string       `json:"description"`
		Permissions []Permission `json:"permissions"`
	}

	return json.Marshal(roleJson{
		ID:          p.ID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Permissions: p.PermissionModels,
	})
}

func (p *Role) FromRequest(r *http.Request) error {
	p.Story = r.Context().Value("story").(string)
	return nil
}

func (p *Role) Validate() error {
	if p.Name == "" {
		return quality.ErrValidation.WithDetail("name is required")
	}
	if p.StoryID == 0 && p.Story == "" {
		return quality.ErrValidation.WithDetail("story is required")
	}
	if len(p.PermissionModels) > 0 {
		for _, p := range p.PermissionModels {
			err := p.Validate()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Role) Process() error {
	p.Slug = generateSlug(p.Name)
	return nil
}

func CreateRoleAction(p actions.Payload, container actions.Container) error {
	var err error
	db := adapter.(*sql.DB)

	store := NewRoleStore(db)
	err = p.Process()
	if err != nil {
		return err
	}

	err = store.Create(context.Background(), p)
	if err != nil {
		return err
	}

	err = store.AddPermissionRelationships(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}

func UpdateRoleAction(p actions.Payload, container actions.Container) error {
	var err error
	db := adapter.(*sql.DB)

	store := NewRoleStore(db)
	err = p.Process()
	if err != nil {
		return err
	}

	err = store.Update(context.Background(), p)
	if err != nil {
		return err
	}

	// Use sync to ensure permission relationships match the updated role
	err = store.SyncPermissionRelationships(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRoleAction(p actions.Payload, container actions.Container) error {
	store := NewRoleStore(container.DB())
	return store.Delete(context.Background(), p)
}
