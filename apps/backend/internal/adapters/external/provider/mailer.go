package provider

import (
	"errors"
	"fmt"
	"html"
	"net/smtp"
	"os"
	"strings"

	"go.uber.org/zap"
)

type Mailer interface {
	Send(to, subject, html string) error
}

var ErrMailerNotConfigured = errors.New("mailer is not configured")

type SMTPMailer struct{ logger *zap.SugaredLogger }

func NewSMTPMailer(logger *zap.SugaredLogger) *SMTPMailer {
	return &SMTPMailer{logger: logger.Named("Mailer")}
}

func (m *SMTPMailer) Send(to, subject, html string) error {
	host, port, from := os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT"), os.Getenv("MAIL_FROM")
	if host == "" || port == "" || from == "" {
		m.logger.Warnw("email delivery unavailable (SMTP is not configured)", "to", to)
		return ErrMailerNotConfigured
	}
	message := []byte(fmt.Sprintf("MIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\nFrom: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, html))
	var auth smtp.Auth
	if user := os.Getenv("SMTP_USER"); user != "" {
		auth = smtp.PlainAuth("", user, os.Getenv("SMTP_PASSWORD"), host)
	}
	return smtp.SendMail(host+":"+port, auth, from, []string{to}, message)
}

type OfficeInvitation struct {
	CouncillorName string
	Party          string
	City           string
	State          string
}

func OfficeInvitationHTML(url string, invitation OfficeInvitation) string {
	name := html.EscapeString(strings.TrimSpace(invitation.CouncillorName))
	party := html.EscapeString(strings.TrimSpace(invitation.Party))
	city := html.EscapeString(strings.TrimSpace(invitation.City))
	state := html.EscapeString(strings.ToUpper(strings.TrimSpace(invitation.State)))
	region := strings.Trim(strings.Join([]string{city, state}, "/"), "/")
	partyLine := ""
	if party != "" {
		partyLine = fmt.Sprintf(`<p style="margin:4px 0 0;color:#617166;font-size:14px">%s</p>`, party)
	}
	regionLine := ""
	if region != "" {
		regionLine = fmt.Sprintf(`<p style="margin:4px 0 0;color:#617166;font-size:14px">%s</p>`, region)
	}
	return fmt.Sprintf(`<!doctype html><html lang="pt-BR"><body style="margin:0;font-family:Arial,sans-serif;background:#f6f4ed;padding:32px;color:#1f3d2a"><main style="max-width:520px;margin:auto;background:#fff;padding:32px;border-radius:16px"><p style="margin:0 0 18px;color:#698d22;font-size:12px;font-weight:bold;letter-spacing:.08em;text-transform:uppercase">Convite para equipe</p><h1 style="margin:0;font-size:28px">Você foi convidado para o gabinete de %s</h1><div style="margin:22px 0;padding:16px;border-radius:12px;background:#f1f7dc"><strong>Gabinete de %s</strong>%s%s</div><p style="line-height:1.6">Finalize seu cadastro para atuar, junto ao gabinete, nas demandas da região.</p><a href="%s" style="display:inline-block;background:#b7d84b;color:#1f3d2a;padding:14px 20px;border-radius:8px;font-weight:bold;text-decoration:none">Finalizar cadastro</a><p style="margin:22px 0 0;color:#617166;font-size:13px;line-height:1.5">Este convite é pessoal e expira em 72 horas.</p></main></body></html>`, name, name, partyLine, regionLine, html.EscapeString(url))
}
