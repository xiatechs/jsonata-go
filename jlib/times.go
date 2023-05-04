package jlib

import (
	"fmt"
	"encoding/json"
)

type DateTimeObj struct {
	DateId string
	TimestampUTC string
	DateUTC string
	TimeZone string
	DateTimeLocal string
	DateLocal string
	HourId string
	Millis string
	TimeZoneOffset string
	YearMonth string
	YearWeek string
}

// TimeFunc todo provide detailed comment
func TimeFunc(inputTimeStamp, rawTimeGrain, timezone string, isRawUTC, timezoneOffset bool) (interface{}, error) {
	timeObj := DateTimeObj{}

	bytes, _ := json.Marshal(timeObj)

	var output interface{}

	err := json.Unmarshal(bytes, &output)
	if err != nil {
		return output, fmt.Errorf("unescape json unmarshal error: %v", err)
	}

	return output, nil
}