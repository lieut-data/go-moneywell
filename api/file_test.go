package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lieut-data/go-moneywell/api"
)

func TestOpenDocument(t *testing.T) {
	t.Parallel()

	_, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)

	_, err = api.OpenDocument("Test.moneywell/StoreContent/persistentStore")
	assert.NoError(t, err)

	_, err = api.OpenDocument("NoSuchFile.moneywell")
	assert.Error(t, err)

	_, err = api.OpenDocument("NoSuchFile.moneywell/StoreContent/persistentStore")
	assert.Error(t, err)
}
