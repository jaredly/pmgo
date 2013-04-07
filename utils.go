package main

import "time"
import "log"

type DateError struct {
    Value string
    Default time.Time
}

func (err *DateError) Error() string {
    return "Unable to parse " + err.Value + ". Using default date."
}

// Parse the date
//
// Arguments:
//
// - date: "mm-dd-yyyy" or "mm-dd-yy" or "today" or "yesterday"
// - defaultt: the time to return by default (of parsing fails)
// - today:     the time to treat as "today" (if today or yesterday are used)
//
// Returns: `time.Time`
//
// @tested
func parseDate(date string, defaultt time.Time, today time.Time) (time.Time, error) {
	if date == "" {
		return defaultt, nil
	}
	if date == "today" {
		return today, nil
	}
	if date == "yesterday" {
		der, _ := time.ParseDuration("-24h")
		return today.Add(der), nil
	}
	if date == "tomorrow" {
		der, _ := time.ParseDuration("24h")
		return today.Add(der), nil
	}
	var err error
	var tm time.Time
  // here are all the kinds of dates we accept
	formats := []string{"01-02-2006", "01-02-06", "1-02-06", "01-2-06",
                      "1-2-06", "1-2-2006", "01-2-2006", "1-02-2006"}
	for _, format := range formats {
		tm, err = time.Parse(format, date)
		if err == nil {
			return tm, nil
		}
	}
  var terr DateError
  terr.Value = date
  terr.Default = defaultt
	return today, terr
}

