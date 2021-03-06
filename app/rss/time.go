package rss

import (
	"strings"
	"time"
)

func parseTime(s string) (time.Time, error) {
	formats := []string{
		"Mon, _2 Jan 2006 15:04:05 MST",
		"Mon, _2 Jan 2006 15:04:05 Z",
		"Mon, _2 Jan 2006 15:04:05 UT",
		"Mon, _2 Jan 2006 15:04 MST",
		"Mon, _2 Jan 2006 15:04:05 -0700",
		"Mon, _2 Jan 2006 15:04:05",
		"_2 Jan 2006 15:04:05 -0700",
		"2006-01-02 15:04:05",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
	}

	s = strings.TrimSpace(s)

	var e error
	var t time.Time

	for _, format := range formats {
		t, e = time.Parse(format, s)
		if e == nil {
			return t, e
		}
	}

	return time.Time{}, e
}
