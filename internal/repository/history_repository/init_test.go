package history_repository

import (
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestHistoryQueriesBindNamedParameters(t *testing.T) {
	for _, query := range []string{queryGetByUserID, queryGetRecentActivities} {
		bound, _, err := sqlx.Named(query, map[string]any{"user_id": 1, "limit": 10, "offset": 0})
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(bound, ":") || !strings.Contains(bound, "CAST(ARRAY[] AS text[])") {
			t.Fatalf("invalid named query: %s", bound)
		}
	}
}
