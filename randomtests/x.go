package main

import (
	"fmt"
	"time"
)

func isWeekend(d time.Time) bool {
	return d.Weekday() == time.Saturday || d.Weekday() == time.Sunday
}

func main() {
	now := time.Now()
	d := time.Date(now.Year(), now.Month(), 17, 0 /*the hour */, 0 /* the minutes */, 0, 0, time.UTC)
	fmt.Println(d)
	fmt.Println(d.Weekday())
	fmt.Println(isWeekend(d))

	fmt.Println("~~~~~~~~~~~~>")
	fmt.Println(time.Saturday)
}
