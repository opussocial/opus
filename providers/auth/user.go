package auth

import (
	"context"

	// "strconv"
	"encoding/json"
	"net/http"
	"net/mail"

	// "database/sql"
	"time"

	"github.com/opussocialcontent/opus-go/actions"
	"github.com/opussocialcontent/opus-go/quality"
)

type User struct {
	AuthBasePayload
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MarshalJSON implements json.Marshaler - controls what gets output to JSON
func (u User) MarshalJSON() ([]byte, error) {
	// Define exactly what fields you want in the output
	type userJson struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
		// Token is intentionally omitted from JSON output
	}

	return json.Marshal(userJson{
		ID:    u.ID,
		Email: u.Email,
	})
}

// UnmarshalJSON implements json.Unmarshaler - controls what gets read from JSON
func (u *User) UnmarshalJSON(data []byte) error {
	// Define what fields you accept from input
	type userInput struct {
		Email string `json:"email"`
		// ID and Token are ignored from input
	}

	var input userInput
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}

	u.Email = input.Email
	// u.ID and u.Token remain unchanged (or set to defaults)
	return nil
}

func (p *User) FromRequest(r *http.Request) error {
	id, ok := r.Context().Value("userID").(uint)
	if !ok {
		return nil
	}
	p.ID = id
	return nil
}

func (p *User) Validate() error {
	if p.ID == 0 {
		return quality.ErrValidation.WithDetail("id is required")
	}
	if p.Email != "" {
		if _, err := mail.ParseAddress(p.Email); err != nil {
			return quality.ErrValidation.WithDetail("invalid email format")
		}
	}
	return nil
}

func (p *User) Process() error {
	return nil
}

func ShowUserAction(p actions.Payload, container actions.Container) error {
	store := NewUserStore(container.DB())
	return store.GetByID(context.Background(), p)
}

func UpdateUserAction(p actions.Payload, container actions.Container) error {
	store := NewUserStore(container.DB())
	return store.UpdateUser(context.Background(), p)
}

func RemoveAccountAction(p actions.Payload, container actions.Container) error {
	store := NewUserStore(container.DB())
	return store.RemoveAccount(context.Background(), p)
}
