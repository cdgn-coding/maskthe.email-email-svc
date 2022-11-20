package services

import (
	"email-svc/src/application/events"
	"email-svc/src/business/repositories"
)

type SendEmail struct {
	outboundEmail repositories.OutboundEmail
}

func NewSendEmail(outboundEmail repositories.OutboundEmail) *SendEmail {
	return &SendEmail{outboundEmail: outboundEmail}
}

func (r SendEmail) Execute(email *events.ReceivedEmail) error {
	err := r.outboundEmail.Send(email)
	return err
}
