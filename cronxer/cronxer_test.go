package cronxer

import (
	//	"strings"
	"testing"
)

func TestCronParser_Parse(t *testing.T) {
	parser := New()

	cases := []struct {
		cronString   string
		expected     string
		expectedErr  bool
		errorMessage string
	}{
		{
			cronString: "* * * * * /usr/bin/command",
			expected: `Minute        *
Hour          *
Day of Month  *
Month         *
Day of Week   *
Command       /usr/bin/command
`,
			expectedErr:  false,
			errorMessage: "",
		},
		{
			cronString: "*/15 2-4,6-8 1,3,5,7,9 1-12 0-6 /usr/bin/command",
			expected: `Minute        00 15 30 45
Hour          02 03 04 06 07 08
Day of Month  01 03 05 07 09
Month         01 02 03 04 05 06 07 08 09 10 11 12
Day of Week   00 01 02 03 04 05 06
Command       /usr/bin/command
`,
			expectedErr:  false,
			errorMessage: "",
		},
		{
			cronString: "*/10 * * * *",
			expected: `Minute        00 10 20 30 40 50
Hour          *
Day of Month  *
Month         *
Day of Week   *
`,
			expectedErr:  false,
			errorMessage: "",
		},
		{
			cronString: "0 1 * * * /usr/bin/command",
			expected: `Minute        00
Hour          01
Day of Month  *
Month         *
Day of Week   *
Command       /usr/bin/command
`,
			expectedErr:  false,
			errorMessage: "",
		},
		{
			cronString:   "*/30 1-24 foo * *",
			expected:     "",
			expectedErr:  true,
			errorMessage: "string value in cron field",
		},
		{
			cronString:   "* * * *",
			expected:     "",
			expectedErr:  true,
			errorMessage: "incomplete cron string",
		},
		{
			cronString: "1-10/2 3-5/2 7-10/2 11-12/2 0-6/2 /usr/bin/command",
			expected: `Minute        01 03 05 07 09
Hour          03 05
Day of Month  07 09
Month         11
Day of Week   00 02 04 06
Command       /usr/bin/command
`,
			expectedErr:  false,
			errorMessage: "",
		},
		{
			cronString:   "*/30 1-24 * * *",
			expected:     "",
			expectedErr:  true,
			errorMessage: "invalid minute field",
		},
	}

	for _, c := range cases {
		actual, err := parser.Parse(c.cronString)

		if (err != nil) != c.expectedErr {
			t.Errorf("Test case failed for cron string '%s'. Expected error: %v, Actual error: %v", c.cronString, c.expectedErr, err)
		}

		if actual != c.expected {
			t.Errorf("Test case failed for cron string '%s'. Expected output:\n%s\nActual output:\n%s", c.cronString, c.expected, actual)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	cp := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := cp.Parse("1-10/2 3-5/2 7-10/2 11-12/2 0-6/2 /usr/bin/command")
		if err != nil {
			b.Errorf("parsing failed: %v", err)
		}
	}
}