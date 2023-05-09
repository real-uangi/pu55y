package date

import (
	"time"
)

const (
	Layout string = "2006-01-02 15:04:05"
	Area          = "Asia/Shanghai"
)

var (
	location, _ = time.LoadLocation(Area)
)

func CurrentDateString() string {
	return FormatDate(time.Now())
}

func FormatDate(t time.Time) string {
	return t.In(location).Format(Layout)
}
