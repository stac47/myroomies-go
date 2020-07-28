package models

import (
	"time"
)

type Expense struct {
	Id          string    `bson:"_id,omitempty"`
	PayerLogin  string    `bson:"payer_login"`
	Date        time.Time `bson:"date"`
	Amount      float64   `bson:"amount"`
	Recipient   string    `bson:"recipient"`
	Description string    `bson:"description"`
}
