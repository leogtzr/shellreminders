package main

import "github.com/muesli/termenv"

// Reminder ...
type Reminder struct {
	Name      string
	EveryWhen int
}

// Configuration ...
type Configuration struct {
	termProfile        termenv.Profile
	colorConfiguration map[string]string
}
