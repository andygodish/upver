package main

import (
	"fmt"
	"os"

	"github.com/andygodish/upver/internal/config"
	"github.com/andygodish/upver/internal/plan"
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

	p, err := plan.Build(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

	fmt.Println("Current version:", p.CurrentVersion)
	fmt.Println("Upstream tag:", p.UpstreamTag)
	fmt.Println("BUMP_MODE:", p.BumpMode)
	fmt.Println("NEW_VERSION:", p.NewVersion)

	if planOnly {
		return
	}
}
