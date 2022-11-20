package consumers

import (
	"email-svc/src/application/events"
	"email-svc/src/application/services"
	"email-svc/src/infrastructure/configuration"
	"encoding/json"
	"errors"
	"fmt"
)

type EmailConsumer struct {
	logger    configuration.Logger
	sendEmail *services.SendEmail
}

func NewEmailConsumer(logger configuration.Logger, redirectEmail *services.SendEmail) *EmailConsumer {
	return &EmailConsumer{logger: logger, sendEmail: redirectEmail}
}

var PayloadNotValid = errors.New("payload not valid")

func (e *EmailConsumer) Invoke(payload string) error {
	command := &events.ReceivedEmail{}
	err := json.Unmarshal([]byte(payload), command)

	if err != nil {
		e.logger.Error(fmt.Sprintf("Unable to parse payload. %s", payload))
		return PayloadNotValid
	}

	err = e.sendEmail.Execute(command)

	if err != nil {
		e.logger.Error(fmt.Sprintf("Unable to send email. cause: %v", err))
		return err
	}

	return nil
}
