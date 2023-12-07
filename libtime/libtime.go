package libtime

import (
	"os"
	"time"
)

// Config -
type Config struct {
	Timezone *time.Location
}

var initialize = false
var cfg = Config{}

// loadConfig -
func loadConfig() {
	if !initialize {
		tz, err := time.LoadLocation(os.Getenv("app_timezone"))
		if err == nil {
			cfg.Timezone = tz
		} else {
			tz, _ = time.LoadLocation("Asia/Jakarta")
			cfg.Timezone = tz
		}

		initialize = true
	}
}

// Now - Return now based on environment app timezone, not recommend usage to store data. Use only for readable logging
func Now() time.Time {
	loadConfig()
	return time.Now().In(cfg.Timezone)
}

// NowToString - Return Now() to format yyyy-mm-dd hh:ii:ss
func NowToString() string {
	return Now().Format("2006-01-02 15:04:05")
}

// NullFormat - Return string or nil for sql.Nulltime
func NullFormat(input interface{}, tz string) interface{} {
	v, ok := input.(map[string]interface{})
	if ok {
		b, ok := v["Valid"].(bool)
		if ok && b {
			if s, ok := v["Time"].(string); ok {
				t, err := time.Parse(time.RFC3339, s)
				if err == nil {
					loc, _ := time.LoadLocation(tz)
					return t.In(loc).Format(time.RFC3339)
				}
			}
		}
	}
	return nil
}

// DateFormat - Return string or nil for sql.Nulltime
func DateFormat(input interface{}, tz string) interface{} {
	v, ok := input.(map[string]interface{})
	if ok {
		b, ok := v["Valid"].(bool)
		if ok && b {
			if s, ok := v["Time"].(string); ok {
				t, err := time.Parse(time.RFC3339, s)
				if err == nil {
					loc, _ := time.LoadLocation(tz)
					return t.In(loc).Format("2006-01-02")
				}
			}
		}
	}
	return nil
}

// DateTimeFormat - Return string or nil for sql.Nulltime
func DateTimeFormat(input interface{}, tz string) interface{} {
	v, ok := input.(map[string]interface{})
	if ok {
		b, ok := v["Valid"].(bool)
		if ok && b {
			if s, ok := v["Time"].(string); ok {
				t, err := time.Parse(time.RFC3339, s)
				if err == nil {
					loc, _ := time.LoadLocation(tz)
					return t.In(loc).Format("2006-01-02 15:04:05")
				}
			}
		}
	}
	return nil
}

// RFC3339Format - Return RFC3339 string or nil for sql.Nulltime
func RFC3339Format(input interface{}, tz string) interface{} {
	v, ok := input.(map[string]interface{})
	if ok {
		b, ok := v["Valid"].(bool)
		if ok && b {
			if s, ok := v["Time"].(string); ok {
				t, err := time.Parse(time.RFC3339, s)
				if err == nil {
					loc, _ := time.LoadLocation(tz)
					return t.In(loc).Format(time.RFC3339)
				}
			}
		}
	}
	return nil
}

// Based on https://github.com/bearbin/go-age

// AgeAt - gets the age of an entity at a certain time.
func AgeAt(birthDate time.Time, now time.Time) int {
	// Get the year number change since the user's birth.
	years := now.Year() - birthDate.Year()

	// If the date is before the date of birth, then not that many years have elapsed.
	birthDay := getAdjustedBirthDay(birthDate, now)
	if now.YearDay() < birthDay {
		years--
	}

	return years
}

// Age - Shorthand for AgeAt(birthDate, time.Now()), and carries the same usage and limitations.
func Age(birthDate time.Time) int {
	return AgeAt(birthDate, time.Now())
}

// getAdjustedBirthDay - Gets the adjusted date of birth to work around leap year differences.
func getAdjustedBirthDay(birthDate time.Time, now time.Time) int {
	birthDay := birthDate.YearDay()
	currentDay := now.YearDay()
	if isLeap(birthDate) && !isLeap(now) && birthDay >= 60 {
		return birthDay - 1
	}
	if isLeap(now) && !isLeap(birthDate) && currentDay >= 60 {
		return birthDay + 1
	}
	return birthDay
}

// isLeap - Works out if a time.Time is in a leap year.
func isLeap(date time.Time) bool {
	year := date.Year()
	if year%400 == 0 {
		return true
	} else if year%100 == 0 {
		return false
	} else if year%4 == 0 {
		return true
	}
	return false
}
