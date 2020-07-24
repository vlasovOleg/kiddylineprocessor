package sqlstore

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

// TestDB ...
func TestDB(t *testing.T) (*sql.DB, func(...string)) {
	t.Helper()
	databaseURL := "host=localhost dbname=kiddylineprocessor_test sslmode=disable"
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	db.Ping()

	return db, func(tables ...string) {
		if len(tables) > 0 {
			if _, err := db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", "))); err != nil {
				t.Fatal(err)
			}
		}

		db.Close()
	}
}
