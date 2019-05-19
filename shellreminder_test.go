package main

import (
	"io/ioutil"
	"testing"
	"time"
)

func TestReminderFileParsing(t *testing.T) {
	reminderFileContent := `Santander Platino;18
Promotions;13;counter
	`
	err := ioutil.WriteFile("/tmp/rmnd.txt", []byte(reminderFileContent), 0644)
	if err != nil {
		t.Fatal("Error generating test reminder file")
	}

	reminders, err := parseRemindersFromFile("/tmp/rmnd.txt")
	if err != nil {
		t.Fatal("error parsing reminder file")
	}

	expectedRemindersCount := 2
	if len(reminders) != expectedRemindersCount {
		t.Fatalf("error parsing reminder file, it should have had %d records, got=%d", expectedRemindersCount, len(reminders))
	}
}

func generateBaseDateTime() time.Time {
	now := time.Now()
	then := time.Date(2018, now.Month(), now.Day(), 0 /*the hour */, 0 /* the minutes */, 0, 0, time.UTC)
	return then
}

func TestReminderRecordParsing(t *testing.T) {
	tests := []struct {
		input        string
		expectedDays int
		expectedName string
	}{
		{`Santander Platino;18`, 18, "Santander Platino"},
		{`Promotions;13;counter`, 13, "Promotions"},
	}

	for _, tt := range tests {
		reminder, err := extractReminderFromText(tt.input)
		if err != nil {
			t.Fatalf("error parsing '%s' record", tt.input)
		}
		if reminder.Name != tt.expectedName {
			t.Fatalf("got=%s as name, expected=%s", reminder.Name, tt.expectedName)
		}
		if reminder.EveryWhen != tt.expectedDays {
			t.Fatalf("got=%d as days, expected=%d", reminder.EveryWhen, tt.expectedDays)
		}
	}
}

func TestReminderString(t *testing.T) {
	// 'Santander Platino' day 28 of each month
	r := Reminder{Name: "Santander Platino", EveryWhen: 28}
	expectedReminderToStringValue := "'Santander Platino' day 28 of each month"
	if expectedReminderToStringValue != r.String() {
		t.Fatalf("got=[%q] as string, expected=[%q]", r.String(), expectedReminderToStringValue)
	}
}
