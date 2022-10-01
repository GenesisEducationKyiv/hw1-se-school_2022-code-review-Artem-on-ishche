package models

const CurrencySeparator = "-"

type CurrencyPair struct {
	Base  Currency
	Quote Currency
}

func (pair *CurrencyPair) String() string {
	return pair.Base.Name + CurrencySeparator + pair.Quote.Name
}

func NewCurrencyPair(base, quote Currency) CurrencyPair {
	return CurrencyPair{
		Base:  base,
		Quote: quote,
	}
}
