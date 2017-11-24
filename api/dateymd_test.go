package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDateymd(t *testing.T) {
	testCases := []struct {
		Description string
		Input       int

		ExpectedTime  time.Time
		ExpectedError bool
	}{
		{
			"0",
			0,
			time.Time{},
			false,
		},
		{
			"negative value",
			-1,
			time.Time{},
			true,
		},
		{
			"invalid year",
			0102,
			time.Time{},
			true,
		},
		{
			"invalid month",
			20061301,
			time.Time{},
			true,
		},
		{
			"invalid day",
			20060145,
			time.Time{},
			true,
		},
		{
			"January 1, 2017",
			20170101,
			time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
			false,
		},
		{
			"March 1, 2006",
			20060301,
			time.Date(2006, 3, 1, 0, 0, 0, 0, time.UTC),
			false,
		},
		{
			"August 31, 2020",
			20200831,
			time.Date(2020, 8, 31, 0, 0, 0, 0, time.UTC),
			false,
		},
		{
			"December 31, 2017",
			20171231,
			time.Date(2017, 12, 31, 0, 0, 0, 0, time.UTC),
			false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Description, func(t *testing.T) {
			t.Parallel()

			actualTime, actualErr := parseDateymd(testCase.Input)

			assert.Equal(t, testCase.ExpectedTime, actualTime)
			if testCase.ExpectedError {
				assert.Error(t, actualErr)
			} else {
				assert.NoError(t, actualErr)
			}
		})
	}
}
