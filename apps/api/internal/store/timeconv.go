package store

import (
	"fmt"
	"time"
)

// rfc3339Milli is the exact format every timestamp column uses (R2): RFC
// 3339, UTC, millisecond precision, literal "Z" offset.
const rfc3339Milli = "2006-01-02T15:04:05.000Z"

// parseStoreTime converts a column's TEXT timestamp back into time.Time —
// the conversion lives at the store's edge; RFC 3339 TEXT never leaves this
// package (data-access-entity-framework, per the confirmed contract).
func parseStoreTime(s string) (time.Time, error) {
	t, err := time.Parse(rfc3339Milli, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("store: parse time %q: %w", s, err)
	}
	return t, nil
}

// nowStoreTime stamps a column that has no schema-level default (last_seen_at
// on project_paths) with the same convention every other timestamp uses.
func nowStoreTime() string {
	return time.Now().UTC().Format(rfc3339Milli)
}
