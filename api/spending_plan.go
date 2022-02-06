package api

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"github.com/lieut-data/go-moneywell/api/money"
)

// SpendingPlan represents spending plan events associated with a bucket in a MoneyWell document.
// This is a subset of the ZACTIVITY table, where both transactions, favourite transactions and
// spending plan amounts are recorded. Not all spending plan columns are exported.
//
// The MoneyWell SQLite schema for the ZACTIVITY table is as follows:
//  > .schema ZACTIVITY
//  CREATE TABLE ZACTIVITY (
//      Z_PK INTEGER PRIMARY KEY,
//      Z_ENT INTEGER,
//      Z_OPT INTEGER,
//      ZDATEYMD INTEGER,
//      ZISBUCKETOPTIONAL INTEGER,
//      ZISFLAGGED INTEGER,
//      ZISLOCKED INTEGER,
//      ZSEQUENCE INTEGER,
//      ZSOURCE INTEGER,
//      ZTYPE INTEGER,
//      ZSPLITPARENT INTEGER,
//      Z3_SPLITPARENT INTEGER,
//      ZTRANSFERACCOUNT INTEGER,
//      ZTRANSFERSIBLING INTEGER,
//      Z3_TRANSFERSIBLING INTEGER,
//      ZARETAGSSAVED INTEGER,
//      ZISAMOUNTSAVED INTEGER,
//      ZISBUCKETSAVEDINTEGER,
//      ZISMEMOSAVED INTEGER,
//      ZISTYPESAVED INTEGER,
//      ZACCOUNT INTEGER,
//      ZBUCKET INTEGER,
//      ZGENERATECOUNT INTEGER,
//      ZHASVARIABLEAMOUNT INTEGER,
//      ZACCOUNT1 INTEGER,
//      ZBUCKET1 INTEGER,
//      ZRECURRENCERULE INTEGER,
//      ZISPERCENTAGE INTEGER,
//      ZFILLRECURRENCERULE INTEGER,
//      ZSPENDINGPLAN INTEGER,
//      ZDATERECONCILEDYMD INTEGER,
//      ZISLASTIMPORT INTEGER,
//      ZISQUARANTINED INTEGER,
//      ZISREPEATING INTEGER,
//      ZSPLITINDEX INTEGER,
//      ZSTATUS INTEGER,
//      ZACCOUNT2 INTEGER,
//      ZBUCKET2 INTEGER,
//      ZINVESTMENTSECURITYID INTEGER,
//      ZDATELASTUSED TIMESTAMP,
//      ZLASTGENERATEDDATE TIMESTAMP,
//      ZAMOUNT DECIMAL,
//      ZCHECKREF DECIMAL,
//      ZLOCALIZEDAMOUNT DECIMAL,
//      ZSALEAMOUNT DECIMAL,
//      ZCOMMISSION DECIMAL,
//      ZFEES DECIMAL,
//      ZLOCALIZEDAMOUNT1 DECIMAL,
//      ZSALEAMOUNT1 DECIMAL,
//      ZSHAREPRICE DECIMAL,
//      ZSHARES DECIMAL,
//      ZSHARESSPLITDENOMINATOR DECIMAL,
//      ZSHARESSPLITNUMERATOR DECIMAL,
//      ZSPLITTOTAL DECIMAL,
//      ZTAXES DECIMAL,
//      ZAMOUNTSTRING VARCHAR,
//      ZCHECKREFSTRING VARCHAR,
//      ZMEMO VARCHAR,
//      ZPAYEE VARCHAR,
//      ZTICDSSYNCID VARCHAR,
//      ZUNIQUEID VARCHAR,
//      ZSALECURRENCYCODE VARCHAR,
//      ZEXTERNALID VARCHAR,
//      ZLOCALIZEDAMOUNTSTRING VARCHAR,
//      ZORIGINALMEMO VARCHAR,
//      ZORIGINALPAYEE VARCHAR,
//      ZRECEIPTFILENAME VARCHAR,
//      ZSALEAMOUNTSTRING VARCHAR,
//      ZSALECURRENCYCODE1 VARCHAR
//  );
type SpendingPlan struct {
	PrimaryKey         int64
	Date               time.Time
	Name               string
	Amount             money.Money
	Bucket             int64
	RecurrenceRule     int64
	FillRecurrenceRule int64
}

// GetSpendingPlan fetches the set of spending plan events in a MoneyWell document.
func GetSpendingPlan(database *sql.DB) ([]SpendingPlan, error) {
	rows, err := database.Query(`
            SELECT
		za.Z_PK,
		za.ZDATEYMD,
		za.ZPAYEE,
		CAST(ROUND(za.ZAMOUNT * 100) AS INTEGER),
		COALESCE(za.ZBUCKET, za.ZBUCKET1, za.ZBUCKET2),
		zb.ZCURRENCYCODE,
		za.ZRECURRENCERULE,
		za.ZFILLRECURRENCERULE
            FROM
                ZACTIVITY za
	    JOIN
		ZBUCKET zb ON (zb.Z_PK = COALESCE(za.ZBUCKET, za.ZBUCKET1, za.ZBUCKET2))
	    WHERE
		za.Z_ENT = ?
	    ORDER BY
		za.Z_PK ASC
        `, ActivityTypeSpendingPlan)

	if err != nil {
		return nil, errors.Wrap(err, "failed to query spending plan")
	}
	defer rows.Close()

	spendingPlan := []SpendingPlan{}

	var primaryKey, amountRaw int64
	var dateymd int
	var name string
	var bucket sql.NullInt64
	var currencyCode sql.NullString
	var recurrenceRule, fillRecurrenceRule sql.NullInt64
	for rows.Next() {
		err := rows.Scan(
			&primaryKey,
			&dateymd,
			&name,
			&amountRaw,
			&bucket,
			&currencyCode,
			&recurrenceRule,
			&fillRecurrenceRule,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan spending plan")
		}

		date, err := parseDateymd(dateymd)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse spending plan date")
		}

		spendingPlan = append(spendingPlan, SpendingPlan{
			PrimaryKey: primaryKey,
			Date:       date,
			Name:       name,
			Amount: money.Money{
				Currency: currencyCode.String,
				Amount:   amountRaw,
			},
			Bucket:             bucket.Int64,
			RecurrenceRule:     recurrenceRule.Int64,
			FillRecurrenceRule: fillRecurrenceRule.Int64,
		})
	}

	return spendingPlan, nil
}
