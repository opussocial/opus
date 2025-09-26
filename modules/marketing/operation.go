package marketing

import "gitlab.com/pedrokoblitz/opus-go/actions"

type Operation struct {
	ID    uint
	Name  string
	Slug  string
	Scope string
}

func (h *Operation) Validate() error { return nil }
func (h *Operation) Process() error  { return nil }

func CreateOperationAction(p actions.Payload, container actions.Container) error { return nil }

func DeleteOperationAction(p actions.Payload, container actions.Container) error { return nil }
