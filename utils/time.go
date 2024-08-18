package utils

import "time"

// Use "02-01-2006" for dd-mm-yyyy
// Use Asia/Kolkata for ITC time
func FormatTimeByLocation(timeObj time.Time, location string, formatType string) string {
	loc, _ := time.LoadLocation(location)
	timeString := timeObj.In(loc).Format(formatType)
	return timeString
}
