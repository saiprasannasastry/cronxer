package cronxer

import (
	"fmt"
	"strconv"
	"strings"
)

func validateCronString(cronString string) error {
	fields := strings.Fields(cronString)

	if len(fields) != 5 && len(fields) != 6 {
		return fmt.Errorf("invalid cron string: %s", cronString)
	}

	if err := validateField(fields[0], 0, 59); err != nil {
		return fmt.Errorf("invalid minute field: %v", err)
	}

	if err := validateField(fields[1], 0, 23); err != nil {
		return fmt.Errorf("invalid hour field: %v", err)
	}

	if err := validateField(fields[2], 1, 31); err != nil {
		return fmt.Errorf("invalid day of month field: %v", err)
	}

	if err := validateField(fields[3], 1, 12); err != nil {
		return fmt.Errorf("invalid month field: %v", err)
	}

	if err := validateField(fields[4], 0, 6); err != nil {
		return fmt.Errorf("invalid day of week field: %v", err)
	}

	if len(fields) == 6 {
		if err := validateCommandField(fields[5]); err != nil {
			return fmt.Errorf("invalid command field: %v", err)
		}
	}

	return nil
}

func validateField(field string, min, max int) error {
	if field == "*" {
		return nil
	}

	parts := strings.Split(field, ",")
	for _, part := range parts {
		if strings.Contains(part, "/") {
			stepParts := strings.Split(part, "/")
			if len(stepParts) != 2 {
				return fmt.Errorf("invalid step value in field: %s", field)
			}

			rangeField := stepParts[0]
			step, err := strconv.Atoi(stepParts[1])
			if err != nil || step <= 0 || step > max-min+1 {
				return fmt.Errorf("invalid step value in field: %s (step must be a positive integer less than or equal to %d)", field, max-min+1)
			}

			if err := validateRange(rangeField, min, max); err != nil {
				return err
			}
		} else if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return fmt.Errorf("invalid range value in field: %s", field)
			}

			start, err := strconv.Atoi(rangeParts[0])
			if err != nil || start < min || start > max {
				return fmt.Errorf("invalid range value in field: %s (start must be a number between %d and %d)", field, min, max)
			}

			end, err := strconv.Atoi(rangeParts[1])
			if err != nil || end < min || end > max || end < start {
				return fmt.Errorf("invalid range value in field: %s (end must be a number between %d and %d)", field, min, max)
			}
		} else {
			value, err := strconv.Atoi(part)
			if err != nil || value < min || value > max {
				return fmt.Errorf("invalid value in field: %s (value must be a number between %d and %d)", field, min, max)
			}
		}
	}

	return nil
}

func validateRange(rangeStr string, min, max int) error {
	if rangeStr == "*" {
		return nil
	}
	value, err := strconv.Atoi(rangeStr)
	if err != nil || value < min || value > max {
		return fmt.Errorf("invalid range value in field: %s", rangeStr)
	}

	return nil
}

func validateCommandField(field string) error {
	// Additional validation logic for the command field can be added here if needed.
	// ask the interviewer if * is required
	return nil
}
