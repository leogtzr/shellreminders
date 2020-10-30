package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/muesli/termenv"
	"github.com/nexmo-community/nexmo-go"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"
)

func colorForMessages() map[string]string {
	return map[string]string{
		"red":    redHexColor,
		"yellow": yellowHexColor,
		"green":  greenHexColor,
	}
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

func createOutputText(
	cmdArgs []string, msg string, remainingDays, warningRemainingDays int, config *ColorConfiguration) string {
	cmd := exec.Command(shellPresenterCommand, append(cmdArgs[:], msg)...)
	cmdOut, err := cmd.Output()
	if err != nil {
		return msg
	}
	return withColor(string(cmdOut), remainingDays, warningRemainingDays, config)
}

func withColor(msg string, remainingDays, warningRemainingDays int, config *ColorConfiguration) string {
	if (remainingDays <= warningRemainingDays) && (remainingDays > 0) {
		return termenv.String(msg).Foreground(config.termProfile.Color(yellowHexColor)).String()
	} else if remainingDays == 0 {
		return termenv.String(msg).Foreground(config.termProfile.Color(redHexColor)).String()
	}
	return termenv.String(msg).Foreground(config.termProfile.Color(greenHexColor)).String()
}

func createMessage(next, now time.Time, r Reminder) (string, int) {
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

	return msg, remainingDays
}

func daysBetween(a, b time.Time) int {
	return a.YearDay() - b.YearDay()
}

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

	notify := false
	if len(records) > minNumberOfRecordsInFile {
		notify = true
	}

	w, err := strconv.Atoi(when)
	if err != nil {
		return Reminder{}, errors.New("not enough records in row")
	}

	return Reminder{Name: name, EveryWhen: w, Notify: notify}, nil
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

func getRemindersFilePath(remindersDirectory string) string {
	remindersFile := path.Join(remindersDirectory, "reminders")
	return remindersFile
}

func readConfig(filename, configPath string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath(configPath)
	v.SetConfigType("env")
	err := v.ReadInConfig()
	return v, err
}

func notifySMS(msg string, r *Reminder, envConfig *viper.Viper) error {
	auth := nexmo.NewAuthSet()
	auth.SetAPISecret(envConfig.GetString("api_key"), envConfig.GetString("api_secret"))

	client := nexmo.NewClient(http.DefaultClient, auth)

	smsContent := nexmo.SendSMSRequest{
		From: "447700900004",
		To:   envConfig.GetString("to_phone"),
		Text: msg,
	}

	_, _, err := client.SMS.SendSMS(smsContent)
	return err
}

func notifyEmail(msg string, r *Reminder, envConfig *viper.Viper) error {
	from := mail.NewEmail("Leonidas", "leonidas@root.com")
	subject := fmt.Sprintf("REMINDER -> %s", msg)
	to := mail.NewEmail("Leo Gtz", envConfig.GetString("email_to"))

	message := mail.NewSingleEmail(from, subject, to, msg, msg)
	client := sendgrid.NewSendClient(envConfig.GetString("sendgrid_api_key"))
	_, err := client.Send(message)
	return err
}

func buildHash(reminderName string) string {
	today := time.Now()
	text := fmt.Sprintf("%s%d%s%d", reminderName, today.Day(), today.Month(), today.Year())
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func dirExists(dirPath string) bool {
	if _, err := os.Stat(dirPath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

func createDirectory(dirPath string) error {
	if !dirExists(dirPath) {
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func exists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}

func getColorConfig() ColorConfiguration {
	p := termenv.ColorProfile()
	colors := colorForMessages()

	colorConfig := ColorConfiguration{
		colorConfiguration: colors,
		termProfile:        p,
	}

	return colorConfig
}

func notify(msg string, r *Reminder, envConfig *viper.Viper) error {
	err := notifySMS(msg, r, envConfig)
	if err != nil {
		return err
	}
	err = notifyEmail(msg, r, envConfig)
	if err != nil {
		return err
	}
	return nil
}
