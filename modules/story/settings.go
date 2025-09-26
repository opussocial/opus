package story

import (
	"time"
    "net/http"
	"gitlab.com/pedrokoblitz/opus-go/quality"
)

type SettingResults []Setting 

func (p *SettingResults) Validate() error {
	return nil
}

type Setting struct {
	ID          uint         `json:"id"`
	UserID     uint         `json:"userId"`
	StoryID     uint         `json:"storyId"`
	Name        string       `json:"name"`
	Slug        string       `json:"slug"`
	Value       string       `json:"value"`
	Scope       string       `json:"scope"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

func (p *Setting) FromRequest(r *http.Request) error {
    return nil
}

func (p *Setting) Validate() error {
	if p.Name == "" {
        return quality.ErrValidation.WithDetail("name is required")
	} 
	if p.Scope == "" {
        return quality.ErrValidation.WithDetail("scope is required")
	} 
	return nil
}

func (p *Setting) Process() error {
	p.Slug = generateSlug(p.Name)
	return nil
}
