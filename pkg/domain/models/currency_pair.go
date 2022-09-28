package models

import "fmt"

type CurrencyPair struct {
	Base  Currency
	Quote Currency
}

func (pair *CurrencyPair) String() string {
	return fmt.Sprintf("%s-%s", pair.Base.Name, pair.Quote.Name)
}

func NewCurrencyPair(base, quote Currency) CurrencyPair {
	return CurrencyPair{
		Base:  base,
		Quote: quote,
	}
}
