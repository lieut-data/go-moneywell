package doctor

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/lieut-data/go-moneywell/api"
	"github.com/lieut-data/go-moneywell/api/money"
)

const (
	// ProblemNotFullySplit identifies a split transaction whose bucketed children do not sum
	// to the transaction amount. This leads to an imbalance that doesn't show up as
	// unassigned.
	ProblemNotFullySplit = 1
	// ProblemSplitParentAssignedBucket identifies a split transaction that is itself assigned
	// a bucket. Only the children should be assigned to buckets. This problem has never been
	// observed in an actual MoneyWell document.
	ProblemSplitParentAssignedBucket = 2
	// ProblemTransferInsideCashFlowAssignedBucket identifies a transfer that should not be
	// assigned a bucket since both accounts are inside the cash flow.
	ProblemTransferInsideCashFlowAssignedBucket = 3
	// ProblemTransferOutsideCashFlowAssignedBucket identifies a transfer that should not be
	// assigned a bucket since both accounts are outside the cash flow.
	ProblemTransferOutsideCashFlowAssignedBucket = 4
	// ProblemTransferOutOfCashFlowMissingBucket identifies a transfer that should be assigned
	// a bucket since it moves money out of the cash flow.
	ProblemTransferOutOfCashFlowMissingBucket = 5
	// ProblemTransferFromCashFlowAssignedBucket identifies a transfer that should not be
	// assigned a bucket since it receives money from inside the cash flow.
	ProblemTransferFromCashFlowAssignedBucket = 6
	// ProblemBucketOptionalInsideCashFlow identifies a transaction marked as bucket optional
	// that should not be.
	ProblemBucketOptionalInsideCashFlow = 7
	// ProblemMissingBucketInsideCashFlow identifies a transaction missing an assigned bucket.
	ProblemMissingBucketInsideCashFlow = 8
)

// ProblematicTranscations represents a transaction diagnosed with a potential problem.
type ProblematicTransaction struct {
	Transaction int64
	Problem     int
	Description string
}

