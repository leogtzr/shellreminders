package main

import (
	"bufio"
	"bytes"
	"errors"
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
	lessThanDays               = 8
	recordFileSeparator        = ";"
)

var cmdArgs = [6]string{"-f", "smblock", "-w", "900", "-F", "border"}

func existsFileOrDirectory(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

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

func getRemindersFile() (string, error) {
	remindersDir := path.Join(os.Getenv("HOME"), shellReminderMainDirectory)
	if !existsFileOrDirectory(remindersDir) {
		return "", fmt.Errorf("%s does not exists", remindersDir)
	}

	remindersFile := path.Join(remindersDir, "reminders")
	if !existsFileOrDirectory(remindersFile) {
		return "", fmt.Errorf("%s file does not exists", remindersFile)
	}
	return remindersFile, nil
}

func main() {

	remindersFile, err := getRemindersFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
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
	remainingDays := daysBetween(next, now)
	if remainingDays == 0 {
		msg = fmt.Sprintf("'%s' TODAY! (%s)", r.Name, formatDate(&now))
	} else if remainingDays < lessThanDays {
		if isWeekend(&next) {
			if remainingDays == 1 {
				msg = fmt.Sprintf("'%s' in %d day (WEEKEND) (%s)", r.Name, remainingDays, formatDate(&next))
			} else {
				msg = fmt.Sprintf("'%s' in %d days (WEEKEND) (%s)", r.Name, remainingDays, formatDate(&next))
			}
		} else {
			if remainingDays == 1 {
				msg = fmt.Sprintf("'%s' in %d day (%s)", r.Name, remainingDays, formatDate(&next))
			} else {
				msg = fmt.Sprintf("'%s' in %d days (%s)", r.Name, remainingDays, formatDate(&next))
			}
		}
	}

	return msg
}

func daysBetween(a, b time.Time) int {

	// fmt.Println("a", a.YearDay())
	// fmt.Println("b", b.YearDay())

	// if a.After(b) {
	// 	a, b = b, a
	// }

	// days := -a.YearDay()
	// for year := a.Year(); year < b.Year(); year++ {
	// 	days += time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()
	// }
	// days += b.YearDay()
	// return days
	return a.YearDay() - b.YearDay()
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
