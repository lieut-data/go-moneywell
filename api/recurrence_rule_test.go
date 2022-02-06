package api_test

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lieut-data/go-moneywell/api"
)

func TestRecurrenceRuleEqual(t *testing.T) {
	t.Run("self", func(t *testing.T) {
		var r *api.RecurrenceRule
		assert.True(t, r.Equals(r))
	})

	t.Run("nil", func(t *testing.T) {
		var r api.RecurrenceRule
		assert.False(t, r.Equals(nil))
	})

	t.Run("exactly equal", func(t *testing.T) {
		r1 := &api.RecurrenceRule{
			PrimaryKey:         1,
			EndDate:            time.Now(),
			FirstDayOfTheWeek:  1,
			OccurrenceCount:    2,
			RecurrenceInterval: 3,
			RecurrenceType:     4,
			DaysOfTheMonth:     []int64{1, 2, 3},
			DaysOfTheWeek:      []int64{4, 5, 6},
			MonthsOfTheYear:    []int64{7, 8, 9},
			OnThe: api.RecurrenceRuleOnThe{
				DayOfTheWeek: 1,
				WeekNumber:   10,
			},
			WeekdaysOfTheMonth: []int64{1, 3, 6},
		}

		r2 := &api.RecurrenceRule{
			PrimaryKey:         1,
			EndDate:            r1.EndDate,
			FirstDayOfTheWeek:  1,
			OccurrenceCount:    2,
			RecurrenceInterval: 3,
			RecurrenceType:     4,
			DaysOfTheMonth:     []int64{1, 2, 3},
			DaysOfTheWeek:      []int64{4, 5, 6},
			MonthsOfTheYear:    []int64{7, 8, 9},
			OnThe: api.RecurrenceRuleOnThe{
				DayOfTheWeek: 1,
				WeekNumber:   10,
			},
			WeekdaysOfTheMonth: []int64{1, 3, 6},
		}

		assert.True(t, r1.Equals(r2))
		assert.True(t, r2.Equals(r1))
	})

	t.Run("different primary keys, still equal", func(t *testing.T) {
		r1 := &api.RecurrenceRule{
			PrimaryKey:         1,
			EndDate:            time.Now(),
			FirstDayOfTheWeek:  1,
			OccurrenceCount:    2,
			RecurrenceInterval: 3,
			RecurrenceType:     4,
			DaysOfTheMonth:     []int64{1, 2, 3},
			DaysOfTheWeek:      []int64{4, 5, 6},
			MonthsOfTheYear:    []int64{7, 8, 9},
			OnThe: api.RecurrenceRuleOnThe{
				DayOfTheWeek: 1,
				WeekNumber:   10,
			},
			WeekdaysOfTheMonth: []int64{1, 3, 6},
		}

		r2 := &api.RecurrenceRule{
			PrimaryKey:         2,
			EndDate:            r1.EndDate,
			FirstDayOfTheWeek:  1,
			OccurrenceCount:    2,
			RecurrenceInterval: 3,
			RecurrenceType:     4,
			DaysOfTheMonth:     []int64{1, 2, 3},
			DaysOfTheWeek:      []int64{4, 5, 6},
			MonthsOfTheYear:    []int64{7, 8, 9},
			OnThe: api.RecurrenceRuleOnThe{
				DayOfTheWeek: 1,
				WeekNumber:   10,
			},
			WeekdaysOfTheMonth: []int64{1, 3, 6},
		}

		assert.True(t, r1.Equals(r2))
		assert.True(t, r2.Equals(r1))
	})

	t.Run("different end dates", func(t *testing.T) {
		r1 := &api.RecurrenceRule{
			PrimaryKey:         1,
			EndDate:            time.Now(),
			FirstDayOfTheWeek:  1,
			OccurrenceCount:    2,
			RecurrenceInterval: 3,
			RecurrenceType:     4,
			DaysOfTheMonth:     []int64{1, 2, 3},
			DaysOfTheWeek:      []int64{4, 5, 6},
			MonthsOfTheYear:    []int64{7, 8, 9},
			OnThe: api.RecurrenceRuleOnThe{
				DayOfTheWeek: 1,
				WeekNumber:   10,
			},
			WeekdaysOfTheMonth: []int64{1, 3, 6},
		}

		r2 := &api.RecurrenceRule{
			PrimaryKey:         2,
			EndDate:            r1.EndDate.Add(1 * time.Second),
			FirstDayOfTheWeek:  1,
			OccurrenceCount:    2,
			RecurrenceInterval: 3,
			RecurrenceType:     4,
			DaysOfTheMonth:     []int64{1, 2, 3},
			DaysOfTheWeek:      []int64{4, 5, 6},
			MonthsOfTheYear:    []int64{7, 8, 9},
			OnThe: api.RecurrenceRuleOnThe{
				DayOfTheWeek: 1,
				WeekNumber:   10,
			},
			WeekdaysOfTheMonth: []int64{1, 3, 6},
		}

		assert.False(t, r1.Equals(r2))
		assert.False(t, r2.Equals(r1))
	})
}

