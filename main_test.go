package main

import (
	"testing"
)

func TestGetTicketIdFromListLine(t *testing.T) {
	type testCase struct {
		in       string
		expected string
	}

	var testCases = []testCase{
		{"BLAH-1234 wibble wobble", "BLAH-1234"},
		{"OPS-1     wibble wobble", "OPS-1"},
	}

	for _, v := range testCases {
		actual := getTicketIdFromListLine(v.in)
		if v.expected != actual {
			t.Fatalf("expected %q, actual %q", v.expected, actual)
		}
	}
}
