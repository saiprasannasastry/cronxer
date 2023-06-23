package cronxer

import "testing"

func TestCronParser_Parse(t *testing.T) {
	cronParser := New()

	// 	// Valid cron string
	// 	cronString := "* * * * * /usr/bin/command"
	// 	expectedOutput := `Minute         *
	// Hour           *
	// Day of Month   *
	// Month          *
	// Day of Week    *
	// Command        /usr/bin/command
	// `
	// 	output, err := cronParser.Parse(cronString)
	// 	if err != nil {
	// 		t.Errorf("Error parsing valid cron string: %v", err)
	// 	}
	// 	if output != expectedOutput {
	// 		t.Errorf("Invalid output for valid cron string\ngot:\n%s\nexpected:\n%s", output, expectedOutput)
	// 	}

	// 	// Invalid cron string
	// 	cronString = "invalid cron string"
	// 	_, err = cronParser.Parse(cronString)
	// 	if err == nil {
	// 		t.Error("Expected error for invalid cron string, but got nil")
	// 	}

	// 	// Cron string missing command
	// 	cronString = "* * * * *"
	// 	expectedOutput = `Minute         *
	// Hour           *
	// Day of Month   *
	// Month          *
	// Day of Week    *
	// Command
	// `
	// 	output, err = cronParser.Parse(cronString)
	// 	if err != nil {
	// 		t.Errorf("Error parsing cron string missing command: %v", err)
	// 	}
	// 	if output != expectedOutput {
	// 		t.Errorf("Invalid output for cron string missing command\ngot:\n%s\nexpected:\n%s", output, expectedOutput)
	// 	}

	// Cron string with step values
	cronString := "*/15 */2 1-10/2 * * /usr/bin/command"
	expectedOutput := `
Minute         00 15 30 45
Hour           00 02 04 06 08 10 12 14 16 18 20 22
Day of Month   01 03 05 07 09
Month          *
Day of Week    *
Command        /usr/bin/command
	`
	output, err := cronParser.Parse(cronString)
	if err != nil {
		t.Errorf("Error parsing cron string with step values: %v", err)
	}
	if output != expectedOutput{
		t.Errorf("Invalid output for cron string with step values\ngot:\n%s\nexpected:\n%s", output, expectedOutput)
	}
}
