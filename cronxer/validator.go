package cronxer

import (
	"fmt"
	"strconv"
	"strings"
)

// https://crontab.guru/
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
			rangeAndStep := strings.Split(part, "/")
			if len(rangeAndStep) != 2 {
				return fmt.Errorf("invalid step value in field: %s", field)
			}

			rangeField := rangeAndStep[0]
			stepField := rangeAndStep[1]

			if err := validateStepField(stepField, min, max); err != nil {
				return err
			}

			if err := validateRangeField(rangeField, min, max, stepField); err != nil {
				return err
			}
		} else if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")

			if len(rangeParts) != 2 {
				return fmt.Errorf("invalid range value in field: %s", rangeParts)
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

func validateStepField(stepField string, min, max int) error {
	step, err := strconv.Atoi(stepField)
	if err != nil || step <= 0 || step > max-min+1 {
		return fmt.Errorf("invalid step value in field: %s (step must be a positive integer less than or equal to %d)", stepField, max-min+1)
	}

	return nil
}

func validateRangeField(rangeField string, min, max int, stepField string) error {
	if rangeField == "*" {
		return nil
	}

	rangeParts := strings.Split(rangeField, "-")
	if len(rangeParts) != 2 {
		return fmt.Errorf("invalid range value in field: %s", rangeField)
	}

	start, err := strconv.Atoi(rangeParts[0])
	if err != nil || start < min || start > max {
		return fmt.Errorf("invalid range value in field: %s (start must be a number between %d and %d)", rangeField, min, max)
	}

	end, err := strconv.Atoi(rangeParts[1])
	if err != nil || end < min || end > max || end < start {
		return fmt.Errorf("invalid range value in field: %s (end must be a number between %d and %d)", rangeField, min, max)
	}

	//This is an addidtional validation which is commented
	//but https://crontab.guru/#1-10/12_3-5/2_7-10/2_11-12/2_0-6/2 says its a valid cron/
	// "1-10/12 2-4/2 7-10/3 10-12/2 0-6/3 /usr/bin/command" in general makes no sense to have step greater than start
	// step, _ := strconv.Atoi(stepField)
	// if step >end{
	// 	return fmt.Errorf("invalid step range in field %s (step must be a number between %d and %d)",rangeField,start,end)
	// }

	return nil
}
func validateCommandField(field string) error {

	return nil
}
