package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/muesli/termenv"
)

func TestReminderFileParsing(t *testing.T) {
	reminderFileContent := `Santander Platino;18
Promotions;13;true
	`

	const remindersTmpFilePath = "/tmp/rmnd.txt"

	_, err := parseRemindersFromFile("does_not_exist")
	if err == nil {
		t.Error("File not found, it should have failed.")
	}

	err = ioutil.WriteFile(remindersTmpFilePath, []byte(reminderFileContent), 0644)
	if err != nil {
		t.Fatal("Error generating test reminder file")
	}

	reminders, err := parseRemindersFromFile(remindersTmpFilePath)
	if err != nil {
		t.Fatal("error parsing reminder file")
	}

	expectedRemindersCount := 2
	if len(reminders) != expectedRemindersCount {
		t.Fatalf("error parsing reminder file, it should have had %d records, got=%d", expectedRemindersCount, len(reminders))
	}

	err = os.RemoveAll(remindersTmpFilePath)
	if err != nil {
		t.Errorf("unexpedted error: [%s]", err)
	}
}

func generateBaseDateTime() time.Time {
	now := time.Now()
	then := time.Date(2018, now.Month(), now.Day(), 0 /*the hour */, 0 /* the minutes */, 0, 0, time.UTC)
	return then
}

