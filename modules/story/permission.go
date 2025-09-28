package story

import (
	"time"
    "net/http"
	// "gitlab.com/pedrokoblitz/opus-go/actions"
	"gitlab.com/pedrokoblitz/opus-go/quality"
)

type Permission struct {
	ID               uint           `json:"id"`
	Name             string         `json:"name"`
	Action           string         `json:"action"`
	CreatedAt        time.Time      `json:"createdAt"`
}

func (p *Permission) FromRequest(r *http.Request) error {
    return nil
}

func (p *Permission) Validate() error {
	if p.Name == "" {
        return quality.ErrValidation.WithDetail("name is required")
	} 
	if p.Action == "" {
        return quality.ErrValidation.WithDetail("action is required")
	} 
	return nil
}

func (p *Permission) Process() error {
	return nil
}

// type PermissionRequirement struct {
// 	UserID uint
// 	UserSlug string
// 	ProfileNodeID uint
// 	StoryID uint
// 	StorySlug string
// 	RoleID uint
// 	RoleSlug string
// 	DefinitionID uint
// 	DefinitionSlug string
// 	Action string
// }

// func (p *PermissionRequirement) Validate() error {
// 	return nil
// }

// func (p *PermissionRequirement) Default() error {
// 	return nil
// }

// func (p *PermissionRequirement) Process() error {
// 	return nil
// }

