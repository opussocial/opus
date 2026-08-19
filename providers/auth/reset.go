package auth

import (
	// "fmt"
	"context"
	"net/http"
	"net/mail"
	"unicode/utf8"

	// "database/sql"

	"github.com/opussocialcontent/opus-go/actions"
	"github.com/opussocialcontent/opus-go/quality"
)

type ResetRequest struct {
	AuthBasePayload
}

func (p *ResetRequest) FromRequest(r *http.Request) error {
	// p.ID = r.Context().GetValue("userID").(uint)
	return nil
}

func (p *ResetRequest) Validate() error {
	if p.Email == "" {
		return quality.ErrValidation.WithDetail("email is required")
	}
	if _, err := mail.ParseAddress(p.Email); err != nil {
		return quality.ErrValidation.WithDetail("invalid email format")
	}
	return nil
}

func (p *ResetRequest) Process() error {
	p.Token = GenerateRandomString(16)
	return nil
}

type Reset struct {
	AuthBasePayload
}

func (p *Reset) FromRequest(r *http.Request) error {
	token := r.URL.Query().Get("reset")
	p.Token = token
	return nil
}

func (p *Reset) Validate() error {
	if utf8.RuneCountInString(p.Password) < 8 {
		return quality.ErrValidation.WithDetail("password must be at least 8 characters")
	}
	return nil
}

func (p *Reset) Process() error {
	password, err := EncryptPassword(p.Password)
	p.Password = password
	return err
}

func ResetPasswordAction(p actions.Payload, container actions.Container) error {
	store := NewResetStore(container.DB())
	err := p.Process()
	if err != nil {
		return err
	}
	// Reset the password first (the query uses the token to find the user),
	// then delete the token so it can't be reused.
	err = store.ResetPassword(context.Background(), p)
	if err != nil {
		return err
	}
	return store.DeleteResetToken(context.Background(), p)
}

func RequestPasswordResetAction(p actions.Payload, container actions.Container) error {
	store := NewResetStore(container.DB())
	req := p.(*ResetRequest)
	// Look up the user by email to obtain their ID
	id, err := store.GetUserIDByEmail(context.Background(), req.Email)
	if err != nil {
		return err
	}
	req.ID = id
	err = p.Process()
	if err != nil {
		return err
	}
	return store.CreateResetToken(context.Background(), p)
}
