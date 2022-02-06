package api_test

import (
	// "fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lieut-data/go-moneywell/api"
	"github.com/lieut-data/go-moneywell/api/money"
)

func TestGetSpendingPlan(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	spendingPlan, err := api.GetSpendingPlan(database)
	assert.NoError(t, err)

	expectedSpendingPlan := []api.SpendingPlan{
		{
			PrimaryKey:         21,
			Date:               time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (Never -> Every Event Date)",
			Amount:             money.Money{Amount: 1000, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     0,
			FillRecurrenceRule: 0,
		},

		{
			PrimaryKey:         22,
			Date:               time.Date(2018, 4, 2, 0, 0, 0, 0, time.UTC),
			Name:               "Hobbies (Every Day - Every Event Date)",
			Amount:             money.Money{Amount: 2000, Currency: "CAD"},
			Bucket:             27,
			RecurrenceRule:     15,
			FillRecurrenceRule: 0,
		},

		{
			PrimaryKey:         23,
			Date:               time.Date(2018, 4, 3, 0, 0, 0, 0, time.UTC),
			Name:               "Mortgage/Rent (Every Week -> Every Day)",
			Amount:             money.Money{Amount: 3000, Currency: "CAD"},
			Bucket:             2,
			RecurrenceRule:     16,
			FillRecurrenceRule: 41,
		},

		{
			PrimaryKey:         24,
			Date:               time.Date(2018, 4, 12, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (Every 3 Months -> Every 2 months)",
			Amount:             money.Money{Amount: 12000, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     17,
			FillRecurrenceRule: 34,
		},

		{
			PrimaryKey:         25,
			Date:               time.Date(2018, 4, 7, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (Every Month -> Every 4 weeks)",
			Amount:             money.Money{Amount: 7000, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     24,
			FillRecurrenceRule: 30,
		},

		{
			PrimaryKey:         26,
			Date:               time.Date(2018, 4, 11, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (Every 2 Months -> Every month on the 1st and 16th)",
			Amount:             money.Money{Amount: 11000, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     23,
			FillRecurrenceRule: 35,
		},

		{
			PrimaryKey:         27,
			Date:               time.Date(2018, 4, 14, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (Every Year -> Every 6 months)",
			Amount:             money.Money{Amount: 14000, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     28,
			FillRecurrenceRule: 29,
		},

		{
			PrimaryKey:         28,
			Date:               time.Date(2018, 4, 6, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (Every 4 Weeks -> Every 3 weeks)",
			Amount:             money.Money{Amount: 6000, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     22,
			FillRecurrenceRule: 40,
		},

		{
			PrimaryKey:         29,
			Date:               time.Date(2018, 4, 8, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (Every Month; 15th & 31st -> Every month)",
			Amount:             money.Money{Amount: 8000, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     19,
			FillRecurrenceRule: 37,
		},

		{
			PrimaryKey:         30,
			Date:               time.Date(2018, 4, 10, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (Every Month; 1st & 16th -> Every month on the 1st and 15th)",
			Amount:             money.Money{Amount: 10000, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     21,
			FillRecurrenceRule: 31,
		},

		{
			PrimaryKey:         31,
			Date:               time.Date(2018, 4, 5, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (Every 3 Weeks -> Every 2 weeks)",
			Amount:             money.Money{Amount: 5000, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     25,
			FillRecurrenceRule: 32,
		},

		{
			PrimaryKey:         32,
			Date:               time.Date(2018, 4, 4, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (Every 2 Weeks -> Every week)",
			Amount:             money.Money{Amount: 4000, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     20,
			FillRecurrenceRule: 33,
		},

		{
			PrimaryKey:         33,
			Date:               time.Date(2018, 4, 9, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (Every Month; 1st & 15th -> Every month on the 15th and 31st)",
			Amount:             money.Money{Amount: 9000, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     18,
			FillRecurrenceRule: 39,
		},

		{
			PrimaryKey:         34,
			Date:               time.Date(2018, 4, 13, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (Every 6 Months -> Every 3 months)",
			Amount:             money.Money{Amount: 13000, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     26,
			FillRecurrenceRule: 38,
		},

		{
			PrimaryKey:         35,
			Date:               time.Date(2018, 4, 15, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (Every 2 Years -> Every year)",
			Amount:             money.Money{Amount: 15000, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     27,
			FillRecurrenceRule: 36,
		},

		{
			PrimaryKey:         36,
			Date:               time.Date(2018, 4, 29, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (Every day, 10 times)",
			Amount:             money.Money{Amount: 100000, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     42,
			FillRecurrenceRule: 0,
		},

		{
			PrimaryKey:         37,
			Date:               time.Date(2018, 4, 30, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (Every week, until June 1, 2018)",
			Amount:             money.Money{Amount: 200000, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     43,
			FillRecurrenceRule: 0,
		},

		{
			PrimaryKey:         38,
			Date:               time.Date(2018, 4, 28, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (S/M/W/F)",
			Amount:             money.Money{Amount: 50000, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     45,
			FillRecurrenceRule: 0,
		},

		{
			PrimaryKey:         39,
			Date:               time.Date(2018, 5, 1, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (T/T/S)",
			Amount:             money.Money{Amount: 9999, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     47,
			FillRecurrenceRule: 0,
		},

		{
			PrimaryKey:         40,
			Date:               time.Date(2020, 5, 2, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (2nd day)",
			Amount:             money.Money{Amount: 222, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     48,
			FillRecurrenceRule: 0,
		},

		{
			PrimaryKey:         41,
			Date:               time.Date(2020, 5, 3, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (3rd weekday)",
			Amount:             money.Money{Amount: 333, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     49,
			FillRecurrenceRule: 0,
		},

		{
			PrimaryKey:         42,
			Date:               time.Date(2020, 5, 4, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (Every 2 days)",
			Amount:             money.Money{Amount: 200, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     52,
			FillRecurrenceRule: 0,
		},

		{
			PrimaryKey:         43,
			Date:               time.Date(2020, 5, 5, 0, 0, 0, 0, time.UTC),
			Name:               "Groceries (1st weekend day)",
			Amount:             money.Money{Amount: 3100, Currency: "CAD"},
			Bucket:             13,
			RecurrenceRule:     55,
			FillRecurrenceRule: 0,
		},
	}

	assert.Equal(t, expectedSpendingPlan, spendingPlan)
}

func TestGetSpendingPlanRecurrenceRuleDescriptions(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	spendingPlan, err := api.GetSpendingPlan(database)
	assert.NoError(t, err)

	recurrenceRules, err := api.GetRecurrenceRulesMap(database)
	assert.NoError(t, err)

	expectedDescriptions := map[int64]string{
		21: "Never",
		22: "Every day",
		23: "Every week",
		24: "Every 3 months",
		25: "Every month",
		26: "Every 2 months",
		27: "Every year",
		28: "Every 4 weeks",
		29: "Every month on the 15th and 31st",
		30: "Every month on the 1st and 16th",
		31: "Every 3 weeks",
		32: "Every 2 weeks",
		33: "Every month on the 1st and 15th",
		34: "Every 6 months",
		35: "Every 2 years",
		36: "Every day, ending after 11 times",
		37: "Every week, ending on 2018-06-02",
		38: "Every week on Sunday, Monday, Wednesday and Friday",
		39: "Every week on Tuesday, Thursday and Saturday",
		40: "Every month on the 2nd day",
		41: "Every month on the 3rd week day",
		42: "Every 2 days",
		43: "Every month on the 1st weekend day",
	}

	descriptions := make(map[int64]string)
	for _, spendingPlanEvent := range spendingPlan {
		recurrenceRule := recurrenceRules[spendingPlanEvent.RecurrenceRule]

		descriptions[spendingPlanEvent.PrimaryKey] = api.DescribeRecurrenceRule(recurrenceRule)
	}

	assert.Equal(t, expectedDescriptions, descriptions)
}

func TestGetSpendingPlanFillRecurrenceRuleDescriptions(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	spendingPlan, err := api.GetSpendingPlan(database)
	assert.NoError(t, err)

	recurrenceRules, err := api.GetRecurrenceRulesMap(database)
	assert.NoError(t, err)

	expectedDescriptions := map[int64]string{
		21: "Every Event Date",
		22: "Every Event Date",
		23: "Every day",
		24: "Every 2 months",
		25: "Every 4 weeks",
		26: "Every month on the 1st and 16th",
		27: "Every 6 months",
		28: "Every 3 weeks",
		29: "Every month",
		30: "Every month on the 1st and 15th",
		31: "Every 2 weeks",
		32: "Every week",
		33: "Every month on the 15th and 31st",
		34: "Every 3 months",
		35: "Every year",
		36: "Every Event Date",
		37: "Every Event Date",
		38: "Every Event Date",
		39: "Every Event Date",
		40: "Every Event Date",
		41: "Every Event Date",
		42: "Every Event Date",
		43: "Every Event Date",
	}

	descriptions := make(map[int64]string)
	for _, spendingPlanEvent := range spendingPlan {
		recurrenceRule := recurrenceRules[spendingPlanEvent.FillRecurrenceRule]

		descriptions[spendingPlanEvent.PrimaryKey] = api.DescribeFillRecurrenceRule(recurrenceRule)
	}

	assert.Equal(t, expectedDescriptions, descriptions)
}
