package service

import (
	"testing"
	"time"
)

func TestCronParser_EveryMinute(t *testing.T) {
	parser := NewCronParser()
	from := time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local)

	next, err := parser.Parse("* * * * *", from)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	expected := from.Add(1 * time.Minute)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestCronParser_SpecificMinute(t *testing.T) {
	parser := NewCronParser()
	from := time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local)

	// 每小时的第15分钟
	next, err := parser.Parse("15 * * * *", from)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	expected := time.Date(2024, 1, 15, 11, 15, 0, 0, time.Local)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestCronParser_SpecificTime(t *testing.T) {
	parser := NewCronParser()
	from := time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local)

	// 每天凌晨2点
	next, err := parser.Parse("0 2 * * *", from)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	expected := time.Date(2024, 1, 16, 2, 0, 0, 0, time.Local)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestCronParser_Range(t *testing.T) {
	parser := NewCronParser()
	from := time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local)

	// 工作日9-17点每小时
	next, err := parser.Parse("0 9-17 * * 1-5", from)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if next.Hour() < 9 || next.Hour() > 17 {
		t.Errorf("hour out of range: %d", next.Hour())
	}
	if next.Weekday() == time.Saturday || next.Weekday() == time.Sunday {
		t.Errorf("should be weekday: %v", next.Weekday())
	}
}

func TestCronParser_Step(t *testing.T) {
	parser := NewCronParser()
	from := time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local)

	// 每5分钟
	next, err := parser.Parse("*/5 * * * *", from)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if next.Minute()%5 != 0 {
		t.Errorf("minute not multiple of 5: %d", next.Minute())
	}
}

func TestCronParser_List(t *testing.T) {
	parser := NewCronParser()
	from := time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local)

	// 每小时的第0, 15, 30, 45分钟
	next, err := parser.Parse("0,15,30,45 * * * *", from)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	validMinutes := map[int]bool{0: true, 15: true, 30: true, 45: true}
	if !validMinutes[next.Minute()] {
		t.Errorf("invalid minute: %d", next.Minute())
	}
}

func TestCronParser_Daily(t *testing.T) {
	parser := NewCronParser()
	from := time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local)

	next, err := parser.Parse("@daily", from)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	expected := time.Date(2024, 1, 16, 0, 0, 0, 0, time.Local)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestCronParser_Hourly(t *testing.T) {
	parser := NewCronParser()
	from := time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local)

	next, err := parser.Parse("@hourly", from)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	expected := from.Add(1 * time.Hour)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestCronParser_Every5Minutes(t *testing.T) {
	parser := NewCronParser()
	from := time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local)

	next, err := parser.Parse("@every 5m", from)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	expected := from.Add(5 * time.Minute)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestCronParser_Every2Hours(t *testing.T) {
	parser := NewCronParser()
	from := time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local)

	next, err := parser.Parse("@every 2h", from)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	expected := from.Add(2 * time.Hour)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestCronParser_InvalidExpression(t *testing.T) {
	parser := NewCronParser()
	from := time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local)

	tests := []struct {
		expr string
		desc string
	}{
		{"", "empty"},
		{"* * *", "too few fields"},
		{"* * * * * *", "too many fields"},
		{"60 * * * *", "minute out of range"},
		{"* 25 * * *", "hour out of range"},
		{"* * 32 * *", "day out of range"},
		{"* * * 13 *", "month out of range"},
		{"* * * * 8", "weekday out of range"},
		{"@unknown", "unknown special"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := parser.Parse(tt.expr, from)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestCronParser_ValidateCron(t *testing.T) {
	parser := NewCronParser()

	valid := []string{
		"* * * * *",
		"0 * * * *",
		"0 0 * * *",
		"0 0 1 * *",
		"0 0 * 1 *",
		"0 0 * * 0",
		"*/5 * * * *",
		"0 9-17 * * 1-5",
		"0,30 * * * *",
		"@daily",
		"@hourly",
		"@every 5m",
	}

	for _, expr := range valid {
		if err := parser.ValidateCron(expr); err != nil {
			t.Errorf("expected valid: %s, got error: %v", expr, err)
		}
	}

	invalid := []string{
		"",
		"* * *",
		"60 * * * *",
		"@unknown",
	}

	for _, expr := range invalid {
		if err := parser.ValidateCron(expr); err == nil {
			t.Errorf("expected invalid: %s", expr)
		}
	}
}

func TestCronParser_ParseInterval(t *testing.T) {
	parser := NewCronParser()

	tests := []struct {
		input    string
		expected time.Duration
		wantErr  bool
	}{
		{"30s", 30 * time.Second, false},
		{"5m", 5 * time.Minute, false},
		{"1h", 1 * time.Hour, false},
		{"2d", 48 * time.Hour, false},
		{"", 0, true},
		{"abc", 0, true},
		{"5x", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			duration, err := parser.ParseInterval(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
			}
			if !tt.wantErr && duration != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, duration)
			}
		})
	}
}

func TestCronParser_WeekdaySunday(t *testing.T) {
	parser := NewCronParser()
	from := time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local) // 周一

	// 每周日
	next, err := parser.Parse("0 0 * * 0", from)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if next.Weekday() != time.Sunday {
		t.Errorf("expected Sunday, got %v", next.Weekday())
	}
}

func TestCronParser_Monthly(t *testing.T) {
	parser := NewCronParser()
	from := time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local)

	next, err := parser.Parse("@monthly", from)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	expected := time.Date(2024, 2, 1, 0, 0, 0, 0, time.Local)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestCronParser_Weekly(t *testing.T) {
	parser := NewCronParser()
	from := time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local) // 周一

	next, err := parser.Parse("@weekly", from)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if next.Weekday() != time.Sunday {
		t.Errorf("expected Sunday, got %v", next.Weekday())
	}
}

func TestContainsInt(t *testing.T) {
	tests := []struct {
		slice []int
		val   int
		want  bool
	}{
		{[]int{1, 2, 3}, 2, true},
		{[]int{1, 2, 3}, 4, false},
		{[]int{}, 1, false},
	}

	for _, tt := range tests {
		if got := containsInt(tt.slice, tt.val); got != tt.want {
			t.Errorf("containsInt(%v, %d) = %v, want %v", tt.slice, tt.val, got, tt.want)
		}
	}
}
