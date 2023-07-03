package cronxer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
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
	timeFormat       = "Jan _2, 2006 15:04"
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

// expand fields expands all cases of cron string and returns a resulting array of fields
func expandField(field string, min, max int) []string {
	if field == "*" {
		return []string{"*"}
	}

	parts := strings.Split(field, ",")
	result := make([]string, 0)

	for _, part := range parts {
		expandedValues, err := expandPart(part, min, max)
		if err != nil {
			return nil
		}
		result = append(result, expandedValues...)
	}

	return result
}

// expandPart expands a single part of the cron field and returns the resulting values
func expandPart(part string, min, max int) ([]string, error) {
	if strings.Contains(part, "/") && strings.Contains(part, "-") {
		// Handle the case where both range and step are present
		rangeAndStepParts := strings.Split(part, "/")
		rangePart := rangeAndStepParts[0]
		stepPart := rangeAndStepParts[1]

		rangeParts := strings.Split(rangePart, "-")
		start, err := strconv.Atoi(rangeParts[0])
		if err != nil {
			return nil, err
		}
		end, err := strconv.Atoi(rangeParts[1])
		if err != nil {
			return nil, err
		}

		step, err := strconv.Atoi(stepPart)
		if err != nil {
			return nil, err
		}

		return generateRange(start, end, step), nil
	} else if strings.Contains(part, "/") {
		step, err := strconv.Atoi(strings.Split(part, "/")[1])
		if err != nil {
			return nil, err
		}
		return generateRange(min, max, step), nil
	} else if strings.Contains(part, "-") {
		rangeParts := strings.Split(part, "-")
		start, err := strconv.Atoi(rangeParts[0])
		if err != nil {
			return nil, err
		}
		end, err := strconv.Atoi(rangeParts[1])
		if err != nil {
			return nil, err
		}
		return generateRange(start, end, 1), nil
	} else {
		value, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		if value >= min && value <= max {
			return []string{fmt.Sprintf("%02d", value)}, nil
		}
	}

	return nil, errors.New("invalid cron field")
}

// generateRange generates values in the given range between start and end
func generateRange(start, end, step int) []string {
	result := make([]string, 0)

	for current := start; current <= end; current += step {
		value := fmt.Sprintf("%02d", current)
		result = append(result, value)
	}

	return result
}

// generateTable takes in schedules and returns a string builder
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

func (cp *CronParser) GetNextCronJobs(cronString string, n int) ([]string, error) {

	fields := strings.Fields(cronString)

	if len(fields) < numFields-1 || len(fields) > numFields {
		return nil, errors.New("invalid cron string")
	}
	err := validateCronString(cronString)
	if err != nil {
		return nil, err
	}

	schedules := make([][]string, numFields-1)
	schedules[minutesPos] = expandField(fields[minutesPos], 0, 59)
	schedules[hoursPos] = expandField(fields[hoursPos], 0, 23)
	schedules[daysOfMonthPos] = expandField(fields[daysOfMonthPos], 1, 31)
	schedules[monthsPos] = expandField(fields[monthsPos], 1, 12)
	schedules[daysOfWeekPos] = expandField(fields[daysOfWeekPos], 0, 6)

	currentTime := time.Now()
	nextJobs := make([]string, 0, n)
	for i := 0; i < n; i++ {
		nextJob, err := getNextJob(schedules, currentTime)
		if err != nil {
			return nil, err
		}
		nextJobFormatter := nextJob.Format(timeFormat)
		nextJobs = append(nextJobs, nextJobFormatter)

		currentTime = nextJob.Add(time.Minute)
	}
	return nextJobs, nil
}

func getNextJob(schedules [][]string, currentTime time.Time) (time.Time, error) {
	for {
		currentMinute := strconv.Itoa(currentTime.Minute())
		currentHour := strconv.Itoa(currentTime.Hour())
		currentDayofTheMonth := strconv.Itoa(currentTime.Day())
		if currentDayofTheMonth == "28" {
		}
		currentMonth := strconv.Itoa(int(currentTime.Month()))
		currentWeek := strconv.Itoa(int(currentTime.Weekday()))

		scheduleMinute := schedules[minutesPos]
		scheduleHour := schedules[hoursPos]
		scheduleDayOftheMonth := schedules[daysOfMonthPos]
		scheduleMonth := schedules[monthsPos]
		scheduleWeek := schedules[daysOfWeekPos]

		if isScheduled(scheduleMinute, currentMinute) &&
			isScheduled(scheduleHour, currentHour) &&
			isScheduled(scheduleDayOftheMonth, currentDayofTheMonth) && isScheduled(scheduleMonth, currentMonth) && isScheduled(scheduleWeek, currentWeek) {
			return currentTime, nil
		}

		currentTime = currentTime.Add(time.Minute)

	}
}
func isScheduled(schedule []string, value string) bool {

	for _, sched := range schedule {
		if strings.HasPrefix(sched, "0") {
			sched = sched[1:]
		}
		if sched == value || sched == "*" {
			return true
		}
	}
	return false
}
