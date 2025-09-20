package marketing

import (
	"time"
	// "gitlab.com/pedrokoblitz/opus-go/internal/payloads"
)

type Payoff struct {
	ID uint
	OperationID uint
	Points uint
	CalculatedAt time.Time
}

func (h *Payoff) Validate() error { 
	return nil
}

func (h *Payoff) Process() error { return nil }
