package journal

import "testing"

func TestParseAndRenderRoundTrip(t *testing.T) {
	content := "# 2026-06-12\n\n09:00\nMulai PBI data retention.\n\n09:47\nKepikiran cek AWS billing.\nTahan dulu. Lanjut retention.\n"

	j, err := Parse("2026-06-12", content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	got := Render(j)
	if got != content {
		t.Fatalf("Render() mismatch\nwant:\n%s\n got:\n%s", content, got)
	}
}

func TestInsertKeepsChronologicalOrder(t *testing.T) {
	j := Journal{
		Date: "2026-06-12",
		Entries: []Entry{
			{Time: "09:00", Body: "First"},
			{Time: "10:30", Body: "Third"},
		},
	}

	j.Insert(Entry{Time: "09:45", Body: "Second"})

	got := []string{j.Entries[0].Time, j.Entries[1].Time, j.Entries[2].Time}
	want := []string{"09:00", "09:45", "10:30"}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Insert() order mismatch: got %v want %v", got, want)
		}
	}
}

func TestInsertPreservesSameTimestampOrder(t *testing.T) {
	j := Journal{
		Date: "2026-06-12",
		Entries: []Entry{
			{Time: "09:00", Body: "First"},
			{Time: "09:00", Body: "Second"},
		},
	}

	j.Insert(Entry{Time: "09:00", Body: "Third"})

	if j.Entries[2].Body != "Third" {
		t.Fatalf("expected same timestamp insert at the end, got %#v", j.Entries)
	}
}

func TestRenderTimeline(t *testing.T) {
	j := Journal{
		Date:    "2026-06-12",
		Entries: []Entry{{Time: "09:00", Body: "Mulai PBI data retention."}},
	}

	got := RenderTimeline(j)
	want := "2026-06-12\n\n09:00\nMulai PBI data retention.\n"

	if got != want {
		t.Fatalf("RenderTimeline() mismatch\nwant:\n%s\n got:\n%s", want, got)
	}
}

func TestParseAllowsTimeLikeBodyLine(t *testing.T) {
	content := "# 2026-06-12\n\n09:00\nCatatan schema:\n09:45 bukan entry baru di sini.\n\n10:00\nBalik kerja.\n"

	j, err := Parse("2026-06-12", content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(j.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(j.Entries))
	}

	if j.Entries[0].Body != "Catatan schema:\n09:45 bukan entry baru di sini." {
		t.Fatalf("unexpected first body: %q", j.Entries[0].Body)
	}
}
