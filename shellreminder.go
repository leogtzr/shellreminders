package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Reminder ...
type Reminder struct {
	Name      string
	EveryWhen int
	Type      ReminderType
}

// ReminderType ...
type ReminderType int

const (
	// RecurrentReminder ...
	RecurrentReminder ReminderType = 0
	// Counter ...
	Counter ReminderType = 1

	// ShellReminderMainDirectory ...
	ShellReminderMainDirectory = "/.shellreminder"
	// RemindersFile ...
	RemindersFile            = ShellReminderMainDirectory + "/reminders"
	minNumberOfRecordsInFile = 2

	shellPresenterCommand = "toiletx"

	minDaysToShowInReminders = -3
)

func existsFileOrDirectory(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func extractReminderFromText(text string) (Reminder, error) {
	if !strings.Contains(text, ";") {
		return Reminder{}, fmt.Errorf("[%s] with wrong format", text)
	}
	records := strings.Split(strings.TrimSpace(text), ";")

	if len(records) < minNumberOfRecordsInFile {
		return Reminder{}, fmt.Errorf("not enough records in row ----> [%s]", text)
	}

	name := records[0]
	if len(strings.TrimSpace(name)) == 0 {
		return Reminder{}, fmt.Errorf("reminder name cannot be empty")
	}
	when := records[1]
	if len(strings.TrimSpace(when)) == 0 {
		return Reminder{}, fmt.Errorf("reminder time cannot be empty")
	}

	w, err := strconv.Atoi(when)
	if err != nil {
		return Reminder{}, fmt.Errorf("second record should be numeric")
	}

	reminderType := RecurrentReminder
	if len(records) > minNumberOfRecordsInFile {
		if strings.TrimSpace(strings.ToLower(records[2])) == "counter" {
			reminderType = Counter
		} else {
			return Reminder{}, fmt.Errorf("counter is the only explicit reminder type allowed for now [%s]", records[2])
		}
	}
	return Reminder{Name: name, EveryWhen: w, Type: reminderType}, nil
}

func shouldIgnoreLineInFile(line string) bool {
	return len(line) == 0 || strings.HasPrefix(line, "#")
}

func parseRemindersFromFile(filePath string) ([]Reminder, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	input := bufio.NewScanner(f)
	reminders := make([]Reminder, 0)

	for input.Scan() {
		line := strings.TrimSpace(input.Text())
		if shouldIgnoreLineInFile(line) {
			continue
		}
		reminder, err := extractReminderFromText(line)
		if err != nil {
			return []Reminder{}, err
		}
		reminders = append(reminders, reminder)
	}

	return reminders, nil
}

func (r Reminder) String() string {
	var out bytes.Buffer
	out.WriteString("'")
	out.WriteString(r.Name)
	out.WriteString("'")
	out.WriteString(" every ")
	out.WriteString(fmt.Sprintf("%d", r.EveryWhen))
	out.WriteString(" days. ")

	if r.Type == RecurrentReminder {
		out.WriteString("[recurrent reminder]")
	} else {
		out.WriteString("[counter]")
	}

	return out.String()
}

func isWeekend(d time.Time) bool {
	return d.Weekday() == time.Saturday || d.Weekday() == time.Sunday
}

func buildTime(r *Reminder) time.Time {
	now := time.Now()
	d := time.Date(now.Year(), now.Month(), r.EveryWhen, 0, 0, 0, 0, time.UTC)
	return d
}

func simpleDateFormat(t *time.Time) string {
	return fmt.Sprintf("%d/%d/%d", t.Year(), t.Month(), t.Day())
}

func buildReminderMessage(reminderName string, remainingDays int, r *Reminder) string {
	var out bytes.Buffer

	if remainingDays == 0 {
		out.WriteString(fmt.Sprintf("Remaining days for '%s' : TODAY!", r.Name))
	} else {
		out.WriteString(fmt.Sprintf("Remaining days for '%s' :%d ", r.Name, remainingDays))
	}

	d := buildTime(r)
	if isWeekend(d) {
		out.WriteString("| WARNING! ")
		out.WriteString(d.Weekday().String())
		out.WriteString(" (")
		out.WriteString(simpleDateFormat(&d))
		out.WriteString(")")
	}
	return out.String()
}

func sortRemindersByDay(reminders *[]Reminder) {
	sort.Slice(*reminders,
		func(i, j int) bool {
			return (*reminders)[i].EveryWhen > (*reminders)[j].EveryWhen
		},
	)
}

func main() {

	// Check if the .shellreminder directory exists ...
	if !existsFileOrDirectory(os.Getenv("HOME") + ShellReminderMainDirectory) {
		fmt.Fprintf(os.Stderr, "%s does not exists\n", os.Getenv("HOME")+ShellReminderMainDirectory)
		os.Exit(1)
	}

	if !existsFileOrDirectory(os.Getenv("HOME") + RemindersFile) {
		fmt.Fprintf(os.Stderr, "%s file does not exists", (os.Getenv("HOME") + RemindersFile))
		os.Exit(1)
	}

	reminders, err := parseRemindersFromFile(os.Getenv("HOME") + RemindersFile)
	if err != nil {
		panic(err)
	}

	today := time.Now()

	sortRemindersByDay(&reminders)

	for _, r := range reminders {
		remainingDays := r.EveryWhen - today.Day()

		if remainingDays < minDaysToShowInReminders {
			continue
		}

		msg := buildReminderMessage(r.Name, remainingDays, &r)
		cmdArgs := []string{"-f", "term", "-F", "border", msg}

		if cmdOut, err := exec.Command(shellPresenterCommand, cmdArgs...).Output(); err != nil {
			fmt.Println(msg)
		} else {
			fmt.Print(string(cmdOut))
		}
	}

}
