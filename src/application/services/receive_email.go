package services

import (
	"email-svc/src/application/events"
	"email-svc/src/infrastructure/configuration"
	"encoding/json"
)

type ReceiveEmail struct {
	logger         configuration.Logger
	emailPublisher events.Publisher
}

func NewReceiveEmail(emailPublisher events.Publisher, logger configuration.Logger) *ReceiveEmail {
	return &ReceiveEmail{emailPublisher: emailPublisher, logger: logger}
}

func (r ReceiveEmail) Execute(email *events.ReceivedEmail) error {
	r.logger.Info("Marshaling email to create queue message")
	messageString, err := json.Marshal(email)
	if err != nil {
		return err
	}

	r.logger.Info("Dispatching email to queue. ", string(messageString))
	return r.emailPublisher.Dispatch(string(messageString))
}
