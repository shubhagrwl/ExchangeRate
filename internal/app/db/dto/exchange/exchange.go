package exchange

import (
	"time"
)

const (
	TABLE_NAME           = "exchange"
	COLUMN_EXCHANGE_DATE = "exchange_date"
	COLUMN_BASE          = "base"
	COLUMN_USD           = "usd"
	COLUMN_EUR           = "eur"
	COLUMN_GBP           = "gbp"
	COLUMN_CREATED_AT    = "created_at"
)

type Exchange struct {
	Id           int       `json:"id"`
	ExchangeDate string    `json:"exchange_date"`
	Base         string    `json:"base"`
	Usd          float64   `json:"usd"`
	Eur          float64   `json:"eur"`
	Gbp          float64   `json:"gbp"`
	CreatedAt    time.Time `json:"created_at"`
}
