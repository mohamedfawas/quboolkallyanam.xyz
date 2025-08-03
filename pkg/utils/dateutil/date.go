package dateutil

import "time"

func ParseDate(dateStr string) (*time.Time, error) {
	layout := "2006-01-02"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return nil, err
	}
	return &date, nil
}
