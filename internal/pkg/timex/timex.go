package timex

import (
	"fmt"
	"time"
)

type JsonTime time.Time

var pattern = "2006-01-02 15:04:05"

func (jt JsonTime) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf(`"%s"`, time.Time(jt).Format(pattern))

	return []byte(s), nil
}

func (jt *JsonTime) UnmarshalJSON(b []byte) error {
	data := string(b)
	if data == "" {
		return nil
	}

	parseTime, err := time.Parse(`"` + pattern + `"`, data)
	if err != nil {
		return err
	}

	*jt = JsonTime(parseTime)
	return nil
}