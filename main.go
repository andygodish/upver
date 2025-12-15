package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/andygodish/upver/internal/config"
	"github.com/andygodish/upver/internal/plan"
)

func main() {
	planOnly := flag.Bool("plan", false, "print computed plan and exit (default behavior for now)")
	configPath := flag.String("config", "upver.yaml", "path to upver config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
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

	if *planOnly {
		return
	}
}
