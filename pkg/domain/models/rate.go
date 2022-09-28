package models

import "fmt"

type ExchangeRate struct {
	CurrencyPair CurrencyPair
	Price        float64
}

func (rate *ExchangeRate) String() string {
	return fmt.Sprintf("%s rate is %v", rate.CurrencyPair.String(), rate.Price)
}

func NewExchangeRate(pair CurrencyPair, price float64) *ExchangeRate {
	return &ExchangeRate{CurrencyPair: pair, Price: price}
}
