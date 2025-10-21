package evaluator

import (
	"testing"
	"time"

	"github.com/vpaulo/seda/object"
)

func TestTimeNow(t *testing.T) {
	input := `Time.now()`
	result := testEval(input)

	timeObj, ok := result.(*object.Time)
	if !ok {
		t.Fatalf("result is not Time. got=%T (%+v)", result, result)
	}

	// Check that the time is recent (within last 5 seconds)
	now := time.Now()
	diff := now.Sub(timeObj.Value)
	if diff < 0 || diff > 5*time.Second {
		t.Errorf("Time.now() returned time that is not current. diff=%v", diff)
	}
}

func TestTimeDateParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Time
	}{
		{
			`Time.date("DD-MM-YYYY", "20-01-2000")`,
			time.Date(2000, 1, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			`Time.date("YYYY-MM-DD", "2025-10-14")`,
			time.Date(2025, 10, 14, 0, 0, 0, 0, time.UTC),
		},
		{
			`Time.date("YYYY-MM-DD HH:mm:ss", "2025-10-14 15:30:45")`,
			time.Date(2025, 10, 14, 15, 30, 45, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		result := testEval(tt.input)

		timeObj, ok := result.(*object.Time)
		if !ok {
			t.Fatalf("result is not Time. got=%T (%+v)", result, result)
		}

		if !timeObj.Value.Equal(tt.expected) {
			t.Errorf("Time.date() parsed incorrectly. expected=%v, got=%v",
				tt.expected, timeObj.Value)
		}
	}
}

func TestTimeUnix(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`Time.unix(0)`, 0},                   // Unix epoch
		{`Time.unix(1609459200)`, 1609459200}, // 2021-01-01 00:00:00 UTC
		{`Time.unix(1728921045)`, 1728921045}, // 2024-10-14 15:30:45 UTC
	}

	for _, tt := range tests {
		result := testEval(tt.input)

		timeObj, ok := result.(*object.Time)
		if !ok {
			t.Fatalf("result is not Time. got=%T (%+v)", result, result)
		}

		if timeObj.Value.Unix() != tt.expected {
			t.Errorf("Time.unix() incorrect. expected=%d, got=%d",
				tt.expected, timeObj.Value.Unix())
		}
	}
}