func TestReminderRecordParsing(t *testing.T) {

	_, err := extractReminderFromText("some record")
	if err == nil {
		t.Errorf("It should have failed while parsing ... ")
	}

	_, err = extractReminderFromText(";;")
	if err == nil {
		t.Errorf("Fields cannot be empty")
	}

	_, err = extractReminderFromText("record1;")
	if err == nil {
		t.Errorf("Second field cannot be empty")
	}

	_, err = extractReminderFromText("fiel1;hola")
	if err == nil {
		t.Errorf("Second fields must be numeric")
	}

	tests := []struct {
		input        string
		expectedDays int
		expectedName string
		notify       bool
	}{
		{`Santander Platino;18`, 18, "Santander Platino", false},
		{`Promotions;13;true`, 13, "Promotions", true},
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
		if reminder.Notify != tt.notify {
			t.Fatalf("got=%t as notify, expected=%t", reminder.Notify, tt.notify)
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

func TestShouldIgnoreEmptyLine(t *testing.T) {
	line := ""
	if !shouldIgnoreLineInFile(line) {
		t.Errorf("line is empty, it should be ignored.")
	}
	line = "# This is a comment"
	if !shouldIgnoreLineInFile(line) {
		t.Errorf("line is a comment, it should be ignored.")
	}
}

func TestExistsFileOrDirectory(t *testing.T) {
	if !existsFileOrDirectory("/dev/null") {
		t.Errorf("file should exist")
	}
}

func TestIsWeekend(t *testing.T) {
	// 05/19/2018 is Sunday
	d := time.Date(2018, time.May, 19, 0, 0, 0, 0, time.UTC)
	if !isWeekend(&d) {
		t.Errorf("Date should be identified as Weekend.")
	}
}

func TestFormatDate(t *testing.T) {
	d := time.Date(2018, time.May, 19, 0, 0, 0, 0, time.UTC)
	expectedFormattedDateResult := "2018/05/19 Saturday"
	result := formatDate(&d)
	if expectedFormattedDateResult != result {
		t.Errorf("got [%q], expected [%q]", result, expectedFormattedDateResult)
	}
}

func TestSortRemindersByDay(t *testing.T) {
	reminders := []Reminder{
		Reminder{
			Name:      "A",
			EveryWhen: 4,
		},
		Reminder{
			Name:      "B",
			EveryWhen: 2,
		},
		Reminder{
			Name:      "C",
			EveryWhen: 9,
		},
	}

	sortRemindersByDay(&reminders)

	for i := 0; i < len(reminders)-1; i++ {
		if reminders[i].EveryWhen < reminders[i+1].EveryWhen {
			t.Errorf("%d should come first than %d (%s) and (%s)",
				reminders[i].EveryWhen, reminders[i+1].EveryWhen, reminders[i], reminders[i+1])
		}
	}

}

func TestCreateOutputText(t *testing.T) {
	cmdArgs := [2]string{"-f", "term"}
	msg := "'Santander Platino' in 2 days (2019/05/28 Tuesday)"
	const expectedOutputHexString = "1b5b33383b323b3231393b3137313b3132316d2753616e74616e64657220506c6174696e6f2720696e203220646179732028323031392f30352f32382054756573646179290a1b5b306d"
	const expectedOutputHexStringWithPresenterCommandNotFound = "2753616e74616e64657220506c6174696e6f2720696e203220646179732028323031392f30352f3238205475657364617929"

	p := termenv.ColorProfile()
	colors := colorForMessages()
	config := Configuration{
		colorConfiguration: colors,
		termProfile:        p,
	}

	output := createOutputText(cmdArgs[:], msg, 2, warningRemainingDays, &config)
	output = strings.TrimSpace(output)
	output = strings.TrimSuffix(output, "\n")
	output = fmt.Sprintf("%x", output)

	if output != expectedOutputHexString {
		t.Errorf("got=[%x], expected=[%x]", output, expectedOutputHexString)
	}

	cmdArgs = [2]string{"-f", "foundNotFound"}
	output = createOutputText(cmdArgs[:], msg, 2, warningRemainingDays, &config)
	output = fmt.Sprintf("%x", output)
	if output != expectedOutputHexStringWithPresenterCommandNotFound {
		t.Errorf("got=[%s], want=[%s]", output, expectedOutputHexStringWithPresenterCommandNotFound)
	}
}

func TestNextReminderRecurrentDate(t *testing.T) {

	currentDate := time.Date(2019, 5, 26, 0, 0, 0, 0, time.UTC)
	everyWhen := 27

	next := nextReminderRecurrentDate(currentDate, everyWhen)
	if (currentDate.Month() + 1) <= next.Month() {
		t.Errorf("Date [%q] should be after [%q]", currentDate, next)
	}

	everyWhen = 25
	next = nextReminderRecurrentDate(currentDate, everyWhen)
	if next.Month() <= currentDate.Month() {
		t.Errorf("Date [%q] should be after [%q]", next, currentDate)
	}

	everyWhen = 26
	next = nextReminderRecurrentDate(currentDate, everyWhen)
	if next.Month() != currentDate.Month() {
		t.Errorf("Date [%q] should be equal [%q], %d = %d", next, currentDate, next.Month(), currentDate.Month())
	}

}

func TestCreateMessage(t *testing.T) {
	now := time.Date(2019, 5, 26, 0, 0, 0, 0, time.UTC)
	next := time.Date(2019, 5, 28, 0, 0, 0, 0, time.UTC)

	r := Reminder{Name: "Hello", EveryWhen: 28}
	msg, _ := createMessage(next, now, r)

	expectedMsg := "'Hello' in 2 days (2019/05/28 Tuesday)"
	if msg != expectedMsg {
		t.Errorf("got=[%s], want=[%s]", msg, expectedMsg)
	}

	next = time.Date(2019, 5, 26, 0, 0, 0, 0, time.UTC)
	expectedMsg = "'Hello' TODAY! (2019/05/26 Sunday)"
	msg, _ = createMessage(next, now, r)
	if msg != expectedMsg {
		t.Errorf("got=[%s], want=[%s]", msg, expectedMsg)
	}

	next = time.Date(2019, 5, 27, 0, 0, 0, 0, time.UTC)
	expectedMsg = "'Hello' in 1 day (2019/05/27 Monday)"
	msg, _ = createMessage(next, now, r)
	if msg != expectedMsg {
		t.Errorf("got=[%s], want=[%s]", msg, expectedMsg)
	}

	now = time.Date(2019, 5, 24, 0, 0, 0, 0, time.UTC)
	next = time.Date(2019, 5, 26, 0, 0, 0, 0, time.UTC)

	expectedMsg = "'Hello' in 2 days (WEEKEND) (2019/05/26 Sunday)"
	msg, _ = createMessage(next, now, r)
	if msg != expectedMsg {
		t.Errorf("got=[%s], want=[%s]", msg, expectedMsg)
	}
}

func TestGetRemindersFile(t *testing.T) {
	homeEnvBackup := os.Getenv("HOME")
	os.Setenv("HOME", "/tmp")
	_, err := getRemindersFile()
	if err == nil {
		t.Errorf("reminders file should be in home directory")
	}

	os.Setenv("HOME", homeEnvBackup)
	_, err = getRemindersFile()
	if err != nil {
		t.Errorf("Not able to get reminders file")
	}

}

func Test_withColor(t *testing.T) {
	type test struct {
		msg                                 string
		remainingDays, warningRemainingDays int
		config                              Configuration
		want                                string
	}

	p := termenv.ColorProfile()
	colors := colorForMessages()
	config := Configuration{
		colorConfiguration: colors,
		termProfile:        p,
	}

	tests := []test{
		{
			msg:                  "hola",
			remainingDays:        2,
			warningRemainingDays: 3,
			config:               config,
			want:                 "1b5b33383b323b3231393b3137313b3132316d686f6c611b5b306d",
		},

		{
			msg:                  "abcdef",
			remainingDays:        0,
			warningRemainingDays: 0,
			config:               config,
			want:                 "1b5b33383b323b3233323b3133313b3133366d6162636465661b5b306d",
		},

		{
			msg:                  "otro",
			remainingDays:        10,
			warningRemainingDays: 5,
			config:               config,
			want:                 "1b5b33383b323b3136383b3230343b3134306d6f74726f1b5b306d",
		},
	}

	for _, tt := range tests {
		got := withColor(tt.msg, tt.remainingDays, tt.warningRemainingDays, &tt.config)
		gotHex := fmt.Sprintf("%x", got)
		if gotHex != tt.want {
			t.Errorf("got=[hex: '%s'], want=[%s]", gotHex, tt.want)
		}
	}
}
