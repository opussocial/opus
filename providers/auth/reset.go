package auth

import (
	// "fmt"
	"context"
	"net/mail"
    "net/http"
	// "database/sql"
	// "unicode/utf8"
	"gitlab.com/pedrokoblitz/opus-go/internal/payloads"
	"gitlab.com/pedrokoblitz/opus-go/internal/adapters"
	"gitlab.com/pedrokoblitz/opus-go/internal/quality"
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
	// if utf8.RuneCountInString(p.Password) < 8 {
	// 	return quality.ErrValidation.WithDetail("password must be at least 8 characters")
	// }
	return nil
}

func (p *Reset) Process() error {
	password, err := EncryptPassword(p.Password)
	p.Password = password
	return err
}


func ResetPasswordAction(p payloads.Payload, adapter interface{}) error {
	store := NewResetStore(adapter.(adapters.DatabaseAdapter))
	err := p.Process()
	if err != nil {
		return err
	}
	err = store.DeleteResetToken(context.Background(), p)
	if err != nil {
		return err
	}
	return store.ResetPassword(context.Background(), p)
}

func RequestPasswordResetAction(p payloads.Payload, adapter interface{}) error {
	store := NewResetStore(adapter.(adapters.DatabaseAdapter))
	err := p.Process()
	if err != nil {
		return err
	}
	return store.CreateResetToken(context.Background(), p)
}
