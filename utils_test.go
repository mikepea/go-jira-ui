package main

import (
	"testing"
)

func TestFindTicketIdInString(t *testing.T) {
	var match string
	match = findTicketIdInString("  relates: BLAH-123[Done]  ")
	if match != "BLAH-123" {
		t.Fatalf("expected BLAH-123, got %s", match)
	}
	match = findTicketIdInString("  wibble: xxBLAH-123[Done]  ")
	if match != "" {
		t.Fatalf("expected %q, got %q", "", match)
	}
}

func TestWrapText(t *testing.T) {
	input := []string{
		"",
		"wibble:        hello",
		"longfield:     1234567890123456789012345678901234567890",
		"1234567890123456789012345678901234567890",
		"12345678901234567890123456789012345678901234567890",
		"body: |",
		"   hello there I am a line that is longer than 40 chars yes I am oh aye.",
	}
	expected := []string{
		"",
		"wibble:        hello",
		"longfield:     1234567890123456789012345678901234567890",
		"1234567890123456789012345678901234567890",
		"1234567890123456789012345678901234567890",
		"1234567890",
		"body: |",
		"   hello there I am a line that is longe",
		"r than 40 chars yes I am oh aye.",
	}
	match := WrapText(input, 40)
	for i, _ := range expected {
		if i > len(match)-1 {
			t.Fatalf("expected %d lines, got %d", len(expected), len(match))
		} else if match[i] != expected[i] {
			t.Fatalf("line %d - expected %q, got %q", i, expected[i], match[i])
		}
	}
}
