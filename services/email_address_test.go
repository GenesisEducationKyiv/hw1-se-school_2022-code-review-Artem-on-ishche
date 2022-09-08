package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestThatValidEmailsDontProduceErrors(t *testing.T) {
	for _, address := range validEmailAddresses {
		_, err := NewEmailAddress(address)
		assert.Nil(t, err, fmt.Sprintf("Returned an error on a valid address: %s", address))
	}
}

func TestThatInvalidEmailsProduceErrors(t *testing.T) {
	for _, address := range invalidEmailAddresses {
		_, err := NewEmailAddress(address)
		assert.NotNil(t, err, fmt.Sprintf("Didn't return an error on an invalid address: %s", address))
	}
}
