package models

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/pkg/domain/models"
)

type emailTest struct {
	title, body string
}

var emailTests = []emailTest{
	{"title", "body"},
	{"hey", "there"},
	{"Attention!", "Thanks for your attention."},
	{"!!!BIGGEST SALE OF THE YEAR!!!", "gotcha"},
	{"BTC Quote UAH rate", "The current rate is 831311.2341 UAH for 1 BTC"},
}

func TestNewEmail_WithCustomData_ConstructsCorrectEmails(t *testing.T) {
	for _, test := range emailTests {
		email := models.NewEmail(test.title, test.body)

		assert.Equal(t, test.body, email.Body)
		assert.Equal(t, test.title, email.Title)
	}
}
