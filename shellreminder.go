package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/muesli/termenv"
)

func main() {

	/*
		api_key=37c08302
		api_secret=JIaVw2MowOPfQkX2
	*/
	envConfig, err := readConfig("shellreminders.env", os.Getenv("HOME"), map[string]interface{}{
		"api_key":    os.Getenv("NEXMO_API_KEY"),
		"api_secret": os.Getenv("NEXMO_API_SECRET"),
		"to_phone":   os.Getenv("NOTIFY_PHONE"),
	})

	remindersFile, err := getRemindersFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	reminders, err := parseRemindersFromFile(remindersFile)
	if err != nil {
		panic(err)
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
			// send sms ...
			err = notify(msg, &r, envConfig)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}
