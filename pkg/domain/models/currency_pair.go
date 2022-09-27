package models

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
