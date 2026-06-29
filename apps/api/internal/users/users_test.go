package users

import (
	"testing"
	"time"
)

func TestISOTimeMatchesJSToISOString(t *testing.T) {
	ts := time.Date(2026, time.January, 2, 3, 4, 5, 123_000_000, time.UTC)
	got := ISOTime(ts)
	want := "2026-01-02T03:04:05.123Z"
	if got != want {
		t.Fatalf("ISOTime = %q, want %q", got, want)
	}
}

func TestISOTimeConvertsToUTC(t *testing.T) {
	loc := time.FixedZone("UTC+2", 2*60*60)
	ts := time.Date(2026, time.January, 1, 2, 0, 0, 0, loc)
	got := ISOTime(ts)
	want := "2026-01-01T00:00:00.000Z"
	if got != want {
		t.Fatalf("ISOTime = %q, want %q", got, want)
	}
}
