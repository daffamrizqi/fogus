package journal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	DateLayout = "2006-01-02"
	TimeLayout = "15:04"
)

type Entry struct {
	Time string
	Body string
}

type Journal struct {
	Date    string
	Entries []Entry
}

type Store struct {
	baseDir string
}

func NewStore() (Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Store{}, err
	}

	return Store{baseDir: filepath.Join(home, ".fogus", "journal")}, nil
}

func (s Store) PathForDate(date time.Time) string {
	return filepath.Join(s.baseDir, date.Format(DateLayout)+".md")
}

func (s Store) Load(date time.Time) (Journal, error) {
	path := s.PathForDate(date)

	body, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Journal{Date: date.Format(DateLayout)}, nil
		}
		return Journal{}, err
	}

	j, err := Parse(date.Format(DateLayout), string(body))
	if err != nil {
		return Journal{}, fmt.Errorf("parse %s: %w", path, err)
	}

	return j, nil
}

func (s Store) Save(j Journal) error {
	if err := os.MkdirAll(s.baseDir, 0o755); err != nil {
		return err
	}

	path := filepath.Join(s.baseDir, j.Date+".md")
	return os.WriteFile(path, []byte(Render(j)), 0o644)
}

func (j *Journal) Insert(entry Entry) {
	index := len(j.Entries)
	for i, existing := range j.Entries {
		if compareTimes(entry.Time, existing.Time) < 0 {
			index = i
			break
		}
	}

	j.Entries = append(j.Entries, Entry{})
	copy(j.Entries[index+1:], j.Entries[index:])
	j.Entries[index] = entry
}

func Parse(expectedDate, content string) (Journal, error) {
	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	trimmed := strings.TrimSpace(normalized)
	if trimmed == "" {
		return Journal{Date: expectedDate}, nil
	}

	lines := strings.Split(trimmed, "\n")
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "# ") {
		return Journal{}, errors.New("missing journal date header")
	}

	date := strings.TrimSpace(strings.TrimPrefix(lines[0], "# "))
	if date != expectedDate {
		return Journal{}, fmt.Errorf("journal header date %q does not match %q", date, expectedDate)
	}

	j := Journal{Date: date}
	for i := 1; i < len(lines); {
		if strings.TrimSpace(lines[i]) == "" {
			i++
			continue
		}

		if !isTimeLine(lines[i]) {
			return Journal{}, fmt.Errorf("expected time at line %d", i+1)
		}

		entryTime := lines[i]
		i++

		bodyLines := make([]string, 0)
		for i < len(lines) {
			if strings.TrimSpace(lines[i]) == "" {
				if i+1 < len(lines) && isTimeLine(lines[i+1]) {
					break
				}
				bodyLines = append(bodyLines, "")
				i++
				continue
			}

			bodyLines = append(bodyLines, lines[i])
			i++
		}

		body := strings.TrimSpace(strings.Join(bodyLines, "\n"))
		if body == "" {
			return Journal{}, fmt.Errorf("empty body for entry at %s", entryTime)
		}

		j.Entries = append(j.Entries, Entry{Time: entryTime, Body: body})
	}

	return j, nil
}

func Render(j Journal) string {
	var b strings.Builder
	b.WriteString("# ")
	b.WriteString(j.Date)
	b.WriteString("\n")

	if len(j.Entries) == 0 {
		return b.String()
	}

	b.WriteString("\n")
	for i, entry := range j.Entries {
		if i > 0 {
			b.WriteString("\n")
		}

		b.WriteString(entry.Time)
		b.WriteString("\n")
		b.WriteString(strings.TrimSpace(entry.Body))
		b.WriteString("\n")
	}

	return b.String()
}

func RenderTimeline(j Journal) string {
	var b strings.Builder
	b.WriteString(j.Date)
	b.WriteString("\n\n")

	for i, entry := range j.Entries {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(entry.Time)
		b.WriteString("\n")
		b.WriteString(entry.Body)
		b.WriteString("\n")
	}

	return b.String()
}

func compareTimes(left, right string) int {
	if left < right {
		return -1
	}
	if left > right {
		return 1
	}
	return 0
}

func isTimeLine(line string) bool {
	if len(line) != len(TimeLayout) {
		return false
	}

	_, err := time.Parse(TimeLayout, line)
	return err == nil
}
