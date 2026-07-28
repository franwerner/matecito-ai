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
	return formatStoreTime(time.Now())
}

// formatStoreTime renders t in the same TEXT convention every timestamp
// column uses (R2) — the inverse of parseStoreTime. Used where the caller
// supplies the instant itself (events.occurred_at can be an earlier, caller-
// given time — R6's late-arriving event) instead of "now".
func formatStoreTime(t time.Time) string {
	return t.UTC().Format(rfc3339Milli)
}
