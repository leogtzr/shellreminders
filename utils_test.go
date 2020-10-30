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
			want:         "a04369b08863c07411df0c19a9396f2a",
		},
		test{
			reminderName: "test1",
			want:         "c3fb3ec91f9ed1471864731aa4dc5fa0",
		},
	}

	for _, tt := range tests {
		if got := buildHash(tt.reminderName); got != tt.want {
			t.Errorf("got=[%s], expected=[%s] for '%s' reminder name", got, tt.want, tt.reminderName)
		}
	}
}
