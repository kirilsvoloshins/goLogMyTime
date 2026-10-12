package date

import (
	"time"
)

func DateToString(date1 time.Time) string {
	return date1.Format("2006-01-02")
}

func GetTimestamp() int64 {
	return time.Now().UnixMilli()
}

func MinutesToHoursMinutes(minutesSrc int) (int, int) {
	hours := int(minutesSrc / 60)
	minutes := minutesSrc - hours*60
	return hours, minutes
}
