package doctor_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lieut-data/go-moneywell/api"
	"github.com/lieut-data/go-moneywell/internal/doctor"
)

func TestDiagnose(t *testing.T) {
	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)

	settings, err := api.GetSettings(database)
	assert.NoError(t, err)

	accounts, err := api.GetAccounts(database)
	assert.NoError(t, err)

	transactions, err := api.GetTransactions(database)
	assert.NoError(t, err)

	actualProblematicTransactions, err := doctor.GetProblematicTransactions(
		settings,
		accounts,
		transactions,
	)
	assert.NoError(t, err)

	expectedProblematicTransactions := []doctor.ProblematicTransaction{
		{
			Transaction: 3,
			Problem:     doctor.ProblemNotFullySplit,
			Description: "transaction[3] on 2017-11-19 against Inside Cash Flow #1 for -$100.01 CAD (Not Fully Split Transaction) is not fully split (off by -$0.01 CAD)",
		},
		{
			Transaction: 15,
			Problem:     doctor.ProblemTransferInsideCashFlowAssignedBucket,
			Description: "transfer[15] on 2017-11-19 against Inside Cash Flow #1 for -$50.00 CAD (Transfer inside cash flow assigned bucket) between accounts in the cash flow should not be assigned to a bucket",
		},
		{
			Transaction: 25,
			Problem:     doctor.ProblemBucketOptionalInsideCashFlow,
			Description: "transaction[25] on 2017-11-19 against Inside Cash Flow #1 for -$0.01 CAD (Bucket Optional) should not be marked as bucket optional in a cash flow account",
		},
		{
			Transaction: 25,
			Problem:     doctor.ProblemMissingBucketInsideCashFlow,
			Description: "transaction[25] on 2017-11-19 against Inside Cash Flow #1 for -$0.01 CAD (Bucket Optional) is not assigned to a bucket",
		},
		{
			Transaction: 27,
			Problem:     doctor.ProblemMissingBucketInsideCashFlow,
			Description: "transaction[27] on 2017-11-19 against Inside Cash Flow #1 for -$5.00 CAD (Missing bucket) is not assigned to a bucket",
		},
		{
			Transaction: 20,
			Problem:     doctor.ProblemTransferOutOfCashFlowMissingBucket,
			Description: "transfer[20] on 2017-11-19 against Inside Cash Flow #2 for -$25.00 CAD (Transfer money out of cash flow) from account inside cash flow to account outside cash flow should be assigned to a bucket",
		},
		{
			Transaction: 17,
			Problem:     doctor.ProblemTransferOutsideCashFlowAssignedBucket,
			Description: "transfer[17] on 2017-11-19 against Outside Cash Flow #1 for -$50.00 CAD (Transfer outside cash flow assigned bucket) between accounts outside the cash flow should not be assigned to a bucket",
		},
		{
			Transaction: 21,
			Problem:     doctor.ProblemTransferFromCashFlowAssignedBucket,
			Description: "transfer[21] on 2017-11-19 against Outside Cash Flow #2 for $25.00 CAD (Transfer money out of cash flow) from account outside cash flow to account inside cash flow should not be assigned to a bucket",
		},
	}

	assert.Equal(t, expectedProblematicTransactions, actualProblematicTransactions)
}
