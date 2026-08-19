package actions

import (
	"time"
	"net/http"
	"github.com/opussocialcontent/opus-go/quality"	
)

type Validatable interface {
    Validate() error
}

type Processable interface {
    Process() error
}

type RequestAware interface {
    FromRequest(r *http.Request) error
}

// Enhanced Payload interface
type Payload interface {
    Validatable
    Processable
	RequestAware
}

type DefaultPayload struct {
	ID uint `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
	Description string `json:"name"`
	Scope string `json:"scope"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (p *DefaultPayload) Process() error { return nil }

func (p *DefaultPayload) FromRequest(r *http.Request) error { return nil }

func (p *DefaultPayload) Validate() error {
	if p.ID == 0 {
		return quality.ErrValidation.WithDetail("ID is required")
	}
	return nil
}

type DefaultID struct {
	ID uint `json:"id"`
}

func (p *DefaultID) FromRequest(r *http.Request) error {
	return nil
}

func (p *DefaultID) Process() error { return nil }

func (p *DefaultID) Validate() error {
	if p.ID == 0 {
		return quality.ErrValidation.WithDetail("ID is required")
	}
	return nil
}
