package marketing

import (
	"time"
	"gitlab.com/pedrokoblitz/opus-go/internal/payloads"
)

type Lead struct {
	ID uint
	UserID uint
	Name string
	Email string
	Phone string
	Region string
	Source string
	Cookie string
	SignUpDate time.Time
	CreateAt time.Time
}

func (h *Lead) Validate() error { return nil }
func (h *Lead) Process() error { return nil }

//
func CreateLeadAction(p payloads.Payload, adapter interface{}) error { return nil }

//
func DeleteLeadAction(p payloads.Payload, adapter interface{}) error { return nil }
