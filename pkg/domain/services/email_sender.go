package services

import (
	"gses2.app/api/pkg/domain/models"
)

type EmailSender interface {
	SendEmails(email models.Email, receiverAddresses []string) error
}
