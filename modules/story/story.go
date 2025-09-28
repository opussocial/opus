 package story

import (
    // "fmt"
    "time"
    "encoding/json"
	"context"
    "net/http"
 
	"gitlab.com/pedrokoblitz/opus-go/actions"
	"gitlab.com/pedrokoblitz/opus-go/quality"
)

/**
 * 
 * APP PAYLOADS
 * 
 * 
 **/
type StoryResults struct {
	Filters map[string]interface{} `json:"filters"`
	Results []Story `json:"results'"`
}

func (p StoryResults) MarshalJSON() ([]byte, error) {
    return json.Marshal(p.Results)
}

func (p *StoryResults) FromRequest(r *http.Request) error {
    return nil
}

func (p *StoryResults) Validate() error { return nil }

func (p *StoryResults) Process() error { return nil }

type Story struct {
	ID       uint      `json:"id"`
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	Description     string    `json:"description"`
	Active   bool      `json:"active"`
	Relationships []string `json:"relationships"`
	RelationshipModels []GraphRelationship `json:"relationshipSlugs"`
	Settings []string `json:"settings"`
	SettingModels []Setting `json:"settingModels"`
	Roles []string `json:"roles"`
	RoleModels []Role `json:"roleModels"`
	Definitions []string `json:"definitions"`
	DefinitionModels []Definition `json:"definitionModels"`
	CreatedAt time.Time `json:"createdAt"`
 	UpdatedAt time.Time `json:"updatedAt"`
}

// MarshalJSON implements json.Marshaler - controls what gets output to JSON
func (p Story) MarshalJSON() ([]byte, error) {
    // Define exactly what fields you want in the output
    type storyJson struct {
		ID       uint      `json:"id"`
		Name     string    `json:"name"`
		Slug     string    `json:"slug"`
		Description     string    `json:"description"`
		Active   bool      `json:"active"`
		Relationships []GraphRelationship `json:"relationships"`
		Settings []Setting `json:"settings"`
		Roles []Role `json:"roles"`
		Definitions []Definition `json:"definitions"`
    }
    
    return json.Marshal(storyJson{
    	ID: p.ID,
    	Name: p.Name,
    	Slug: p.Slug,
    	Description: p.Description,
    	Active: p.Active,
    	Relationships: p.RelationshipModels,
    	Settings: p.SettingModels,
    	Roles: p.RoleModels,
    	Definitions: p.DefinitionModels,
    })
}

func (p *Story) FromRequest(r *http.Request) error {
	slug, ok := r.Context().Value("story").(string)
	if (!ok) {
		return nil
	}
	p.Slug = slug
    return nil
}

func (p *Story) Validate() error {
	if p.Name == "" {
        return quality.ErrValidation.WithDetail("name is required")
	} 
	if p.Slug == "" {
        return quality.ErrValidation.WithDetail("slug is required")
	} 
	return nil
}

func (p *Story) Process() error {
	var newRoles []Role
	for _, r := range p.RoleModels {
		err := r.Process()
		if err != nil {
			return err
		}
		newRoles = append(newRoles, r)
	}
	p.RoleModels = newRoles

	var newDefinitions []Definition
	for _, r := range p.DefinitionModels {
		err := r.Process()
		if err != nil {
			return err
		}
		newDefinitions = append(newDefinitions, r)
	}
	p.DefinitionModels = newDefinitions

	var newSettings []Setting
	for _, r := range p.SettingModels {
		err := r.Process()
		if err != nil {
			return err
		}
		newSettings = append(newSettings, r)
	}
	p.SettingModels = newSettings

	var newRelationships []GraphRelationship
	for _, r := range p.RelationshipModels {
		err := r.Process()
		if err != nil {
			return err
		}
		newRelationships = append(newRelationships, r)
	}
	p.RelationshipModels = newRelationships

	return nil
}

/**
 * 
 * STORY ACTIONS
 * 
 * 
 **/
func CreateStoryAction(p actions.Payload, container actions.Container) error {
	var err error
	err = p.Process()
	if err != nil {
		return err
	}
	db := container.DB()
	if err != nil {
		return err
	}

	store := NewStoryStore(db)
	if err != nil {
		return err
	}
	err = store.Create(context.Background(), p)
	if err != nil {
		return err
	}
	err = store.AddDefinitions(context.Background(), p)
	if err != nil {
		return err
	}
	err = store.AddRoles(context.Background(), p)
	if err != nil {
		return err
	}
	err = store.AddSettings(context.Background(), p)
	if err != nil {
		return err
	}
	err = store.AddRelationships(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}

func UpdateStoryAction(p actions.Payload, container actions.Container) error {
	var err error
	// TODO: CHECK IF NECESSARY
	// err = p.Process()
	// if err != nil {
	// 	return err
	// }

	db := container.DB()
	if err != nil {
		return err
	}

	store := NewStoryStore(db)
	err = store.Update(context.Background(), p)
	if err != nil {
		return err
	}
	err = store.AddDefinitions(context.Background(), p)
	if err != nil {
		return err
	}
	err = store.AddRoles(context.Background(), p)
	if err != nil {
		return err
	}
	err = store.AddSettings(context.Background(), p)
	if err != nil {
		return err
	}
	err = store.AddRelationships(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}

func ShowStoryAction(p actions.Payload, container actions.Container) error {
	var err error
	store := NewStoryStore(container.DB())
	err = store.GetBySlug(context.Background(), p)
	if err != nil {
		return err
	}
	err = store.GetRoles(context.Background(), p)
	if err != nil {
		return err
	}
	err = store.GetDefinitions(context.Background(), p)
	if err != nil {
		return err
	}
	err = store.GetSettings(context.Background(), p)
	if err != nil {
		return err
	}
	// err = store.GetRelationships(context.Background(), p)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func IndexStoryAction(p actions.Payload, container actions.Container) error {
	var err error
	store := NewStoryStore(container.DB())
	err = store.Index(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}

func DeleteStoryAction(p actions.Payload, container actions.Container) error {
	store := NewStoryStore(container.DB())
	return store.Delete(context.Background(), p)
}
