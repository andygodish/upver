package version

import (
	"fmt"
	"strconv"
	"strings"
)

type Parsed struct {
	Base    string
	SeqName string
	SeqNum  int
}

func Parse(v string) (Parsed, error) {
	// expected: <base>-<seqName>.<seqNum>   e.g. 1.0-jam.2
	parts := strings.SplitN(v, "-", 2)
	if len(parts) != 2 {
		return Parsed{}, fmt.Errorf("invalid version %q: missing '-'", v)
	}
	base := parts[0]
	rest := parts[1]

	dot := strings.LastIndex(rest, ".")
	if dot == -1 {
		return Parsed{}, fmt.Errorf("invalid version %q: missing '.'", v)
	}

	seqName := rest[:dot]
	seqNumStr := rest[dot+1:]

	n, err := strconv.Atoi(seqNumStr)
	if err != nil {
		return Parsed{}, fmt.Errorf("invalid version %q: seq num not int: %w", v, err)
	}

	return Parsed{Base: base, SeqName: seqName, SeqNum: n}, nil
}