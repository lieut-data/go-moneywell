package api_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lieut-data/go-moneywell/api"
	"github.com/lieut-data/go-moneywell/api/money"
)

func TestGetBucketTransfers(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)

	bucketTransfers, err := api.GetBucketTransfers(database)
	assert.NoError(t, err)

	expectedBucketTransfers := []api.BucketTransfer{
		{
			PrimaryKey:   2,
			Date:         time.Date(2017, 11, 19, 0, 0, 0, 0, time.UTC),
			TransferType: 0,
			Amount:       money.Money{Currency: "CAD", Amount: 100 * 100},
			Bucket:       27,
			TargetBucket: 2,
		},
		{
			PrimaryKey:   5,
			Date:         time.Date(2017, 11, 19, 0, 0, 0, 0, time.UTC),
			TransferType: 0,
			Amount:       money.Money{Currency: "CAD", Amount: 250 * 100},
			Bucket:       13,
			TargetBucket: 3,
		},
		{
			PrimaryKey:   3,
			Date:         time.Date(2017, 11, 19, 0, 0, 0, 0, time.UTC),
			TransferType: 1,
			Amount:       money.Money{Currency: "CAD", Amount: -100 * 100},
			Bucket:       2,
			TargetBucket: 27,
		},
		{
			PrimaryKey:   4,
			Date:         time.Date(2017, 11, 19, 0, 0, 0, 0, time.UTC),
			TransferType: 1,
			Amount:       money.Money{Currency: "CAD", Amount: -250 * 100},
			Bucket:       3,
			TargetBucket: 13,
		},
		{
			PrimaryKey:   1,
			Date:         time.Date(2017, 11, 19, 0, 0, 0, 0, time.UTC),
			TransferType: 0,
			Amount:       money.Money{Currency: "CAD", Amount: 650 * 100},
			Bucket:       2,
			TargetBucket: 3,
		},
		{
			PrimaryKey:   6,
			Date:         time.Date(2017, 11, 19, 0, 0, 0, 0, time.UTC),
			TransferType: 1,
			Amount:       money.Money{Currency: "CAD", Amount: -650 * 100},
			Bucket:       3,
			TargetBucket: 2,
		},
	}

	assert.Equal(t, expectedBucketTransfers, bucketTransfers)
}
