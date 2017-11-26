package api

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/lieut-data/go-moneywell/api/money"
)

// Account represents an cash account, savings account, chequing account, line of credit, credit
// card or other type of account in a MoneyWell document. An account correlates 1:1 with a record
// in the ZACCOUNT table. Not all columns are exported.
//
// The MoneyWell SQLite schema for the ZACCOUNT table is as follows:
//  > .schema ZACCOUNT
//  CREATE TABLE ZACCOUNT (
//      Z_PK INTEGER PRIMARY KEY,
//      Z_ENT INTEGER,
//      Z_OPT INTEGER,
//      ZALLOWAUTORECONCILE INTEGER,
//      ZALLOWONLINEBANKING INTEGER,
//      ZALLOWONLINEBILLPAY INTEGER,
//      ZBALANCEDATEYMD INTEGER,
//      ZINCLUDEINCASHFLOW INTEGER,
//      ZINCLUDEINNETWORTH INTEGER,
//      ZINTRODUCTORYDATEYMD INTEGER,
//      ZISBUCKETOPTIONAL INTEGER,
//      ZISDEBT INTEGER,
//      ZISHIDDEN INTEGER,
//      ZISSELECTED INTEGER,
//      ZISTAXDEFERRED INTEGER,
//      ZONLINEAVAILABLEDATEYMD INTEGER,
//      ZONLINEENDDATEYMD INTEGER,
//      ZONLINELEDGERDATEYMD INTEGER,
//      ZONLINESTARTDATEYMD INTEGER,
//      ZSEQUENCE INTEGER,
//      ZSTARTINGDATEYMD INTEGER,
//      ZTYPE INTEGER,
//      ZACCOUNTGROUP INTEGER,
//      ZONLINECONNECTION INTEGER,
//      ZBALANCE DECIMAL,
//      ZCREDITLIMIT DECIMAL,
//      ZINTERESTRATE DECIMAL,
//      ZINTRODUCTORYRATE DECIMAL,
//      ZMINIMUMBALANCE DECIMAL,
//      ZONLINEAVAILABLEBALANCE DECIMAL,
//      ZONLINELEDGERBALANCE DECIMAL,
//      ZACCOUNTNUMBER VARCHAR,
//      ZCURRENCYCODE VARCHAR,
//      ZIMPORTFORMAT VARCHAR,
//      ZMEMO VARCHAR,
//      ZNAME VARCHAR,
//      ZROUTINGNUMBER VARCHAR,
//      ZTICDSSYNCID VARCHAR,
//      ZUNIQUEID VARCHAR
//  );
type Account struct {
	PrimaryKey        int64
	Name              string
	Balance           money.Money
	IsBucketOptional  bool
	IncludeInCashFlow bool
	CurrencyCode      string
	AccountGroup      int64
}

// GetAccounts fetches the set of accounts in a MoneyWell document, sorted by the display order
// as MoneyWell itself would render.
func GetAccounts(database *sql.DB) ([]Account, error) {
	rows, err := database.Query(`
            SELECT 
                za.Z_PK, 
                za.ZNAME,
                CAST(ROUND(za.ZBALANCE * 100) AS INTEGER),
                za.ZISBUCKETOPTIONAL,
                za.ZINCLUDEINCASHFLOW,
                za.ZCURRENCYCODE,
                za.ZACCOUNTGROUP
            FROM 
                ZACCOUNT za
            LEFT JOIN
                ZACCOUNTGROUP zag ON ( zag.Z_PK = za.ZACCOUNTGROUP )
            ORDER BY
                zag.ZSEQUENCE ASC,
                za.ZSEQUENCE ASC
        `)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query accounts")
	}
	defer rows.Close()

	accounts := []Account{}

	var primaryKey int64
	var name string
	var balanceRaw, accountGroup sql.NullInt64
	var isBucketOptional, includeInCashFlow int
	var currencyCode sql.NullString
	for rows.Next() {
		err := rows.Scan(
			&primaryKey,
			&name,
			&balanceRaw,
			&isBucketOptional,
			&includeInCashFlow,
			&currencyCode,
			&accountGroup,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan account")
		}

		accounts = append(accounts, Account{
			PrimaryKey: primaryKey,
			Name:       name,
			Balance: money.Money{
				Currency: currencyCode.String,
				Amount:   balanceRaw.Int64,
			},
			IsBucketOptional:  isBucketOptional > 0,
			IncludeInCashFlow: includeInCashFlow > 0,
			CurrencyCode:      currencyCode.String,
			AccountGroup:      accountGroup.Int64,
		})
	}

	return accounts, nil
}

// GetAccountsMap gets a map from the account primary key to the account.
func GetAccountsMap(database *sql.DB) (map[int64]Account, error) {
	accounts, err := GetAccounts(database)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	accountsMap := make(map[int64]Account, len(accounts))
	for _, account := range accounts {
		accountsMap[account.PrimaryKey] = account
	}

	return accountsMap, nil
}

// GetAccountBalance uses the given transactions to compute the balance of an account at the given
// time.
func GetAccountBalance(account Account, transactions []Transaction) money.Money {
	balance := money.Money{}

	for _, transaction := range transactions {
		// Ignore voided and pending transactions.
		switch transaction.Status {
		case TransactionStatusVoided:
			fallthrough
		case TransactionStatusPending:
			continue
		}

		// Split transactions result in both a parent transaction and children
		// transactions. Exclude the children and count only the parent transaction to
		// avoid double summing.
		if transaction.Account == account.PrimaryKey && transaction.SplitParent == 0 {
			balance = balance.Add(transaction.Amount)
		}
	}

	return balance
}
