package adapters

import (
    "gopkg.in/mail.v2"
)

// EmailConfig holds SMTP server configuration
type EmailConfig struct {
    SMTPHost     string
    SMTPPort     int
    SenderEmail  string
    SenderPass   string
}

// EmailAdapter manages email sending
type EmailAdapter struct {
    config EmailConfig
    dialer *mail.Dialer
}

// NewEmailAdapter creates a new EmailAdapter with the given config
func NewEmailAdapter(cfg EmailConfig) *EmailAdapter {
    d := mail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SenderEmail, cfg.SenderPass)
    d.StartTLSPolicy = mail.MandatoryStartTLS
    return &EmailAdapter{
        config: cfg,
        dialer: d,
    }
}

// SendWelcomeEmail sends a welcome email with HTML content
func (es *EmailAdapter) SendHtmlEmail(to, subject, message string) error {
    m := mail.NewMessage()
    m.SetHeader("From", es.config.SenderEmail)
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", message)

    return es.dialer.DialAndSend(m)
}

// SendAlertEmail sends a plain text alert email
func (es *EmailAdapter) SendTextEmail(to, subject, message string) error {
    m := mail.NewMessage()
    m.SetHeader("From", es.config.SenderEmail)
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/plain", message)

    return es.dialer.DialAndSend(m)
}
