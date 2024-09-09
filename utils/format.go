package utils

import (
	"fmt"
	"strings"
	"time"
)

const DateTimeLayout1 = "2006-01-02T15:04:05Z07:00"       // Format with time and no milliseconds
const DateTimeLayout2 = "2006-01-02"                      // Date-only format
const DateTimeLayout3 = "2006-01-02T15:04:05.000Z07:00"   // Format with milliseconds
const DateTimeLayout4 = "2006-01-02T15:04:05.000Z"

func UnmarshalDateTime(src []byte, dst *time.Time) error {
	// Convert byte slice to string and trim spaces and quotes
	srcStr := strings.TrimSpace(string(src))
	srcStr = strings.Trim(srcStr, "\"")

	// Switch case based on string length and presence of milliseconds or timezone
	switch {
	case len(srcStr) == len(DateTimeLayout3):
		// Layout with milliseconds and timezone
		parsedTime, err := time.Parse(DateTimeLayout3, srcStr)
		if err != nil {
			return fmt.Errorf("unable to parse date-time (layout: %s): %w", DateTimeLayout3, err)
		}
		*dst = parsedTime
		return nil

	case len(srcStr) == len(DateTimeLayout4):
		// Layout with milliseconds but no timezone
		parsedTime, err := time.Parse(DateTimeLayout4, srcStr)
		if err != nil {
			return fmt.Errorf("unable to parse date-time (layout: %s): %w", DateTimeLayout4, err)
		}
		*dst = parsedTime
		return nil

	case len(srcStr) == len(DateTimeLayout1):
		// Layout without milliseconds, with timezone
		parsedTime, err := time.Parse(DateTimeLayout1, srcStr)
		if err != nil {
			return fmt.Errorf("unable to parse date-time (layout: %s): %w", DateTimeLayout1, err)
		}
		*dst = parsedTime
		return nil

	case len(srcStr) == len(DateTimeLayout2):
		// Date-only layout
		parsedTime, err := time.Parse(DateTimeLayout2, srcStr)
		if err != nil {
			return fmt.Errorf("unable to parse date-time (layout: %s): %w", DateTimeLayout2, err)
		}
		*dst = parsedTime
		return nil

	default:
		// If no known layout matches
		return fmt.Errorf("unrecognized date-time format: %s", srcStr)
	}
}