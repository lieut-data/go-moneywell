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
