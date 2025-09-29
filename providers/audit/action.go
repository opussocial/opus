package audit

import (
	"time"
)

type ActionRecord struct {
	Name string
	Payload string
	PayloadData []byte
	Status string
	CreatedAt time.Time
	UpdatedAt time.Time
}
