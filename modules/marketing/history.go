package marketing

import (
	"time"
	"gitlab.com/pedrokoblitz/opus-go/internal/payloads"
)

type History struct {
	ID uint
	OperationID uint
	Scope string
	ExecutedAt time.Time
}

func (h *History) Validate() error { return nil }
func (h *History) Process() error { return nil }

//
func CreateHistoryAction(p payloads.Payload, adapter interface{}) error { return nil }

//
func DeleteHistoryAction(p payloads.Payload, adapter interface{}) error { return nil }