func TestGetRecurrenceRules(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	recurrenceRules, err := api.GetRecurrenceRules(database)
	assert.NoError(t, err)

	expectedRecurrenceRules := []api.RecurrenceRule{
		{
			PrimaryKey:         1,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
		},
		{
			PrimaryKey:         2,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeDaily,
		},
		{
			PrimaryKey:         3,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
			DaysOfTheMonth:     []int64{1, 16},
		},
		{
			PrimaryKey:         4,
			RecurrenceInterval: 3,
			RecurrenceType:     api.RecurrenceTypeWeekly,
		},
		{
			PrimaryKey:         5,
			RecurrenceInterval: 4,
			RecurrenceType:     api.RecurrenceTypeWeekly,
		},
		{
			PrimaryKey:         6,
			RecurrenceInterval: 3,
			RecurrenceType:     api.RecurrenceTypeMonthly,
		},
		{
			PrimaryKey:         7,
			RecurrenceInterval: 2,
			RecurrenceType:     api.RecurrenceTypeMonthly,
		},
		{
			PrimaryKey:         8,
			RecurrenceInterval: 6,
			RecurrenceType:     api.RecurrenceTypeMonthly,
		},
		{
			PrimaryKey:         9,
			RecurrenceInterval: 2,
			RecurrenceType:     api.RecurrenceTypeWeekly,
		},
		{
			PrimaryKey:         10,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeYearly,
		},
		{
			PrimaryKey:         11,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
			DaysOfTheMonth:     []int64{1, 15},
		},
		{
			PrimaryKey:         12,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
			DaysOfTheMonth:     []int64{15, 31},
		},
		{
			PrimaryKey:         13,
			RecurrenceInterval: 2,
			RecurrenceType:     api.RecurrenceTypeYearly,
		},
		{
			PrimaryKey:         14,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeWeekly,
		},
		{
			PrimaryKey:         15,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeDaily,
		},

		api.RecurrenceRule{
			PrimaryKey:         16,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeWeekly,
		},

		api.RecurrenceRule{
			PrimaryKey:         17,
			RecurrenceInterval: 3,
			RecurrenceType:     api.RecurrenceTypeMonthly,
		},

		api.RecurrenceRule{
			PrimaryKey:         18,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
			DaysOfTheMonth:     []int64{1, 15},
		},

		api.RecurrenceRule{
			PrimaryKey:         19,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
			DaysOfTheMonth:     []int64{15, 31},
		},

		api.RecurrenceRule{
			PrimaryKey:         20,
			RecurrenceInterval: 2,
			RecurrenceType:     api.RecurrenceTypeWeekly,
		},

		api.RecurrenceRule{
			PrimaryKey:         21,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
			DaysOfTheMonth:     []int64{1, 16},
		},

		api.RecurrenceRule{
			PrimaryKey:         22,
			RecurrenceInterval: 4,
			RecurrenceType:     api.RecurrenceTypeWeekly,
		},

		api.RecurrenceRule{
			PrimaryKey:         23,
			RecurrenceInterval: 2,
			RecurrenceType:     api.RecurrenceTypeMonthly,
		},

		api.RecurrenceRule{
			PrimaryKey:         24,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
		},

		api.RecurrenceRule{
			PrimaryKey:         25,
			RecurrenceInterval: 3,
			RecurrenceType:     api.RecurrenceTypeWeekly,
		},

		api.RecurrenceRule{
			PrimaryKey:         26,
			RecurrenceInterval: 6,
			RecurrenceType:     api.RecurrenceTypeMonthly,
		},

		api.RecurrenceRule{
			PrimaryKey:         27,
			RecurrenceInterval: 2,
			RecurrenceType:     api.RecurrenceTypeYearly,
		},

		api.RecurrenceRule{
			PrimaryKey:         28,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeYearly,
		},

		api.RecurrenceRule{
			PrimaryKey:         29,
			RecurrenceInterval: 6,
			RecurrenceType:     api.RecurrenceTypeMonthly,
		},

		api.RecurrenceRule{
			PrimaryKey:         30,
			RecurrenceInterval: 4,
			RecurrenceType:     api.RecurrenceTypeWeekly,
		},

		api.RecurrenceRule{
			PrimaryKey:         31,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
			DaysOfTheMonth:     []int64{1, 15},
		},

		api.RecurrenceRule{
			PrimaryKey:         32,
			RecurrenceInterval: 2,
			RecurrenceType:     api.RecurrenceTypeWeekly,
		},

		api.RecurrenceRule{
			PrimaryKey:         33,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeWeekly,
		},

		api.RecurrenceRule{
			PrimaryKey:         34,
			RecurrenceInterval: 2,
			RecurrenceType:     api.RecurrenceTypeMonthly,
		},

		api.RecurrenceRule{
			PrimaryKey:         35,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
			DaysOfTheMonth:     []int64{1, 16},
		},

		api.RecurrenceRule{
			PrimaryKey:         36,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeYearly,
		},

		api.RecurrenceRule{
			PrimaryKey:         37,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
		},

		api.RecurrenceRule{
			PrimaryKey:         38,
			RecurrenceInterval: 3,
			RecurrenceType:     api.RecurrenceTypeMonthly,
		},

		api.RecurrenceRule{
			PrimaryKey:         39,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
			DaysOfTheMonth:     []int64{15, 31},
		},

		api.RecurrenceRule{
			PrimaryKey:         40,
			RecurrenceInterval: 3,
			RecurrenceType:     api.RecurrenceTypeWeekly,
		},

		api.RecurrenceRule{
			PrimaryKey:         41,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeDaily,
		},

		api.RecurrenceRule{
			PrimaryKey:         42,
			OccurrenceCount:    11,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeDaily,
		},

		api.RecurrenceRule{
			PrimaryKey:         43,
			EndDate:            time.Date(2018, 6, 2, 0, 0, 0, 0, time.UTC),
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeWeekly,
		},

		api.RecurrenceRule{
			PrimaryKey:         44,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeWeekly,
			DaysOfTheWeek: []int64{
				api.DayOfTheWeekSunday,
				api.DayOfTheWeekMonday,
				api.DayOfTheWeekWednesday,
				api.DayOfTheWeekFriday,
			},
		},

		api.RecurrenceRule{
			PrimaryKey:         45,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeWeekly,
			DaysOfTheWeek: []int64{
				api.DayOfTheWeekSunday,
				api.DayOfTheWeekMonday,
				api.DayOfTheWeekWednesday,
				api.DayOfTheWeekFriday,
			},
		},

		api.RecurrenceRule{
			PrimaryKey:         46,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeWeekly,
			DaysOfTheWeek: []int64{
				api.DayOfTheWeekTuesday,
				api.DayOfTheWeekThursday,
				api.DayOfTheWeekSaturday,
			},
		},

		api.RecurrenceRule{
			PrimaryKey:         47,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeWeekly,
			DaysOfTheWeek: []int64{
				api.DayOfTheWeekTuesday,
				api.DayOfTheWeekThursday,
				api.DayOfTheWeekSaturday,
			},
		},

		api.RecurrenceRule{
			PrimaryKey:         48,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
			DaysOfTheMonth:     []int64{10},
			OnThe: api.RecurrenceRuleOnThe{
				DayOfTheWeek: -1,
				WeekNumber:   2,
			},
		},

		api.RecurrenceRule{
			PrimaryKey:         49,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
			DaysOfTheMonth:     []int64{10},
			OnThe: api.RecurrenceRuleOnThe{
				DayOfTheWeek: -2,
				WeekNumber:   3,
			},
		},

		api.RecurrenceRule{
			PrimaryKey:         50,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
			DaysOfTheMonth:     []int64{10},
			OnThe: api.RecurrenceRuleOnThe{
				DayOfTheWeek: -2,
				WeekNumber:   3,
			},
		},

		api.RecurrenceRule{
			PrimaryKey:         51,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
			DaysOfTheMonth:     []int64{10},
			OnThe: api.RecurrenceRuleOnThe{
				DayOfTheWeek: -1,
				WeekNumber:   2,
			},
		},

		api.RecurrenceRule{
			PrimaryKey:         52,
			RecurrenceInterval: 2,
			RecurrenceType:     api.RecurrenceTypeDaily,
		},

		api.RecurrenceRule{
			PrimaryKey:         53,
			RecurrenceInterval: 2,
			RecurrenceType:     api.RecurrenceTypeDaily,
		},

		api.RecurrenceRule{
			PrimaryKey:         54,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
			DaysOfTheMonth:     []int64{10},
			OnThe: api.RecurrenceRuleOnThe{
				DayOfTheWeek: -3,
				WeekNumber:   1,
			},
		},

		api.RecurrenceRule{
			PrimaryKey:         55,
			RecurrenceInterval: 1,
			RecurrenceType:     api.RecurrenceTypeMonthly,
			DaysOfTheMonth:     []int64{10},
			OnThe: api.RecurrenceRuleOnThe{
				DayOfTheWeek: -3,
				WeekNumber:   1,
			},
		},
	}

	assert.Equal(t, expectedRecurrenceRules, recurrenceRules)
}

