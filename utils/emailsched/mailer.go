package emailsched

import (
	"fmt"
	"github-graph-drawer/config"
	"html/template"
	"net/smtp"
)

var (
	templates *template.Template
)

func init() {
	templates = template.Must(template.ParseGlob("./templates/email/*"))
}

func sendEmail(subject, content, to string) error {
	receiver := []string{to}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	_subject := "Subject: " + subject
	_to := "To: " + to
	_from := fmt.Sprintf("From: Baraa from GH Graph Drawer <%s>", config.Config().Mailer.User)
	body := []byte(fmt.Sprintf("%s\n%s\n%s\n%s\n%s", _from, _to, _subject, mime, content))

	addr := config.Config().Mailer.Host + ":" + config.Config().Mailer.Port
	auth := smtp.PlainAuth("", config.Config().Mailer.User, config.Config().Mailer.Password, config.Config().Mailer.Host)

	return smtp.SendMail(addr, auth, config.Config().Mailer.User, receiver, body)
}
