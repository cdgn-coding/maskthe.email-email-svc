package sendgrid

import (
	"email-svc/src/business/entities"
	"fmt"
	"net/http"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendgridOutboundEmail struct {
	client *sendgrid.Client
}

func NewSendgridOutboundEmail(client *sendgrid.Client) *SendgridOutboundEmail {
	return &SendgridOutboundEmail{
		client: client,
	}
}

func (s *SendgridOutboundEmail) Send(email *entities.Email) error {
	from := mail.NewEmail("", email.From)
	to := mail.NewEmail("", email.To)
	message := mail.NewSingleEmail(from, email.Subject, to, email.PlainText, email.HTML)
	resp, err := s.client.Send(message)

	if err != nil {
		return err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("error sending email. Status Code: %d. Body: %s", resp.StatusCode, resp.Body)
	}

	return nil
}
