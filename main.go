package main

import (
	"fmt"
	"os"

	"github.com/andygodish/upver/internal/config"
	"github.com/andygodish/upver/internal/extract"
	"github.com/andygodish/upver/internal/version"
)

var planOnly bool

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

	currentVersion, err := extract.ByRegex(b, cfg.Version.Pattern, cfg.Version.Group)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: extract version:", err)
		os.Exit(1)
	}

	upBytes, err := os.ReadFile(cfg.Upstream.File)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: read upstream file:", err)
		os.Exit(1)
	}

	upstreamTag, err := extract.ByRegex(upBytes, cfg.Upstream.Pattern, cfg.Upstream.Group)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: extract upstream tag:", err)
		os.Exit(1)
	}

	pv, err := version.Parse(currentVersion)
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
