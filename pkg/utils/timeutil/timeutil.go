package timeutil

import (
	"log"
	"time"
)

const ISTLocationName = "Asia/Kolkata"

var istLocation *time.Location

func init() {
	var err error
	istLocation, err = time.LoadLocation(ISTLocationName)
	if err != nil {
		log.Fatalf("failed to load IST time zone: %v", err)
	}
}

func ToIST(t time.Time) time.Time {
	return t.In(istLocation)
}

func NowIST() time.Time {
	return time.Now().In(istLocation)
}
