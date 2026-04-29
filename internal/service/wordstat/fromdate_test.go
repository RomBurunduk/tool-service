package wordstat

import (
	"testing"
	"time"
)

func TestFromDateLastThreeMonths(t *testing.T) {
	// Wednesday 2026-04-28 UTC; anchor 2026-01-28 (Wed); Monday on or before = 2026-01-26
	now := time.Date(2026, 4, 28, 15, 0, 0, 0, time.UTC)
	got := FromDateLastThreeMonths(now)
	want := time.Date(2026, 1, 26, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestFromDateLastThreeMonths_AlreadyMonday(t *testing.T) {
	now := time.Date(2026, 4, 27, 0, 0, 0, 0, time.UTC) // Monday
	got := FromDateLastThreeMonths(now)
	want := time.Date(2026, 1, 26, 0, 0, 0, 0, time.UTC) // Jan 27 Wed -> Mon Jan 26
	if !got.Equal(want) {
		t.Fatalf("got %v want %v", got, want)
	}
}
