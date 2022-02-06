package cli

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/pkg/errors"

	"github.com/lieut-data/go-moneywell/api"

	_ "github.com/mattn/go-sqlite3"
)

func ListAccountGroups(database *sql.DB, verbose bool) error {
	accountGroups, err := api.GetAccountGroups(database)
	if err != nil {
		return errors.Wrap(err, "failed to fetch account groups")
	}

	for _, accountGroup := range accountGroups {
		primaryKey := ""
		if verbose {
			primaryKey = fmt.Sprintf(" [%d]", accountGroup.PrimaryKey)
		}

		fmt.Printf("%s%s\n", accountGroup.Name, primaryKey)
	}

	return nil
}

func ListAccounts(database *sql.DB, verbose bool) error {
	accountGroups, err := api.GetAccountGroups(database)
	if err != nil {
		return errors.Wrap(err, "failed to fetch account groups")
	}

	mappedAccountGroups := make(map[int64]api.AccountGroup, len(accountGroups))
	for _, accountGroup := range accountGroups {
		mappedAccountGroups[accountGroup.PrimaryKey] = accountGroup
	}

	accounts, err := api.GetAccounts(database)
	if err != nil {
		return errors.Wrap(err, "failed to fetch accounts")
	}

	transactions, err := api.GetTransactions(database)
	if err != nil {
		return errors.Wrap(err, "failed to get transactions")
	}

	var lastAccountGroup int64
	for _, account := range accounts {
		if lastAccountGroup == 0 || lastAccountGroup != account.AccountGroup {
			if account.AccountGroup > 0 {
				accountGroup := mappedAccountGroups[account.AccountGroup]
				if !verbose {
					fmt.Printf("%s\n", accountGroup.Name)
				} else {
					fmt.Printf(
						"%s (%d)\n",
						accountGroup.Name,
						accountGroup.PrimaryKey,
					)
				}
			}
			lastAccountGroup = account.AccountGroup
		}

		balance := api.GetAccountBalance(account, transactions)

		var indent string
		if account.AccountGroup > 0 {
			indent = "    "
		}

		primaryKey := ""
		if verbose {
			primaryKey = fmt.Sprintf(" [%d]", account.PrimaryKey)
		}

		fmt.Printf("%s%s (%s)%s\n", indent, account.Name, balance, primaryKey)
	}

	return nil
}

func ListBucketGroups(database *sql.DB, verbose bool) error {
	bucketGroups, err := api.GetBucketGroups(database)
	if err != nil {
		return errors.Wrap(err, "failed to fetch bucket groups")
	}

	var lastBucketType int64
	for _, bucketGroup := range bucketGroups {
		if lastBucketType == 0 || lastBucketType != bucketGroup.Type {
			if bucketGroup.Type == api.BucketGroupTypeIncome {
				fmt.Printf("Income\n")
			} else if bucketGroup.Type == api.BucketGroupTypeExpense {
				fmt.Printf("Expense\n")
			}
			lastBucketType = bucketGroup.Type
		}

		primaryKey := ""
		if verbose {
			primaryKey = fmt.Sprintf(" [%d]", bucketGroup.PrimaryKey)
		}

		fmt.Printf("    %s%s\n", bucketGroup.Name, primaryKey)
	}

	return nil
}

