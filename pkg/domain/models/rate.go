package models

import (
	"fmt"
	"time"
)

type ExchangeRate struct {
	CurrencyPair CurrencyPair
	Price        float64
	Timestamp    time.Time
}

func (rate *ExchangeRate) String() string {
	return fmt.Sprintf("%s rate is %v", rate.CurrencyPair.String(), rate.Price)
}

func NewExchangeRate(pair CurrencyPair, price float64, timestamp time.Time) *ExchangeRate {
	return &ExchangeRate{CurrencyPair: pair, Price: price, Timestamp: timestamp}
}
