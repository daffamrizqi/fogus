package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"fogus/journal"
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		printUsage(stderr)
		return errors.New("missing command")
	}

	switch args[0] {
	case "log":
		return runLog(args[1:], stdin, stdout)
	case "today":
		return runToday(args[1:], stdout)
	case "summary":
		return runSummary(args[1:], stdout)
	case "help", "--help", "-h":
		printUsage(stdout)
		return nil
	default:
		printUsage(stderr)
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "fogus <command> [options]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  log      Add a journal entry")
	fmt.Fprintln(w, "  today    Show the day timeline")
	fmt.Fprintln(w, "  summary  Show the day summary")
}

func runLog(args []string, stdin io.Reader, stdout io.Writer) error {
	fs := flag.NewFlagSet("log", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	dateFlag := fs.String("date", "", "target date in YYYY-MM-DD")
	atFlag := fs.String("at", "", "target time in HH:MM")
	editorFlag := fs.Bool("editor", false, "open editor for multi-line input")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *editorFlag && len(fs.Args()) > 0 {
		return errors.New("--editor cannot be used with inline log text")
	}

	entryDate, err := parseDateFlag(*dateFlag)
	if err != nil {
		return err
	}

	entryTime, err := parseTimeFlag(*atFlag, time.Now())
	if err != nil {
		return err
	}

	body, err := readLogBody(fs.Args(), *editorFlag, stdin, stdout, entryTime)
	if err != nil {
		return err
	}

	store, err := journal.NewStore()
	if err != nil {
		return err
	}

	j, err := store.Load(entryDate)
	if err != nil {
		return err
	}

	j.Insert(journal.Entry{Time: entryTime, Body: body})

	if err := store.Save(j); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "Saved to %s\n", store.PathForDate(entryDate))
	return nil
}

func runToday(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("today", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	dateFlag := fs.String("date", "", "target date in YYYY-MM-DD")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if len(fs.Args()) != 0 {
		return errors.New("today does not accept positional arguments")
	}

	targetDate, err := parseDateFlag(*dateFlag)
	if err != nil {
		return err
	}

	store, err := journal.NewStore()
	if err != nil {
		return err
	}

	j, err := store.Load(targetDate)
	if err != nil {
		return err
	}

	if len(j.Entries) == 0 {
		fmt.Fprintf(stdout, "No entries for %s\n", j.Date)
		return nil
	}

	fmt.Fprint(stdout, journal.RenderTimeline(j))
	return nil
}

func runSummary(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("summary", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	dateFlag := fs.String("date", "", "target date in YYYY-MM-DD")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if len(fs.Args()) != 0 {
		return errors.New("summary does not accept positional arguments")
	}

	targetDate, err := parseDateFlag(*dateFlag)
	if err != nil {
		return err
	}

	store, err := journal.NewStore()
	if err != nil {
		return err
	}

	j, err := store.Load(targetDate)
	if err != nil {
		return err
	}

	if len(j.Entries) == 0 {
		fmt.Fprintf(stdout, "No entries for %s\n", j.Date)
		return nil
	}

	fmt.Fprint(stdout, renderSummary(j))
	return nil
}

func parseDateFlag(value string) (time.Time, error) {
	if value == "" {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), nil
	}

	date, err := time.ParseInLocation(journal.DateLayout, value, time.Local)
	if err != nil || date.Format(journal.DateLayout) != value {
		return time.Time{}, fmt.Errorf("invalid --date %q: expected YYYY-MM-DD", value)
	}

	return date, nil
}

func parseTimeFlag(value string, now time.Time) (string, error) {
	if value == "" {
		return now.Format(journal.TimeLayout), nil
	}

	parsed, err := time.Parse(journal.TimeLayout, value)
	if err != nil || parsed.Format(journal.TimeLayout) != value {
		return "", fmt.Errorf("invalid --at %q: expected HH:MM", value)
	}

	return value, nil
}

func renderSummary(j journal.Journal) string {
	var b strings.Builder
	b.WriteString(j.Date)
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("%d entries\n", len(j.Entries)))
	b.WriteString(fmt.Sprintf("First log: %s\n", j.Entries[0].Time))
	b.WriteString(fmt.Sprintf("Last log: %s\n\n", j.Entries[len(j.Entries)-1].Time))
	b.WriteString("Entries\n")

	for _, entry := range j.Entries {
		firstLine, extraLines := summarizeBody(entry.Body)
		b.WriteString("- ")
		b.WriteString(entry.Time)
		b.WriteString(" ")
		b.WriteString(firstLine)
		if !hasTerminalPunctuation(firstLine) {
			b.WriteString(".")
		}
		if extraLines > 0 {
			b.WriteString(fmt.Sprintf(" (+%d %s)", extraLines, pluralize("line", extraLines)))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func summarizeBody(body string) (string, int) {
	lines := strings.Split(strings.ReplaceAll(body, "\r\n", "\n"), "\n")
	trimmed := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		trimmed = append(trimmed, strings.TrimSpace(line))
	}

	if len(trimmed) == 0 {
		return "", 0
	}

	return trimmed[0], len(trimmed) - 1
}

func pluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}

func hasTerminalPunctuation(line string) bool {
	return strings.HasSuffix(line, ".") || strings.HasSuffix(line, "!") || strings.HasSuffix(line, "?")
}
