package thaidate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const beOffset = 543

// ToBE converts a standard Gregorian time.Time to a Thai Buddhist Era (BE) string
// formatted as YYYY-MM-DD.
func ToBE(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	beYear := t.Year() + beOffset
	return fmt.Sprintf("%04d-%02d-%02d", beYear, t.Month(), t.Day())
}

// ToBERFC3339 converts a standard Gregorian time.Time to a Thai Buddhist Era (BE) string
// formatted as RFC3339.
func ToBERFC3339(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	beYear := t.Year() + beOffset
	// Format as RFC3339 but replace the year
	gregorianStr := t.Format(time.RFC3339)
	return fmt.Sprintf("%04d%s", beYear, gregorianStr[4:])
}

// ParseBE parses a Thai Buddhist Era (BE) string (YYYY-MM-DD or RFC3339)
// into a standard Gregorian time.Time.
func ParseBE(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Time{}, nil
	}

	parts := strings.SplitN(dateStr, "-", 2)
	if len(parts) < 2 {
		return time.Time{}, fmt.Errorf("invalid BE date format: %s", dateStr)
	}

	beYear, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year in BE date: %w", err)
	}

	gregorianYear := beYear - beOffset
	gregorianStr := fmt.Sprintf("%04d-%s", gregorianYear, parts[1])

	// Try parsing as RFC3339 first
	if t, err := time.Parse(time.RFC3339, gregorianStr); err == nil {
		return t, nil
	}

	// Try parsing as YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", gregorianStr); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("could not parse BE date string: %s", dateStr)
}
