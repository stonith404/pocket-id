package utils

import "time"

func FormatDateForDb(time time.Time) string {
	const layout = "2006-01-02 15:04:05.000-07:00"
	return time.Format(layout)
}
