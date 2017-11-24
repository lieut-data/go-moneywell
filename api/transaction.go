package api

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"github.com/lieut-data/go-moneywell/api/money"
)

// Transaction represents a transaction in a MoneyWell document. This is a subset of the
// ZACTIVITY table, where both transactions, favourite transactions and spending plan amounts
// are recorded. Not all transaction columns are exported.
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
//
// By analysis, Z_ENT appears to represent the type of activity in question:
//  4: Placeholders for favourite transactions
//  6: Spending plan amounts
//  7: Transactions
type Transaction struct {
	PrimaryKey       int64
	Date             time.Time
	TransactionType  int
	Amount           money.Money
	Bucket           int64
	Account          int64
	TransferAccount  int64
	TransferSibling  int64
	SplitParent      int64
	IsSplit          bool
	IsBucketOptional bool
	Payee            string
	Memo             string
}

const (
	TransactionTypeDeposit    = 0
	TransactionTypeWithdrawal = 1
	TransactionTypeCheck      = 2
)

func (t *Transaction) IsTransfer() bool {
	return t.TransferAccount > 0 || t.TransferSibling > 0
}

func (t *Transaction) GetDate() time.Time {
	return t.Date
}

func (t *Transaction) GetAmount() money.Money {
	return t.Amount
}

func (t *Transaction) GetBucket() int64 {
	return t.Bucket
}

// GetTransactions fetches the set of transactions in a MoneyWell document, sorted by the display
// order as MoneyWell itself would render (and showing oldest to newest).
func GetTransactions(database *sql.DB) ([]Transaction, error) {
	rows, err := database.Query(`
            SELECT 
                za.Z_PK, 
                za.ZDATEYMD,
                za.ZTYPE,
                CAST(ROUND(za.ZAMOUNT * 100) AS INTEGER),
                COALESCE(za.ZBUCKET, za.ZBUCKET1, za.ZBUCKET2),
                COALESCE(za.ZACCOUNT, za.ZACCOUNT1, za.ZACCOUNT2),
                COALESCE(zat.ZACCOUNT, zat.ZACCOUNT1, zat.ZACCOUNT2),
                COALESCE(za.ZTRANSFERSIBLING, za.Z3_TRANSFERSIBLING),
                COALESCE(za.ZSPLITPARENT, za.Z3_SPLITPARENT),
                (
                    SELECT 
                        1 
                    FROM 
                        ZACTIVITY za2 
                    WHERE 
                        -- Only ZSPLITPARENT has an index, so assume that's the only one that
                        -- matters.
                        za2.ZSPLITPARENT = za.Z_PK 
                    LIMIT 1
                ) IS NOT NULL AS IsSplit,
                za.ZISBUCKETOPTIONAL,
                za.ZPAYEE,
                za.ZMEMO,
                zac.ZCURRENCYCODE
            FROM 
                ZACTIVITY za
            JOIN
                ZACCOUNT zac ON (zac.Z_PK = za.ZACCOUNT2)
            LEFT JOIN
                ZACCOUNTGROUP zacg ON ( zacg.Z_PK = zac.ZACCOUNTGROUP )
            LEFT JOIN
                ZACTIVITY zat ON ( zat.Z_PK = za.ZTRANSFERSIBLING )
            WHERE
                za.Z_ENT = 7
            ORDER BY
                za.ZDATEYMD ASC,
                zacg.ZSEQUENCE ASC,
                zac.ZSEQUENCE ASC,
                za.ZTYPE ASC
        `)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query transactions")
	}
	defer rows.Close()

	transactions := []Transaction{}

	var primaryKey, amountRaw int64
	var dateymd, transactionType int
	var bucket, account, transferAccount, transferSibling, splitParent sql.NullInt64
	var isSplit, isBucketOptional bool
	var payee, memo, currencyCode sql.NullString
	for rows.Next() {
		err := rows.Scan(
			&primaryKey,
			&dateymd,
			&transactionType,
			&amountRaw,
			&bucket,
			&account,
			&transferAccount,
			&transferSibling,
			&splitParent,
			&isSplit,
			&isBucketOptional,
			&payee,
			&memo,
			&currencyCode,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan transaction")
		}

		date, err := parseDateymd(dateymd)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse transaction date")
		}

		transactions = append(transactions, Transaction{
			PrimaryKey:      primaryKey,
			Date:            date,
			TransactionType: transactionType,
			Amount: money.Money{
				Currency: currencyCode.String,
				Amount:   amountRaw,
			},
			Bucket:           bucket.Int64,
			Account:          account.Int64,
			TransferAccount:  transferAccount.Int64,
			TransferSibling:  transferSibling.Int64,
			SplitParent:      splitParent.Int64,
			IsSplit:          isSplit,
			IsBucketOptional: isBucketOptional,
			Payee:            payee.String,
			Memo:             memo.String,
		})
	}

	return transactions, nil
}
