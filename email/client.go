package email

import (
	"errors"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	log "github.com/sirupsen/logrus"

	"github.com/todanni/api/config"
)

const (
	ProjectInviteEmailSubject = "ToDanni Project Invitation"
)

var (
	Sender = mail.NewEmail("ToDanni Notification", "no-reply@todanni.com")
)

type SenderClient interface {
	SendProjectInvitationEmail(email ProjectInviteEmail) error
	SendDashboardInvitationEmail() error
}

type emailClient struct {
	client *sendgrid.Client
}

func NewEmailClient(config config.Config) SenderClient {
	client := sendgrid.NewSendClient(config.SendGridAPIKey)
	return &emailClient{
		client: client,
	}
}

func (e *emailClient) SendProjectInvitationEmail(email ProjectInviteEmail) error {
	to := mail.NewEmail(email.RecipientName, email.RecipientEmail)

	// TODO: Figure out how to use templates
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"

	message := mail.NewSingleEmail(Sender, ProjectInviteEmailSubject, to, plainTextContent, htmlContent)

	response, err := e.client.Send(message)
	if err != nil {
		log.Error(err)
		return errors.New("couldn't send email")
	}

	log.Info(response)
	return nil
}

func (e *emailClient) SendDashboardInvitationEmail() error {
	//TODO implement me
	panic("implement me")
}
