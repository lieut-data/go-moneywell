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
	defer database.Close()

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

func TestGetBucketGroupsMap(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	bucketGroups, err := api.GetBucketGroupsMap(database)
	assert.NoError(t, err)

	expectedBucketGroups := map[int64]api.BucketGroup{
		2: {
			PrimaryKey: 2,
			Type:       api.BucketGroupTypeIncome,
			Name:       "Salary",
		},
		4: {
			PrimaryKey: 4,
			Type:       api.BucketGroupTypeExpense,
			Name:       "Discretionary",
		},
		3: {
			PrimaryKey: 3,
			Type:       api.BucketGroupTypeExpense,
			Name:       "Bills",
		},
	}

	assert.Equal(t, expectedBucketGroups, bucketGroups)
}
