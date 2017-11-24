package main

import (
	"flag"
	"fmt"

	"github.com/lieut-data/go-moneywell/api"
	"github.com/lieut-data/go-moneywell/internal/cli"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var verbose bool
	var moneywellPath, list, tag, bucket, account string
	flag.BoolVar(&verbose, "verbose", false, "be more verbose")
	flag.StringVar(&moneywellPath, "file", "", "the path to the MoneyWell document")
	flag.StringVar(&list, "list", "", "list the given entity")
	flag.StringVar(&account, "account", "", "the bucket by which to filter transactions")
	flag.StringVar(&bucket, "bucket", "", "the bucket by which to filter transactions")
	flag.StringVar(&tag, "tag", "", "the tag by which to filter transactions")

	flag.Parse()

	if moneywellPath == "" {
		fmt.Println("required: path to moneywell document")
		return
	}

	database, err := api.OpenDocument(moneywellPath)
	if err != nil {
		fmt.Printf("failed to open database: %v\n", err)
		return
	}

	switch list {
	case "account-groups":
		err = cli.ListAccountGroups(database, verbose)
	case "accounts":
		err = cli.ListAccounts(database, verbose)
	case "bucket-groups":
		err = cli.ListBucketGroups(database, verbose)
	case "buckets":
		err = cli.ListBuckets(database, verbose)
	case "tags":
		err = cli.ListTags(database, verbose)
	case "transactions":
		err = cli.ListTransactions(database, account, bucket, tag, verbose)
	}

	if err != nil {
		fmt.Printf("cli failed: %v\n", err)
		return
	}

	if err := database.Close(); err != nil {
		fmt.Printf("failed to close database: %v\n", err)
		return
	}
}
