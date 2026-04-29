package wordstat

import "time"

// FromDateLastThreeMonths returns UTC midnight Monday on or before (now - 3 calendar months).
func FromDateLastThreeMonths(now time.Time) time.Time {
	anchor := now.UTC().AddDate(0, -3, 0)
	monday := truncateToMondayUTC(anchor)
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC)
}

func truncateToMondayUTC(d time.Time) time.Time {
	d = d.UTC()
	// Go: Sunday=0, Monday=1, ... Saturday=6
	// Days to subtract to reach Monday of same week:
	off := (int(d.Weekday()) + 6) % 7
	return d.AddDate(0, 0, -off)
}
