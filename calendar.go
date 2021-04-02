package business

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type calendarConf struct {
	WorkingDays []string `yaml:"working_days"`
	Holidays    []string `yaml:"holidays"`
}

// Choices is wrapper object to return multiple Choice objects
type Calendar struct {
	DaysNames map[string]time.Weekday
	Ordinals  map[string]string

	WorkingDays []time.Weekday
	Holidays    []time.Time
}

func (c *Calendar) String() string {
	bs, _ := json.Marshal(c)
	return string(bs)
}

// NewCalendar instantiate new Calendar object
func NewCalendar() *Calendar {
	return &Calendar{
		DaysNames: map[string]time.Weekday{
			"sunday":    time.Sunday,
			"monday":    time.Monday,
			"tuesday":   time.Tuesday,
			"wednesday": time.Wednesday,
			"thursday":  time.Thursday,
			"friday":    time.Friday,
			"saturday":  time.Saturday,
		},
		Ordinals: map[string]string{
			"1st,": "01", "2nd,": "02", "3rd,": "03", "4th,": "04", "5th,": "05",
			"6th,": "06", "7th,": "07", "8th,": "08", "9th,": "09", "10th,": "10",
			"11th,": "11", "12th,": "12", "13th,": "13", "14th,": "14", "15th,": "15",
			"16th,": "16", "17th,": "17", "18th,": "18", "19th,": "19", "20th,": "20",
			"21st,": "21", "22nd,": "22", "23rd,": "23", "24th,": "24", "25th,": "25",
			"26th,": "26", "27th,": "27", "28th,": "28", "29th,": "29", "30th,": "30",
			"31st,": "31",
		},
	}
}

// Relative endpoint: POST /payments
func (c *Calendar) Load(calendar string) (err error) {

	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("./data/%s.yml", calendar))
	if err != nil {
		return
	}

	config := calendarConf{}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return
	}

	// update working days
	for _, workday := range config.WorkingDays {
		if d, ok := c.DaysNames[workday]; ok {
			c.WorkingDays = append(c.WorkingDays, d)
		}
	}

	// and load in holidays
	for _, holiday := range config.Holidays {

		// parse the ordinal form data string
		holdate, parseErr := c.parseDate(holiday)

		if parseErr != nil {
			return parseErr
		}

		// add to our hol array
		c.Holidays = append(c.Holidays, holdate)
	}

	return
}

// Parse format "August 29th, 2011"
func (c *Calendar) parseDate(s string) (time.Time, error) {
	i := strings.IndexByte(s, ' ') + 1
	if i > 0 {
		j := i + strings.IndexByte(s[i:], ' ')
		if j > i {
			if c, ok := c.Ordinals[s[i:j]]; ok {
				s = s[:i] + c + s[j:]
			}
		}
	}
	return time.Parse("January 02 2006", s)
}

// Return true if the date given is identified as a holiday
func (c *Calendar) IsHoliday(date time.Time) bool {

	for _, holiday := range c.Holidays {
		if date.Equal(holiday) {
			return true
		}
	}

	return false
}

// Return true if the date given is a working day
func (c *Calendar) IsWorkingDay(date time.Time) bool {

	for _, workday := range c.WorkingDays {
		if date.Weekday() == workday {
			return true
		}
	}

	return false
}

// Return true if the date given is a working day and not a holiday.
func (c *Calendar) IsBusinessDay(date time.Time) bool {

	if c.IsHoliday(date) {
		return false
	}

	return c.IsWorkingDay(date)
}

// Roll forward to the next business day. If the date given is a business day, that day will be returned.
func (c *Calendar) RollForward(date time.Time) time.Time {

	curdate := date

	for !c.IsBusinessDay(curdate) {
		curdate = curdate.AddDate(0, 0, 1)
	}

	return curdate
}

// Roll backward to the previous business day. If the date given is a business day, that day will be returned.
func (c *Calendar) RollBackward(date time.Time) time.Time {

	curdate := date

	for !c.IsBusinessDay(curdate) {
		curdate = curdate.AddDate(0, 0, -1)
	}

	return curdate
}

// Roll backward to the previous business day regardless of whether the given date is a business day or not.
func (c *Calendar) NextBusinessDay(date time.Time) time.Time {

	curdate := date.AddDate(0, 0, 1)

	return c.RollForward(curdate)
}

// Roll backward to the previous business day regardless of whether the given date is a business day or not.
func (c *Calendar) PreviousBusinessDay(date time.Time) time.Time {

	curdate := date.AddDate(0, 0, -1)

	return c.RollBackward(curdate)
}

// Count the number of business days between two dates.
// This method counts from start of from_date to start of to_date. So,
//        business_days_between(mon, weds) = 2 (assuming no holidays)
func (c *Calendar) BusinessDaysBetween(from time.Time, to time.Time) int {

	var days int = 0
	var direction int = 0

	//
	curdate := from

	//
	if to.Equal(from) {
		return 0
	} else if from.Before(to) {
		// we're moving forward
		direction = 1
	} else {
		// we're moving backward
		direction = -1
	}

	for !to.Equal(curdate) {
		if c.IsBusinessDay(curdate) {
			days += 1
		}
		curdate = curdate.AddDate(0, 0, direction)
	}

	return days * direction
}

// Add or subtract a number of business days to a date.
// If a non-business day is given, counting will start from the next business day. So:
//            monday + 1 = tuesday
//            friday + 1 = monday
//            sunday + 1 = tuesday
//            friday - 1 = thursday
//            monday - 1 = friday
//            sunday - 1 = thursday
func (c *Calendar) AddBusinessDays(date time.Time, delta int) time.Time {

	curdate := date

	if delta == 0 {
		return curdate
	} else if delta < 0 {
		curdate = c.RollBackward(curdate)
	} else {
		curdate = c.RollForward(curdate)
	}

	for i := 0; i < abs(delta); i++ {
		if delta < 0 {
			curdate = c.PreviousBusinessDay(curdate)
		} else {
			curdate = c.NextBusinessDay(curdate)
		}
	}

	return curdate
}

// Get the business day of the month for a given input date.
func (c *Calendar) GetBusinessDay(date time.Time) int {

	firstdate := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, date.Nanosecond(), date.Location())

	return c.BusinessDaysBetween(firstdate, date.AddDate(0, 0, 1))
}

// abs helper function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
