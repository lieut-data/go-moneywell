package api

import (
	"database/sql"

	"github.com/pkg/errors"
)

// TransactionTag represents a tag assigned to a given transaction in a MoneyWell document. A
// transaction tag correlates 1:1 with a record in the Z_3TAGS table. Not all columns are exported.
//
// The MoneyWell SQLite schema for the Z_3TAGS table is as follows:
//  > .schema Z_3ZTAGS
//  CREATE TABLE Z_3TAGS (
//      Z_3ACTIVITIES INTEGER,
//      Z_24TAGS INTEGER
//  );
type TransactionTag struct {
	Transaction int64
	Tag         int64
}

// GetTransactionTags fetches the set of transaction tags in a MoneyWell document.
func GetTransactionTags(database *sql.DB) ([]TransactionTag, error) {
	rows, err := database.Query(`
            SELECT 
                zt.Z_3ACTIVITIES,
                zt.Z_24TAGS
            FROM 
                Z_3TAGS zt
            ORDER BY
                zt.Z_3ACTIVITIES ASC,
                zt.Z_24TAGS ASC
        `)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query transaction tags")
	}
	defer rows.Close()

	transactionTags := []TransactionTag{}

	var transaction, tag int64
	for rows.Next() {
		err := rows.Scan(&transaction, &tag)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan transaction tag")
		}

		transactionTags = append(transactionTags, TransactionTag{
			Transaction: transaction,
			Tag:         tag,
		})
	}

	return transactionTags, nil
}

// GetTransactionTagMap fetches a map from transaction to a set of tags.
func GetTransactionTagMap(database *sql.DB) (map[int64][]int64, error) {
	transactionTags, err := GetTransactionTags(database)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	transactionTagMap := make(map[int64][]int64)
	for _, transactionTag := range transactionTags {
		transactionTagMap[transactionTag.Transaction] = append(
			transactionTagMap[transactionTag.Transaction],
			transactionTag.Tag,
		)
	}

	return transactionTagMap, nil
}

// GetTagTransactionMap fetches a map from tag to a set of transactions.
func GetTagTransactionMap(database *sql.DB) (map[int64][]int64, error) {
	transactionTags, err := GetTransactionTags(database)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tagTransactionMap := make(map[int64][]int64)
	for _, transactionTag := range transactionTags {
		tagTransactionMap[transactionTag.Tag] = append(
			tagTransactionMap[transactionTag.Tag],
			transactionTag.Transaction,
		)
	}

	return tagTransactionMap, nil
}
