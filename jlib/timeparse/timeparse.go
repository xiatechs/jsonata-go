package timeparse

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	strtime "github.com/ncruces/go-strftime"
)

// DateDim is the date dimension object returned from the timeparse function
type DateDim struct {
	// Other
	TimeZone       string `json:"TimeZone"`       // lite
	TimeZoneOffset string `json:"TimeZoneOffset"` // lite
	YearMonth      int    `json:"YearMonth"`      // int
	YearWeek       int    `json:"YearWeek"`       // int
	YearIsoWeek    int    `json:"YearIsoWeek"`    // int
	YearDay        int    `json:"YearDay"`        // int
	DateID         string `json:"DateId"`         // lite
	DateKey        int    `json:"DateKey"`        // lite
	DateTimeKey    int    `json:"DateTimeKey"`    // lite
	HourID         string `json:"HourId"`
	HourKey        int    `json:"HourKey"`
	Millis         int    `json:"Millis"`   // lite
	RawValue       string `json:"RawValue"` // lite

	// UTC
	UTC     string `json:"UTC"`     // lite
	DateUTC string `json:"DateUTC"` // lite
	HourUTC int    `json:"HourUTC"`

	// Local
	Local     string `json:"Local"`     // lite
	DateLocal string `json:"DateLocal"` // lite
	HourLocal int    `json:"HourLocal"`
}

// TimeDateDimensions generates a JSON object dependent on input source timestamp, input source format and input source timezone
// using golang time formats
func TimeDateDimensions(inputSrcTs, inputSrcFormat, inputSrcTz, requiredTz string) (*DateDim, error) {
	// first we create a time location based on the input source timezone location
	inputLocation, err := time.LoadLocation(inputSrcTz)
	if err != nil {
		return nil, err
	}

	// Since the source timestamp is implied to be in local time ("Europe/London"),
	// we parse it with the location set to Europe/London
	inputTime, err := time.ParseInLocation(inputSrcFormat, inputSrcTs, inputLocation)
	if err != nil {
		return nil, err
	}

	// then we create a time location based on the output timezone location
	outputLocation, err := time.LoadLocation(requiredTz)
	if err != nil {
		return nil, err
	}

	// here we translate te input time into a local time based on the output location
	localTime := inputTime.In(outputLocation)

	// convert the parsed time into a UTC time for UTC calculations
	utcTime := localTime.UTC()

	// now we have inputTime (the time parsed from the input location)
	// the local time which is the inputTime converted to the output location
	// and the UTC time which is the local time converted into UTC time

	// UTC TIME values
	utcAsYearMonthDay := utcTime.Format("2006-01-02")

	// Input time stamp TIME values (we confirmed there need to be a seperate set of UTC values)
	dateID := localTime.Format("20060102")

	// here we get the year and week for golang standard library ISOWEEK for the local time
	year, week := localTime.ISOWeek()

	yearDay, err := strconv.Atoi(localTime.Format("2006") + localTime.Format("002"))
	if err != nil {
		return nil, err
	}

	hourKeyStr := localTime.Format("2006010215")

	yearMondayWeek, _:= strconv.Atoi(strtime.Format(`%Y%W`, localTime))

	// the ISO yearweek
	yearIsoWeekInt, err := strconv.Atoi(fmt.Sprintf("%d%02d", year, week))
	if err != nil {
		return nil, err
	}

	// year month
	yearMonthInt, err := strconv.Atoi(localTime.Format("200601"))
	if err != nil {
		return nil, err
	}

	// the date key
	dateKeyInt, err := strconv.Atoi(dateID)
	if err != nil {
		return nil, err
	}

	// hack - golang time library doesn't handle millis correct unless you do the below
	dateTimeID := localTime.Format("20060102150405.000")

	dateTimeID = strings.ReplaceAll(dateTimeID, ".", "")

	dateTimeKeyInt, err := strconv.Atoi(dateTimeID)
	if err != nil {
		return nil, err
	}

	// the hours
	hourKeyInt, err := strconv.Atoi(hourKeyStr)
	if err != nil {
		return nil, err
	}

	localTimeStamp := localTime.Format("2006-01-02T15:04:05.000-07:00")
	offsetStr := localTime.Format("-07:00")
	// construct the date dimension structure
	dateDim := &DateDim{
		RawValue:       inputSrcTs,
		TimeZoneOffset: offsetStr,
		YearWeek:       yearMondayWeek,
		YearDay:        yearDay,
		YearIsoWeek:    yearIsoWeekInt,
		YearMonth:      yearMonthInt,
		Millis:         int(localTime.UnixMilli()),
		HourLocal:      localTime.Hour(),
		HourKey:        hourKeyInt,
		HourID:         "Hours_" + hourKeyStr,
		DateLocal:      localTime.Format("2006-01-02"),
		TimeZone:       localTime.Location().String(),
		Local:          localTimeStamp,
		DateKey:        dateKeyInt,
		DateTimeKey:    dateTimeKeyInt,
		DateID:         "Dates_" + dateID,
		DateUTC:        utcAsYearMonthDay,
		UTC:            utcTime.Format("2006-01-02T15:04:05.000Z"),
		HourUTC:        utcTime.Hour(),
	}

	return dateDim, nil
}

// WeekNumberZeroBasedImproved calculates the week number of the year for the given date
// considering weeks starting from Monday and numbering from 0.
func WeekNumberZeroBasedImproved(t time.Time) (int, error) {
	startOfYear := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	firstDayOfYear := startOfYear.Weekday()

	offset := (int(firstDayOfYear) + 6) % 7

	yearDay := t.YearDay()
	week := (yearDay - 1 - offset) / 7
	year := 0

	if week < 0 {
		dec31LastYear := time.Date(t.Year()-1, 12, 31, 0, 0, 0, 0, t.Location())
		return WeekNumberZeroBasedImproved(dec31LastYear)
	}

	if t.Month() == time.December && week >= 52 {
		jan1NextYear := time.Date(t.Year()+1, 1, 1, 0, 0, 0, 0, t.Location())
		if jan1NextYear.Weekday() <= time.Thursday {
			if yearDay+int(time.Thursday-jan1NextYear.Weekday()) > 365+(time.Date(t.Year(), 2, 29, 0, 0, 0, 0, t.Location()).YearDay()/366) {
				year++
				week = 0
				return strconv.Atoi(fmt.Sprintf("%04d%02d", year, week))
			}
		}
	}

	year = t.Year()
	return strconv.Atoi(fmt.Sprintf("%04d%02d", year, week))
}
