package marketing

import (
	"time"

	"gitlab.com/pedrokoblitz/opus-go/actions"
)

type Lead struct {
	ID         uint
	UserID     uint
	Name       string
	Email      string
	Phone      string
	Region     string
	Source     string
	Cookie     string
	SignUpDate time.Time
	CreateAt   time.Time
}

func (h *Lead) Validate() error { return nil }
func (h *Lead) Process() error  { return nil }

func CreateLeadAction(p actions.Payload, container actions.Container) error { return nil }

func DeleteLeadAction(p actions.Payload, container actions.Container) error { return nil }
