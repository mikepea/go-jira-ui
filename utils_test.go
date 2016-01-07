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
