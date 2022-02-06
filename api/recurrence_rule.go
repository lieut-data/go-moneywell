package api

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
	"howett.net/plist"
)

// RecurrenceRule represents the conditions under which a spending plan event or its corresponding
// fill event are repeated.
//
// The MoneyWell SQLite schema for the ZRECURRENCERULE table is as follows:
// > .schema ZRECURRENCERULE
// CREATE TABLE ZRECURRENCERULE (
//    Z_PK INTEGER PRIMARY KEY,
//    Z_ENT INTEGER,
//    Z_OPT INTEGER,
//    ZENDDATEYMD INTEGER,
//    ZFIRSTDAYOFTHEWEEK INTEGER,
//    ZOCCURRENCECOUNT INTEGER,
//    ZRECURRENCEINTERVAL INTEGER,
//    ZRECURRENCETYPE INTEGER,
//    ZACTIVITY INTEGER,
//    Z5_ACTIVITY INTEGER,
//    ZEVENT INTEGER,
//    ZTICDSSYNCID VARCHAR,
//    ZUNIQUEID VARCHAR,
//    ZDAYSOFTHEMONTH BLOB,
//    ZDAYSOFTHEWEEK BLOB,
//    ZMONTHSOFTHEYEAR BLOB,
//    ZNTHWEEKDAYSOFTHEMONTH BLOB
// );
type RecurrenceRuleOnThe struct {
	DayOfTheWeek int64
	WeekNumber   int64
}

type RecurrenceRule struct {
	PrimaryKey         int64
	EndDate            time.Time
	FirstDayOfTheWeek  int64
	OccurrenceCount    int64
	RecurrenceInterval int64
	RecurrenceType     int64
	DaysOfTheMonth     []int64
	DaysOfTheWeek      []int64
	MonthsOfTheYear    []int64
	OnThe              RecurrenceRuleOnThe
	WeekdaysOfTheMonth []int64
}

func (r *RecurrenceRule) Equals(other *RecurrenceRule) bool {
	if r == other {
		return true
	}
	if r == nil || other == nil {
		return false
	}

	return r.EndDate == other.EndDate &&
		r.FirstDayOfTheWeek == other.FirstDayOfTheWeek &&
		r.OccurrenceCount == other.OccurrenceCount &&
		r.RecurrenceInterval == other.RecurrenceInterval &&
		r.RecurrenceType == other.RecurrenceType &&
		reflect.DeepEqual(r.DaysOfTheMonth, other.DaysOfTheMonth) &&
		reflect.DeepEqual(r.DaysOfTheWeek, other.DaysOfTheWeek) &&
		reflect.DeepEqual(r.MonthsOfTheYear, other.MonthsOfTheYear) &&
		r.OnThe.DayOfTheWeek == other.OnThe.DayOfTheWeek &&
		r.OnThe.WeekNumber == other.OnThe.WeekNumber &&
		reflect.DeepEqual(r.WeekdaysOfTheMonth, other.WeekdaysOfTheMonth)
}

const (
	RecurrenceTypeDaily   = 0
	RecurrenceTypeWeekly  = 1
	RecurrenceTypeMonthly = 2
	RecurrenceTypeYearly  = 3

	WeekNumberNone   = 0
	WeekNumberFirst  = 1
	WeekNumberSecond = 2
	WeekNumberThird  = 3
	WeekNumberFourth = 4
	WeekNumberLast   = -1

	DayOfTheWeekNone       = 0
	DayOfTheWeekSunday     = 1
	DayOfTheWeekMonday     = 2
	DayOfTheWeekTuesday    = 3
	DayOfTheWeekWednesday  = 4
	DayOfTheWeekThursday   = 5
	DayOfTheWeekFriday     = 6
	DayOfTheWeekSaturday   = 7
	DayOfTheWeekDay        = -1
	DayOfTheWeekWeekday    = -2
	DayOfTheWeekWeekendday = -3
)

