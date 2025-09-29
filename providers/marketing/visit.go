package marketing

import (
	"time"

	"gitlab.com/pedrokoblitz/opus-go/actions"
)

type Visit struct {
	ID          uint
	LeadID      uint
	VisitDate   time.Time
	PageUrl     string
	UtmSource   string
	UtmMedium   string
	UtmCampaign string
	ReferrerUrl string
	UserAgent   string
	IpAddress   string
}

func (h *Visit) Validate() error { return nil }
func (h *Visit) Process() error  { return nil }

func CreateVisitAction(p actions.Payload, container actions.Container) error { return nil }

func DeleteVisitAction(p actions.Payload, container actions.Container) error { return nil }
