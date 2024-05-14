package mail

import (
	"encoding/json"
	"io/ioutil"
	"net/smtp"
	"strings"

	"github.com/nontechno/link"
	log "github.com/sirupsen/logrus"
)

const (
	Separator = "\u0000"
)

type EmailConfig struct {
	Username string `json:"user"`
	Password string `json:"pass"`
	SMTP     string `json:"smtp"`
	Port     string `json:"port"`
	Message  string `json:"message"`
	From     string `json:"from"`
	Subject  string `json:"subject"`
}

func (ec EmailConfig) Address() string {
	return ec.SMTP + ":" + ec.Port
}

var (
	emailConfig EmailConfig
	getLog      func() *log.Entry
)

func SendEmail(to []string, txt string) {
	if len(emailConfig.Username) == 0 {
		if data, err := ioutil.ReadFile("./email.config"); err == nil {
			if err := json.Unmarshal(data, &emailConfig); err != nil {
				// no point to go on, since we have no config
				getLog().Error("failed to parse config", err)
				return
			}
		} else {
			getLog().Error("failed to find email config", err)
			return
		}
	}

	message := emailConfig.Message
	subject := emailConfig.Subject
	if len(txt) > 0 {
		if parts := strings.Split(txt, Separator); len(parts) > 1 {
			subject = parts[0]
			message = parts[1]
		} else {
			message = txt
		}
	}

	msg := []byte("To:" + strings.Join(to, ";") +
		"\r\nFrom: " + emailConfig.From +
		"\r\nSubject: " + subject +
		"\r\nContent-Type: text/plain\r\n\r\n" +
		message)

	// Authentication.
	auth := smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.SMTP)

	// Sending email.
	err := smtp.SendMail(emailConfig.Address(), auth, emailConfig.Username, to, msg)
	if err != nil {
		getLog().Error("failed to send email", err)
		return
	}

	getLog().Info("email sent")
}

func init() {
	link.Register(SendEmail, "send.mail")
	link.Link(&getLog, "get.log", nil)
}
