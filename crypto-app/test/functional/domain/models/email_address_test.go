package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/pkg/domain/models"
)

var validEmailAddresses = []string{
	"artem.mykytyshyn@gmail.com",
	"meornotme@gmail.com",
	"some1l234aldf02l@gmail.com",
	"mikemichaels@domain.org",
	"text@text.text",
	"artem.mykytyshyn@ukma.edu.ua",
	"a.a.a.a.a.a.a.1.10.1.212.1.131,131@a.a.a.a",
	"a@a.a",
}

var invalidEmailAddresses = []string{
	"artem@mykytyshyn@gmail.com",
	"meornotme@gmailcom",
	"some.1l23.4aldf02l.gmail.com",
	".text@text.text",
	"text@text.text.",
	"artem.mykytyshyn@ukma..ua",
	"a@a",
	"who@.",
}

func TestNewEmailAddress_WithValidInputs_DoNotProduceErrors(t *testing.T) {
	for _, address := range validEmailAddresses {
		_, err := models.NewEmailAddress(address)
		assert.Nil(t, err, fmt.Sprintf("Returned an error on a valid address: %s", address))
	}
}

func TestNewEmailAddress_WithInvalidInputs_ProduceErrors(t *testing.T) {
	for _, address := range invalidEmailAddresses {
		_, err := models.NewEmailAddress(address)
		assert.NotNil(t, err, fmt.Sprintf("Didn't return an error on an invalid address: %s", address))
	}
}
