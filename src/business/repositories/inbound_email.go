package repositories

import (
	"email-svc/src/business/entities"
)

type OutboundEmail interface {
	Send(email *entities.Email) error
}
