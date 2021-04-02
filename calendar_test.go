package business

import (
	"fmt"
	"testing"
	"time"
)

var holidaytests = []struct {
	in  string
	out bool
}{
	{"January 01 2020", true},    // new years day
	{"April 10 2020", true},      // Good Friday
	{"December 25 2020", true},   // Christmas Day
	{"June 23 2020", false},      //
	{"September 14 2020", false}, //
}

func TestIsHoliday(t *testing.T) {

	calendar := NewCalendar()
	result := calendar.Load("bacs")

	if result != nil {
		t.Errorf("error loading bacs calendar: %q", result)
	}

	for _, tt := range holidaytests {
		t.Run(tt.in, func(t *testing.T) {

			holdate, _ := time.Parse("January 02 2006", tt.in)

			result := calendar.IsHoliday(holdate)

			// christmas day
			if result != tt.out {
				t.Errorf("error indentifying holiday got %v, want %v", result, tt.out)
			}
		})
	}
}

var workdaytests = []struct {
	in  string
	out bool
}{
	{"January 01 2020", true},    // new years day
	{"April 10 2020", true},      // Good Friday
	{"December 25 2020", true},   // Christmas Day
	{"December 26 2020", false},  //
	{"June 23 2020", true},       //
	{"September 26 2020", false}, //
}

func TestIsWorkingDay(t *testing.T) {

	calendar := NewCalendar()
	result := calendar.Load("bacs")

	if result != nil {
		t.Errorf("error loading bacs calendar: %q", result)
	}

	for _, tt := range workdaytests {
		t.Run(tt.in, func(t *testing.T) {

			holdate, _ := time.Parse("January 02 2006", tt.in)

			result := calendar.IsWorkingDay(holdate)

			// christmas day
			if result != tt.out {
				t.Errorf("error indentifying weekday got %v, want %v", result, tt.out)
			}
		})
	}
}

var businessdaytests = []struct {
	in  string
	out bool
}{
	{"January 01 2020", false},  // new years day
	{"April 10 2020", false},    // Good Friday
	{"December 25 2020", false}, // Christmas Day
	{"December 26 2020", false}, //
	{"June 23 2020", true},      //
	{"September 25 2020", true}, //
}

func TestIsBusinessDay(t *testing.T) {

	calendar := NewCalendar()
	result := calendar.Load("bacs")

	if result != nil {
		t.Errorf("error loading bacs calendar: %q", result)
	}

	for _, tt := range businessdaytests {
		t.Run(tt.in, func(t *testing.T) {

			holdate, _ := time.Parse("January 02 2006", tt.in)

			result := calendar.IsBusinessDay(holdate)

			// christmas day
			if result != tt.out {
				t.Errorf("error indentifying business day got %v, want %v", result, tt.out)
			}
		})
	}
}

var rollfowardtests = []struct {
	in  string
	out string
}{
	{"January 01 2020", "January 02 2020"},     // new years day
	{"April 10 2020", "April 14 2020"},         // Good Friday
	{"December 25 2020", "December 29 2020"},   // Christmas Day
	{"December 26 2020", "December 29 2020"},   //
	{"June 23 2020", "June 23 2020"},           //
	{"September 26 2020", "September 28 2020"}, //
}

func TestRollForward(t *testing.T) {

	calendar := NewCalendar()
	result := calendar.Load("bacs")

	if result != nil {
		t.Errorf("error loading bacs calendar: %q", result)
	}

	for _, tt := range rollfowardtests {
		t.Run(tt.in, func(t *testing.T) {

			curdate, _ := time.Parse("January 02 2006", tt.in)
			targetdate, _ := time.Parse("January 02 2006", tt.out)

			result := calendar.RollForward(curdate)

			// christmas day
			if !result.Equal(targetdate) {
				t.Errorf("error indentifying forward day got %v, want %v", result, targetdate)
			}
		})
	}
}

var rollbackwardtests = []struct {
	in  string
	out string
}{
	{"January 01 2020", "December 31 2019"},    // new years day
	{"April 10 2020", "April 09 2020"},         // Good Friday
	{"December 25 2020", "December 24 2020"},   // Christmas Day
	{"December 26 2020", "December 24 2020"},   //
	{"June 23 2020", "June 23 2020"},           //
	{"September 26 2020", "September 25 2020"}, //
}

func TestRoleBackward(t *testing.T) {

	calendar := NewCalendar()
	result := calendar.Load("bacs")

	if result != nil {
		t.Errorf("error loading bacs calendar: %q", result)
	}

	for _, tt := range rollbackwardtests {
		t.Run(tt.in, func(t *testing.T) {

			curdate, _ := time.Parse("January 02 2006", tt.in)
			targetdate, _ := time.Parse("January 02 2006", tt.out)

			result := calendar.RollBackward(curdate)

			// christmas day
			if !result.Equal(targetdate) {
				t.Errorf("error indentifying backward day got %v, want %v", result, targetdate)
			}
		})
	}
}

var nextbusinesstests = []struct {
	in  string
	out string
}{
	{"January 01 2020", "January 02 2020"},     // new years day
	{"April 10 2020", "April 14 2020"},         // Good Friday
	{"December 25 2020", "December 29 2020"},   // Christmas Day
	{"December 26 2020", "December 29 2020"},   // Boxing Day
	{"June 23 2020", "June 24 2020"},           //
	{"September 26 2020", "September 28 2020"}, //
}