func ListBuckets(database *sql.DB, verbose bool) error {
	bucketGroups, err := api.GetBucketGroups(database)
	if err != nil {
		return errors.Wrap(err, "failed to fetch bucket groups")
	}

	mappedBucketGroups := make(map[int64]api.BucketGroup, len(bucketGroups))
	for _, bucketGroup := range bucketGroups {
		mappedBucketGroups[bucketGroup.PrimaryKey] = bucketGroup
	}

	settings, err := api.GetSettings(database)
	if err != nil {
		return errors.Wrap(err, "failed to get settings")
	}

	buckets, err := api.GetBuckets(database)
	if err != nil {
		return errors.Wrap(err, "failed to fetch buckets")
	}

	transactions, err := api.GetTransactions(database)
	if err != nil {
		return errors.Wrap(err, "failed to get transactions")
	}

	bucketTransfers, err := api.GetBucketTransfers(database)
	if err != nil {
		return errors.Wrap(err, "failed to get bucket transfers")
	}

	var lastBucketGroup, lastBucketGroupType int64
	for _, bucket := range buckets {
		if lastBucketGroup == 0 || lastBucketGroup != bucket.BucketGroup {
			if bucket.BucketGroup > 0 {
				bucketGroup := mappedBucketGroups[bucket.BucketGroup]

				if lastBucketGroupType == 0 || lastBucketGroupType != bucketGroup.Type {
					if bucketGroup.Type == 1 {
						fmt.Printf("Income\n")
					} else if bucketGroup.Type == 2 {
						fmt.Printf("Expense\n")
					}
					lastBucketGroupType = bucketGroup.Type
				}

				if !verbose {
					fmt.Printf("    %s\n", bucketGroup.Name)
				} else {
					fmt.Printf(
						"    %s (%d)\n",
						bucketGroup.Name,
						bucketGroup.PrimaryKey,
					)
				}
			}
			lastBucketGroup = bucket.BucketGroup
		}

		events, err := api.GetBucketEvents(bucket, transactions, bucketTransfers)
		if err != nil {
			return errors.Wrap(err, "failed to get bucket events")
		}

		balance, err := api.GetBucketBalance(bucket, events, settings)
		if err != nil {
			return errors.Wrap(err, "failed to get bucket balance")
		}

		var indent string
		if bucket.BucketGroup > 0 {
			indent = "    "
		}

		primaryKey := ""
		if verbose {
			primaryKey = fmt.Sprintf(" [%d]", bucket.PrimaryKey)
		}

		fmt.Printf("    %s%s (%s)%s\n", indent, bucket.Name, balance, primaryKey)
	}

	return nil
}

func ListTags(database *sql.DB, verbose bool) error {
	tags, err := api.GetTags(database)
	if err != nil {
		return errors.Wrap(err, "failed to fetch tags")
	}

	for _, tag := range tags {
		primaryKey := ""
		if verbose {
			primaryKey = fmt.Sprintf(" [%d]", tag.PrimaryKey)
		}

		fmt.Printf("%s%s\n", tag.Name, primaryKey)
	}

	return nil
}

func ListTransactions(
	database *sql.DB,
	accountFilter,
	bucketFilter,
	tagFilter string,
	verbose bool,
) error {
	transactions, err := api.GetTransactions(database)
	if err != nil {
		return errors.Wrap(err, "failed to fetch transactions")
	}

	accountsMap, err := api.GetAccountsMap(database)
	if err != nil {
		return errors.Wrap(err, "failed to fetch accounts map")
	}

	bucketsMap, err := api.GetBucketsMap(database)
	if err != nil {
		return errors.Wrap(err, "failed to fetch buckets map")
	}

	tagsMap, err := api.GetTagsMap(database)
	if err != nil {
		return errors.Wrap(err, "failed to fetch tags map")
	}

	transactionTagMap, err := api.GetTransactionTagMap(database)
	if err != nil {
		return errors.Wrap(err, "failed to fetch transaction tag map")
	}

	for _, transaction := range transactions {
		primaryKey := ""
		if verbose {
			primaryKey = fmt.Sprintf(" [%d]", transaction.PrimaryKey)
		}

		bucket := bucketsMap[transaction.Bucket]
		if len(bucketFilter) > 0 && bucket.Name != bucketFilter {
			continue
		}

		account := accountsMap[transaction.Account]
		if len(accountFilter) > 0 && account.Name != accountFilter {
			continue
		}

		transactionTags := transactionTagMap[transaction.PrimaryKey]
		if len(tagFilter) > 0 {
			found := false
			for _, tag := range transactionTags {
				if tagFilter == tagsMap[tag].Name {
					found = true
					continue
				}
			}

			if !found {
				continue
			}
		}

		memo := ""
		if len(transaction.Memo) > 0 {
			memo = fmt.Sprintf(" (%s)", transaction.Memo)
		}

		bucketName := ""
		if !account.IncludeInCashFlow && transaction.IsBucketOptional {
			bucketName = "(optional)"
		} else if transaction.IsTransfer() {
			bucketName = "(transfer)"
		} else {
			bucketName = bucket.Name
		}

		accountName := ""
		if transaction.IsTransfer() {
			switch transaction.TransactionType {
			case api.TransactionTypeDeposit:
				accountName = fmt.Sprintf(
					"%s from %s",
					account.Name,
					accountsMap[transaction.TransferAccount].Name,
				)
			case api.TransactionTypeWithdrawal:
				accountName = fmt.Sprintf(
					"%s to %s",
					account.Name,
					accountsMap[transaction.TransferAccount].Name,
				)
			}
		}

		fmt.Printf(
			"%s%s\t%s\t%s\t%s\t%s%s\n",
			transaction.Payee,
			memo,
			bucketName,
			accountName,
			transaction.Date.Format("Jan 2, 2006"),
			transaction.Amount,
			primaryKey,
		)
	}

	return nil
}

