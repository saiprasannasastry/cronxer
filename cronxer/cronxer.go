package cronxer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	minutesPos = iota
	hoursPos
	daysOfMonthPos
	monthsPos
	daysOfWeekPos
	commandPos
	numFields = 6
)

const (
	timeFormat       = "Jan _2 15:04"
	minutesField     = "Minute"
	hoursField       = "Hour"
	daysOfMonthField = "Day of Month"
	monthsField      = "Month"
	daysOfWeekField  = "Day of Week"
	commandField     = "Command"
)

var (
	instance *CronParser
)

// CronParser represents the singleton Cron Parser.
type CronParser struct{}

// GetInstance returns the singleton instance of CronParser.
func New() *CronParser {
	if instance == nil {
		instance = &CronParser{}
	}
	return instance
}

// Parse parses the given cron string and returns the output as a string.
func (cp *CronParser) Parse(cronString string) (string, error) {
	fields := strings.Fields(cronString)

	if len(fields) < numFields-1 || len(fields) > numFields {
		return "", errors.New("invalid cron string")
	}
	err := validateCronString(cronString)
	if err != nil {
		return "", err
	}

	schedules := make([][]string, numFields-1)
	schedules[minutesPos] = expandField(fields[minutesPos], 0, 59)
	schedules[hoursPos] = expandField(fields[hoursPos], 0, 23)
	schedules[daysOfMonthPos] = expandField(fields[daysOfMonthPos], 1, 31)
	schedules[monthsPos] = expandField(fields[monthsPos], 1, 12)
	schedules[daysOfWeekPos] = expandField(fields[daysOfWeekPos], 0, 6)

	command := ""
	if len(fields) == numFields {
		command = fields[commandPos]
	}

	headers := []string{minutesField, hoursField, daysOfMonthField, monthsField, daysOfWeekField}

	return generateTable(headers, schedules, command), nil
}

func expandField(field string, min, max int) []string {
	if field == "*" {
		return []string{"*"}
		//if we need all days use the below 
		//return generateRange(min, max, 1)
	}

	parts := strings.Split(field, ",")
	result := []string{}
	for _, part := range parts {
		if strings.Contains(part, "/") && strings.Contains(part, "-") {
			// Handle the case where both range and step are present
			rangeAndStepParts := strings.Split(part, "/")
			rangePart := rangeAndStepParts[0]
			stepPart := rangeAndStepParts[1]

			rangeParts := strings.Split(rangePart, "-")
			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil
			}
			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil
			}

			step, err := strconv.Atoi(stepPart)
			if err != nil {
				return nil
			}

			result = append(result, generateRange(start, end, step)...)
		} else if strings.Contains(part, "/") {
			step, err := strconv.Atoi(strings.Split(part, "/")[1])
			if err != nil {
				return nil
			}
			result = append(result, generateRange(min, max, step)...)
		} else if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil
			}
			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil
			}
			result = append(result, generateRange(start, end, 1)...)
		} else {
			value, err := strconv.Atoi(part)
			if err != nil {
				return nil
			}
			if value >= min && value <= max {
				result = append(result, fmt.Sprintf("%02d", value))
			}
		}
	}

	return result
}

func generateRange(start, end, step int) []string {
	result := []string{}
	for i := start; i <= end; i += step {
		result = append(result, fmt.Sprintf("%02d", i))
	}
	return result
}

func generateTable(headers []string, schedules [][]string, command string) string {
	maxHeaderLength := getMaxHeaderLength(headers)
	var output strings.Builder

	for i := 0; i < numFields-1; i++ {
		output.WriteString(fmt.Sprintf("%-*s%s\n", maxHeaderLength+2, headers[i], strings.Join(schedules[i], " ")))
	}

	if command != "" {
		output.WriteString(fmt.Sprintf("%-*s%s\n", maxHeaderLength+2, commandField, command))
	}

	return output.String()
}

func getMaxHeaderLength(headers []string) int {
	maxHeaderLength := 0
	for _, header := range headers {
		if len(header) > maxHeaderLength {
			maxHeaderLength = len(header)
		}
	}
	return maxHeaderLength
}
