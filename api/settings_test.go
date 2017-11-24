package api_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lieut-data/go-moneywell/api"
)

func TestGetSettings(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)

	settings, err := api.GetSettings(database)
	assert.NoError(t, err)

	assert.Equal(t, time.Date(2017, 11, 01, 0, 0, 0, 0, time.UTC), settings.CashFlowStartDate)
	assert.Equal(t, time.Time{}, settings.LastFillBucketsDate)
	assert.Equal(t, "", settings.AttachmentPath)
}
