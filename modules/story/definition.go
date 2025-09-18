package story

import (
	"time"
	"context"
	"encoding/json"
    "net/http"
	"gitlab.com/pedrokoblitz/opus-go/internal/quality"
	"gitlab.com/pedrokoblitz/opus-go/internal/adapters"
	"gitlab.com/pedrokoblitz/opus-go/internal/payloads"
)

type Definition struct {
	ID               uint    `json:"id"`

	Name             string  `json:"name"`
	Slug             string  `json:"slug"`
	Description      string  `json:"description"`

	Icon             string  `json:"icon"`
	Color            string  `json:"color"`

	Story            string `json:"story"`
	StoryID          uint   `json:"storyId"`
	StoryModel       Story  `json:"storyModel"`

	SchemaModels     []DefinitionSchema
	Schemas          []string
	
	StatusModels     []DefinitionStatus
	Status           []string

	PossibleParents  []string

	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`

	HasTaxonomy      bool  `json:"hasTaxonomy"`
	HasInteractions  bool  `json:"hasInteractions"`
	BelongsToGraph   bool  `json:"belongsToGraph"`
}

// MarshalJSON implements json.Marshaler - controls what gets output to JSON
func (p Definition) MarshalJSON() ([]byte, error) {
    // Define exactly what fields you want in the output
    type defJson struct {
		ID       uint      `json:"id"`
		Name     string    `json:"name"`
		Slug     string    `json:"slug"`
		Description     string    `json:"description"`
		Icon             string  `json:"icon"`
		Color            string  `json:"color"`
		StoryID          uint   `json:"storyId"`
		Schemas []DefinitionSchema `json:"schemas"`
		Status     []DefinitionStatus `json:"status"`
		PossibleParents []string `json:"parents"`
    }
    
    return json.Marshal(defJson{
    	ID: p.ID,
    	Name: p.Name,
    	Slug: p.Slug,
    	Description: p.Description,
    	Icon: p.Icon,
    	Color: p.Color,
    	StoryID: p.StoryID,
    	Status: p.StatusModels,
    	Schemas: p.SchemaModels,
    	PossibleParents: p.PossibleParents,
    })
}

func (p *Definition) FromRequest(r *http.Request) error {
	p.Story = r.Context().Value("story").(string)
    return nil
}

func (p *Definition) Validate() error {
	if p.Name == "" {
        return quality.ErrValidation.WithDetail("name is required")
	} 
	// if <condition> {
    //     return quality.ErrValidation.WithDetail("<message>")
	// } 
	return nil
}

func (p *Definition) Process() error {
	p.Slug = generateSlug(p.Name)
	return nil
}

type DefinitionHierarchy struct {
	ParentDefinitionID uint   `json:"parentDefinitionID,omitempty"`
	Parent             string `json:"parentDefinition,omitempty"`
	ChildDefinitionID  uint   `json:"childDefinitionID,omitempty"`
	Child              string `json:"childDefinition,omitempty"`
}

func (p *DefinitionHierarchy) FromRequest(r *http.Request) error {
    return nil
}

func (p *DefinitionHierarchy) Validate() error {
	// if <condition> {
    //     return quality.ErrValidation.WithDetail("<message>")
	// } 
	return nil
}

func (p *DefinitionHierarchy) Process() error {
	return nil
}

type DefinitionSchema struct {
	ID       uint   `json:"id"`
	Name     string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	Schema   string `json:"schema,omitempty"`
}

func (p *DefinitionSchema) FromRequest(r *http.Request) error {
    return nil
}

func (p *DefinitionSchema) Validate() error {
	if p.Name == "" {
        return quality.ErrValidation.WithDetail("name is required")
	} 
	if p.Schema == "" {
        return quality.ErrValidation.WithDetail("schema is required")
	} 
	// if <condition> {
    //     return quality.ErrValidation.WithDetail("<message>")
	// } 

	return nil
}

func (p *DefinitionSchema) Process() error {
	return nil
}

type DefinitionStatus struct {
	ID             uint         `json:"id"`
	
	Definition   string         `json:"definition"`
	DefinitionID   uint         `json:"definitionId"`
	
	Name           string       `json:"name"`
	Slug           string       `json:"slug"`

	Icon           string       `json:"icon"`
	Color          string       `json:"color"`

	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
}

func (p *DefinitionStatus) Validate() error {
	if p.Name == "" {
        return quality.ErrValidation.WithDetail("name is required")
	} 
	if p.Definition == "" || p.DefinitionID == 0 {
        return quality.ErrValidation.WithDetail("definition is required")
	} 
	// if <condition> {
    //     return quality.ErrValidation.WithDetail("<message>")
	// } 
	return nil
}

func (p *DefinitionStatus) Process() error {
	p.Slug = generateSlug(p.Name)
	return nil
}

func CreateDefinitionAction(p payloads.Payload, adapter interface{}) error {
	var err error
	err = p.Process()
	if err != nil {
		return err
	}
	
	db := adapter.(adapters.DatabaseAdapter)
	store := NewDefinitionStore(db)
	err = store.Create(context.Background(), p)
	if err != nil {
		return err
	}
	err = store.AddStatus(context.Background(), p)
	if err != nil {
		return err
	}
	err = store.AddSchemaRelationships(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDefinitionAction(p payloads.Payload, adapter interface{}) error {
	var err error
	// err = p.Process()
	// if err != nil {
	// 	return err
	// }

	db := adapter.(adapters.DatabaseAdapter)
	if err != nil {
		return err
	}
	store := NewDefinitionStore(db)

	err = store.Update(context.Background(), p)
	if err != nil {
		return err
	}
	err = store.SyncStatus(context.Background(), p)
	if err != nil {
		return err
	}
	err = store.SyncSchemaRelationships(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDefinitionAction(p payloads.Payload, adapter interface{}) error {
	store := NewDefinitionStore(adapter.(adapters.DatabaseAdapter))
	return store.Delete(context.Background(), p)
}
