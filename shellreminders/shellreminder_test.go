package shellreminders

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
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
		{
			Name:      "A",
			EveryWhen: 4,
		},
		{
			Name:      "B",
			EveryWhen: 2,
		},
		{
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
	type testCase struct {
		baseDir string
		want    string
	}

	tests := []testCase{
		{
			baseDir: "/home/leo",
			want:    "/home/leo/reminders",
		},
		{
			baseDir: "/tmp",
			want:    "/tmp/reminders",
		},
	}

	for _, tc := range tests {
		if got := getRemindersFilePath(tc.baseDir); got != tc.want {
			t.Errorf("got=[%s], want=[%s]", got, tc.want)
		}
	}
}
