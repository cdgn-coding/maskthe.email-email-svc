package services

import (
	"email-svc/src/application/events"
	"encoding/json"
)

type ReceiveEmail struct {
	emailPublisher events.Publisher
}

func NewReceiveEmail(emailPublisher events.Publisher) *ReceiveEmail {
	return &ReceiveEmail{emailPublisher: emailPublisher}
}

func (r ReceiveEmail) Execute(email *events.ReceivedEmail) error {
	messageString, err := json.Marshal(email)
	if err != nil {
		return err
	}
	return r.emailPublisher.Dispatch(string(messageString))
}
