package models

import "strings"

type Currency struct {
	Name string
}

func NewCurrency(name string) Currency {
	return Currency{strings.ToUpper(name)}
}

func (c Currency) IsEmpty() bool {
	return c == Currency{}
}
