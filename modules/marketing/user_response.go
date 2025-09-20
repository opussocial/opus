package marketing

import (
	"time"
	// "gitlab.com/pedrokoblitz/opus-go/internal/payloads"
)

type UserResponse struct {
	ID uint
	LeadID uint
	OperationID uint
	Scope string
	RespondedAt time.Time
}

func (h *UserResponse) Validate() error { return nil }
func (h *UserResponse) Process() error { return nil }
