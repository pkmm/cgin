package util

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// JSONTime format json time field by myself
//type JSONTime struct {
//	time.Time
//}
type JSONTime time.Time

func (t *JSONTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"` + "2006-01-02 15:04:05" + `"`, string(data), time.Local)
	*t = JSONTime(now)
	return err
}

// MarshalJSON on JSONTime format Time field with Y-m-d H:i:s
func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	var tmp = time.Time(t)
	if tmp.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tmp, nil
}

// Scan value of time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
