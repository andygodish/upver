package version

import (
	"fmt"
	"strconv"
	"strings"
)

type Parsed struct {
	Base    string
	SeqName string // if empty, implies number-only versioning, ex: 0.1.2-3
	SeqNum  int
}

func Parse(v string) (Parsed, error) {
	sep := strings.LastIndex(v, "-")
	if sep == -1 {
		return Parsed{}, fmt.Errorf("invalid version %q: missing '-'", v)
	}

	base := v[:sep]
	suffix := v[sep+1:]
	if base == "" || suffix == "" {
		return Parsed{}, fmt.Errorf("invalid version %q: empty base or suffix", v)
	}

	// Case A: <base>-<seqName>.<seqNum>  (e.g. 1.0-jam.2)
	if dot := strings.LastIndex(suffix, "."); dot != -1 {
		seqName := suffix[:dot]
		seqNumStr := suffix[dot+1:]

		n, err := strconv.Atoi(seqNumStr)
		if err != nil {
			return Parsed{}, fmt.Errorf("invalid version %q: seq num not int: %w", v, err)
		}

		return Parsed{Base: base, SeqName: seqName, SeqNum: n}, nil
	}

	// Case B: <base>-<seqNum> (e.g. 2.19.1-0)
	n, err := strconv.Atoi(suffix)
	if err != nil {
		return Parsed{}, fmt.Errorf("invalid version %q: suffix not int: %w", v, err)
	}

	return Parsed{Base: base, SeqName: "", SeqNum: n}, nil
}