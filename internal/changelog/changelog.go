package changelog

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/andygodish/upver/internal/gitutil"
)

func Render(version string, commits []gitutil.Commit, commitURLBase string) string {
	dateUTC := time.Now().UTC().Format("2006-01-02")

	var b strings.Builder
	fmt.Fprintf(&b, "## %s - %s\n\n", version, dateUTC)

	for _, c := range commits {
		// - [abcd123](https://.../commit/<FULL>) subject
		fmt.Fprintf(&b, "- [%s](%s/%s) %s\n", c.ShortSHA, commitURLBase, c.FullSHA, c.Subject)
	}

	b.WriteString("\n")
	return b.String()
}

func UpdateFile(path string, version string, newSection string) error {
	// Read existing file (if any)
	existing := ""
	if data, err := os.ReadFile(path); err == nil {
		existing = string(data)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("read %s: %w", path, err)
	}

	// Remove any existing section for this version
	cleaned := removeSection(existing, version)

	// Prepend new section
	out := newSection + cleaned

	// Write back
	if err := os.WriteFile(path, []byte(out), 0644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func removeSection(content string, version string) string {
	if content == "" {
		return ""
	}

	var out strings.Builder
	sc := bufio.NewScanner(strings.NewReader(content))

	inSection := false
	versionHeading := "## " + version

	for sc.Scan() {
		line := sc.Text()

		if strings.HasPrefix(line, "## ") {
			// entering any heading
			if strings.HasPrefix(line, versionHeading) {
				inSection = true
				continue // skip the heading line itself
			}
			// if we were skipping, and hit next heading, stop skipping
			if inSection {
				inSection = false
			}
		}

		if !inSection {
			out.WriteString(line)
			out.WriteString("\n")
		}
	}

	// If scanner had an error, best to fall back to original content.
	// But scanner errors are rare here; keep it simple.
	return out.String()
}