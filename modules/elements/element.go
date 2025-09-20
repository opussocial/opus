package elements

import (
	"time"
	"context"
	"gitlab.com/pedrokoblitz/opus-go/internal/quality"
	"gitlab.com/pedrokoblitz/opus-go/internal/payloads"
	"gitlab.com/pedrokoblitz/opus-go/internal/adapters"
	"gitlab.com/pedrokoblitz/opus-go/modules/story"
)

type Element struct {
	ID       uint   `json:"id"`
	Slug     string `json:"slug"`

	Level    uint   `json:"level"`
	Active   bool   `json:"active"`
	Locked   bool   `json:"locked"`

	ParentID *uint  `json:"parentID,omitempty"`
	Parent *Element  `json:"parent,omitempty"`
	Children []Element  `json:"children,omitempty"`

	DefinitionID uint           `json:"definitionID"`
	DefinitionModel   story.Definition `json:"definition"`
	Definition   string `json:"definition"`
	StatusID uint           `json:"statusID"`
	Status string           `json:"status"`
	StatusModel   story.DefinitionStatus `json:"statusModel"`

	CreatedByProfile    string          `json:"createdByProfile"`
	CreatedByProfileID    uint          `json:"createdByProfileID"`
	CreatedByProfileModel      *Element          `json:"createdByProfileModel"`

	// SCHEMAS
	File          *FileSchema         `json:"file,omitempty"`
	Text     	  *TextSchema     `json:"text,omitempty"`
	Person        *PersonSchema      `json:"person,omitempty"`
	ContactPoint  *ContactPointSchema      `json:"contactPoint,omitempty"`
	PostalAddress *PostalAddressSchema      `json:"postalAddress,omitempty"`
	TimeTracking  *TimeTrackingSchema `json:"timeTracking,omitempty"`
	WebResource   *WebResourceSchema      `json:"webResource,omitempty"`

	// END SCHEMAS
	Interactions []Interaction `json:"interactions,omitempty"`

	// RENAME TERMS TO KEYWORDS
	// Keywords see : https://schema.org/keywords
	// Keywords or tags used to describe this content. Multiple entries in a keywords list are typically delimited by commas.
	// types : Text
	Keywords []string `json:"keywords,omitempty"`

	// MainEntityOfPage see : https://schema.org/mainEntityOfPage
	// Indicates a page (or other CreativeWork) for which this thing is the main entity being described. See background notes (see: https://schema.org/docs/datamodel.html#mainEntityBackground) for details. Inverse property: mainEntity (see: https://schema.org/mainEntity).
	// types : CreativeWork URL
	MainEntity string `json:"mainEntity,omitempty"`

	// Name see : https://schema.org/name
	// The name of the item.
	// types : Text
	Name string `json:"name,omitempty"`

	// AlternateName see : https://schema.org/alternateName
	// An alias for the item.
	// types : Text
	AlternateName string `json:"alternateName,omitempty"`

	// Description see : https://schema.org/description
	// A description of the item.
	// types : Text
	Description string `json:"description,omitempty"`

	// SameAs see : https://schema.org/sameAs
	// URL of a reference Web page that unambiguously indicates the item&#39;s identity. E.g. the URL of the item&#39;s Wikipedia page, Wikidata entry, or official website.
	// types : URL
	SameAs string `json:"sameAs,omitempty"`

	// Url see : https://schema.org/url
	// URL of the item.
	// types : URL
	Url string `json:"url,omitempty"`

	// Version see : https://schema.org/version
	// The version of the CreativeWork embodied by a specified resource.
	// types : Number Text
	// Version uint `json:"version,omitempty"`

	// InLanguage see : https://schema.org/inLanguage
	// The language of the content or performance or used in an action. Please use one of the language codes from the IETF BCP 47 standard (see: https://schema.orghttp://tools.ietf.org/html/bcp47). See also availableLanguage (see: https://schema.org/availableLanguage). Supersedes language (see: https://schema.org/language).
	// types : Language Text
	// InLanguage string `json:"inLanguage,omitempty"`

	// IsAccessibleForFree see : https://schema.org/isAccessibleForFree
	// A flag to signal that the item, event, or place is accessible for free. Supersedes free (see: https://schema.org/free).
	// types : Boolean
	// IsAccessibleForFree bool `json:"isAccessibleForFree,omitempty"`

	// IsBasedOn see : https://schema.org/isBasedOn
	// A resource that was used in the creation of this resource. This term can be repeated for multiple sources. For example, http://example.com/great-multiplication-intro.html. Supersedes isBasedOnUrl (see: https://schema.org/isBasedOnUrl).
	// types : URL
	// IsBasedOn string `json:"isBasedOn,omitempty"`

	// DateCreated see : https://schema.org/dateCreated
	// The date on which the CreativeWork was created or the item was added to a DataFeed.
	// types : Date DateTime
	CreatedAt time.Time `js		DefinitionID: 1,
on:"dateCreated,omitempty"`

	// DateModified see : https://schema.org/dateModified
	// The date on which the CreativeWork was most recently modified or when the item&#39;s entry was modified within a DataFeed.
	// types : Date DateTime
	UpdatedAt time.Time `json:"dateModified,omitempty"`

	// DatePublished see : https://schema.org/datePublished
	// Date of first broadcast/publication.
	// types : Date
	PublishedAt time.Time `json:"datePublished,omitempty"`

	// Expires see : https://schema.org/expires
	// Date the content expires and is no longer useful or available. For example a VideoObject (see: https://schema.org/VideoObject) or NewsArticle (see: https://schema.org/NewsArticle) whose availability or relevance is time-limited, or a ClaimReview (see: https://schema.org/ClaimReview) fact check whose publisher wants to indicate that it may no longer be relevant (or helpful to highlight) after some date.
	// types : Date
	ExpiresAt time.Time `json:"expires,omitempty"`

	DeletedAt time.Time // NO JSON FOR THIS `json:"deletedAt,omitempty"`
}

func (p *Element) Validate() error {
	if p.Name == "" {
        return quality.ErrValidation.WithDetail("name is required")
	} 
	if p.DefinitionID == 0 && p.DefinitionModel.ID == 0 && p.Definition == "" {
        return quality.ErrValidation.WithDetail("definition is required")
	}
	if  p.CreatedByProfileID == 0 {
        return quality.ErrValidation.WithDetail("profile is required")
	}
	return nil
}

func (p *Element) Process() error {
	return nil
}

func CreateElementAction(p payloads.Payload, adapter interface{}) error {
	store := NewElementStore(adapter.(adapters.DatabaseAdapter))
	return store.Create(context.Background(), p)
}

func DeleteElementAction(p payloads.Payload, adapter interface{}) error {
	store := NewElementStore(adapter.(adapters.DatabaseAdapter))
	return store.Delete(context.Background(), p)
}
