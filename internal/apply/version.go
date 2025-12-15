package apply

import (
	"fmt"
	"os"
	"regexp"

	"github.com/andygodish/upver/internal/config"
)

func Version(cfg config.VersionConfig, newVersion string) error {
	b, err := os.ReadFile(cfg.File)
	if err != nil {
		return fmt.Errorf("read %s: %w", cfg.File, err)
	}

	re, err := regexp.Compile(cfg.Pattern)
	if err != nil {
		return fmt.Errorf("compile version regex: %w", err)
	}

	// Find match indices and capture-group indices
	idx := re.FindSubmatchIndex(b)
	if idx == nil {
		return fmt.Errorf("version pattern did not match in %s", cfg.File)
	}

	// idx layout:
	// [ fullStart, fullEnd, g1Start, g1End, g2Start, g2End, ... ]
	group := cfg.Group
	if group <= 0 {
		group = 1
	}
	gi := group * 2
	if gi+1 >= len(idx) {
		return fmt.Errorf("group %d out of range (have %d capture groups)", group, (len(idx)/2)-1)
	}

	start, end := idx[gi], idx[gi+1]
	if start < 0 || end < 0 {
		return fmt.Errorf("group %d did not participate in the match", group)
	}

	updated := make([]byte, 0, len(b)-(end-start)+len(newVersion))
	updated = append(updated, b[:start]...)
	updated = append(updated, []byte(newVersion)...)
	updated = append(updated, b[end:]...)

	// Preserve existing file mode if possible
	mode := os.FileMode(0644)
	if st, statErr := os.Stat(cfg.File); statErr == nil {
		mode = st.Mode()
	}

	if err := os.WriteFile(cfg.File, updated, mode); err != nil {
		return fmt.Errorf("write %s: %w", cfg.File, err)
	}

	return nil
}