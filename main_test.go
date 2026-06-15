package main

import (
	"testing"

	"fogus/journal"
)

func TestRenderSummary(t *testing.T) {
	j := journal.Journal{
		Date: "2026-06-12",
		Entries: []journal.Entry{
			{Time: "09:00", Body: "Mulai PBI data retention."},
			{Time: "09:47", Body: "Kepikiran cek AWS billing.\nTahan dulu. Lanjut retention."},
		},
	}

	got := renderSummary(j)
	want := "2026-06-12\n2 entries\nFirst log: 09:00\nLast log: 09:47\n\nEntries\n- 09:00 Mulai PBI data retention.\n- 09:47 Kepikiran cek AWS billing. (+1 line)\n"

	if got != want {
		t.Fatalf("renderSummary() mismatch\nwant:\n%s\n got:\n%s", want, got)
	}
}

func TestParseDateFlagInvalid(t *testing.T) {
	if _, err := parseDateFlag("2026/06/12"); err == nil {
		t.Fatal("expected invalid date error")
	}
}

func TestParseTimeFlagInvalid(t *testing.T) {
	if _, err := parseTimeFlag("9:47", testNow()); err == nil {
		t.Fatal("expected invalid time error")
	}
}
