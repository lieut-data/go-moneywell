package doctor

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/lieut-data/go-moneywell/api"

	_ "github.com/mattn/go-sqlite3"
)

// Diagnose analyzes the given MoneyWell document for potential issues.
func Diagnose(moneywellPath string) error {
	database, err := api.OpenDocument(moneywellPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", moneywellPath)
	}

	settings, err := api.GetSettings(database)
	if err != nil {
		return errors.Wrap(err, "failed to get settings")
	}

	accounts, err := api.GetAccounts(database)
	if err != nil {
		return errors.Wrap(err, "failed to get accounts")
	}

	transactions, err := api.GetTransactions(database)
	if err != nil {
		return errors.Wrap(err, "failed to get transactions")
	}

	problematicTransactions, err := GetProblematicTransactions(
		settings,
		accounts,
		transactions,
	)
	if err != nil {
		return errors.Wrap(err, "failed to query for problematic transactions")
	}

	for _, problematicTransaction := range problematicTransactions {
		fmt.Printf("WARNING: %s\n", problematicTransaction.Description)
	}

	return nil
}