func ListRecurrenceRules(
	database *sql.DB,
	verbose bool,
) error {
	recurrenceRules, err := api.GetRecurrenceRules(database)
	if err != nil {
		return errors.Wrap(err, "failed to fetch recurrence rules")
	}

	sort.Sort(api.RecurrenceRuleSort(recurrenceRules))

	for _, recurrenceRule := range recurrenceRules {
		primaryKey := ""
		if verbose {
			primaryKey = fmt.Sprintf(" [%d]", recurrenceRule.PrimaryKey)
		}

		fmt.Printf("%s%s\n", api.DescribeRecurrenceRule(recurrenceRule), primaryKey)
	}

	return nil
}

func ListSpendingPlanEvents(
	database *sql.DB,
	bucketFilter string,
	verbose bool,
) error {
	spendingPlanEvents, err := api.GetSpendingPlan(database)
	if err != nil {
		return errors.Wrap(err, "failed to fetch spending plan")
	}

	recurrenceRulesMap, err := api.GetRecurrenceRulesMap(database)
	if err != nil {
		return errors.Wrap(err, "failed to fetch recurrence rules")
	}

	bucketsMap, err := api.GetBucketsMap(database)
	if err != nil {
		return errors.Wrap(err, "failed to fetch buckets map")
	}

	bucketGroups, err := api.GetBucketGroups(database)
	if err != nil {
		return errors.Wrap(err, "failed to fetch bucket groups")
	}

	mappedBucketGroups := make(map[int64]api.BucketGroup, len(bucketGroups))
	for _, bucketGroup := range bucketGroups {
		mappedBucketGroups[bucketGroup.PrimaryKey] = bucketGroup
	}

	bucketGroupTypes := []int64{api.BucketGroupTypeIncome, api.BucketGroupTypeExpense}
	for _, bucketGroupType := range bucketGroupTypes {
		headerPrinted := false

		for _, event := range spendingPlanEvents {
			primaryKey := ""
			if verbose {
				primaryKey = fmt.Sprintf(" [%d]", event.PrimaryKey)
			}

			recurrenceRule := recurrenceRulesMap[event.RecurrenceRule]
			fillRecurrenceRule := recurrenceRulesMap[event.FillRecurrenceRule]

			bucket := bucketsMap[event.Bucket]
			if len(bucketFilter) > 0 && bucket.Name != bucketFilter {
				continue
			}

			bucketGroup := mappedBucketGroups[bucket.BucketGroup]
			if bucketGroup.Type != bucketGroupType {
				continue
			}

			if !headerPrinted {
				if bucketGroup.Type == api.BucketGroupTypeIncome {
					fmt.Printf("Income\n")
				} else if bucketGroup.Type == api.BucketGroupTypeExpense {
					fmt.Printf("Expense\n")
				}

				headerPrinted = true
			}

			fmt.Printf(
				"    %s\t%s\t%s\t%s\t%s\t%s%s\n",
				event.Name,
				event.Date,
				bucket.Name,
				event.Amount,
				api.DescribeRecurrenceRule(recurrenceRule),
				api.DescribeFillRecurrenceRule(fillRecurrenceRule),
				primaryKey,
			)
		}
	}

	return nil
}
