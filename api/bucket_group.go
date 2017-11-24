package api

import (
	"database/sql"

	"github.com/pkg/errors"
)

// BucketGroup represents an income or expense bucket group in a MoneyWell document. A bucket group
// correlates 1:1 with a record in the ZBUCKETGROUP table. Not all columns are exported.
//
// The MoneyWell SQLite schema for the ZBUCKETGROUP table is as follows:
//  > .schema ZBUCKETGROUP
//  CREATE TABLE ZBUCKETGROUP (
//      Z_PK INTEGER PRIMARY KEY,
//      Z_ENT INTEGER,
//      Z_OPT INTEGER,
//      ZSEQUENCE INTEGER,
//      ZTYPE INTEGER,
//      ZNAME VARCHAR,
//      ZTICDSSYNCID VARCHAR,
//      ZUNIQUEID VARCHAR
//  );
type BucketGroup struct {
	PrimaryKey int64
	Type       int64
	Name       string
}

const (
	BucketGroupTypeIncome  = 1
	BucketGroupTypeExpense = 2
)

// GetBucketGroups fetches the set of bucket groups in a MoneyWell document.
func GetBucketGroups(database *sql.DB) ([]BucketGroup, error) {
	rows, err := database.Query(`
            SELECT 
                zbg.Z_PK, 
                zbg.ZTYPE,
                zbg.ZNAME
            FROM 
                ZBUCKETGROUP zbg
            ORDER BY
                -- Sort income bucket groups before expense bucket groups
                zbg.ZTYPE ASC,
                -- Sort by the bucket group sequence
                zbg.ZSEQUENCE ASC
        `)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query bucket groups")
	}
	defer rows.Close()

	bucketGroups := []BucketGroup{}

	var primaryKey, bucketGroupType int64
	var name string
	for rows.Next() {
		err := rows.Scan(&primaryKey, &bucketGroupType, &name)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan bucket group")
		}

		bucketGroups = append(bucketGroups, BucketGroup{
			PrimaryKey: primaryKey,
			Type:       bucketGroupType,
			Name:       name,
		})
	}

	return bucketGroups, nil
}
