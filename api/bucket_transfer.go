package api

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"github.com/lieut-data/go-moneywell/api/money"
)

// BucketTransfer represents a transfer between buckets in a MoneyWell document.
//
// The MoneyWell SQLite schema for the ZBUCKETRANSFER table is as follows:
//  > .schema ZBUCKETTRANSFER {
//  CREATE TABLE ZBUCKETTRANSFER (
//     Z_PK INTEGER PRIMARY KEY,
//     Z_ENT INTEGER,
//     Z_OPT INTEGER,
//     ZDATEYMD INTEGER,
//     ZISUSERCREATED INTEGER,
//     ZSEQUENCE INTEGER,
//     ZTYPE INTEGER,
//     ZACCOUNT INTEGER,
//     ZBUCKET INTEGER,
//     ZEVENT INTEGER,
//     ZTRANSFERSIBLING INTEGER,
//     ZAMOUNT DECIMAL,
//     ZMEMO VARCHAR,
//     ZTICDSSYNCID VARCHAR,
//     ZUNIQUEID VARCHAR
//  );
type BucketTransfer struct {
	PrimaryKey   int
	Date         time.Time
	TransferType int
	Amount       money.Money
	Bucket       int64
	TargetBucket int64
}

// GetDate implements the Event interface to return the bucket transfer date.
func (bt *BucketTransfer) GetDate() time.Time {
	return bt.Date
}

// GetAmount implements the Event interface to return the bucket transfer amount.
func (bt *BucketTransfer) GetAmount() money.Money {
	return bt.Amount
}

// GetAmount implements the Event interface to return the bucket in question.
func (bt *BucketTransfer) GetBucket() int64 {
	return bt.Bucket
}

// GetBucketTransfers fetches the set of bucket transfers in a MoneyWell document.
func GetBucketTransfers(database *sql.DB) ([]BucketTransfer, error) {
	rows, err := database.Query(`
            SELECT 
                zbt.Z_PK, 
                zbt.ZDATEYMD,
                zbt.ZTYPE,
                CAST(ROUND(zbt.ZAMOUNT * 100) AS INTEGER),
                zbt.ZBUCKET,
                zbt2.ZBUCKET,
                zb.ZCURRENCYCODE
            FROM 
                ZBUCKETTRANSFER zbt
            JOIN
                ZBUCKET zb ON ( zb.Z_PK = zbt.ZBUCKET )
            JOIN
                ZBUCKETTRANSFER zbt2 ON ( zbt2.Z_PK = zbt.ZTRANSFERSIBLING )
            ORDER BY
                zbt.ZDATEYMD ASC,
                zbt.ZSEQUENCE ASC
        `)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query bucket transfers")
	}
	defer rows.Close()

	bucketTransfers := []BucketTransfer{}

	var primaryKey, dateymd, transferType int
	var amountRaw, bucket, targetBucket int64
	var currencyCode string
	for rows.Next() {
		err := rows.Scan(
			&primaryKey,
			&dateymd,
			&transferType,
			&amountRaw,
			&bucket,
			&targetBucket,
			&currencyCode,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan account")
		}

		date, err := parseDateymd(dateymd)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse bucket transfer date")
		}

		bucketTransfers = append(bucketTransfers, BucketTransfer{
			PrimaryKey:   primaryKey,
			Date:         date,
			TransferType: transferType,
			Amount: money.Money{
				Currency: currencyCode,
				Amount:   amountRaw,
			},
			Bucket:       bucket,
			TargetBucket: targetBucket,
		})
	}

	return bucketTransfers, nil
}
