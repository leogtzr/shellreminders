package main

import "github.com/muesli/termenv"

// Reminder ...
type Reminder struct {
	Name      string
	EveryWhen int
	Notify    bool
}

// ColorConfiguration ...
type ColorConfiguration struct {
	termProfile        termenv.Profile
	colorConfiguration map[string]string
}
