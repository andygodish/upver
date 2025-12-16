package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/andygodish/upver/internal/config"
	"github.com/andygodish/upver/internal/plan"
	"github.com/andygodish/upver/internal/apply"
	"github.com/andygodish/upver/internal/changelog"
	"github.com/andygodish/upver/internal/gitutil"
)

func main() {
	applyFlag := flag.Bool("apply", false, "apply changes (write updated version to file)")
	changelogFlag := flag.Bool("changelog", false, "update changelog file based on git history")
	configPath := flag.String("config", "upver.yaml", "path to upver config file")
	planOnly := flag.Bool("plan", false, "print computed plan and exit (default behavior for now)")
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

	// changelog generation (always for now; we’ll gate behind a flag next)
	lastTag, err := gitutil.LastTag()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

	rangeSpec := "HEAD"
	if lastTag != "" {
		rangeSpec = fmt.Sprintf("%s..HEAD", lastTag)
	}

	commits, err := gitutil.Commits(rangeSpec)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

	// commit URL base: for now, require it in config later; temporary placeholder:
	commitURLBase := "" // we’ll fill this next step
	_ = commitURLBase

	_ = commits

	fmt.Println("Current version:", p.CurrentVersion)
	fmt.Println("Upstream tag:", p.UpstreamTag)
	fmt.Println("BUMP_MODE:", p.BumpMode)
	fmt.Println("NEW_VERSION:", p.NewVersion)

	if *changelogFlag {
		lastTag, err := gitutil.LastTag()
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			os.Exit(1)
		}

		rangeSpec := "HEAD"
		if lastTag != "" {
			rangeSpec = fmt.Sprintf("%s..HEAD", lastTag)
		}

		commits, err := gitutil.Commits(rangeSpec)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			os.Exit(1)
		}

		commitURLBase, err := gitutil.CommitURLBase(os.Getenv("CI_SERVER_URL"), os.Getenv("CI_PROJECT_PATH"))
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			os.Exit(1)
		}

		section := changelog.Render(p.NewVersion, commits, commitURLBase)

		if err := changelog.UpdateFile(cfg.Changelog.File, p.NewVersion, section); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			os.Exit(1)
		}
	}

	if *planOnly {
		return
	}

	if *applyFlag {
		if err := apply.Version(cfg.Version, p.NewVersion); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			os.Exit(1)
		}
	}
}
