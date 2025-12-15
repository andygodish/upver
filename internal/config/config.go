package config

import (
	"fmt"
	"os"

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
	File    string `yaml:"file"`
	Pattern string `yaml:"pattern"`
	Group   int    `yaml:"group"`
}

type ChangelogConfig struct {
	File string `yaml:"file"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	// minimal validation (fail fast)
	// Version
	if cfg.Version.File == "" || cfg.Version.Pattern == "" {
		return nil, fmt.Errorf("config: version.file and version.pattern are required")
	}
	if cfg.Version.Group == 0 {
		cfg.Version.Group = 1
	}

	// Upstream
	if cfg.Upstream.File == "" || cfg.Upstream.Pattern == "" {
		return nil, fmt.Errorf("config: upstream.file and upstream.pattern are required")
	}
	if cfg.Upstream.Group == 0 {
		cfg.Upstream.Group = 1
	}

	// Changelog
	if cfg.Changelog.File == "" {
		return nil, fmt.Errorf("config: changelog.file is required")
	}

	return &cfg, nil
}
