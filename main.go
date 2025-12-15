package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/andygodish/upver/internal/config"
)

var planOnly bool

func extractByRegex(b []byte, pattern string, group int) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("compile regex: %w", err)
	}

	m := re.FindSubmatch(b)
	if m == nil {
		return "", fmt.Errorf("pattern did not match")
	}

	if group < 0 || group >= len(m) {
		return "", fmt.Errorf("group %d out of range (matches=%d)", group, len(m)-1)
	}

	return string(m[group]), nil
}

type ParsedVersion struct {
	Base    string
	SeqName string
	SeqNum  int
}

func parseVersion(v string) (ParsedVersion, error) {
	// expected: <base>-<seqName>.<seqNum>   e.g. 1.0-jam.2
	parts := strings.SplitN(v, "-", 2)
	if len(parts) != 2 {
		return ParsedVersion{}, fmt.Errorf("invalid version %q: missing '-'", v)
	}
	base := parts[0]
	rest := parts[1]

	dot := strings.LastIndex(rest, ".")
	if dot == -1 {
		return ParsedVersion{}, fmt.Errorf("invalid version %q: missing '.'", v)
	}

	seqName := rest[:dot]
	seqNumStr := rest[dot+1:]

	n, err := strconv.Atoi(seqNumStr)
	if err != nil {
		return ParsedVersion{}, fmt.Errorf("invalid version %q: seq num not int: %w", v, err)
	}

	return ParsedVersion{Base: base, SeqName: seqName, SeqNum: n}, nil
}

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "--plan" {
			planOnly = true
		}
	}

	cfg, err := config.Load("upver.yaml")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

	b, err := os.ReadFile(cfg.Version.File)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: read version file:", err)
		os.Exit(1)
	}

	currentVersion, err := extractByRegex(b, cfg.Version.Pattern, cfg.Version.Group)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: extract version:", err)
		os.Exit(1)
	}

	upBytes, err := os.ReadFile(cfg.Upstream.File)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: read upstream file:", err)
		os.Exit(1)
	}

	upstreamTag, err := extractByRegex(upBytes, cfg.Upstream.Pattern, cfg.Upstream.Group)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: extract upstream tag:", err)
		os.Exit(1)
	}

	pv, err := parseVersion(currentVersion)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: parse current version:", err)
		os.Exit(1)
	}

	bumpMode := "seq"
	newBase := pv.Base
	newSeqNum := pv.SeqNum + 1

	if upstreamTag != pv.Base {
		bumpMode = "semver"
		newBase = upstreamTag
		newSeqNum = 0
	}

	newVersion := fmt.Sprintf("%s-%s.%d", newBase, pv.SeqName, newSeqNum)

	fmt.Println("Current version:", currentVersion)
	fmt.Println("Upstream tag:", upstreamTag)
	fmt.Println("BUMP_MODE:", bumpMode)
	fmt.Println("NEW_VERSION:", newVersion)

	if planOnly {
		return
	}
}
