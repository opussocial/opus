package auth

import (
	"context"
	"fmt"

	// "database/sql"
	"net/http"
	"net/mail"
	"unicode/utf8"

	"github.com/opussocialcontent/opus-go/actions"
	"github.com/opussocialcontent/opus-go/quality"
)

type SignUp struct {
	AuthBasePayload
}

func (p *SignUp) FromRequest(r *http.Request) error {
	return nil
}

func (p *SignUp) Validate() error {
	if p.Email == "" {
		return quality.ErrValidation.WithDetail("email is required")
	}
	if _, err := mail.ParseAddress(p.Email); err != nil {
		return quality.ErrValidation.WithDetail("invalid email format")
	}
	if utf8.RuneCountInString(p.Password) < 8 {
		return quality.ErrValidation.WithDetail("password must be at least 8 characters")
	}
	return nil
}

func (p *SignUp) Process() error {
	password, err := EncryptPassword(p.Password)
	p.Password = password
	if err != nil {
		return err
	}
	p.Token = GenerateRandomString(16)
	return nil
}

func ConfirmAction(p actions.Payload, container actions.Container) error {
	store := NewSignUpStore(container.DB())
	err := store.Confirm(context.Background(), p)
	if err != nil {
		return err
	}
	err = store.DeleteConfirmationToken(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}

func SignUpAction(p actions.Payload, container actions.Container) error {
	var err error
	store := NewSignUpStore(container.DB())
	err = p.Process()
	if err != nil {
		return err
	}
	err = store.SignUp(context.Background(), p)
	if err != nil {
		return err
	}
	return store.CreateConfirmationToken(context.Background(), p)
}

func WelcomeEmailAction(p actions.Payload, container actions.Container) error {
	emailService := container.Email()
	signUp := p.(*SignUp)
	message := fmt.Sprintf(`
        <html>
        <body>
            <h1>Welcome, %s!</h1>
            <p>Thank you for signing up to our web application. We are excited to have you on board.</p>
            <p>If you have any questions, feel free to reply to this email.</p>
        </body>
        </html>
    `, "username")

	err := emailService.SendHtmlEmail(signUp.Email, "Welcome to Opus", message)
	if err != nil {
		return quality.ErrExecution.WithDetail("failed to send welcome email: " + err.Error())
	}
	fmt.Println("Welcome email sent successfully")
	return nil
}
