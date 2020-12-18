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
		{
			reminderName: "hola",
			want:         "ff6f3a046811c67b2c3e0b6ebd534a17",
		},
		{
			reminderName: "test1",
			want:         "3d5d13dafdf63b89659bc8deae6996c5",
		},
	}

	for _, tt := range tests {
		if got := buildHash(tt.reminderName); got != tt.want {
			t.Errorf("got=[%s], expected=[%s] for '%s' reminder name", got, tt.want, tt.reminderName)
		}
	}
}
