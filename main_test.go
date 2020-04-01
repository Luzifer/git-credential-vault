package main

import (
	"strings"
	"testing"
)

func TestReadInput(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		"protocol=https",
		"host=github.com",
		"", // Defined as end-of-input
		"shouldnt=bethere",
	}, "\n"))

	vars, err := readInput(input)
	if err != nil {
		t.Fatalf("Read input errored: %s", err)
	}

	if l := len(vars); l != 2 {
		t.Errorf("Unexpected number of values: %d != 2", l)
	}

	if h := vars["host"]; h != "github.com" {
		t.Errorf("Unexpected value %s: %s != %s", "host", h, "github.com")
	}

	if p := vars["protocol"]; p != "https" {
		t.Errorf("Unexpected value %s: %s != %s", "protocol", p, "https")
	}
}
