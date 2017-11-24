package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lieut-data/go-moneywell/api"
	"github.com/lieut-data/go-moneywell/api/money"
)

func TestGetBuckets(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)

	buckets, err := api.GetBuckets(database)
	assert.NoError(t, err)

	expectedBuckets := []api.Bucket{
		{
			PrimaryKey:      3,
			Type:            1,
			BucketGroup:     2,
			Name:            "Salary",
			StartingBalance: money.Money{Currency: "CAD", Amount: 0},
			CurrencyCode:    "",
		},
		{
			PrimaryKey:      13,
			Type:            2,
			BucketGroup:     0,
			Name:            "Groceries",
			StartingBalance: money.Money{Currency: "CAD", Amount: 0},
			CurrencyCode:    "",
		},
		{
			PrimaryKey:      27,
			Type:            2,
			BucketGroup:     4,
			Name:            "Hobbies",
			StartingBalance: money.Money{Currency: "CAD", Amount: 0},
			CurrencyCode:    "",
		},
		{
			PrimaryKey:      2,
			Type:            2,
			BucketGroup:     3,
			Name:            "Mortgage/Rent",
			StartingBalance: money.Money{Currency: "CAD", Amount: 0},
			CurrencyCode:    "",
		},
	}

	assert.Equal(t, expectedBuckets, buckets)
}

func TestGetBucketsMap(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)

	buckets, err := api.GetBucketsMap(database)
	assert.NoError(t, err)

	expectedBuckets := map[int64]api.Bucket{
		2: {
			PrimaryKey:      2,
			Type:            2,
			BucketGroup:     3,
			Name:            "Mortgage/Rent",
			StartingBalance: money.Money{Currency: "CAD", Amount: 0},
			CurrencyCode:    "",
		},
		3: {
			PrimaryKey:      3,
			Type:            1,
			BucketGroup:     2,
			Name:            "Salary",
			StartingBalance: money.Money{Currency: "CAD", Amount: 0},
			CurrencyCode:    "",
		},
		13: {
			PrimaryKey:      13,
			Type:            2,
			BucketGroup:     0,
			Name:            "Groceries",
			StartingBalance: money.Money{Currency: "CAD", Amount: 0},
			CurrencyCode:    "",
		},
		27: {
			PrimaryKey:      27,
			Type:            2,
			BucketGroup:     4,
			Name:            "Hobbies",
			StartingBalance: money.Money{Currency: "CAD", Amount: 0},
			CurrencyCode:    "",
		},
	}

	assert.Equal(t, expectedBuckets, buckets)
}

func TestGetBucketEvents(t *testing.T) {
	t.Skip("not yet implemented")
}

func TestGetBucketBalance(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)

	settings, err := api.GetSettings(database)
	assert.NoError(t, err)

	buckets, err := api.GetBucketsMap(database)
	assert.NoError(t, err)

	transactions, err := api.GetTransactions(database)
	assert.NoError(t, err)

	bucketTransfers, err := api.GetBucketTransfers(database)
	assert.NoError(t, err)

	testCases := []struct {
		Description     string
		Bucket          int64
		ExpectedBalance money.Money
	}{
		{
			"Salary",
			3,
			money.Money{Currency: "CAD", Amount: 100 * 100},
		},
		{
			"Groceries",
			13,
			money.Money{Currency: "CAD", Amount: -200 * 100},
		},
		{
			"Mortgage/Rent",
			2,
			money.Money{Currency: "CAD", Amount: 50 * 100},
		},
		{
			"Hobbies",
			27,
			money.Money{Currency: "CAD", Amount: 100 * 100},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Description, func(t *testing.T) {
			t.Parallel()

			events, err := api.GetBucketEvents(
				buckets[testCase.Bucket],
				transactions,
				bucketTransfers,
			)
			assert.NoError(t, err)

			balance, err := api.GetBucketBalance(
				buckets[testCase.Bucket],
				events,
				settings,
			)
			assert.NoError(t, err)

			assert.Equal(t, testCase.ExpectedBalance, balance)
		})
	}
}
