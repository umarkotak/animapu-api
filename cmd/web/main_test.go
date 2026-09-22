package main

import "testing"

func TestIsMigrateUp(t *testing.T) {
	if !isMigrateUp([]string{"animapu-api", "migrate", "up"}) || isMigrateUp([]string{"animapu-api"}) {
		t.Fatal("migration command detection failed")
	}
}