func parsePlistAsInts(plistStr *string) ([]int64, interface{}, error) {
	if plistStr == nil {
		return nil, nil, nil
	}

	var ints []int64

	decoder := plist.NewDecoder(strings.NewReader(*plistStr))
	var target interface{}
	err := decoder.Decode(&target)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to decode plist")
	}

	if targetMap, ok := target.(map[string]interface{}); !ok {
		return nil, nil, errors.New("failed to decode plist as map")
	} else if objects, ok := targetMap["$objects"]; !ok {
		return nil, nil, errors.New("failed to decode plist as map, accessing $objects")
	} else if objectsArray, ok := objects.([]interface{}); !ok {
		return nil, nil, errors.New("failed to decode plist as map, accessing $objects as array")
	} else {
		for _, object := range objectsArray {
			if dayOfTheWeek, ok := object.(uint64); ok {
				ints = append(ints, int64(dayOfTheWeek))
			}
		}
	}

	return ints, target, nil
}

func parsePlistAsWeekdaysOfTheMonth(plistStr *string) (int64, int64, interface{}, error) {
	if plistStr == nil {
		return 0, 0, nil, nil
	}

	decoder := plist.NewDecoder(strings.NewReader(*plistStr))
	var target interface{}
	err := decoder.Decode(&target)
	if err != nil {
		return 0, 0, nil, errors.Wrap(err, "failed to decode plist")
	}

	if targetMap, ok := target.(map[string]interface{}); !ok {
		return 0, 0, target, errors.New("failed to decode plist as map")
	} else if objects, ok := targetMap["$objects"]; !ok {
		return 0, 0, target, errors.New("failed to decode plist as map, accessing $objects")
	} else if objectsArray, ok := objects.([]interface{}); !ok {
		return 0, 0, target, errors.New("failed to decode plist as map, accessing $objects as array")
	} else {
		for _, object := range objectsArray {
			if objectMap, ok := object.(map[string]interface{}); ok {
				if cls := objectMap["$class"]; cls.(plist.UID) == 2 {
					var dayOfTheWeek, weekNumber int64

					if dayOfTheWeekUint64, ok := objectMap["dayOfTheWeek"].(uint64); ok {
						dayOfTheWeek = int64(dayOfTheWeekUint64)
					} else if dayOfTheWeekInt64, ok := objectMap["dayOfTheWeek"].(int64); ok {
						dayOfTheWeek = dayOfTheWeekInt64
					}
					if weekNumberUint64, ok := objectMap["weekNumber"].(uint64); ok {
						weekNumber = int64(weekNumberUint64)
					} else if weekNumberInt64, ok := objectMap["weekNumber"].(int64); ok {
						weekNumber = weekNumberInt64
					}

					return dayOfTheWeek, weekNumber, target, nil
				}
			}
		}
	}

	return 0, 0, target, nil
}

// GetRecurrenceRules fetches the set of recurrence rules in a MoneyWell document.
func GetRecurrenceRules(database *sql.DB) ([]RecurrenceRule, error) {
	rows, err := database.Query(`
            SELECT 
		zr.Z_PK,
		zr.ZENDDATEYMD,
		zr.ZFIRSTDAYOFTHEWEEK,
		zr.ZOCCURRENCECOUNT,
		zr.ZRECURRENCEINTERVAL,
		zr.ZRECURRENCETYPE,
		zr.ZDAYSOFTHEMONTH,
		zr.ZDAYSOFTHEWEEK,
		zr.ZMONTHSOFTHEYEAR,
		zr.ZNTHWEEKDAYSOFTHEMONTH
            FROM 
                ZRECURRENCERULE zr
	    ORDER BY
		zr.Z_PK ASC
        `)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query recurrence rules")
	}
	defer rows.Close()

	recurrenceRules := []RecurrenceRule{}

	var primaryKey, firstDayOfTheWeek, occurrenceCount, recurrenceInterval, recurrenceType int64
	var daysOfTheMonth, daysOfTheWeek, monthsOfTheYear, weekdaysOfTheMonth *string

	var endDateymd sql.NullInt64

	for rows.Next() {
		err := rows.Scan(
			&primaryKey,
			&endDateymd,
			&firstDayOfTheWeek,
			&occurrenceCount,
			&recurrenceInterval,
			&recurrenceType,
			&daysOfTheMonth,
			&daysOfTheWeek,
			&monthsOfTheYear,
			&weekdaysOfTheMonth,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan recurrence rule")
		}

		endDate, err := parseDateymd(int(endDateymd.Int64))
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse recurrence rule end date")
		}

		recurrenceRule := RecurrenceRule{
			PrimaryKey:         primaryKey,
			EndDate:            endDate,
			FirstDayOfTheWeek:  firstDayOfTheWeek,
			OccurrenceCount:    occurrenceCount,
			RecurrenceInterval: recurrenceInterval,
			RecurrenceType:     recurrenceType,
		}

		recurrenceRule.DaysOfTheMonth, _, err = parsePlistAsInts(daysOfTheMonth)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode days of the month")
		}

		recurrenceRule.DaysOfTheWeek, _, err = parsePlistAsInts(daysOfTheWeek)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode days of the week")
		}

		recurrenceRule.MonthsOfTheYear, _, err = parsePlistAsInts(monthsOfTheYear)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode months of the year")
		}

		recurrenceRule.OnThe.DayOfTheWeek, recurrenceRule.OnThe.WeekNumber, _, err = parsePlistAsWeekdaysOfTheMonth(weekdaysOfTheMonth)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode weekdays of the month")
		}

		recurrenceRules = append(recurrenceRules, recurrenceRule)
	}

	return recurrenceRules, nil
}

