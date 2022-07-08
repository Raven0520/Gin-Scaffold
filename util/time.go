package util

import "time"

const (
	TZDateTimeFormat        = "2006-01-02T15:04:05Z"
	UTCDateTimeFormat       = "2006-01-02T15:04:05Z07:00"
	UTCDateTimeMinuteFormat = "2006-01-02T15:04Z07:00"
	DateTimeFormat          = "2006-01-02 15:04:05"
	DateFormat              = "2006-01-02"
	DateMinuteFormat        = "2006-01-02 15:04"
	HourFormat              = "15:04"
	DateTimeStringFormat    = "20060102150405"
)

// StartEndFromString 时间字符串 转换为时间
func StartEndFromString(s, e, format string) (start, end time.Time, err error) {
	start, err = time.Parse(format, s)
	if err != nil {
		return
	}
	end, err = time.Parse(format, e)
	if err != nil {
		return
	}
	return
}

// TimeToString 时间转换为字符串
func TimeToString(t time.Time, f string) string {
	return t.Format(f)
}

// TimeToStringInLocation 时间转换为字符串 带时区
func TimeToStringInLocation(t time.Time, location *time.Location, f string) string {
	return t.In(location).Format(f)
}
