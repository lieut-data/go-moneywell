# go-moneywell
[![GoDoc](https://godoc.org/github.com/lieut-data/go-moneywell/api?status.svg)](http://godoc.org/github.com/lieut-data/go-moneywell/api)

A Go (golang) API for programmatic access to a [MoneyWell](https://moneywellapp.com/) document, 
along with various tools leveraging this API to improve common workflows in using MoneyWell.

This project is in no way affiliated with [No Thirst Software](http://nothirst.com/company/).

## Getting Started

### Prerequisites

You need to [install Go](https://golang.org/doc/install). If you're using MoneyWell, you're 
also using OSX, so use the [Homebrew](https://brew.sh/) package manager and simply:

    brew install go

### Installation

To install go-moneywell for the command line tools, use `go install`:

    go install github.com/lieut-data/go-moneywell/cmd/moneywellcli
    go install github.com/lieut-data/go-moneywell/cmd/moneywelldoctor

To install go-moneywell for the api, use `go get`:

    go get github.com/lieut-data/go-moneywell/api

### Running

To analyze a MoneyWell document for problems, use [moneywelldoctor](#moneywelldoctor):

    moneywelldoctor Finances.moneywell

To dump various data from a MoneyWell document, use [moneywellcli](#moneywellcli):

    moneywellcli -file Finances.moneywell -list accounts
    moneywellcli -file Finances.moneywell -list buckets
    moneywellcli -file Finances.moneywell -list tags
    moneywellcli -file Finances.moneywell -list account-groups
    moneywellcli -file Finances.moneywell -list bucket-groups
    moneywellcli -file Finances.moneywell -list transactions
    moneywellcli -file Finances.moneywell -list recurrence-rules
    moneywellcli -file Finances.moneywell -list spending-plan

Optionally filter transactions by account, bucket or tag:

    moneywellcli -file Finances.moneywell -list transactions -account "Chequing"
    moneywellcli -file Finances.moneywell -list transactions -bucket "Salary"
    moneywellcli -file Finances.moneywell -list transactions -tag "family_vacation_2017"

## Command-line Tools

### [moneywelldoctor](cmd/moneywelldoctor)

In short, moneywelldoctor is designed to solve this problem:

![Account/bucket imbalance, but no unassigned transactions.](internal/doctor/why-moneywelldoctor.png?raw=true)

MoneyWell always warns when the account and bucket totals are mismatched. Most of the time, this
is simply because there are some newly entered transactions against accounts in the cash flow that
have not been assigned a bucket. MoneyWell provides a helpful `Unassigned` smart filter to track
these down. 

Unfortunately, there are other ways to create such an imbalance, all of which moneywelldoctor
is designed to detect:
* A split transaction whose children do not sum to the transaction amount.
* A transfer between accounts both inside or outside the cash flow that is incorrectly assigned a 
bucket.
* A transfer between an account inside the cash flow and an account outside the cash flow that
is missing a bucket.
* A transaction incorrectly marked as bucket optional.

Finding these issues previously involved a "binary search" through Time Machine to discover which
transaction introduced the imbalance, or giving up and resetting the cash flow start date. Given
the path to a `*.moneywell` document, `moneywelldoctor` will instead pin down exactly what
transactions are at fault:

    WARNING: transaction[3] on 2017-11-19 against Chequing for -$100.01 CAD (Cash Rebate) is not fully split (off by -$0.01 CAD)
    WARNING: transfer[15] on 2017-11-19 against Cash Account for -$50.00 CAD (Withdrawal for buying movie tickets) between accounts in the cash flow should not be assigned to a bucket

Note that `moneywelldoctor` will not make any changes to the given MoneyWell document. Any
transactions identified must be then fixed within MoneyWell itself.

### [moneywellcli](cmd/moneywellcli)

The bundled `moneywellcli` is largely an exercise of the [api](api)
package. It supports various operations, dumping the resulting data structures to STDOUT:

    moneywellcli -file Finances.moneywell -list accounts
    moneywellcli -file Finances.moneywell -list buckets
    moneywellcli -file Finances.moneywell -list tags
    moneywellcli -file Finances.moneywell -list account-groups
    moneywellcli -file Finances.moneywell -list bucket-groups
    moneywellcli -file Finances.moneywell -list transactions
    moneywellcli -file Finances.moneywell -list transactions -account "Chequing"
    moneywellcli -file Finances.moneywell -list transactions -bucket "Salary"
    moneywellcli -file Finances.moneywell -list transactions -tag "family_vacation_2017"
    moneywellcli -file Finances.moneywell -list recurrence-rules
    moneywellcli -file Finances.moneywell -list spending-plan
    moneywellcli -file Finances.moneywell -list spending-plan -bucket "Tech"

The API to this command line tool is subject to change. A future revision will likely support CSV 
encoding for export to spreadsheets along with JSON encoding for integration with other scripts.

## Packages

### [api](api)

The command line tools above depend on the primary export of this repository: a read-only,
programmatic API into a MoneyWell document. A MoneyWell document is itself a Core Data
SQLite database. You can inspect the contents of this database by installing sqlite (if necessary):

    brew install sqlite

and then opening the `persistentStore` database therein:

    sqlite3 api/Test.moneywell/StoreContent/persistentStore

Within the sqlite shell, you can list the tables using `.tables`, and run queries such as:

    SELECT * FROM ZBUCKET;

This API is effectively a low-level wrapper around such queries, taking into account knowledge
about how MoneyWell arranges its data model. The API is currently read-only and compatible with
documents readable to at least MoneyWell 3.0.5. 

For detailed documentation on the API, see the generated [GoDoc](https://godoc.org/github.com/lieut-data/go-moneywell/api).

Note that nothing precludes a future version of this API from writing to the MoneyWell document.
A logical place to start would be tag management: MoneyWell makes it hard to see all tags and
remove or rename tags in bulk.

#### TODO

Some random notes on things in the API yet to be tackled:
- [ ] Add support for querying favourite transactions.
- [ ] Fix order of hidden buckets relative to other buckets.
- [ ] Fix sorting of accounts outside of account groups.
- [ ] Fix sorting of buckets outside of bucket groups.
- [ ] Make the bucket event order fully match MoneyWell's own listing within a given day.
- [ ] Determine why the sqlite3 driver leaves persistentStore-wal and persistentStore-shm around.
- [ ] Update moneywelldoctor to warn about scheduled transactions in the past

## Authors

* Jesse Hallam ([keybase.io/lieutdata](https://keybase.io/lieutdata))

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
