package services

import "strings"

type Currency struct {
	Name string
}

func NewCurrency(name string) Currency {
	return Currency{strings.ToUpper(name)}
}

type CurrencyPair struct {
	From Currency
	To   Currency
}

func NewCurrencyPair(from, to Currency) CurrencyPair {
	return CurrencyPair{
		From: from,
		To:   to,
	}
}
