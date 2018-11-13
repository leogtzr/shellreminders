package main

import (
	"io/ioutil"
	"testing"
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

func TestReminderRecordParsing(t *testing.T) {
	tests := []struct {
		input                string
		expectedDays         int
		expectedName         string
		expectedReminderType ReminderType
	}{
		{`Santander Platino;18`, 18, "Santander Platino", RecurrentReminder},
		{`Promotions;13;counter`, 13, "Promotions", Counter},
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
		if reminder.Type != tt.expectedReminderType {
			t.Fatalf("got=%+v as type, expected=%+v", reminder.Type, tt.expectedReminderType)
		}
	}
}
