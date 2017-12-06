package api

import (
	"database/sql"

	"github.com/pkg/errors"
)

// AccountGroup represents an account group in a MoneyWell document. An account group correlates
// 1:1 with a record in the ZACCOUNTGROUP table. Not all columns are exported.
//
// The MoneyWell SQLite schema for the ZACCOUNTGROUP table is as follows:
//  > .schema ZACCOUNTGROUP
//  CREATE TABLE ZACCOUNTGROUP (
//      Z_PK INTEGER PRIMARY KEY,
//      Z_ENT INTEGER,
//      Z_OPT INTEGER,
//      ZSEQUENCE INTEGER,
//      ZNAME VARCHAR,
//      ZTICDSSYNCID VARCHAR,
//      ZUNIQUEID VARCHAR
//  );
type AccountGroup struct {
	PrimaryKey int64
	Name       string
}

// GetAccountGroups fetches the set of accounts in a MoneyWell document, sorted by the display
// order as MoneyWell itself would render.
func GetAccountGroups(database *sql.DB) ([]AccountGroup, error) {
	rows, err := database.Query(`
            SELECT 
                zag.Z_PK, 
                zag.ZNAME
            FROM 
                ZACCOUNTGROUP zag
            ORDER BY
                -- Sort by the account group sequence
                zag.ZSEQUENCE ASC
        `)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query account groups")
	}
	defer rows.Close()

	accountGroups := []AccountGroup{}

	var primaryKey int64
	var name string
	for rows.Next() {
		err := rows.Scan(&primaryKey, &name)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan account group")
		}

		accountGroups = append(accountGroups, AccountGroup{
			PrimaryKey: primaryKey,
			Name:       name,
		})
	}

	return accountGroups, nil
}

func GetAccountGroupsMap(database *sql.DB) (map[int64]AccountGroup, error) {
	accountGroups, err := GetAccountGroups(database)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	accountGroupsMap := make(map[int64]AccountGroup, len(accountGroups))
	for _, accountGroup := range accountGroups {
		accountGroupsMap[accountGroup.PrimaryKey] = accountGroup
	}

	return accountGroupsMap, nil
}
