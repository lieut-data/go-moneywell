package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lieut-data/go-moneywell/api"
)

func TestGetTransactionTags(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	transactionTags, err := api.GetTransactionTags(database)
	assert.NoError(t, err)

	expectedTransactionTags := []api.TransactionTag{
		{2, 1},
		{4, 2},
		{4, 4},
		{5, 3},
		{5, 4},
	}

	assert.Equal(t, expectedTransactionTags, transactionTags)
}

func TestGetTransactionTagMap(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	transactionTagMap, err := api.GetTransactionTagMap(database)
	assert.NoError(t, err)

	expectedTransactionTagMap := map[int64][]int64{
		2: []int64{1},
		5: []int64{3, 4},
		4: []int64{2, 4},
	}

	assert.Equal(t, expectedTransactionTagMap, transactionTagMap)
}

func TestGetTagTransactionMap(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	transactionTagMap, err := api.GetTagTransactionMap(database)
	assert.NoError(t, err)

	expectedTagTransactionMap := map[int64][]int64{
		1: []int64{2},
		3: []int64{5},
		2: []int64{4},
		4: []int64{4, 5},
	}

	assert.Equal(t, expectedTagTransactionMap, transactionTagMap)
}
