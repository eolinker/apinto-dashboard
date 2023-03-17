package common

import "time"

const formatTime = "2006-01-02 15:04:05"

func TimeToStr(t time.Time) string {
	return t.Format(formatTime)
}

func StrToTime(t string) time.Time {
	location, _ := time.ParseInLocation(formatTime, t, time.Local)
	return location
}
