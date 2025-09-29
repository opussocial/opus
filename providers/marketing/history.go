package marketing

import (
	"time"

	"gitlab.com/pedrokoblitz/opus-go/actions"
)

type History struct {
	ID          uint
	OperationID uint
	Scope       string
	ExecutedAt  time.Time
}

func (h *History) Validate() error { return nil }
func (h *History) Process() error  { return nil }

func CreateHistoryAction(p actions.Payload, container actions.Container) error { return nil }

func DeleteHistoryAction(p actions.Payload, container actions.Container) error { return nil }