func TestGetRecurrenceRulesMap(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	recurrenceRules, err := api.GetRecurrenceRules(database)
	assert.NoError(t, err)

	recurrenceRulesMap, err := api.GetRecurrenceRulesMap(database)
	assert.NoError(t, err)

	for primaryKey := range recurrenceRulesMap {
		found := false
		for _, recurrenceRule := range recurrenceRules {
			if primaryKey == recurrenceRule.PrimaryKey {
				assert.Equal(t, recurrenceRule, recurrenceRulesMap[primaryKey])
				found = true
			}
		}

		assert.True(t, found, "found recurrence rule in map not in slice")
	}
}

func TestDescribeRecurrenceRule(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	recurrenceRules, err := api.GetRecurrenceRulesMap(database)
	assert.NoError(t, err)

	expectedRecurrenceRuleDescriptions := map[int64]string{
		1:  "Every month",
		2:  "Every day",
		3:  "Every month on the 1st and 16th",
		4:  "Every 3 weeks",
		5:  "Every 4 weeks",
		6:  "Every 3 months",
		7:  "Every 2 months",
		8:  "Every 6 months",
		9:  "Every 2 weeks",
		10: "Every year",
		11: "Every month on the 1st and 15th",
		12: "Every month on the 15th and 31st",
		13: "Every 2 years",
		14: "Every week",
		15: "Every day",
		16: "Every week",
		17: "Every 3 months",
		18: "Every month on the 1st and 15th",
		19: "Every month on the 15th and 31st",
		20: "Every 2 weeks",
		21: "Every month on the 1st and 16th",
		22: "Every 4 weeks",
		23: "Every 2 months",
		24: "Every month",
		25: "Every 3 weeks",
		26: "Every 6 months",
		27: "Every 2 years",
		28: "Every year",
		29: "Every 6 months",
		30: "Every 4 weeks",
		31: "Every month on the 1st and 15th",
		32: "Every 2 weeks",
		33: "Every week",
		34: "Every 2 months",
		35: "Every month on the 1st and 16th",
		36: "Every year",
		37: "Every month",
		38: "Every 3 months",
		39: "Every month on the 15th and 31st",
		40: "Every 3 weeks",
		41: "Every day",
		42: "Every day, ending after 11 times",
		43: "Every week, ending on 2018-06-02",
		44: "Every week on Sunday, Monday, Wednesday and Friday",
		45: "Every week on Sunday, Monday, Wednesday and Friday",
		46: "Every week on Tuesday, Thursday and Saturday",
		47: "Every week on Tuesday, Thursday and Saturday",
	}

	for id, expectedDescription := range expectedRecurrenceRuleDescriptions {
		t.Run(fmt.Sprintf("rule %d", id), func(t *testing.T) {
			assert.Equal(t, expectedDescription, api.DescribeRecurrenceRule(recurrenceRules[id]))
		})
	}
}

