package utils

import (
	"fmt"
	"time"
)

// Define the layout for marshaling/unmarshaling. Modify this based on your API's date-time format.
const DateTimeLayout = "2006-01-02T15:04:05Z07:00"

// MarshalDateTime converts a time.Time to a string formatted according to the GraphQL schema.
func MarshalDateTime(t *time.Time) ([]byte, error) {
	if t == nil {
		return nil, fmt.Errorf("time value is nil")
	}
	return []byte(t.Format(DateTimeLayout)), nil
}

// UnmarshalDateTime converts a string from the GraphQL response to a time.Time.
func UnmarshalDateTime(src []byte, dst *time.Time) error {
	// Convert byte slice to string
	srcStr := string(src)

	// Try parsing the string with time component
	parsedTime, err := time.Parse(DateTimeLayout, srcStr)
	if err == nil {
		*dst = parsedTime
		return nil
	}

	// If the string doesn't include a time, try parsing it as just a date
	parsedDateOnly, err := time.Parse("2006-01-02", srcStr)
	if err != nil {
		return fmt.Errorf("unable to parse date-time: %w", err)
	}

	*dst = parsedDateOnly
	return nil
}
