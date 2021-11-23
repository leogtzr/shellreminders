package shellreminders

import (
	"testing"
	"time"
)

func Test_colorForMessages(t *testing.T) {
	const expectedNumberOfColors = 3
	expectedColorsInMap := [3]string{
		"red",
		"yellow",
		"green",
	}
	colors := colorForMessages()
	if len(colors) != expectedNumberOfColors {
		t.Errorf("got=[%d], want=[%d]", len(colors), expectedNumberOfColors)
	}

	for _, color := range expectedColorsInMap {
		if _, ok := colors[color]; !ok {
			t.Errorf("color '%s' does not exist in map", color)
		}
	}
}

func Test_buildHash(t *testing.T) {
	type test struct {
		reminderName string
		time         time.Time
		want         string
	}

	tests := []test{
		{
			reminderName: "him",
			time:         time.Date(2021, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
			want:         "26b104add45f3e3c597bbb5050456c52",
		},
		{
			reminderName: "her",
			time:         time.Date(2021, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
			want:         "7fc4a9fe041118c9e4a237a36df8ada9",
		},
	}

	/*
		Since the buildHasH() uses a time/date and we can't use time.Now(), we will always pass
		a time we know...
	*/
	for _, tt := range tests {

		if got := buildHash(tt.reminderName, tt.time); got != tt.want {
			t.Errorf("got=[%s], expected=[%s] for '%s' reminder name", got, tt.want, tt.reminderName)
		}
	}
}
