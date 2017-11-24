package api

import (
	"database/sql"

	"github.com/pkg/errors"
)

// Tag represents a tag annotating a transaction in a MoneyWell document. A tag correlates 1:1
// with a record in the ZTAG table. Not all columns are exported.
//
// The MoneyWell SQLite schema for the ZTAG table is as follows:
//  > .schema ZTAG
//  CREATE TABLE ZTAG (
//      Z_PK INTEGER PRIMARY KEY,
//      Z_ENT INTEGER,
//      Z_OPT INTEGER,
//      ZNAME VARCHAR,
//      ZTICDSSYNCID VARCHAR
//  );
type Tag struct {
	PrimaryKey int64
	Name       string
}

// GetTags fetches the set of tags in a MoneyWell document.
func GetTags(database *sql.DB) ([]Tag, error) {
	rows, err := database.Query(`
            SELECT 
                zt.Z_PK, 
                zt.ZNAME
            FROM 
                ZTAG zt
            ORDER BY
                zt.ZNAME ASC
        `)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query tags")
	}
	defer rows.Close()

	tags := []Tag{}

	var primaryKey int64
	var name string
	for rows.Next() {
		err := rows.Scan(&primaryKey, &name)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan tag")
		}

		tags = append(tags, Tag{
			PrimaryKey: primaryKey,
			Name:       name,
		})
	}

	return tags, nil
}

// GetTagsMap gets a map from the tag primary key to the tag.
func GetTagsMap(database *sql.DB) (map[int64]Tag, error) {
	tags, err := GetTags(database)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tagsMap := make(map[int64]Tag, len(tags))
	for _, tag := range tags {
		tagsMap[tag.PrimaryKey] = tag
	}

	return tagsMap, nil
}
