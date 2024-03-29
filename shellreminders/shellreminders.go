package shellreminders

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type appEnv struct {
	envConfig *viper.Viper
}

func CLI(args []string) int {
	var app appEnv

	err := app.fromArgs(args)
	if err != nil {
		return errorBuildingAppConfigFromArgs
	}

	if err = app.run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)

		return errorRunningApplication
	}

	return 0
}

func (app *appEnv) run() error {
	remindersFile := getRemindersFilePath(path.Join(os.Getenv("HOME"), app.envConfig.GetString("reminders_directory")))
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

	colorConfig := getColorConfig()

	sortRemindersByDay(&reminders)

	now := time.Now()

	for _, r := range reminders {
		next := nextReminderRecurrentDate(now, r.EveryWhen)
		msg, remainingDays := createMessage(next, now, r)

		if len(msg) != 0 && remainingDays > 0 {
			fmt.Println(createOutputText(cmdArgs[:], msg, remainingDays, warningRemainingDays, &colorConfig))
		}

		if r.Notify && remainingDays == 0 {
			hash := buildHash(r.Name, time.Now())
			notifHashFilePath := filepath.Join(notifsDir, hash)

			if !exists(notifHashFilePath) {
				err = notify(msg, app.envConfig)
				if err != nil {
					fmt.Println(err)
				}

				err = ioutil.WriteFile(notifHashFilePath, []byte(r.Name), 0600)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (app *appEnv) fromArgs(args []string) error {
	envConfig, err := readConfig("shellreminders.env", os.Getenv("HOME"), map[string]interface{}{
		"api_key":             os.Getenv("NEXMO_API_KEY"),
		"api_secret":          os.Getenv("NEXMO_API_SECRET"),
		"to_phone":            os.Getenv("NOTIFY_PHONE"),
		"sendgrid_api_key":    os.Getenv("SENDGRID_API_KEY"),
		"reminders_directory": shellReminderMainDirectory,
		"email_to":            "",
	})
	if err != nil {
		return err
	}
	app.envConfig = envConfig

	return nil
}
