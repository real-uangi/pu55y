// Package date @author uangi 2023-05
package date

import (
	"github.com/real-uangi/pu55y/config"
	"strconv"
	"time"
)

type Date int64

var (
	location *time.Location
	layout   *string
)

func CurrentDateString() string {
	return FormatDate(time.Now())
}

func FormatDate(t time.Time) string {
	return t.In(location).Format(*layout)
}

func New() Date {
	return Date(time.Now().UnixMilli())
}

func (d *Date) ToTime() time.Time {
	return time.UnixMilli(int64(*d))
}

func (d *Date) Format() string {
	return d.ToTime().In(location).Format(*layout)
}

func (d *Date) Mill() int64 {
	return int64(*d)
}

func (d *Date) MillString() string {
	return strconv.FormatInt(int64(*d), 10)
}

func Init() {
	conf := &config.GetConfig().Sys.Date
	l, err := time.LoadLocation(conf.Zone)
	if err != nil {
		panic(err)
	}
	location = l
	layout = &conf.Layout
}
