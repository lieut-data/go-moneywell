package api

import (
	"database/sql"
	"sort"

	"github.com/pkg/errors"

	"github.com/lieut-data/go-moneywell/api/money"
)

// Bucket represents an income or expense bucket in a MoneyWell document. A bucket correlates 1:1
// with a record in the ZBUCKET table. Not all columns are exported.
//
// The MoneyWell SQLite schema for the ZBUCKET table is as follows:
//  > .schema ZBUCKET
//  CREATE TABLE ZBUCKET (
//      Z_PK INTEGER PRIMARY KEY,
//      Z_ENT INTEGER,
//      Z_OPT INTEGER,
//      ZINCLUDEINAUTOFILL INTEGER,
//      ZISHIDDEN INTEGER,
//      ZISSELECTED INTEGER,
//      ZISTAXRELATED INTEGER,
//      ZSEQUENCE INTEGER,
//      ZTYPE INTEGER,
//      ZBUCKETGROUP INTEGER,
//      ZOVERFLOWBUCKET INTEGER,
//      ZSOURCEBUCKET INTEGER,
//      ZCURRENCYCODE VARCHAR,
//      ZMEMO VARCHAR,
//      ZNAME VARCHAR,
//      ZTICDSSYNCID VARCHAR,
//      ZUNIQUEID VARCHAR
//  );
type Bucket struct {
	PrimaryKey      int64
	Type            int64
	BucketGroup     int64
	Name            string
	StartingBalance money.Money
	CurrencyCode    string
}

// GetBuckets fetches the set of buckets in a MoneyWell document, sorted by the display order
// as MoneyWell itself would render.
func GetBuckets(database *sql.DB) ([]Bucket, error) {
	rows, err := database.Query(`
            SELECT 
                zb.Z_PK, 
                zb.ZTYPE,
                zb.ZBUCKETGROUP, 
                zb.ZNAME,
                CAST(ROUND(zbsb.ZAMOUNT * 100) AS INTEGER),
                zb.ZCURRENCYCODE
            FROM 
                ZBUCKET zb
            LEFT JOIN
                ZBUCKETGROUP zbg ON ( zbg.Z_PK = zb.ZBUCKETGROUP )
            LEFT JOIN
                ZBUCKETSTARTINGBALANCE zbsb ON ( zbsb.ZBUCKET = zb.Z_PK )
            ORDER BY
                -- Sort income buckets before expense buckets
                zb.ZTYPE ASC,
                -- Sort by the bucket group, if any
                zbg.ZSEQUENCE ASC,
                -- Sort within the bucket group, if any
                zb.ZSEQUENCE ASC
        `)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query buckets")
	}
	defer rows.Close()

	buckets := []Bucket{}

	var primaryKey, bucketType int64
	var name, currencyCode string
	var bucketGroup, startingBalance sql.NullInt64
	for rows.Next() {
		err := rows.Scan(
			&primaryKey,
			&bucketType,
			&bucketGroup,
			&name,
			&startingBalance,
			&currencyCode,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan bucket")
		}

		buckets = append(buckets, Bucket{
			PrimaryKey:  primaryKey,
			Type:        bucketType,
			BucketGroup: bucketGroup.Int64,
			Name:        name,
			StartingBalance: money.Money{
				Currency: currencyCode,
				Amount:   startingBalance.Int64,
			},
		})
	}

	return buckets, nil
}

// GetBucketsMap gets a map from the bucket primary key to the bucket.
func GetBucketsMap(database *sql.DB) (map[int64]Bucket, error) {
	buckets, err := GetBuckets(database)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	bucketsMap := make(map[int64]Bucket, len(buckets))
	for _, bucket := range buckets {
		bucketsMap[bucket.PrimaryKey] = bucket
	}

	return bucketsMap, nil
}

// GetBucketEvents uses the given transactions and bucket transfers to generate a list of events
// against the bucket sorted by date, generally matching the display of a bucket in MoneyWell.
func GetBucketEvents(
	bucket Bucket,
	transactions []Transaction,
	bucketTransfers []BucketTransfer,
) ([]Event, error) {
	events := []Event{}

	for _, transaction := range transactions {
		transaction := transaction

		// Ignore voided and pending transactions.
		switch transaction.Status {
		case TransactionStatusVoided:
			fallthrough
		case TransactionStatusPending:
			continue
		}

		if transaction.Bucket == bucket.PrimaryKey {
			events = append(events, &transaction)
		}
	}

	for _, bucketTransfer := range bucketTransfers {
		bucketTransfer := bucketTransfer
		if bucketTransfer.Bucket == bucket.PrimaryKey {
			events = append(events, &bucketTransfer)
		}
	}

	sort.Sort(ByEventDate(events))

	return events, nil
}

// GetBucketBalance uses the given events to compute the current balance of a bucket.
func GetBucketBalance(bucket Bucket, events []Event, settings Settings) (money.Money, error) {
	balance := bucket.StartingBalance

	for _, event := range events {
		if event.GetDate().Before(settings.CashFlowStartDate) {
			continue
		}
		if event.GetBucket() != bucket.PrimaryKey {
			continue
		}

		balance = balance.Add(event.GetAmount())
	}

	return balance, nil
}
