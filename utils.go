package main

import "time"
import "log"

// Parse the date
//
// Arguments:
//
// - date: "mm-dd-yyyy" or "mm-dd-yy" or "today" or "yesterday"
// - defaultt: the time to return by default (of parsing fails)
// - today:     the time to treat as "today" (if today or yesterday are used)
//
// Reutrns: `time.Time`
func parseDate(date string, defaultt time.Time, today time.Time) time.Time {
	if date == "" {
		return defaultt
	}
	if date == "today" {
		return today
	}
	if date == "yesterday" {
		der, _ := time.ParseDuration("-24h")
		return today.Add(der)
	}
	var err error
	var tm time.Time
	formats := []string{"01-02-2006", "01-02-06", "1-02-06", "01-2-06", "1-2-06", "1-2-2006", "01-2-2006", "1-02-2006"}
	for _, format := range formats {
		tm, err = time.Parse(format, date)
		if err == nil {
			return tm
		}
	}
	log.Printf("Warning: unrecognizable date. Defaulting to today: %q\n", date)
	return today
}
