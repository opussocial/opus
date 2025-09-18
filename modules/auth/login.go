package auth

import (
    "fmt"
    // "time"
    "context"
    "net/mail"
    "net/http"
//  "unicode/utf8"
    
    "gitlab.com/pedrokoblitz/opus-go/internal/payloads"
    "gitlab.com/pedrokoblitz/opus-go/internal/adapters"
    "gitlab.com/pedrokoblitz/opus-go/internal/quality"
)

type Token struct {
    AuthBasePayload
}

func (p *Token) FromRequest(r *http.Request) error {
    return nil
}

func (p *Token) Validate() error {
    if p.Token == "" {
        return quality.ErrValidation.WithDetail("token is required")
    }   
    return nil
}

func (p *Token) Process() error {
    return nil
}

type Login struct {
    AuthBasePayload
}

func (p *Login) FromRequest(r *http.Request) error {
    return nil
}

func (p *Login) Validate() error {
    if p.Email == "" {
        return quality.ErrValidation.WithDetail("email is required")
    }
    if _, err := mail.ParseAddress(p.Email); err != nil {
        return quality.ErrValidation.WithDetail("invalid email format")
    }
    // if utf8.RuneCountInString(p.Password) < 8 {
    //     return quality.ErrValidation.WithDetail("password must be at least 8 characters")
    // }

    p.InputPassword = p.Password
    return nil
}

func (p *Login) Process() error {
    token, err := GenerateBearerToken(p.ID)
    if err != nil {
        return quality.ErrExecution.WithDetail("error generating token")
    }
    p.Token = token
    err = CheckPassword(p.InputPassword, p.Password)
    if err != nil {
        return quality.ErrUnauthorized
    }
    return nil
}

func LoginAction(p payloads.Payload, adapter interface{}) error {
    store := NewLoginStore(adapter.(adapters.DatabaseAdapter))
    err := store.GetByEmail(context.Background(), p)
    if err != nil {
        return err
    }
    err = p.Process()
    if err != nil {
        return err
    }
    return store.CreateAuthToken(context.Background(), p)
}

func BeforeShowHook(p payloads.Payload, adapter interface{}) error {
    fmt.Println("HELLO")
    fmt.Println("HELLO")
    fmt.Println("HELLO")
    fmt.Println("HELLO")
    return nil
}

func LogoutAction(p payloads.Payload, adapter interface{}) error {
    store := NewLoginStore(adapter.(adapters.DatabaseAdapter))
    return store.DeleteAuthToken(context.Background(), p)
}
