package money

import (
	"fmt"
)

// Money represents an amount and currency in a MoneyWell document.
//
// It represents the amount as some number of cents, assuming 100 cents in a dollar, and is not
// useful for currencies that do not fit this mould. Naively assumes all currencies are prefixed
// with "$".
type Money struct {
	Currency string
	Amount   int64
}

// IsZero determines if the value of the money is 0.
func (c Money) IsZero() bool {
	return c.Amount == 0
}

// Add together two instances of Money, adopting the appropriate currency or panicking on a
// mismatch.
func (c Money) Add(other Money) Money {
	currency := c.Currency
	if currency == "" {
		currency = other.Currency
	} else if other.Currency != "" && currency != other.Currency {
		panic(fmt.Sprintf(
			"cannot add currencies of different types %s and %s",
			currency,
			other.Currency,
		))
	}

	return Money{
		Currency: currency,
		Amount:   c.Amount + other.Amount,
	}
}

// Multiply a Money by some constant scaling factor.
func (c Money) Multiply(factor int64) Money {
	return Money{
		Currency: c.Currency,
		Amount:   c.Amount * factor,
	}
}

// String describes the Money in textual form.
func (c Money) String() string {
	currency := ""
	if len(c.Currency) > 0 {
		currency = fmt.Sprintf(" %s", c.Currency)
	}

	if c.Amount < 0 {
		dollars := int64(-c.Amount / 100)
		cents := (-c.Amount) % 100
		return fmt.Sprintf("-$%d.%02d%s", dollars, cents, currency)
	} else {
		dollars := int64(c.Amount / 100)
		cents := c.Amount % 100
		return fmt.Sprintf("$%d.%02d%s", dollars, cents, currency)
	}
}
