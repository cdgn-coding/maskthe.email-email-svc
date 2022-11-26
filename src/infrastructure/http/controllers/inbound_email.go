package controllers

import (
	"email-svc/src/application/services"
	"email-svc/src/business/entities"
	"email-svc/src/infrastructure/configuration"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/sendgrid/sendgrid-go/helpers/inbound"
	"net/http"
)

type InboundEmailController struct {
	receiveEmail *services.ReceiveEmail
	logger       configuration.Logger
}

func NewInboundEmailController(
	receiveEmail *services.ReceiveEmail, logger configuration.Logger) *InboundEmailController {
	return &InboundEmailController{receiveEmail: receiveEmail, logger: logger}
}

var invalidInboundEmailRequest = errors.New("http request is not a valid inbound email")

var inboundEmailRequestNotSecure = errors.New("inbound email not passed security assessments")

var tooMuchRecipients = errors.New("inbound email has too much recipients")

func (controller InboundEmailController) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	email, err := controller.parseInboundEmail(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(email)
	if err != nil {
		controller.logger.Error(request.Context(), "invalid email: %v", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.receiveEmail.Execute(email)
	if err != nil {
		controller.logger.Error(request.Context(), "error receiving email: %v", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (controller InboundEmailController) parseInboundEmail(request *http.Request) (*entities.Email, error) {
	inboundEmail, err := inbound.Parse(request)

	if err != nil {
		return nil, fmt.Errorf("%w. %v", invalidInboundEmailRequest, err)
	}

	err = inboundEmail.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w. %v", inboundEmailRequestNotSecure, err)
	}

	if len(inboundEmail.Envelope.To) > 1 {
		return nil, fmt.Errorf("%w. %v", tooMuchRecipients, err)
	}

	if len(inboundEmail.Envelope.To) == 0 {
		return nil, fmt.Errorf("%w", invalidInboundEmailRequest)
	}

	email := &entities.Email{
		From:      inboundEmail.Envelope.From,
		To:        inboundEmail.Envelope.To[0],
		Subject:   inboundEmail.ParsedValues["subject"],
		PlainText: inboundEmail.TextBody,
		HTML:      inboundEmail.Body["text/html"],
	}

	return email, nil
}
