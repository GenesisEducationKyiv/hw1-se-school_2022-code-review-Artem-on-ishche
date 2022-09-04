package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type emailTest struct {
	title, body string
}

var emailTests = []emailTest{
	{"title", "body"},
	{"hey", "there"},
	{"Attention!", "Thanks for your attention."},
	{"!!!BIGGEST SALE OF THE YEAR!!!", "gotcha"},
	{"BTC to UAH rate", "The current rate is 831311.2341 UAH for 1 BTC"},
}

func TestEmailConstructor(t *testing.T) {
	for _, test := range emailTests {
		email := NewEmail(test.title, test.body)

		assert.Equal(t, test.body, email.Body)
		assert.Equal(t, test.title, email.Title)
	}
}
