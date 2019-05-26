package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Reminder ...
type Reminder struct {
	Name      string
	EveryWhen int
}

const (
	shellReminderMainDirectory = ".shellreminder"
	minNumberOfRecordsInFile   = 2
	shellPresenterCommand      = "toilet"
	minimumDaysAgo             = 2
	lessThanDays               = 7
	recordFileSeparator        = ";"
)

var cmdArgs = [6]string{"-f", "smblock", "-w", "900", "-F", "border"}

func existsFileOrDirectory(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// type errNotEnoughRecordsx struct{ error }

func extractReminderFromText(text string) (Reminder, error) {
	if !strings.Contains(text, recordFileSeparator) {
		return Reminder{}, fmt.Errorf("[%s] with wrong format", text)
	}
	records := strings.Split(strings.TrimSpace(text), ";")

	name := records[0]
	if len(strings.TrimSpace(name)) == 0 {
		return Reminder{}, errors.New("not enough records in row, field1")
	}
	when := records[1]
	if len(strings.TrimSpace(when)) == 0 {
		return Reminder{}, errors.New("not enough records in row, field2")
	}

	w, err := strconv.Atoi(when)
	if err != nil {
		return Reminder{}, errors.New("not enough records in row")
	}

	return Reminder{Name: name, EveryWhen: w}, nil
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
	out.WriteString(" day ")
	out.WriteString(fmt.Sprintf("%d", r.EveryWhen))
	out.WriteString(" of each month")

	return out.String()
}

func isWeekend(d *time.Time) bool {
	return d.Weekday() == time.Saturday || d.Weekday() == time.Sunday
}

func formatDate(t *time.Time) string {
	return fmt.Sprintf("%d/%02d/%02d %s", t.Year(), t.Month(), t.Day(), t.Weekday())
}

func sortRemindersByDay(reminders *[]Reminder) {
	sort.Slice(*reminders,
		func(i, j int) bool {
			return (*reminders)[i].EveryWhen > (*reminders)[j].EveryWhen
		},
	)
}

func main() {

	remindersDir := path.Join(os.Getenv("HOME"), shellReminderMainDirectory)
	if !existsFileOrDirectory(remindersDir) {
		fmt.Fprintf(os.Stderr, "%s does not exists\n", remindersDir)
		os.Exit(1)
	}

	remindersFile := path.Join(remindersDir, "reminders")
	if !existsFileOrDirectory(remindersFile) {
		fmt.Fprintf(os.Stderr, "%s file does not exists\n", remindersFile)
		os.Exit(1)
	}

	reminders, err := parseRemindersFromFile(remindersFile)
	if err != nil {
		panic(err)
	}

	sortRemindersByDay(&reminders)

	now := time.Now()
	for _, r := range reminders {

		msg := ""
		next := nextReminderRecurrentDate(now, r.EveryWhen)
		msg = createMessage(next, now, r)

		if len(msg) != 0 {
			fmt.Println(createOutputText(cmdArgs[:], msg))
		}
	}

}

func createMessage(next, now time.Time, r Reminder) string {
	msg := ""
	remainingDays := int(math.Ceil(next.Sub(now).Hours() / 24.0))
	if remainingDays == 0 {
		msg = fmt.Sprintf("'%s' TODAY! (%s)", r.Name, formatDate(&now))
	} else if remainingDays < lessThanDays {
		if isWeekend(&next) {
			msg = fmt.Sprintf("'%s' in %d days (WEEKEND) (%s)", r.Name, remainingDays, formatDate(&next))
		} else {
			msg = fmt.Sprintf("'%s' in %d days (%s)", r.Name, remainingDays, formatDate(&next))
		}
	}

	return msg
}

func nextReminderRecurrentDate(currentDate time.Time, everyWhen int) time.Time {
	next := currentDate
	if currentDate.Day() == everyWhen {
		next = currentDate
	} else if currentDate.Day() > everyWhen {
		next = time.Date(currentDate.Year(), currentDate.Month()+1, everyWhen, 0, 0, 0, 0, time.UTC)
	} else if currentDate.Day() < everyWhen {
		next = time.Date(currentDate.Year(), currentDate.Month(), everyWhen, 0, 0, 0, 0, time.UTC)
	}
	return next
}

func createOutputText(cmdArgs []string, msg string) string {
	cmd := exec.Command(shellPresenterCommand, append(cmdArgs[:], msg)...)
	cmdOut, err := cmd.Output()
	if err != nil {
		return msg
	}
	return string(cmdOut)
}
