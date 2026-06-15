# fogus

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

`fogus` is a terminal tool for interstitial journaling.

Its goal is simple: help you quickly capture what you are working on, what interrupted your focus, and how your day actually unfolded.

This tool is a good fit for people who frequently context-switch, jump between tasks, or need a lightweight external memory while working.

Or maybe if you're an engineer and have adhd like me LOL. The idea of me for making this tool is to simply keeping track and help me to focus in the middle of this chaotic 
distraction world in a simple way. Hence why i made this tool as simple as possible. Looking forward to add new features for this tool.

## Why fogus

Interstitial journaling is useful because you do not need to wait until the end of the day to remember what happened. You just capture small slices of activity as they happen.

`fogus` tries to keep that process lightweight:

1. terminal-first
2. local-first
3. fast to use
4. human-readable file format
5. no lock-in to a hard-to-read storage format

## Current status

The current version is an MVP that is already usable for personal daily use.

The focus is not on having many features yet, but on making the core flow fast and reliable.

## Current features

1. Log activities with automatic timestamps
2. Add entries for a specific time with `--at HH:MM`
3. Add entries for a specific date with `--date YYYY-MM-DD`
4. Write one-line or multi-line logs
5. Write logs through your editor with `--editor`
6. Automatically save daily journals as Markdown
7. Show the timeline for a given day with `today`
8. Show a neutral daily summary with `summary`
9. Insert entries chronologically instead of only appending at the bottom

## Quickstart

### Setup from scratch

1. Install Go:

```bash
brew install go
```

2. Clone the repository:

```bash
git clone <repo-url>
cd fogus
```

3. Install `fogus` to `~/.local/bin`:

```bash
make install
```

4. Verify the command:

```bash
fogus help
```

5. Try it right away:

```bash
fogus log "Start working"
fogus today
fogus summary
```

### Short setup

If Go is already installed and you are already inside the repo:

```bash
make install
fogus log "Start working"
fogus today
fogus summary
```

## Installation

The binary will be installed to:

```text
~/.local/bin/fogus
```

That directory must be available in your shell `PATH`.

## Storage location

All logs are automatically stored in:

```text
~/.fogus/journal
```

Daily file format:

```text
~/.fogus/journal/YYYY-MM-DD.md
```

Example:

```text
~/.fogus/journal/2026-06-12.md
```

If the folder does not exist yet, `fogus` will create it automatically.

## Commands

```bash
fogus log [--date YYYY-MM-DD] [--at HH:MM] [--editor] [text...]
fogus today [--date YYYY-MM-DD]
fogus summary [--date YYYY-MM-DD]
```

## Basic usage

### Quick capture

```bash
fogus log "Start PBI data retention."
```

`fogus` will automatically:

1. capture the current time
2. determine today's file
3. save the log to `~/.fogus/journal`

### Multi-line from the terminal

```bash
fogus log
```

Then write your log and finish with `Ctrl+D`.

Example:

```text
09:47 - write log, end with Ctrl+D
Thought about checking AWS billing.
Hold that thought. Back to retention.
```

### Multi-line through your editor

```bash
export EDITOR=nvim
fogus log --editor
```

`fogus` will open the editor defined in the `$EDITOR` environment variable.

If `$EDITOR` is not set, the command will fail.

### Backfill a specific time

If you forgot to log at the exact moment:

```bash
fogus log --at 09:45 "Finished analyzing transactions.parquet."
```

The entry will be inserted in the correct chronological position.

### Backfill a specific date

If you want to log something for a previous day:

```bash
fogus log --date 2026-06-12 --at 11:20 "Got stuck on the query finder."
```

### View the timeline

```bash
fogus today
fogus today --date 2026-06-12
```

If no entries exist:

```text
No entries for 2026-06-13
```

### View the summary

```bash
fogus summary
fogus summary --date 2026-06-12
```

The current summary is intentionally neutral. It does not yet group entries into categories like `focus`, `break`, or `blocker`.

## A realistic workflow example

This example is close to the kind of workday where focus gets interrupted, small distractions show up, and you need to return to the main task quickly.

### Start work

```bash
fogus log "Start PBI data retention."
```

### Finish one analysis step

```bash
fogus log --at 09:45 "Finished analyzing transactions.parquet."
```

### A sudden distraction appears

```bash
fogus log
```

Then write:

```text
Thought about checking AWS billing.
Hold that thought. Back to retention.
```

### Take a short break

```bash
fogus log "Coffee break."
```

### Return to the main task

```bash
fogus log "Back to the batch delete script."
```

### Get stuck and need a reference

```bash
fogus log
```

Then write:

```text
Got stuck on the query finder.
Looking up the schema reference.
```

### Review the day so far

```bash
fogus today
fogus summary
```

With a flow like this, you do not need to keep everything in your head. `fogus` becomes a lightweight trail that helps you:

1. see real progress
2. capture interruptions before they disappear
3. return to the right context after getting distracted
4. review work patterns at the end of the day

## Example output

### Timeline

```text
2026-06-12

09:00
Start PBI data retention.

09:45
Finished analyzing transactions.parquet.

09:47
Thought about checking AWS billing.
Hold that thought. Back to retention.
```

### Summary

```text
2026-06-12
6 entries
First log: 09:00
Last log: 11:20

Entries
- 09:00 Start PBI data retention.
- 09:45 Finished analyzing transactions.parquet.
- 09:47 Thought about checking AWS billing. (+1 line)
```

## Journal file format

Example saved file:

```md
# 2026-06-12

09:00
Start PBI data retention.

09:45
Finished analyzing transactions.parquet.

09:47
Thought about checking AWS billing.
Hold that thought. Back to retention.
```

This format is intentionally simple so that it is:

1. easy to read directly
2. easy to back up
3. easy to move into Obsidian or Git
4. still useful even without `fogus`

## Important behavior

1. Entries with earlier timestamps are inserted into the middle of the file when needed
2. Entries with the same timestamp keep input order
3. The journal file is rewritten in a consistent format after each change
4. `today` and `summary` read from the same Markdown source

## Current limitations

1. No tagging yet such as `#break` or `#blocker`
2. No automatic grouped summaries yet
3. No weekly review yet
4. No sync to a specific Obsidian vault yet
5. No full-screen TUI mode yet

## Roadmap

### Near term

1. Lightweight summary grouping such as `focus`, `break`, `blocker`, and `interruption`
2. Lightweight tagging without making input feel heavy
3. Weekly review to spot work patterns over time
4. Cleaner install and distribution workflows

### Mid term

1. Optional integration with a specific Obsidian vault
2. Export or output formats better suited for daily reflection
3. Search and filtering by date or keyword
4. Faster quick-capture prompts or shortcuts

### Longer term

1. TUI mode for journaling and review without leaving the terminal
2. Daily templates or more structured review flows
3. Lightweight insights from journaling patterns without making the tool feel heavy
4. More mature positioning for open source distribution and paid distribution

## For open source or product direction

If `fogus` becomes a public or paid project later, the strongest positioning is probably:

1. a journaling tool for people who frequently context-switch
2. a terminal-native capture tool for knowledge workers
3. a local-first daily work log that stays human-readable

Its real value is not just that it is a CLI note app, but that it helps people return to working context after their attention gets fragmented.

## Development

```bash
make build
make test
make install
make clean
```