// GetRecurrenceRulesMap gets a map from the bucket primary key to the bucket.
func GetRecurrenceRulesMap(database *sql.DB) (map[int64]RecurrenceRule, error) {
	recurrenceRules, err := GetRecurrenceRules(database)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	recurrenceRulesMap := make(map[int64]RecurrenceRule, len(recurrenceRules))
	for _, recurrenceRule := range recurrenceRules {
		recurrenceRulesMap[recurrenceRule.PrimaryKey] = recurrenceRule
	}

	return recurrenceRulesMap, nil
}

func describeDayOfTheWeek(dayOfTheWeek int64) string {
	switch dayOfTheWeek {
	case 1:
		return "Sunday"
	case 2:
		return "Monday"
	case 3:
		return "Tuesday"
	case 4:
		return "Wednesday"
	case 5:
		return "Thursday"
	case 6:
		return "Friday"
	case 7:
		return "Saturday"
	default:
		return "Unknown"
	}
}

func describeNth(day int64) string {
	switch day % 10 {
	case 1:
		return fmt.Sprintf("%dst", day)
	case 2:
		return fmt.Sprintf("%dnd", day)
	case 3:
		return fmt.Sprintf("%drd", day)
	default:
		return fmt.Sprintf("%dth", day)
	}
}

func describeMonth(month int64) string {
	switch month {
	case 1:
		return "January"
	case 2:
		return "February"
	case 3:
		return "March"
	case 4:
		return "April"
	case 5:
		return "May"
	case 6:
		return "June"
	case 7:
		return "July"
	case 8:
		return "August"
	case 9:
		return "September"
	case 10:
		return "October"
	case 11:
		return "November"
	case 12:
		return "December"
	default:
		return "Unknown"
	}
}

func joinWords(words []string) string {
	s := ""
	for i := range words {
		if i == 0 {
			s = words[i]
		} else if i == len(words)-1 {
			s = fmt.Sprintf("%s and %s", s, words[i])
		} else {
			s = fmt.Sprintf("%s, %s", s, words[i])
		}
	}

	return s
}

