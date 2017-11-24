package api

import (
	"time"

	"github.com/lieut-data/go-moneywell/api/money"
)

// Event abstracts a transaction or bucket transfer in a MoneyWell document.
type Event interface {
	GetDate() time.Time
	GetAmount() money.Money
	GetBucket() int64
}

type ByEventDate []Event

func (s ByEventDate) Len() int           { return len(s) }
func (s ByEventDate) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByEventDate) Less(i, j int) bool { return s[j].GetDate().After(s[i].GetDate()) }
