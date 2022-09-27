package models

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/pkg/domain/models"
)

type currencyTest struct {
	provided, expected string
}

var currencyTests = []currencyTest{
	{"uah", "UAH"},
	{"UAH", "UAH"},
	{"uAh", "UAH"},
	{"Uah", "UAH"},
	{"UaH", "UAH"},
	{"cad", "CAD"},
	{"some-random-name", "SOME-RANDOM-NAME"},
	{"cur", "CUR"},
}

func TestThatCurrencyNameIsUppercase(t *testing.T) {
	for _, data := range currencyTests {
		currency := models.NewCurrency(data.provided)

		assert.Equal(t, data.expected, currency.Name)
	}
}