func TestSortRecurrenceRules(t *testing.T) {
	t.Parallel()

	database, err := api.OpenDocument("Test.moneywell")
	assert.NoError(t, err)
	defer database.Close()

	expectedRecurrenceRuleDescriptions := []string{
		"Every day",
		"Every day, ending after 11 times",
		"Every 2 days",
		"Every week",
		"Every week, ending on 2018-06-02",
		"Every week on Sunday, Monday, Wednesday and Friday",
		"Every week on Tuesday, Thursday and Saturday",
		"Every 2 weeks",
		"Every 3 weeks",
		"Every 4 weeks",
		"Every month",
		"Every month on the 1st and 15th",
		"Every month on the 1st and 16th",
		"Every month on the 1st weekend day",
		"Every month on the 3rd week day",
		"Every month on the 2nd day",
		"Every month on the 15th and 31st",
		"Every 2 months",
		"Every 3 months",
		"Every 6 months",
		"Every year",
		"Every 2 years",
	}

	recurrenceRules, err := api.GetRecurrenceRules(database)
	assert.NoError(t, err)

	sort.Sort(api.RecurrenceRuleSort(recurrenceRules))

	sortedRecurrenceRuleDescriptions := make([]string, 0, len(recurrenceRules))
	for i, recurrenceRule := range recurrenceRules {
		if i > 0 && recurrenceRule.Equals(&recurrenceRules[i-1]) {
			continue
		}

		sortedRecurrenceRuleDescriptions = append(
			sortedRecurrenceRuleDescriptions,
			api.DescribeRecurrenceRule(recurrenceRule),
		)
	}

	assert.Equal(t, expectedRecurrenceRuleDescriptions, sortedRecurrenceRuleDescriptions)

}
