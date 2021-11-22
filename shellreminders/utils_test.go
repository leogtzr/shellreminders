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
		want         string
	}

	tests := []test{
		{
			reminderName: "him",
			want:         "1272e1c9194cfd5797c342621ba3f2fb",
		},
		{
			reminderName: "her",
			want:         "58ea1b1fe5e99cc4e8f62aaed7fdc10b",
		},
	}

	for _, tt := range tests {
		if got := buildHash(tt.reminderName, time.Now()); got != tt.want {
			t.Errorf("got=[%s], expected=[%s] for '%s' reminder name", got, tt.want, tt.reminderName)
		}
	}
}
