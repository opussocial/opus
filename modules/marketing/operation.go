package marketing

import (
	"gitlab.com/pedrokoblitz/opus-go/internal/payloads"
)

type Operation struct {
	ID uint
	Name string
	Slug string
	Scope string
}

func (h *Operation) Validate() error { return nil }
func (h *Operation) Process() error { return nil }

//
func CreateOperationAction(p payloads.Payload, adapter interface{}) error { return nil }

//
func DeleteOperationAction(p payloads.Payload, adapter interface{}) error { return nil }
