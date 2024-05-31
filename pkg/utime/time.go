package utime

import "time"

func GetTomorrowZero() time.Time {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	tomorrowZero := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())

	return tomorrowZero
}
