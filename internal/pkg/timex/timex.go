package timex

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

type JsonTime time.Time

var (
	DatetimePattern = "2006-01-02 15:04:05"
	DatePattern     = "2006-01-02"
)

var (
	datetimeRegex, _ = regexp.Compile("[1-2][0-9][0-9][0-9]-[0-1]{0,1}[0-9]-[0-3]{0,1}[0-9]\\s+[0-2][0-9]:[0-5][0-9]:[0-5][0-9]")
	dateRegex, _     = regexp.Compile("[1-2][0-9][0-9][0-9]-[0-1]{0,1}[0-9]-[0-3]{0,1}[0-9]")
)

func (jt JsonTime) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf(`"%s"`, time.Time(jt).Format(DatetimePattern))

	return []byte(s), nil
}

func (jt *JsonTime) UnmarshalJSON(b []byte) error {
	if ret := datetimeRegex.Match(b); ret {
		return jt.parse(b, DatetimePattern)
	}

	if ret := dateRegex.Match(b); ret {
		return jt.parse(b, DatePattern)
	}

	return errors.New("时间格式未匹配")
}

func (jt *JsonTime) parse(b []byte, pattern string) error {
	data := string(b)
	if data == "" {
		return nil
	}

	parseTime, err := time.Parse(`"`+pattern+`"`, data)
	if err != nil {
		return err
	}

	*jt = JsonTime(parseTime)
	return nil
}

// Timestamp 返回自1970年以来的秒数
func (jt JsonTime) Timestamp() int64 {
	return time.Time(jt).Unix()
}

// GetFirstDateOfMonth 获取指定日期所在月第一天的零点
func GetFirstDateOfMonth(d time.Time) time.Time {
	t := d.AddDate(0, 0, -d.Day()+1)
	return getZeroTime(t)
}

// GetLastDateOfMonth 获取指定日志所在月的最后一天的零点
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

func getZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// Parse 将字符串转为time.Time类型
func Parse(str, pattern string) (time.Time, error) {
	return time.ParseInLocation(pattern, str, time.Local)
}

// Format 格式化时间
func Format(date time.Time, pattern string) string {
	return date.Format(pattern)
}
