package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lieut-data/go-moneywell/api"
	"github.com/lieut-data/go-moneywell/api/money"
)

func TestGetAccounts(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	accounts, err := api.GetAccounts(database)
	assert.NoError(t, err)

	expectedAccounts := []api.Account{
		{
			PrimaryKey:        3,
			Name:              "Cash",
			Balance:           money.Money{Currency: "CAD", Amount: 0},
			IsBucketOptional:  false,
			IncludeInCashFlow: true,
			CurrencyCode:      "CAD",
			AccountGroup:      0,
		},
		{
			PrimaryKey:        4,
			Name:              "Cash (USD)",
			Balance:           money.Money{Currency: "USD", Amount: 0},
			IsBucketOptional:  false,
			IncludeInCashFlow: true,
			CurrencyCode:      "USD",
			AccountGroup:      0,
		},
		{
			PrimaryKey:        1,
			Name:              "Chequing Account",
			Balance:           money.Money{Currency: "CAD", Amount: 0},
			IsBucketOptional:  false,
			IncludeInCashFlow: true,
			CurrencyCode:      "CAD",
			AccountGroup:      2,
		},
		{
			PrimaryKey:        2,
			Name:              "Savings",
			Balance:           money.Money{Currency: "CAD", Amount: 100 * 100},
			IsBucketOptional:  false,
			IncludeInCashFlow: false,
			CurrencyCode:      "CAD",
			AccountGroup:      2,
		},
		{
			PrimaryKey:        5,
			Name:              "Line of Credit",
			Balance:           money.Money{Currency: "CAD", Amount: 0},
			IsBucketOptional:  false,
			IncludeInCashFlow: true,
			CurrencyCode:      "CAD",
			AccountGroup:      1,
		},
		{
			PrimaryKey:        6,
			Name:              "Visa",
			Balance:           money.Money{Currency: "CAD", Amount: 0},
			IsBucketOptional:  false,
			IncludeInCashFlow: true,
			CurrencyCode:      "CAD",
			AccountGroup:      1,
		},
	}

	assert.Equal(t, expectedAccounts, accounts)
}

func TestGetAccountsMap(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	accounts, err := api.GetAccountsMap(database)
	assert.NoError(t, err)

	expectedAccounts := map[int64]api.Account{
		3: {
			PrimaryKey:        3,
			Name:              "Cash",
			Balance:           money.Money{Currency: "CAD", Amount: 0},
			IsBucketOptional:  false,
			IncludeInCashFlow: true,
			CurrencyCode:      "CAD",
			AccountGroup:      0,
		},
		4: {
			PrimaryKey:        4,
			Name:              "Cash (USD)",
			Balance:           money.Money{Currency: "USD", Amount: 0},
			IsBucketOptional:  false,
			IncludeInCashFlow: true,
			CurrencyCode:      "USD",
			AccountGroup:      0,
		},
		1: {
			PrimaryKey:        1,
			Name:              "Chequing Account",
			Balance:           money.Money{Currency: "CAD", Amount: 0},
			IsBucketOptional:  false,
			IncludeInCashFlow: true,
			CurrencyCode:      "CAD",
			AccountGroup:      2,
		},
		2: {
			PrimaryKey:        2,
			Name:              "Savings",
			Balance:           money.Money{Currency: "CAD", Amount: 100 * 100},
			IsBucketOptional:  false,
			IncludeInCashFlow: false,
			CurrencyCode:      "CAD",
			AccountGroup:      2,
		},
		5: {
			PrimaryKey:        5,
			Name:              "Line of Credit",
			Balance:           money.Money{Currency: "CAD", Amount: 0},
			IsBucketOptional:  false,
			IncludeInCashFlow: true,
			CurrencyCode:      "CAD",
			AccountGroup:      1,
		},
		6: {
			PrimaryKey:        6,
			Name:              "Visa",
			Balance:           money.Money{Currency: "CAD", Amount: 0},
			IsBucketOptional:  false,
			IncludeInCashFlow: true,
			CurrencyCode:      "CAD",
			AccountGroup:      1,
		},
	}

	assert.Equal(t, expectedAccounts, accounts)
}

func TestGetAccountBalance(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	accounts, err := api.GetAccountsMap(database)
	assert.NoError(t, err)

	transactions, err := api.GetTransactions(database)
	assert.NoError(t, err)

	testCases := []struct {
		Description     string
		Account         int64
		ExpectedBalance money.Money
	}{
		{
			"Cash",
			3,
			money.Money{Currency: "CAD", Amount: 400 * 100},
		},
		{
			"Chequing Account",
			1,
			money.Money{Currency: "CAD", Amount: -350 * 100},
		},
		{
			"Savings",
			2,
			money.Money{Currency: "CAD", Amount: 100 * 100},
		},
		{
			"Line of Credit",
			5,
			money.Money{Currency: "CAD", Amount: 0},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Description, func(t *testing.T) {
			t.Parallel()

			balance := api.GetAccountBalance(
				accounts[testCase.Account],
				transactions,
			)
			assert.NoError(t, err)

			assert.Equal(t, testCase.ExpectedBalance, balance)
		})
	}
}
