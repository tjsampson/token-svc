package datehelpers

import (
	"time"
)

var (
	dateLayouts = []string{"2006-1-2", "2006-01-02", time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC850, time.RFC822Z, time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano, time.Kitchen, time.Stamp, time.StampMilli, time.StampMicro, time.StampNano}
)

func parseDate(date string) (time.Time, error) {
	var result time.Time
	var err error
	for _, dateLayout := range dateLayouts {
		parsedDate, parseErr := time.Parse(dateLayout, date)
		if parseErr == nil {
			result = parsedDate
			break
		} else {
			err = parseErr
		}
	}
	// This will only return the last error...
	return result, err
}

// ConvertDateToMilli takes a date Time object and returs the milliseconds
func ConvertDateToMilli(date time.Time) int64 {
	return date.UnixNano() / int64(time.Millisecond)
}

// ConvertStringDateToMilli pareses a string date with format = "2006-01-02"
// and returns the millisecond Unix time
func ConvertStringDateToMilli(date string) int64 {
	t, err := parseDate(date)
	if err != nil {
		return 0
	}
	return ConvertDateToMilli(t)
}

// ParseDate parses the date
func ParseDate(date string) (time.Time, error) {
	return parseDate(date)
}
