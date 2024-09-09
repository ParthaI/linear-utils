package utils

import (
	"fmt"
	"time"
)

// Define the layouts for marshaling/unmarshaling
const DateTimeLayout1 = "2006-01-02T15:04:05Z07:00"       // Format with time and no milliseconds
const DateTimeLayout2 = "2006-01-02"                      // Date-only format
const DateTimeLayout3 = "2006-01-02T15:04:05.000Z07:00"   // Format with milliseconds and timezone
const DateTimeLayout4 = "2006-01-02T15:04:05.000Z"        // Format with milliseconds but no timezone

// UnmarshalDateTime converts a string from the GraphQL response to a time.Time.
func UnmarshalDateTime(src []byte, dst *time.Time) error {
	// Convert byte slice to string
	srcStr := string(src)

	// Define a list of possible layouts to try
	layouts := []string{
		DateTimeLayout1,
		DateTimeLayout2,
		DateTimeLayout3,
		DateTimeLayout4,
	}

	// Try each layout in sequence
	for _, layout := range layouts {
		parsedTime, err := time.Parse(`"` + layout + `"`, srcStr)
		if err == nil {
			*dst = parsedTime
			return nil
		}
	}

	// If none of the layouts succeeded, return an error
	return fmt.Errorf("unable to parse date-time: no matching layout for input '%s'", srcStr)
}