// GetProblematicTransactions finds transactions with potential problems, typically leading to
// an imbalance between accounts and buckets within MoneyWell.
func GetProblematicTransactions(
	settings api.Settings,
	accounts []api.Account,
	transactions []api.Transaction,
) ([]ProblematicTransaction, error) {
	problematicTransactions := []ProblematicTransaction{}

	for _, transaction := range transactions {
		// Ignore transactions before the cash flow start date. They won't contribute
		// to any current imbalance.
		if transaction.Date.Before(settings.CashFlowStartDate) {
			continue
		}

		// Ignore $0.00 transactions. These won't contribute to an imbalance, and might
		// be used to demarcate initial balances.
		if transaction.Amount.IsZero() {
			continue
		}

		account, err := getAccount(accounts, transaction.Account)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		problematicSplitTransactions, err := checkSplitTransaction(
			account,
			transactions,
			transaction,
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		problematicTransactions = append(
			problematicTransactions,
			problematicSplitTransactions...,
		)

		problematicTransferTransactions, err := checkTransferTransaction(
			accounts,
			account,
			transactions,
			transaction,
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		problematicTransactions = append(
			problematicTransactions,
			problematicTransferTransactions...,
		)

		problematicBucketOptionalTransactions, err := checkBucketOptionalTransaction(
			account,
			transactions,
			transaction,
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		problematicTransactions = append(
			problematicTransactions,
			problematicBucketOptionalTransactions...,
		)

		problematicMissingBucketTransactions, err := checkMissingBucketTransaction(
			account,
			transactions,
			transaction,
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		problematicTransactions = append(
			problematicTransactions,
			problematicMissingBucketTransactions...,
		)
	}

	return problematicTransactions, nil
}

func checkSplitTransaction(
	account api.Account,
	transactions []api.Transaction,
	transaction api.Transaction,
) ([]ProblematicTransaction, error) {
	if !transaction.IsSplit {
		return nil, nil
	}

	problematicTransactions := []ProblematicTransaction{}

	// The split parent in a transaction should not be assigned a bucket.
	if transaction.Bucket != 0 {
		problematicTransactions = append(problematicTransactions, ProblematicTransaction{
			Transaction: transaction.PrimaryKey,
			Problem:     ProblemSplitParentAssignedBucket,
			Description: fmt.Sprintf(
				"%s should not be assigned to a bucket",
				describeTransaction("split parent", account, transaction),
			),
		})
	}

	// The children of a split transaction should sum to the transaction amount. Otherwise,
	// this creates an imbalance that doesn't even show up in the "Unassigned" Smart Bucket
	// within MoneyWell.

	// Find and sum the children
	// TODO: Avoid O(n^2) only if this ever seems slow.
	childBalance := money.Money{}
	for _, child := range transactions {
		if child.SplitParent == transaction.PrimaryKey {
			childBalance = childBalance.Add(child.Amount)
		}
	}

	if transaction.Amount != childBalance {
		problematicTransactions = append(problematicTransactions, ProblematicTransaction{
			Transaction: transaction.PrimaryKey,
			Problem:     ProblemNotFullySplit,
			Description: fmt.Sprintf(
				"%s is not fully split (off by %s)",
				describeTransaction(
					"transaction",
					account,
					transaction,
				),
				transaction.Amount.Add(
					childBalance.Multiply(-1),
				),
			),
		})
	}

	return problematicTransactions, nil
}

func checkTransferTransaction(
	accounts []api.Account,
	account api.Account,
	transactions []api.Transaction,
	transaction api.Transaction,
) ([]ProblematicTransaction, error) {
	if !transaction.IsTransfer() {
		return nil, nil
	}

	// Assume split transactions are checked elsewhere.
	if transaction.IsSplit {
		return nil, nil
	}

	problematicTransactions := []ProblematicTransaction{}

	transferAccount, err := getAccount(accounts, transaction.TransferAccount)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// A transfer between accounts inside the cash flow should not have a bucket assigned.
	if account.IncludeInCashFlow && transferAccount.IncludeInCashFlow && transaction.Bucket != 0 {
		problematicTransactions = append(problematicTransactions, ProblematicTransaction{
			Transaction: transaction.PrimaryKey,
			Problem:     ProblemTransferInsideCashFlowAssignedBucket,
			Description: fmt.Sprintf(
				"%s between accounts in the cash flow should not be assigned to a bucket",
				describeTransaction("transfer", account, transaction),
			),
		})
	}

	// A transfer between accounts outside the cash flow should not have a bucket assigned.
	if !account.IncludeInCashFlow && !transferAccount.IncludeInCashFlow && transaction.Bucket != 0 {
		problematicTransactions = append(problematicTransactions, ProblematicTransaction{
			Transaction: transaction.PrimaryKey,
			Problem:     ProblemTransferOutsideCashFlowAssignedBucket,
			Description: fmt.Sprintf(
				"%s between accounts outside the cash flow should not be assigned to a bucket",
				describeTransaction("transfer", account, transaction),
			),
		})
	}

	// A transfer between accounts with one in the cash flow and one outside should have a
	// bucket assigned only on the account inside the cash flow.
	if account.IncludeInCashFlow && !transferAccount.IncludeInCashFlow && transaction.Bucket == 0 {
		problematicTransactions = append(problematicTransactions, ProblematicTransaction{
			Transaction: transaction.PrimaryKey,
			Problem:     ProblemTransferOutOfCashFlowMissingBucket,
			Description: fmt.Sprintf(
				"%s from account inside cash flow to account outside cash flow should be assigned to a bucket",
				describeTransaction("transfer", account, transaction),
			),
		})
	} else if !account.IncludeInCashFlow && transferAccount.IncludeInCashFlow && transaction.Bucket != 0 {
		problematicTransactions = append(problematicTransactions, ProblematicTransaction{
			Transaction: transaction.PrimaryKey,
			Problem:     ProblemTransferFromCashFlowAssignedBucket,
			Description: fmt.Sprintf(
				"%s from account outside cash flow to account inside cash flow should not be assigned to a bucket",
				describeTransaction("transfer", account, transaction),
			),
		})
	}

	return problematicTransactions, nil
}

func checkBucketOptionalTransaction(
	account api.Account,
	transactions []api.Transaction,
	transaction api.Transaction,
) ([]ProblematicTransaction, error) {
	// If it's not marked as bucket optional, it's not a problem!
	if !transaction.IsBucketOptional {
		return nil, nil
	}

	// Assume split transactions are checked elsewhere.
	if transaction.IsSplit {
		return nil, nil
	}

	// A transaction against an account outside the cash flow won't impact the cash flow.
	if !account.IncludeInCashFlow {
		return nil, nil
	}

	// A transaction marked as bucket optional but that strangely has a bucket assigned seems
	// to occur from time to time normally.
	if transaction.Bucket != 0 {
		return nil, nil
	}

	// Assume transfer transfers are checked elsewhere.
	if transaction.IsTransfer() {
		return nil, nil
	}

	// A transaction should generally not be marked as bucket optional in an account that is
	// part of the of the cash flow.
	return []ProblematicTransaction{
		{
			Transaction: transaction.PrimaryKey,
			Problem:     ProblemBucketOptionalInsideCashFlow,
			Description: fmt.Sprintf(
				"%s should not be marked as bucket optional in a cash flow account",
				describeTransaction("transaction", account, transaction),
			),
		},
	}, nil
}

func checkMissingBucketTransaction(
	account api.Account,
	transactions []api.Transaction,
	transaction api.Transaction,
) ([]ProblematicTransaction, error) {
	// If a bucket is assigned, it's not missing!
	if transaction.Bucket != 0 {
		return nil, nil
	}

	// Assume transfer and split transactions are checked elsewhere.
	if transaction.IsTransfer() || transaction.IsSplit {
		return nil, nil
	}

	// A transaction against an account outside the cash flow won't impact the cash flow.
	if !account.IncludeInCashFlow {
		return nil, nil
	}

	return []ProblematicTransaction{
		{
			Transaction: transaction.PrimaryKey,
			Problem:     ProblemMissingBucketInsideCashFlow,
			Description: fmt.Sprintf(
				"%s is not assigned to a bucket",
				describeTransaction("transaction", account, transaction),
			),
		},
	}, nil
}

func describeTransaction(description string, account api.Account, transaction api.Transaction) string {
	memo := transaction.Memo
	if len(memo) > 1 {
		memo = fmt.Sprintf(" (%s)", memo)
	}
	return fmt.Sprintf(
		"%s[%d] on %s against %s for %s%s",
		description,
		transaction.PrimaryKey,
		transaction.Date.Format("2006-01-02"),
		account.Name,
		transaction.Amount,
		memo,
	)
}

func getAccount(accounts []api.Account, accountPrimaryKey int64) (api.Account, error) {
	// TODO: Avoid O(n^2) only if this ever seems slow.
	for _, account := range accounts {
		if account.PrimaryKey == accountPrimaryKey {
			return account, nil
		}
	}

	return api.Account{}, errors.Errorf("failed to find account %v", accountPrimaryKey)
}