func DescribeRecurrenceRule(recurrenceRule RecurrenceRule) string {
	s := "Unknown"

	switch recurrenceRule.RecurrenceType {
	case RecurrenceTypeDaily:
		if recurrenceRule.RecurrenceInterval == 0 {
			s = "Never"
		} else if recurrenceRule.RecurrenceInterval == 1 {
			s = "Every day"
		} else {
			s = fmt.Sprintf("Every %d days", recurrenceRule.RecurrenceInterval)
		}
	case RecurrenceTypeWeekly:
		if recurrenceRule.RecurrenceInterval == 1 {
			s = "Every week"
		} else {
			s = fmt.Sprintf("Every %d weeks", recurrenceRule.RecurrenceInterval)
		}

		if len(recurrenceRule.DaysOfTheWeek) > 0 {
			var daysOfTheWeek []string
			for _, dayOfTheWeek := range recurrenceRule.DaysOfTheWeek {
				daysOfTheWeek = append(daysOfTheWeek, describeDayOfTheWeek(dayOfTheWeek))
			}
			s = fmt.Sprintf("%s on %s", s, joinWords(daysOfTheWeek))
		}
	case RecurrenceTypeMonthly:
		if recurrenceRule.RecurrenceInterval == 1 {
			s = "Every month"
		} else {
			s = fmt.Sprintf("Every %d months", recurrenceRule.RecurrenceInterval)
		}

		switch recurrenceRule.OnThe.WeekNumber {
		case WeekNumberFirst:
			fallthrough
		case WeekNumberSecond:
			fallthrough
		case WeekNumberThird:
			fallthrough
		case WeekNumberFourth:
			s = fmt.Sprintf("%s on the %s", s, describeNth(recurrenceRule.OnThe.WeekNumber))
		case WeekNumberLast:
			s = fmt.Sprintf("%s on the last", s)
		case WeekNumberNone:
		default:
		}

		switch recurrenceRule.OnThe.DayOfTheWeek {
		case DayOfTheWeekDay:
			s = fmt.Sprintf("%s day", s)
		case DayOfTheWeekWeekday:
			s = fmt.Sprintf("%s week day", s)
		case DayOfTheWeekWeekendday:
			s = fmt.Sprintf("%s weekend day", s)
		case DayOfTheWeekSunday:
			fallthrough
		case DayOfTheWeekMonday:
			fallthrough
		case DayOfTheWeekTuesday:
			fallthrough
		case DayOfTheWeekWednesday:
			fallthrough
		case DayOfTheWeekThursday:
			fallthrough
		case DayOfTheWeekFriday:
			fallthrough
		case DayOfTheWeekSaturday:
			s = fmt.Sprintf("%s %s", s, describeDayOfTheWeek(recurrenceRule.OnThe.DayOfTheWeek))
		case DayOfTheWeekNone:
		default:
		}

		if recurrenceRule.OnThe.DayOfTheWeek == DayOfTheWeekNone && len(recurrenceRule.DaysOfTheMonth) > 0 {
			var daysOfTheMonth []string
			for _, dayOfTheMonth := range recurrenceRule.DaysOfTheMonth {
				daysOfTheMonth = append(daysOfTheMonth, describeNth(dayOfTheMonth))
			}
			s = fmt.Sprintf("%s on the %s", s, joinWords(daysOfTheMonth))
		}
	case RecurrenceTypeYearly:
		if recurrenceRule.RecurrenceInterval == 1 {
			s = "Every year"
		} else {
			s = fmt.Sprintf("Every %d years", recurrenceRule.RecurrenceInterval)
		}

		switch recurrenceRule.OnThe.WeekNumber {
		case WeekNumberFirst:
			s = fmt.Sprintf("%s on the 1st", s)
		case WeekNumberSecond:
			s = fmt.Sprintf("%s on the 2nd", s)
		case WeekNumberThird:
			s = fmt.Sprintf("%s on the 3rd", s)
		case WeekNumberFourth:
			s = fmt.Sprintf("%s on the 4th", s)
		case WeekNumberLast:
			s = fmt.Sprintf("%s on the last", s)
		case WeekNumberNone:
		default:
		}

		switch recurrenceRule.OnThe.DayOfTheWeek {
		case DayOfTheWeekDay:
			s = fmt.Sprintf("%s day", s)
		case DayOfTheWeekWeekday:
			s = fmt.Sprintf("%s week day", s)
		case DayOfTheWeekWeekendday:
			s = fmt.Sprintf("%s weekend day", s)
		case DayOfTheWeekSunday, DayOfTheWeekMonday, DayOfTheWeekTuesday, DayOfTheWeekWednesday, DayOfTheWeekThursday, DayOfTheWeekFriday, DayOfTheWeekSaturday:
			s = fmt.Sprintf("%s %s", s, describeDayOfTheWeek(recurrenceRule.OnThe.DayOfTheWeek))
		case DayOfTheWeekNone:
		default:
		}

		if len(recurrenceRule.MonthsOfTheYear) > 0 {
			if recurrenceRule.OnThe.DayOfTheWeek == DayOfTheWeekNone {
				s = fmt.Sprintf("%s in", s)
			} else {
				s = fmt.Sprintf("%s of", s)
			}

			var monthsOfTheYear []string
			for _, monthOfTheYear := range recurrenceRule.MonthsOfTheYear {
				monthsOfTheYear = append(monthsOfTheYear, describeMonth(monthOfTheYear))
			}
			s = fmt.Sprintf("%s %s", s, joinWords(monthsOfTheYear))
		}
	}

	if recurrenceRule.OccurrenceCount > 0 {
		plural := ""
		if recurrenceRule.OccurrenceCount > 1 {
			plural = "s"
		}

		s = fmt.Sprintf("%s, ending after %d time%s", s, recurrenceRule.OccurrenceCount, plural)
	} else if !recurrenceRule.EndDate.IsZero() {
		s = fmt.Sprintf("%s, ending on %s", s, recurrenceRule.EndDate.Format("2006-01-02"))
	}

	return s
}

