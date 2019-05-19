package main

import (
	"bufio"
	"bytes"
	"fmt"
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

func existsFileOrDirectory(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func extractReminderFromText(text string) (Reminder, error) {
	if !strings.Contains(text, recordFileSeparator) {
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
	return fmt.Sprintf("%d/%02d/%02d", t.Year(), t.Month(), t.Day())
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

	cmdArgs := []string{"-f", "term", "-F", "border"}

	sortRemindersByDay(&reminders)

	now := time.Now()
	for _, r := range reminders {

		msg := ""
		next := now
		if now.Day() == r.EveryWhen {
			next = now
		} else if now.Day() > r.EveryWhen {
			next = time.Date(now.Year(), now.Month()+1, r.EveryWhen, 0, 0, 0, 0, time.UTC)
		} else if now.Day() < r.EveryWhen {
			next = time.Date(now.Year(), now.Month(), r.EveryWhen, 0, 0, 0, 0, time.UTC)
		} else {
			continue
		}

		remainingDays := int(next.Sub(now).Hours() / 24)
		if int(remainingDays) == 0 {
			msg = fmt.Sprintf("'%s' TODAY! (%s)", r.Name, formatDate(&now))
		} else if remainingDays < lessThanDays {
			if isWeekend(&next) {
				msg = fmt.Sprintf("'%s' in less than %d days (WEEKEND) (%s)", r.Name, int(remainingDays), formatDate(&next))
			} else {
				msg = fmt.Sprintf("'%s' in less than %d days (%s)", r.Name, int(remainingDays), formatDate(&next))
			}
		} else {
			continue
		}

		if cmdOut, err := exec.Command(shellPresenterCommand, append(cmdArgs, msg)...).Output(); err != nil {
			fmt.Println(msg)
		} else {
			fmt.Print(string(cmdOut))
		}
	}

}
