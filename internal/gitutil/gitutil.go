package gitutil

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Commit struct {
	FullSHA  string
	ShortSHA string
	Subject  string
}

func runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), errMsg)
	}

	return strings.TrimSpace(stdout.String()), nil
}

func LastTag() (string, error) {
	out, err := runGit("describe", "--tags", "--abbrev=0")
	if err != nil {
		// If there are no tags, git describe errors. Treat that as "no tag".
		// The error message varies, so we match loosely.
		if strings.Contains(err.Error(), "No names found") ||
			strings.Contains(err.Error(), "No tags can describe") {
			return "", nil
		}
		return "", err
	}
	return out, nil
}

func Commits(rangeSpec string) ([]Commit, error) {
	args := []string{
		"log",
		"--first-parent",
		"--no-merges",
		`--pretty=format:%H%x09%h%x09%s`,
	}
	if strings.TrimSpace(rangeSpec) != "" {
		args = append(args, rangeSpec)
	}

	out, err := runGit(args...)
	if err != nil {
		return nil, err
	}
	if out == "" {
		return []Commit{}, nil
	}

	lines := strings.Split(out, "\n")
	commits := make([]Commit, 0, len(lines))

	for _, line := range lines {
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("unexpected git log line: %q", line)
		}
		commits = append(commits, Commit{
			FullSHA:  parts[0],
			ShortSHA: parts[1],
			Subject:  parts[2],
		})
	}

	return commits, nil
}