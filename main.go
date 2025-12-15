package main

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
	File     string `yaml:"file"`
	YamlPath string `yaml:"yamlPath"`
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
	if cfg.Version.File == "" || cfg.Version.YamlPath == "" {
		return nil, fmt.Errorf("config: version.file and version.yamlPath are required")
	}
	if cfg.Upstream.Type == "" || cfg.Upstream.File == "" {
		return nil, fmt.Errorf("config: upstream.type and upstream.file are required")
	}
	if cfg.Changelog.File == "" {
		return nil, fmt.Errorf("config: changelog.file is required")
	}

	return &cfg, nil
}

func main() {
	cfg, err := loadConfig("upver.yaml")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded config:\n%+v\n", *cfg)
}
