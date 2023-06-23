package cronxer

import (
	"testing"
)

func TestValidateCronString(t *testing.T) {
	testCases := []struct {
		cronString string
		expected   bool
	}{
		// Valid cron strings
		{"*/15 * * * *", true},
		{"0 0 1 1 1", true},
		{"0 12 * * 1-5", true},
		{"* * * * *", true},
		{"0 0 * * 1,6", true},
		{"0 0 1 1 0", true},
		{"*/15 */2 1-10/2 * * ", true},

		// Invalid cron strings
		{"", false},                // Empty string
		{"* * * *", false},         // Incomplete cron string
		{"61 * * * *", false},      // Invalid minute value
		{"* 24 * * *", false},      // Invalid hour value
		{"* * 0 * *", false},       // Invalid day of month value
		{"* * * 13 *", false},      // Invalid month value
		{"* * * * 8", false},       // Invalid day of week value
		{"0 0 1 1 8", false},       // Invalid day of week value (8)
		{"0 0 0 0 1", false},       // Invalid day of month value (0)
		{"0 24 * * 1-5", false},    // Invalid hour value (24)
		{"* 0 1-32 * *", false},    // Invalid day of month range (32)
		{"* 0 1,15 * 1,13", false}, // Invalid month value (13)
		{"* 0 1,15 * 1-13", false}, // Invalid month range (13)
		{"* 0 1-15 * 1,15", false}, // Invalid day of week value (15)
		{"*/15 */2 1-10/2 * * ", true},
	}

	for _, tc := range testCases {
		err := validateCronString(tc.cronString)
		if (err == nil) != tc.expected {
			t.Errorf("Failed test case: %s with err: %s", tc.cronString, err)
		}
	}
}