func TestTimeFormatting(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Time.date("YYYY-MM-DD", "2000-01-20").format("DD-MM-YYYY")`, "20-01-2000"},
		{`Time.date("YYYY-MM-DD", "2000-01-20").format("YYYY/MM/DD")`, "2000/01/20"},
		{`Time.date("YYYY-MM-DD HH:mm:ss", "2025-10-14 15:30:45").format("HH:mm:ss")`, "15:30:45"},
	}

	for _, tt := range tests {
		result := testEval(tt.input)

		str, ok := result.(*object.String)
		if !ok {
			t.Fatalf("result is not String. got=%T (%+v)", result, result)
		}

		if str.Value != tt.expected {
			t.Errorf("formatting incorrect. expected=%s, got=%s", tt.expected, str.Value)
		}
	}
}

func TestTimeToString(t *testing.T) {
	input := `Time.date("YYYY-MM-DD", "2000-01-20").to_string()`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	// Should contain the date in ISO 8601 format
	if len(str.Value) == 0 || str.Value[:10] != "2000-01-20" {
		t.Errorf("to_string() incorrect. got=%s", str.Value)
	}
}

func TestTimeUnixMethod(t *testing.T) {
	input := `Time.date("YYYY-MM-DD", "2000-01-20").unix()`
	result := testEval(input)

	num, ok := result.(*object.Number)
	if !ok {
		t.Fatalf("result is not Number. got=%T (%+v)", result, result)
	}

	expected := time.Date(2000, 1, 20, 0, 0, 0, 0, time.UTC).Unix()
	if int64(num.Value) != expected {
		t.Errorf("unix() incorrect. expected=%d, got=%d", expected, int64(num.Value))
	}
}

func TestTimeComponents(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`Time.date("YYYY-MM-DD", "2000-01-20").year()`, 2000},
		{`Time.date("YYYY-MM-DD", "2000-01-20").month()`, 1},
		{`Time.date("YYYY-MM-DD", "2000-01-20").day()`, 20},
		{`Time.date("YYYY-MM-DD HH:mm:ss", "2025-10-14 15:30:45").hour()`, 15},
		{`Time.date("YYYY-MM-DD HH:mm:ss", "2025-10-14 15:30:45").minute()`, 30},
		{`Time.date("YYYY-MM-DD HH:mm:ss", "2025-10-14 15:30:45").second()`, 45},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

func TestTimeArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{
			`var t1 = Time.date("YYYY-MM-DD", "2000-01-01")
			 var t2 = t1.add_seconds(60)
			 t2.diff(t1)`,
			60,
		},
		{
			`var t1 = Time.date("YYYY-MM-DD", "2000-01-01")
			 var t2 = t1.add_minutes(5)
			 t2.diff(t1)`,
			300,
		},
		{
			`var t1 = Time.date("YYYY-MM-DD", "2000-01-01")
			 var t2 = t1.add_hours(2)
			 t2.diff(t1)`,
			7200,
		},
		{
			`var t1 = Time.date("YYYY-MM-DD", "2000-01-01")
			 var t2 = t1.add_days(1)
			 t2.diff(t1)`,
			86400,
		},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

func TestTimeComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			`var t1 = Time.date("YYYY-MM-DD", "2000-01-01")
			 var t2 = Time.date("YYYY-MM-DD", "2000-01-02")
			 t1.is_before(t2)`,
			true,
		},
		{
			`var t1 = Time.date("YYYY-MM-DD", "2000-01-02")
			 var t2 = Time.date("YYYY-MM-DD", "2000-01-01")
			 t1.is_before(t2)`,
			false,
		},
		{
			`var t1 = Time.date("YYYY-MM-DD", "2000-01-02")
			 var t2 = Time.date("YYYY-MM-DD", "2000-01-01")
			 t1.is_after(t2)`,
			true,
		},
		{
			`var t1 = Time.date("YYYY-MM-DD", "2000-01-01")
			 var t2 = Time.date("YYYY-MM-DD", "2000-01-02")
			 t1.is_after(t2)`,
			false,
		},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testBooleanObject(t, result, tt.expected)
	}
}

func TestTimeImmutability(t *testing.T) {
	input := `
		var t1 = Time.date("YYYY-MM-DD", "2000-01-01")
		var unix1 = t1.unix()
		var t2 = t1.add_days(10)
		var unix1_after = t1.unix()
		unix1 == unix1_after
	`

	result := testEval(input)
	testBooleanObject(t, result, true)
}

func TestTimeChaining(t *testing.T) {
	input := `
		var t1 = Time.date("YYYY-MM-DD", "2000-01-01")
		var t2 = t1.add_days(1).add_hours(2).add_minutes(30).add_seconds(45)
		t2.diff(t1)
	`

	expected := float64(86400 + 7200 + 1800 + 45) // 1 day + 2 hours + 30 min + 45 sec
	result := testEval(input)
	testNumberObject(t, result, expected)
}

func TestTimeRoundTrip(t *testing.T) {
	input := `
		var t1 = Time.date("YYYY-MM-DD HH:mm:ss", "2000-01-20 15:30:45")
		var unix_val = t1.unix()
		var t2 = Time.unix(unix_val)
		t1.unix() == t2.unix()
	`

	result := testEval(input)
	testBooleanObject(t, result, true)
}

func TestTimeErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			`Time.now(42)`,
			"Time.now() takes no arguments, got 1",
		},
		{
			`Time.date("YYYY-MM-DD")`,
			"Time.date() takes 2 arguments (format, dateString), got 1",
		},
		{
			`Time.unix()`,
			"Time.unix() takes 1 argument (seconds), got 0",
		},
		{
			`Time.unix("not a number")`,
			"Time.unix() argument must be NUMBER, got STRING",
		},
	}

	for _, tt := range tests {
		result := testEval(tt.input)

		errObj, ok := result.(*object.Error)
		if !ok {
			t.Errorf("Expected error object, got %T", result)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("Expected error message %q, got %q", tt.expectedMessage, errObj.Message)
		}
	}
}
