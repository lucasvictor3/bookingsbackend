package main

import "testing"

func TestRun(t *testing.T) {
	// temporary fix
	_, err := run()

	if err != nil {
		t.Error("failed run()")
	}
}
