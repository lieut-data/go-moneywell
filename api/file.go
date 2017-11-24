package api

import (
	"database/sql"
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"

	_ "github.com/mattn/go-sqlite3"
)

// OpenDocument opens the given MoneyWell document and returns an interface to the SQLite database
// therein.
func OpenDocument(moneywellPath string) (*sql.DB, error) {
	fi, err := os.Stat(moneywellPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to stat %s", moneywellPath)
	}

	switch mode := fi.Mode(); {
	// Assume the `.moneywell` bundle was specified and resolve to the persistentStore therein.
	case mode.IsDir():
		moneywellPath = path.Join(moneywellPath, "StoreContent/persistentStore")
	case mode.IsRegular():
		// Attempt to access the path as-is.
	}

	database, err := sql.Open("sqlite3", fmt.Sprintf("%s?mode=ro", moneywellPath))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open database")
	}

	return database, nil
}
