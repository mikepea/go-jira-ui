package main

import (
	"testing"
)

func TestTicketIdSetting(t *testing.T) {
	sp := new(TicketShowPage)
	if sp.TicketId != "" {
		t.Fatalf("sp.TicketId: expected %q, got %q", "", sp.TicketId)
	}

	sp.TicketId = "ABC-123"
	if sp.TicketId != "ABC-123" {
		t.Fatalf("sp.TicketId: expected %q, got %q", "ABC-123", sp.TicketId)
	}

	if sp.Id() != "ABC-123" {
		t.Fatalf("sp.Id: expected %q, got %q", "ABC-123", sp.Id())
	}

}
