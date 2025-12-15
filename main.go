package main

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version   VersionConfig   `yaml:"version"`
	Upstream  UpstreamConfig  `yaml:"upstream"`
	Changelog ChangelogConfig `yaml:"changelog"`
}

type VersionConfig struct {
	File    string `yaml:"file"`
	Pattern string `yaml:"pattern"`
	Group   int    `yaml:"group"`
}

type UpstreamConfig struct {
	Type        string `yaml:"type"`
	File        string `yaml:"file"`
	ImagePrefix string `yaml:"imagePrefix"`
}

type ChangelogConfig struct {
	File string `yaml:"file"`
}

func loadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	// minimal validation (fail fast)
	if cfg.Version.File == "" || cfg.Version.Pattern == "" {
		return nil, fmt.Errorf("config: version.file and version.pattern are required")
	}
	if cfg.Version.Group == 0 {
		cfg.Version.Group = 1
	}
	if cfg.Upstream.Type == "" || cfg.Upstream.File == "" {
		return nil, fmt.Errorf("config: upstream.type and upstream.file are required")
	}
	if cfg.Changelog.File == "" {
		return nil, fmt.Errorf("config: changelog.file is required")
	}

	return &cfg, nil
}

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

func main() {
	cfg, err := loadConfig("upver.yaml")
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

	fmt.Println("Current version:", currentVersion)
}
