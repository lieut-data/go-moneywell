package api_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lieut-data/go-moneywell/api"
	"github.com/lieut-data/go-moneywell/api/money"
)

func TestGetTransactions(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)

	transactions, err := api.GetTransactions(database)
	assert.NoError(t, err)

	expectedTransactions := []api.Transaction{
		{
			PrimaryKey:       1,
			Date:             time.Date(2017, 11, 1, 0, 0, 0, 0, time.UTC),
			TransactionType:  api.TransactionTypeDeposit,
			Amount:           money.Money{Currency: "CAD", Amount: 0},
			Bucket:           0,
			Account:          1,
			TransferAccount:  0,
			TransferSibling:  0,
			SplitParent:      0,
			IsSplit:          false,
			IsBucketOptional: true,
			Status:           api.TransactionStatusReconciled,
			Payee:            "Initial Balance",
			Memo:             "",
		},
		{
			PrimaryKey:       2,
			Date:             time.Date(2017, 11, 1, 0, 0, 0, 0, time.UTC),
			TransactionType:  api.TransactionTypeDeposit,
			Amount:           money.Money{Currency: "CAD", Amount: 1000 * 100},
			Bucket:           3,
			Account:          1,
			TransferAccount:  0,
			TransferSibling:  0,
			SplitParent:      0,
			IsSplit:          false,
			IsBucketOptional: false,
			Status:           api.TransactionStatusCleared,
			Payee:            "Work",
			Memo:             "",
		},
		{
			PrimaryKey:       11,
			Date:             time.Date(2017, 11, 5, 0, 0, 0, 0, time.UTC),
			TransactionType:  api.TransactionTypeDeposit,
			Amount:           money.Money{Currency: "CAD", Amount: 0},
			Bucket:           0,
			Account:          3,
			TransferAccount:  0,
			TransferSibling:  0,
			SplitParent:      0,
			IsSplit:          false,
			IsBucketOptional: true,
			Status:           api.TransactionStatusReconciled,
			Payee:            "Initial Balance",
			Memo:             "",
		},
		{
			PrimaryKey:       4,
			Date:             time.Date(2017, 11, 5, 0, 0, 0, 0, time.UTC),
			TransactionType:  api.TransactionTypeWithdrawal,
			Amount:           money.Money{Currency: "CAD", Amount: -1 * 350 * 100},
			Bucket:           13,
			Account:          1,
			TransferAccount:  0,
			TransferSibling:  0,
			SplitParent:      0,
			IsSplit:          false,
			IsBucketOptional: false,
			Status:           api.TransactionStatusOpen,
			Payee:            "Grocery Store",
			Memo:             "Chick peas and tuna.",
		},
		{
			PrimaryKey:       5,
			Date:             time.Date(2017, 11, 10, 0, 0, 0, 0, time.UTC),
			TransactionType:  api.TransactionTypeWithdrawal,
			Amount:           money.Money{Currency: "CAD", Amount: -1 * 500 * 100},
			Bucket:           2,
			Account:          1,
			TransferAccount:  0,
			TransferSibling:  0,
			SplitParent:      0,
			IsSplit:          false,
			IsBucketOptional: false,
			Status:           api.TransactionStatusOpen,
			Payee:            "Rent",
			Memo:             "",
		},
		{
			PrimaryKey:       15,
			Date:             time.Date(2017, 11, 12, 0, 0, 0, 0, time.UTC),
			TransactionType:  api.TransactionTypeDeposit,
			Amount:           money.Money{Currency: "CAD", Amount: 400 * 100},
			Bucket:           0,
			Account:          3,
			TransferAccount:  1,
			TransferSibling:  16,
			SplitParent:      0,
			IsSplit:          false,
			IsBucketOptional: true,
			Status:           api.TransactionStatusOpen,
			Payee:            "Split Test",
			Memo:             "",
		},
		{
			PrimaryKey:       12,
			Date:             time.Date(2017, 11, 12, 0, 0, 0, 0, time.UTC),
			TransactionType:  api.TransactionTypeDeposit,
			Amount:           money.Money{Currency: "USD", Amount: 0},
			Bucket:           0,
			Account:          4,
			TransferAccount:  0,
			TransferSibling:  0,
			SplitParent:      0,
			IsSplit:          false,
			IsBucketOptional: true,
			Status:           api.TransactionStatusReconciled,
			Payee:            "Initial Balance",
			Memo:             "",
		},
		{
			PrimaryKey:       13,
			Date:             time.Date(2017, 11, 12, 0, 0, 0, 0, time.UTC),
			TransactionType:  api.TransactionTypeWithdrawal,
			Amount:           money.Money{Currency: "CAD", Amount: -100 * 100},
			Bucket:           13,
			Account:          1,
			TransferAccount:  0,
			TransferSibling:  0,
			SplitParent:      14,
			IsSplit:          false,
			IsBucketOptional: false,
			Status:           api.TransactionStatusOpen,
			Payee:            "Split Test",
			Memo:             "",
		},
		{
			PrimaryKey:       14,
			Date:             time.Date(2017, 11, 12, 0, 0, 0, 0, time.UTC),
			TransactionType:  api.TransactionTypeWithdrawal,
			Amount:           money.Money{Currency: "CAD", Amount: -500 * 100},
			Bucket:           0,
			Account:          1,
			TransferAccount:  0,
			TransferSibling:  0,
			SplitParent:      0,
			IsSplit:          true,
			IsBucketOptional: false,
			Status:           api.TransactionStatusOpen,
			Payee:            "Split Test",
			Memo:             "",
		},
		{
			PrimaryKey:       16,
			Date:             time.Date(2017, 11, 12, 0, 0, 0, 0, time.UTC),
			TransactionType:  api.TransactionTypeWithdrawal,
			Amount:           money.Money{Currency: "CAD", Amount: -400 * 100},
			Bucket:           0,
			Account:          1,
			TransferAccount:  3,
			TransferSibling:  15,
			SplitParent:      14,
			IsSplit:          false,
			IsBucketOptional: true,
			Status:           api.TransactionStatusOpen,
			Payee:            "Split Test",
			Memo:             "",
		},
		{
			PrimaryKey:       8,
			Date:             time.Date(2017, 11, 12, 0, 0, 0, 0, time.UTC),
			TransactionType:  api.TransactionTypeDeposit,
			Amount:           money.Money{Currency: "CAD", Amount: 100 * 100},
			Bucket:           0,
			Account:          2,
			TransferAccount:  0,
			TransferSibling:  0,
			SplitParent:      0,
			IsSplit:          false,
			IsBucketOptional: true,
			Status:           api.TransactionStatusReconciled,
			Payee:            "Initial Balance",
			Memo:             "",
		},
		{
			PrimaryKey:       10,
			Date:             time.Date(2017, 11, 12, 0, 0, 0, 0, time.UTC),
			TransactionType:  api.TransactionTypeDeposit,
			Amount:           money.Money{Currency: "CAD", Amount: 0},
			Bucket:           0,
			Account:          5,
			TransferAccount:  0,
			TransferSibling:  0,
			SplitParent:      0,
			IsSplit:          false,
			IsBucketOptional: true,
			Status:           api.TransactionStatusReconciled,
			Payee:            "Initial Balance",
			Memo:             "",
		},
		{
			PrimaryKey:       9,
			Date:             time.Date(2017, 11, 12, 0, 0, 0, 0, time.UTC),
			TransactionType:  api.TransactionTypeDeposit,
			Amount:           money.Money{Currency: "CAD", Amount: 0},
			Bucket:           0,
			Account:          6,
			TransferAccount:  0,
			TransferSibling:  0,
			SplitParent:      0,
			IsSplit:          false,
			IsBucketOptional: true,
			Status:           api.TransactionStatusReconciled,
			Payee:            "Initial Balance",
			Memo:             "",
		},
		{
			PrimaryKey:       20,
			Date:             time.Date(2017, 11, 25, 0, 0, 0, 0, time.UTC),
			TransactionType:  api.TransactionTypeWithdrawal,
			Amount:           money.Money{Currency: "CAD", Amount: -100 * 100},
			Bucket:           3,
			Account:          1,
			TransferAccount:  0,
			TransferSibling:  0,
			SplitParent:      0,
			IsSplit:          false,
			IsBucketOptional: false,
			Status:           api.TransactionStatusVoided,
			Payee:            "Voided",
			Memo:             "Voided transaction.",
		},
		{
			PrimaryKey:       18,
			Date:             time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC),
			TransactionType:  api.TransactionTypeWithdrawal,
			Amount:           money.Money{Currency: "CAD", Amount: -100 * 100},
			Bucket:           3,
			Account:          1,
			TransferAccount:  0,
			TransferSibling:  0,
			SplitParent:      0,
			IsSplit:          false,
			IsBucketOptional: false,
			Status:           api.TransactionStatusPending,
			Payee:            "Future",
			Memo:             "Future transaction.",
		},
	}

	assert.Equal(t, expectedTransactions, transactions)
}

func TestIsTrnasfer(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)

	transactions, err := api.GetTransactions(database)
	assert.NoError(t, err)

	expectedTransferTransactions := []int64{15, 16}
	actualTransferTransactions := []int64{}
	for _, transaction := range transactions {
		if transaction.IsTransfer() {
			actualTransferTransactions = append(
				actualTransferTransactions,
				transaction.PrimaryKey,
			)
		}
	}

	assert.Equal(t, expectedTransferTransactions, actualTransferTransactions)
}
