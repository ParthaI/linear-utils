package utils

import (
	"fmt"
	"time"
)

// Define the layout for marshaling/unmarshaling. Modify this based on your API's date-time format.
// const DateTimeLayout1 = "2006-01-02T15:04:05Z" // 2023-05-10T19:16:33.208Z
const DateTimeLayout2 = "2006-01-02"
const DateTimeLayout3 = "2006-01-02T15:04:05.000Z"

// 2023-05-10T19:16:33.208Z

// MarshalDateTime converts a time.Time to a string formatted according to the GraphQL schema.
func MarshalDateTime(t *time.Time) ([]byte, error) {
	if t == nil {
		return nil, fmt.Errorf("time value is nil")
	}
	return []byte(t.Format(DateTimeLayout2)), nil
}

// UnmarshalDateTime converts a string from the GraphQL response to a time.Time.
func UnmarshalDateTime(src []byte, dst *time.Time) error {
	// Trim quotes if the string is quoted
	// src = bytes.Trim(src, "\\\"")

	// Convert byte slice to string
	// srcStr := strings.TrimSpace(string(src))
	// srcStr = strings.Trim(srcStr, "\"")
	// if srcStr == "" {
	// 	return nil
	// }

	// Define a list of possible layouts to try
	// layouts := []string{
	// 	DateTimeLayout3,
	// }

	// Try each layout in sequence
	// time.Parse(`"`+dateFormat+`"`, string(b))
	var parsedTime2, parsedTime3 time.Time
	var parseErr2, parseErr3 error
	parsedTime2, parseErr2 = time.Parse(`"`+DateTimeLayout2+`"`, string(src))
	if parseErr2 != nil {
		parsedTime3, parseErr3 = time.Parse(`"`+DateTimeLayout3+`"`, string(src))
	}

	// if parseErr3 != nil && parseErr2 != nil {
	// 	if parseErr3 != nil {
	// 		return parseErr3
	// 	}
	// 	if parseErr2 != nil {
	// 		return parseErr2
	// 	}
	// }
	if parseErr3 == nil {
		*dst = parsedTime3
	}
	if parseErr2 == nil {
		*dst = parsedTime2
	}
	// var err []error
	// for _, layout := range layouts {
	// 	parsedTime, parseErr := time.Parse(layout, srcStr)
	// 	if parseErr == nil {
	// 		*dst = parsedTime
	// 		return nil
	// 	}
	// 	err = append(err, parseErr)
	// }

	// // If none of the layouts succeeded, return the last error
	// if len(err) > 2 {
	// 	return err[0]
	// }

	return nil
}
