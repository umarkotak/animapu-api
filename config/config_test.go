package config

import "testing"

func TestIsAdminEmail(t *testing.T) {
	config = Config{AdminEmails: []string{"admin@example.com", "other@example.com"}}
	if !IsAdminEmail("ADMIN@example.com") || IsAdminEmail("user@example.com") || IsAdminEmail("") {
		t.Fatal("admin email validation failed")
	}
}
