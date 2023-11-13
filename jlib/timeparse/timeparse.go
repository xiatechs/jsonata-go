package timeparse

import (
	"fmt"
	"strconv"
	"time"
)

// DateDim is the date dimension object returned from the timeparse function
type DateDim struct {
	DateID         string `json:"DateId"`
	DateKey        string `json:"DateKey"`
	UTC            string `json:"UTC"`
	DateUTC        string `json:"DateUTC"`
	Parsed         string `json:"Parsed"`
	Local          string `json:"Local"`
	DateLocal      string `json:"DateLocal"`
	HourID         string `json:"HourId"`
	HourKey        string `json:"HourKey"`
	Millis         string `json:"Millis"`
	Hour           string `json:"Hour"`
	TimeZone       string `json:"TimeZone"`
	TimeZoneOffset string `json:"TimeZoneOffset"`
	YearMonth      string `json:"YearMonth"`
	YearWeek       string `json:"YearWeek"`
	YearIsoWeek    string `json:"YearIsoWeek"`
	YearDay        string `json:"YearDay"`

	// TODO add UTC fields
}

// TimeDateDimensions generates a JSON object dependent on input source timestamp, input source format and input source timezone
// using golang time formats
func TimeDateDimensions(inputSrcTs, inputSrcFormat, inputSrcTz string) (interface{}, error) {
	parsedTime, err := getTimeWithLocation(inputSrcTs, inputSrcFormat, inputSrcTz)
	if err != nil {
		return nil, err
	}

	// convert the parsed time into a UTC time for UTC calculations
	parsedTimeUTC := parsedTime.UTC()

	// UTC TIME values

	utcAsYearMonthDay := parsedTimeUTC.Format("2006-01-02")

	// Input time stamp TIME values (we confirmed there need to be a seperate set of UTC values)
	dateID := parsedTime.Format("20060102")

	// Get the time zone offset
	_, offset := parsedTime.Zone()

	// LOCAL TIME VALUES

	// this is the "local time" with offset removed, and added to the time itself
	// i.2 2020-08-01+01:00 --> 2020-08-02
	localTime := parsedTimeUTC.Add(time.Duration(offset) * time.Second)

	year, week := localTime.ISOWeek()

	yearDay := localTime.Format("2006") + localTime.Format("002")

	offsetStr, err := getTimeOffsetString(localTime, parsedTime)
	if err != nil {
		return nil, err
	}

	mondayWeek := getWeekOfYearString(localTime)

	// construct the date dimension structure
	dateDim := DateDim{
		TimeZoneOffset: offsetStr,
		YearWeek:       mondayWeek,
		YearDay:        yearDay,
		YearIsoWeek:    fmt.Sprintf("%d%02d", year, week),
		YearMonth:      localTime.Format("200601"),
		Millis:         strconv.Itoa(int(localTime.UnixMilli())),
		Hour:           strconv.Itoa(localTime.Hour()),
		HourKey:        localTime.Format("2006010215"),
		HourID:         "Hours_" + localTime.Format("2006010215"),
		DateLocal:      localTime.Format("2006-01-02"),
		TimeZone:       parsedTime.Location().String(),
		Local:          parsedTime.Format("2006-01-02T15:04:05.000Z-07:00"),
		Parsed:         parsedTime.Format("2006-01-02T15:04:05.000Z"),
		DateKey:        dateID,
		DateID:         "Dates_" + dateID,
		DateUTC:        utcAsYearMonthDay,
		UTC:            parsedTimeUTC.Format("2006-01-02T15:04:05.000Z"),
	}

	return dateDim, nil
}

func parseDateTimeLocation(d string, layout string, location *time.Location) (time.Time, error) {
	date, err := time.Parse(layout, d)
	if err != nil {
		return date, err
	}

	return date.In(location), nil
}

func getTimeOffsetString(t1, t2 time.Time) (string, error) {
	duration := t1.Sub(t2)

	offsetString := formatOffset(duration)

	return offsetString, nil
}

func formatOffset(diff time.Duration) string {
	sign := "+"

	if diff < 0 {
		sign = "-"
		diff = -diff
	}

	hours := diff / time.Hour

	minutes := (diff % time.Hour) / time.Minute

	return fmt.Sprintf("%s%02d:%02d", sign, hours, minutes)
}

func getWeekOfYearString(date time.Time) string {
	_, week := date.ISOWeek()

	firstWednesday := date.AddDate(0, 0, -int(date.Weekday())+1)
	if firstWednesday.Weekday() != time.Wednesday {
		firstWednesday = firstWednesday.AddDate(0, 0, 7-int(firstWednesday.Weekday())+int(time.Wednesday))
	}

	if date.Weekday() == time.Sunday || date.Before(firstWednesday) {
		week--
	}

	return fmt.Sprintf("%04d%02d", date.Year(), week)
}
