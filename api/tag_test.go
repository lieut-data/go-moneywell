package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lieut-data/go-moneywell/api"
)

func TestGetTags(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	tags, err := api.GetTags(database)
	assert.NoError(t, err)

	expectedTags := []api.Tag{
		{1, "tag1"},
		{2, "tag2"},
		{3, "tag3"},
		{4, "tag4"},
	}

	assert.Equal(t, expectedTags, tags)
}

func TestGetTagsMap(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	tagsMap, err := api.GetTagsMap(database)
	assert.NoError(t, err)

	expectedTagsMap := map[int64]api.Tag{
		1: {1, "tag1"},
		2: {2, "tag2"},
		3: {3, "tag3"},
		4: {4, "tag4"},
	}

	assert.Equal(t, expectedTagsMap, tagsMap)
}
