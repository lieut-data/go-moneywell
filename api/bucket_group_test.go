package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lieut-data/go-moneywell/api"
)

func TestGetBucketGroups(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)

	bucketGroups, err := api.GetBucketGroups(database)
	assert.NoError(t, err)

	expectedBucketGroups := []api.BucketGroup{
		{
			PrimaryKey: 2,
			Type:       api.BucketGroupTypeIncome,
			Name:       "Salary",
		},
		{
			PrimaryKey: 4,
			Type:       api.BucketGroupTypeExpense,
			Name:       "Discretionary",
		},
		{
			PrimaryKey: 3,
			Type:       api.BucketGroupTypeExpense,
			Name:       "Bills",
		},
	}

	assert.Equal(t, expectedBucketGroups, bucketGroups)
}
