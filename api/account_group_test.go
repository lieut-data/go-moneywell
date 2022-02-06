package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lieut-data/go-moneywell/api"
)

func TestGetAccountGroups(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	accountGroups, err := api.GetAccountGroups(database)
	assert.NoError(t, err)

	expectedAccountGroups := []api.AccountGroup{
		{
			PrimaryKey: 2,
			Name:       "Bank",
		},
		{
			PrimaryKey: 1,
			Name:       "Other Bank",
		},
	}

	assert.Equal(t, expectedAccountGroups, accountGroups)
}

func TestGetAccountGroupsMap(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	accountGroups, err := api.GetAccountGroupsMap(database)
	assert.NoError(t, err)

	expectedAccountGroups := map[int64]api.AccountGroup{
		2: {
			PrimaryKey: 2,
			Name:       "Bank",
		},
		1: {
			PrimaryKey: 1,
			Name:       "Other Bank",
		},
	}

	assert.Equal(t, expectedAccountGroups, accountGroups)
}
