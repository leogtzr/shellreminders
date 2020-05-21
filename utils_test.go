package main

import (
	"testing"
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
		test{
			reminderName: "hola",
			want:         "77a69f8b593ff373140666fe1faddcb7",
		},
		test{
			reminderName: "test1",
			want:         "0ebec8ca1bbe2cf1922446edbf9eb9ba",
		},
	}

	for _, tt := range tests {
		if got := buildHash(tt.reminderName); got != tt.want {
			t.Errorf("got=[%s], expected=[%s]", got, tt.want)
		}
	}
}
