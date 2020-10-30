package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/muesli/termenv"
	"github.com/spf13/viper"
)

func run(envConfig *viper.Viper) error {

	remindersFile := getRemindersFilePath(path.Join(os.Getenv("HOME"), envConfig.GetString("reminders_directory")))
	if !existsFileOrDirectory(remindersFile) {
		return fmt.Errorf(`error: '%s' does not exist

Create the file '%s' with a content such as the following:

Cancel UFCFightPass;29;true
Pagar Internet;7;true`, remindersFile, remindersFile)
	}

	notifsDir := path.Join(os.Getenv("HOME"), shellReminderMainDirectory, notificationsDirectory)
	err := createDirectory(notifsDir)
	if err != nil {
		return err
	}

	reminders, err := parseRemindersFromFile(remindersFile)
	if err != nil {
		return err
	}

	p := termenv.ColorProfile()
	colors := colorForMessages()

	config := Configuration{
		colorConfiguration: colors,
		termProfile:        p,
	}

	sortRemindersByDay(&reminders)

	now := time.Now()
	for _, r := range reminders {

		msg := ""
		next := nextReminderRecurrentDate(now, r.EveryWhen)
		msg, remainingDays := createMessage(next, now, r)

		if len(msg) != 0 {
			fmt.Println(createOutputText(cmdArgs[:], msg, remainingDays, warningRemainingDays, &config))
		}

		if r.Notify && remainingDays == 0 {
			hash := buildHash(r.Name)
			notifHashFilePath := filepath.Join(notifsDir, hash)
			if !exists(notifHashFilePath) {
				err = notifySMS(msg, &r, envConfig)
				if err != nil {
					log.Fatal(err)
				}
				err = notifyEmail(msg, &r, envConfig)
				if err != nil {
					fmt.Println(err)
				}
				err = ioutil.WriteFile(notifHashFilePath, []byte(r.Name), 0644)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func main() {

	envConfig, err := readConfig("shellreminders.env", os.Getenv("HOME"), map[string]interface{}{
		"api_key":             os.Getenv("NEXMO_API_KEY"),
		"api_secret":          os.Getenv("NEXMO_API_SECRET"),
		"to_phone":            os.Getenv("NOTIFY_PHONE"),
		"sendgrid_api_key":    os.Getenv("SENDGRID_API_KEY"),
		"reminders_directory": shellReminderMainDirectory,
		"email_to":            "",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if err := run(envConfig); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}
