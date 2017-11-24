package api

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// parseDateymd parses an integer representing a date in a MoneyWell document into a time.Time.
func parseDateymd(dateymd int) (time.Time, error) {
	if dateymd == 0 {
		return time.Time{}, nil
	}

	date, err := time.Parse("20060102", strconv.Itoa(dateymd))
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "failed to parse time %d", dateymd)
	}

	return date, nil
}
