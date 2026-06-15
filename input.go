package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func readLogBody(args []string, useEditor bool, stdin io.Reader, stdout io.Writer, entryTime string) (string, error) {
	if len(args) > 0 {
		body := strings.TrimSpace(strings.Join(args, " "))
		if body == "" {
			return "", errors.New("log body cannot be empty")
		}
		return body, nil
	}

	if useEditor {
		return readFromEditor()
	}

	if stdinFile, ok := stdin.(*os.File); ok {
		if info, err := stdinFile.Stat(); err == nil && (info.Mode()&os.ModeCharDevice) != 0 {
			fmt.Fprintf(stdout, "%s - write log, end with Ctrl+D\n", entryTime)
		}
	}

	body, err := io.ReadAll(bufio.NewReader(stdin))
	if err != nil {
		return "", err
	}

	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" {
		return "", errors.New("log body cannot be empty")
	}

	return trimmed, nil
}

func readFromEditor() (string, error) {
	editor := strings.TrimSpace(os.Getenv("EDITOR"))
	if editor == "" {
		return "", errors.New("$EDITOR is not set")
	}

	tempDir := os.TempDir()
	file, err := os.CreateTemp(tempDir, "fogus-*.md")
	if err != nil {
		return "", err
	}
	path := file.Name()
	if err := file.Close(); err != nil {
		return "", err
	}
	defer os.Remove(path)

	parts := strings.Fields(editor)
	if len(parts) == 0 {
		return "", errors.New("$EDITOR is invalid")
	}

	cmdArgs := append(parts[1:], path)
	cmd := exec.Command(parts[0], cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	body, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", err
	}

	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" {
		return "", errors.New("log body cannot be empty")
	}

	return trimmed, nil
}