func DescribeFillRecurrenceRule(recurrenceRule RecurrenceRule) string {
	if recurrenceRule.RecurrenceType == RecurrenceTypeDaily && recurrenceRule.RecurrenceInterval == 0 {
		return "Every Event Date"
	}

	return DescribeRecurrenceRule(recurrenceRule)
}

func bToP(b bool) *bool {
	return &b
}

func cmpArrays(a, b []int64) *bool {
	l := len(a)
	if len(b) > l {
		l = len(b)
	}

	for k := 0; k < l; k++ {
		if k >= len(a) {
			return bToP(true)
		}
		if k >= len(b) {
			return bToP(false)
		}
		if a[k] != b[k] {
			return bToP(a[k] < b[k])
		}
	}

	return nil
}

type RecurrenceRuleSort []RecurrenceRule

func (s RecurrenceRuleSort) Len() int { return len(s) }
func (s RecurrenceRuleSort) Less(i, j int) bool {
	if s[i].RecurrenceType != s[j].RecurrenceType {
		return s[i].RecurrenceType < s[j].RecurrenceType
	}
	if s[i].RecurrenceInterval != s[j].RecurrenceInterval {
		return s[i].RecurrenceInterval < s[j].RecurrenceInterval
	}

	switch s[i].RecurrenceType {
	case RecurrenceTypeDaily:
	case RecurrenceTypeWeekly:
		cmpDaysOfTheWeek := cmpArrays(s[i].DaysOfTheWeek, s[j].DaysOfTheWeek)
		if cmpDaysOfTheWeek != nil {
			return *cmpDaysOfTheWeek
		}

	case RecurrenceTypeMonthly:
		cmpDaysOfTheMonth := cmpArrays(s[i].DaysOfTheMonth, s[j].DaysOfTheMonth)
		if cmpDaysOfTheMonth != nil {
			return *cmpDaysOfTheMonth
		}

		if s[i].OnThe.DayOfTheWeek != s[j].OnThe.DayOfTheWeek {
			return s[i].OnThe.DayOfTheWeek < s[j].OnThe.DayOfTheWeek
		}

	case RecurrenceTypeYearly:
		cmpMonthsOfTheYear := cmpArrays(s[i].MonthsOfTheYear, s[j].MonthsOfTheYear)
		if cmpMonthsOfTheYear != nil {
			return *cmpMonthsOfTheYear
		}
	}

	if s[i].OccurrenceCount != s[j].OccurrenceCount {
		return s[i].OccurrenceCount < s[j].OccurrenceCount
	}

	if s[i].EndDate != s[j].EndDate {
		return s[i].EndDate.UnixNano() < s[j].EndDate.UnixNano()
	}

	return false
}
func (s RecurrenceRuleSort) Swap(i, j int) { s[j], s[i] = s[i], s[j] }
