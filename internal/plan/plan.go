package plan

import (
	"fmt"
	"os"

	"github.com/andygodish/upver/internal/config"
	"github.com/andygodish/upver/internal/extract"
	"github.com/andygodish/upver/internal/version"
)

type Plan struct {
	CurrentVersion string
	UpstreamTag    string
	BumpMode       string
	NewVersion     string
}

func Build(cfg *config.Config) (*Plan, error) {
	// Read + extract current version
	b, err := os.ReadFile(cfg.Version.File)
	if err != nil {
		return nil, fmt.Errorf("read version file: %w", err)
	}

	currentVersion, err := extract.ByRegex(b, cfg.Version.Pattern, cfg.Version.Group)
	if err != nil {
		return nil, fmt.Errorf("extract version: %w", err)
	}

	// Read + extract upstream tag
	upBytes, err := os.ReadFile(cfg.Upstream.File)
	if err != nil {
		return nil, fmt.Errorf("read upstream file: %w", err)
	}

	upstreamTag, err := extract.ByRegex(upBytes, cfg.Upstream.Pattern, cfg.Upstream.Group)
	if err != nil {
		return nil, fmt.Errorf("extract upstream tag: %w", err)
	}

	// Parse current version and compute bump/new version
	pv, err := version.Parse(currentVersion)
	if err != nil {
		return nil, fmt.Errorf("parse current version: %w", err)
	}

	bumpMode := "seq"
	newBase := pv.Base
	newSeqNum := pv.SeqNum + 1

	if upstreamTag != pv.Base {
		bumpMode = "semver"
		newBase = upstreamTag
		newSeqNum = 0
	}

	var newVersion string
	if pv.SeqName != "" {
		newVersion = fmt.Sprintf("%s-%s.%d", newBase, pv.SeqName, newSeqNum)
	} else {
		newVersion = fmt.Sprintf("%s-%d", newBase, newSeqNum)
	}

	return &Plan{
		CurrentVersion: currentVersion,
		UpstreamTag:    upstreamTag,
		BumpMode:       bumpMode,
		NewVersion:     newVersion,
	}, nil
}