package main

import "testing"

func TestNewAppUsesConfiguredName(t *testing.T) {
	app := NewApp(Config{Name: "CellPilot"})
	if got := app.Name(); got != "CellPilot" {
		t.Fatalf("Name() = %q, want CellPilot", got)
	}
}
