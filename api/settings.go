package api

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// Settings represents the single row from the ZETTINGS table in a MoneyWell document.
//
// The MoneyWell SQLite schema for the ZSETTINGS table is as follows:
//  > .schema ZSETTINGS
//  CREATE TABLE ZSETTINGS (
//      Z_PK INTEGER PRIMARY KEY,
//      Z_ENT INTEGER,
//      Z_OPT INTEGER,
//      ZCASHFLOWSTARTDATEYMD INTEGER,
//      ZLASTFILLBUCKETSDATEYMD INTEGER,
//      ZATTACHMENTPATH VARCHAR,
//      ZTICDSSYNCID VARCHAR
//  );
type Settings struct {
	PrimaryKey          int64
	CashFlowStartDate   time.Time
	LastFillBucketsDate time.Time
	AttachmentPath      string
}

// GetSettings fetches the settings in a MoneyWell document.
func GetSettings(database *sql.DB) (Settings, error) {
	row := database.QueryRow(`
            SELECT 
                zs.Z_PK, 
                zs.ZCASHFLOWSTARTDATEYMD,
                zs.ZLASTFILLBUCKETSDATEYMD,
                zs.ZATTACHMENTPATH
            FROM 
                ZSETTINGS zs
        `)

	var primaryKey int64
	var cashFlowStartDateYMD, lastFillBucketsDateYMD int
	var attachmentPath sql.NullString

	err := row.Scan(
		&primaryKey,
		&cashFlowStartDateYMD,
		&lastFillBucketsDateYMD,
		&attachmentPath,
	)
	if err != nil {
		return Settings{}, errors.Wrap(err, "failed to scan settings")
	}

	cashFlowStartDate, err := parseDateymd(cashFlowStartDateYMD)
	if err != nil {
		return Settings{}, errors.Wrapf(err, "failed to parse cash flow start date")
	}

	lastFillBucketsDate, err := parseDateymd(lastFillBucketsDateYMD)
	if err != nil {
		return Settings{}, errors.Wrapf(err, "failed to parse last fill buckets date")
	}

	return Settings{
		PrimaryKey:          primaryKey,
		CashFlowStartDate:   cashFlowStartDate,
		LastFillBucketsDate: lastFillBucketsDate,
		AttachmentPath:      attachmentPath.String,
	}, nil
}
