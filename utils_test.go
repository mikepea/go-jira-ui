package main

import (
	"encoding/json"
	"testing"
)

func TestCountLabelsFromQueryData(t *testing.T) {
	var data interface{}
	inputJSON := []byte(`{
		  "issues": [
				{ "fields": { "labels": [ "wibble", "bibble" ] } },
				{ "fields": { "labels": [ "wibble", "bibble" ] } },
				{ "fields": { "labels": [ "bibble" ] } }
			]
		}`)

	expected := make(map[string]int)
	expected["wibble"] = 2
	expected["bibble"] = 3

	err := json.Unmarshal(inputJSON, &data)
	if err != nil {
		t.Fatal(err)
	}

	actual := countLabelsFromQueryData(data)
	if expected["wibble"] != actual["wibble"] {
		t.Fatalf("wibble: expected %q, got %q", expected, actual)
	}
	if expected["bibble"] != actual["bibble"] {
		t.Fatalf("bibble: expected %q, got %q", expected, actual)
	}
}

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
		"      {code}   ",
		"      # This is code it should not be wrapped at all herpdy derp",
		"      # weoijwefoi wpeifjwoiejf pwjefoijwefij wefjowiejf wefwefwefijwe",
		"      {code}   ",
		"      {code:bash}   ",
		"      # This is code it should not be wrapped at all herpdy derp",
		"      # weoijwefoi wpeifjwoiejf pwjefoijwefij wefjowiejf wefwefwefijwe",
		"      {code}   ",
		"      {noformat}   ",
		"      # This is noformat it should not be wrapped at all herpdy derp",
		"      # weoijwefoi wpeifjwoiejf pwjefoijwefij wefjowiejf wefwefwefijwe",
		"      {noformat}   ",
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
		"      {code}   ",
		"      # This is code it should not be wrapped at all herpdy derp",
		"      # weoijwefoi wpeifjwoiejf pwjefoijwefij wefjowiejf wefwefwefijwe",
		"      {code}   ",
		"      {code:bash}   ",
		"      # This is code it should not be wrapped at all herpdy derp",
		"      # weoijwefoi wpeifjwoiejf pwjefoijwefij wefjowiejf wefwefwefijwe",
		"      {code}   ",
		"      {noformat}   ",
		"      # This is noformat it should not be wrapped at all herpdy derp",
		"      # weoijwefoi wpeifjwoiejf pwjefoijwefij wefjowiejf wefwefwefijwe",
		"      {noformat}   ",
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