func TestNextBusinessDay(t *testing.T) {

	calendar := NewCalendar()
	result := calendar.Load("bacs")

	if result != nil {
		t.Errorf("error loading bacs calendar: %q", result)
	}

	for _, tt := range nextbusinesstests {
		t.Run(tt.in, func(t *testing.T) {

			curdate, _ := time.Parse("January 02 2006", tt.in)
			targetdate, _ := time.Parse("January 02 2006", tt.out)

			result := calendar.NextBusinessDay(curdate)

			// christmas day
			if !result.Equal(targetdate) {
				t.Errorf("error indentifying next business day got %v, want %v", result, targetdate)
			}
		})
	}
}

var prevbusinesstests = []struct {
	in  string
	out string
}{
	{"January 01 2020", "December 31 2019"},    // new years day
	{"April 10 2020", "April 09 2020"},         // Good Friday
	{"December 25 2020", "December 24 2020"},   // Christmas Day
	{"December 26 2020", "December 24 2020"},   // Boxing Day
	{"June 23 2020", "June 22 2020"},           //
	{"September 26 2020", "September 25 2020"}, //
}

func TestPreviousBusinessDay(t *testing.T) {

	calendar := NewCalendar()
	result := calendar.Load("bacs")

	if result != nil {
		t.Errorf("error loading bacs calendar: %q", result)
	}

	for _, tt := range prevbusinesstests {
		t.Run(tt.in, func(t *testing.T) {

			curdate, _ := time.Parse("January 02 2006", tt.in)
			targetdate, _ := time.Parse("January 02 2006", tt.out)

			result := calendar.PreviousBusinessDay(curdate)

			// christmas day
			if !result.Equal(targetdate) {
				t.Errorf("error indentifying next business day got %v, want %v", result, targetdate)
			}
		})
	}
}

var addbusinesstests = []struct {
	in    string
	delta int
	out   string
}{
	{"January 01 2020", 1, "January 03 2020"},   // new years day
	{"April 10 2020", 2, "April 16 2020"},       // Good Friday
	{"December 25 2020", 2, "December 31 2020"}, // Christmas Day
	{"December 26 2020", 6, "January 07 2021"},  // Boxing Day
	{"June 23 2020", -7, "June 12 2020"},        //
	{"September 26 2020", 200, "July 14 2021"},  //
	{"April 01 2021", 1, "April 06 2021"},       //
}

func TestAddBusinessDays(t *testing.T) {

	calendar := NewCalendar()
	result := calendar.Load("bacs")

	if result != nil {
		t.Errorf("error loading bacs calendar: %q", result)
	}

	for _, tt := range addbusinesstests {
		t.Run(tt.in, func(t *testing.T) {

			curdate, _ := time.Parse("January 02 2006", tt.in)
			targetdate, _ := time.Parse("January 02 2006", tt.out)

			result := calendar.AddBusinessDays(curdate, tt.delta)

			// christmas day
			if !result.Equal(targetdate) {
				t.Errorf("error indentifying next business day got %v, want %v", result, targetdate)
			}
		})
	}
}

var businessdaysbetweentests = []struct {
	to     string
	from   string
	result int
}{
	{"June 02 2014", "June 05 2014", 3},
	{"June 02 2014", "June 09 2014", 5},
	{"May 29 2014", "June 03 2014", 3},
	{"June 09 2014", "June 13 2014", 4},
	{"June 02 2014", "June 13 2014", 9},
	{"June 26 2014", "July 01 2014", 3},
	{"January 01 2020", "January 03 2020", 1},
	{"April 10 2020", "April 16 2020", 2},
	{"December 25 2020", "December 31 2020", 2},
	{"December 26 2020", "January 07 2021", 6},
	{"June 23 2020", "June 12 2020", -7},
	{"September 26 2020", "July 14 2021", 200},
}

func TestBusinessDaysBetween(t *testing.T) {

	calendar := NewCalendar()
	result := calendar.Load("bacs")

	if result != nil {
		t.Errorf("error loading bacs calendar: %q", result)
	}

	for _, tt := range businessdaysbetweentests {
		t.Run(fmt.Sprintf("%v->%v", tt.from, tt.to), func(t *testing.T) {

			from, _ := time.Parse("January 02 2006", tt.from)
			to, _ := time.Parse("January 02 2006", tt.to)

			result := calendar.BusinessDaysBetween(to, from)

			// christmas day
			if result != tt.result {
				t.Errorf("error indentifying next business between %v and %v got %v, want %v", from, to, result, tt.result)
			}
		})
	}
}

var businessdaystests = []struct {
	in     string
	result int
}{
	{"June 02 2014", 1},
	{"May 29 2014", 19},
	{"June 09 2014", 6},
	{"June 26 2014", 19},
	{"January 01 2020", 0},
	{"April 10 2020", 7},
	{"December 25 2020", 18},
	{"December 26 2020", 18},
	{"June 23 2020", 17},
	{"September 26 2020", 19},
}

func TestBusinessDay(t *testing.T) {

	calendar := NewCalendar()
	result := calendar.Load("bacs")

	if result != nil {
		t.Errorf("error loading bacs calendar: %q", result)
	}

	for _, tt := range businessdaystests {
		t.Run(tt.in, func(t *testing.T) {

			from, _ := time.Parse("January 02 2006", tt.in)

			result := calendar.GetBusinessDay(from)

			// christmas day
			if result != tt.result {
				t.Errorf("error indentifying business day for %v got %v, want %v", from, result, tt.result)
			}
		})
	}
}
