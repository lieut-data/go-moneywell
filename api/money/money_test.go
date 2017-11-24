package money_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lieut-data/go-moneywell/api/money"
)

func TestIsZero(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Description    string
		Input          money.Money
		ExpectedIsZero bool
	}{
		{
			"default money",
			money.Money{},
			true,
		},
		{
			"money with only a currency",
			money.Money{Currency: "CAD"},
			true,
		},
		{
			"money with explicitly 0 amount",
			money.Money{Amount: 0},
			true,
		},
		{
			"money with a currency and an explicitly 0 amount",
			money.Money{Currency: "CAD", Amount: 0},
			true,
		},
		{
			"money with a non-zero amount",
			money.Money{Amount: 1},
			false,
		},
		{
			"money with a currency and a non-zero amount",
			money.Money{Currency: "CAD", Amount: 1},
			false,
		},
		{
			"money with a currency and a large non-zero amount",
			money.Money{Currency: "CAD", Amount: 100},
			false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Description, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.ExpectedIsZero, testCase.Input.IsZero())
		})
	}
}

func TestAdd(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Description   string
		M1            money.Money
		M2            money.Money
		ExpectedMoney money.Money
	}{
		{
			"uninitialized monies",
			money.Money{},
			money.Money{},
			money.Money{},
		},
		{
			"amount-only monies",
			money.Money{Amount: 1},
			money.Money{Amount: 2},
			money.Money{Amount: 3},
		},
		{
			"one amount with currency",
			money.Money{Currency: "CAD", Amount: 1},
			money.Money{Amount: 2},
			money.Money{Currency: "CAD", Amount: 3},
		},
		{
			"one larger amount with currency",
			money.Money{Currency: "CAD", Amount: 100},
			money.Money{Amount: 250},
			money.Money{Currency: "CAD", Amount: 350},
		},
		{
			"two amounts with currency",
			money.Money{Currency: "CAD", Amount: 100},
			money.Money{Currency: "CAD", Amount: 250},
			money.Money{Currency: "CAD", Amount: 350},
		},
		{
			"two amounts with currency, one negative",
			money.Money{Currency: "CAD", Amount: 100},
			money.Money{Currency: "CAD", Amount: -250},
			money.Money{Currency: "CAD", Amount: -150},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Description, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.ExpectedMoney, testCase.M1.Add(testCase.M2))
			assert.Equal(t, testCase.ExpectedMoney, testCase.M2.Add(testCase.M1))
		})
	}
}

func TestAddInvalidCurrencies(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		m1 := money.Money{Currency: "CAD", Amount: 100}
		m2 := money.Money{Currency: "USD", Amount: 100}
		m1.Add(m2)
	})

	assert.Panics(t, func() {
		m1 := money.Money{Currency: "CAD", Amount: 100}
		m2 := money.Money{Currency: "USD", Amount: 100}
		m2.Add(m1)
	})
}

func TestMultiply(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Description   string
		Money         money.Money
		Multiplier    int64
		ExpectedMoney money.Money
	}{
		{
			"uninitialized money, zero",
			money.Money{},
			0,
			money.Money{},
		},
		{
			"uninitialized money, unity",
			money.Money{},
			1,
			money.Money{},
		},
		{
			"currency-only money, zero",
			money.Money{Currency: "CAD"},
			0,
			money.Money{Currency: "CAD"},
		},
		{
			"currency-only money, unity",
			money.Money{Currency: "CAD"},
			1,
			money.Money{Currency: "CAD"},
		},
		{
			"no currency, initialized amount, zero",
			money.Money{Amount: 100},
			0,
			money.Money{Amount: 0},
		},
		{
			"no currency, initialized amount, unity",
			money.Money{Amount: 100},
			1,
			money.Money{Amount: 100},
		},
		{
			"no currency, initialized amount, -1",
			money.Money{Amount: 100},
			-1,
			money.Money{Amount: -100},
		},
		{
			"no currency, initialized amount, 2",
			money.Money{Amount: 100},
			2,
			money.Money{Amount: 200},
		},
		{
			"currency, initialized amount, zero",
			money.Money{Currency: "CAD", Amount: 100},
			0,
			money.Money{Currency: "CAD", Amount: 0},
		},
		{
			"currency, initialized amount, unity",
			money.Money{Currency: "CAD", Amount: 100},
			1,
			money.Money{Currency: "CAD", Amount: 100},
		},
		{
			"currency, initialized amount, -1",
			money.Money{Currency: "CAD", Amount: 100},
			-1,
			money.Money{Currency: "CAD", Amount: -100},
		},
		{
			"currency, initialized amount, 2",
			money.Money{Currency: "CAD", Amount: 100},
			2,
			money.Money{Currency: "CAD", Amount: 200},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Description, func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				testCase.ExpectedMoney,
				testCase.Money.Multiply(testCase.Multiplier),
			)
		})
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Description    string
		Money          money.Money
		ExpectedString string
	}{
		{
			"uninitialized money",
			money.Money{},
			"$0.00",
		},
		{
			"currency-only money (CAD)",
			money.Money{Currency: "CAD"},
			"$0.00 CAD",
		},
		{
			"currency-only money (USD)",
			money.Money{Currency: "USD"},
			"$0.00 USD",
		},
		{
			"currency-only money, $0.01",
			money.Money{Amount: 1},
			"$0.01",
		},
		{
			"currency-only money, $0.99",
			money.Money{Amount: 99},
			"$0.99",
		},
		{
			"currency-only money, $1.00",
			money.Money{Amount: 1 * 100},
			"$1.00",
		},
		{
			"currency-only money, $1.99",
			money.Money{Amount: 1*100 + 99},
			"$1.99",
		},
		{
			"currency-only money, -$1.99",
			money.Money{Amount: -1 * (1*100 + 99)},
			"-$1.99",
		},
		{
			"currency-only money, $350.25",
			money.Money{Amount: 350*100 + 25},
			"$350.25",
		},
		{
			"currency-only money, -$350.25",
			money.Money{Amount: -1 * (350*100 + 25)},
			"-$350.25",
		},
		{
			"currency-only money, $0.01 CAD",
			money.Money{Currency: "CAD", Amount: 1},
			"$0.01 CAD",
		},
		{
			"currency-only money, $0.99 CAD",
			money.Money{Currency: "CAD", Amount: 99},
			"$0.99 CAD",
		},
		{
			"currency-only money, $1.00 CAD",
			money.Money{Currency: "CAD", Amount: 1 * 100},
			"$1.00 CAD",
		},
		{
			"currency-only money, $1.99 CAD",
			money.Money{Currency: "CAD", Amount: 1*100 + 99},
			"$1.99 CAD",
		},
		{
			"currency-only money, -$1.99 CAD",
			money.Money{Currency: "CAD", Amount: -1 * (1*100 + 99)},
			"-$1.99 CAD",
		},
		{
			"currency-only money, $350.25 CAD",
			money.Money{Currency: "CAD", Amount: 350*100 + 25},
			"$350.25 CAD",
		},
		{
			"currency-only money, -$350.25 CAD",
			money.Money{Currency: "CAD", Amount: -1 * (350*100 + 25)},
			"-$350.25 CAD",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Description, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.ExpectedString, testCase.Money.String())
		})
	}
}
