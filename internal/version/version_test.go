package version

import "testing"

func TestParseNumberOnlySuffix(t *testing.T) {
	got, err := Parse("0.11.1-3")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if got.Base != "0.11.1" || got.SeqName != "" || got.SeqNum != 3 {
		t.Fatalf("Parse returned %#v", got)
	}
}

func TestParseNamedSequence(t *testing.T) {
	got, err := Parse("1.0-jam.2")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if got.Base != "1.0" || got.SeqName != "jam" || got.SeqNum != 2 {
		t.Fatalf("Parse returned %#v", got)
	}
}

func TestParseBaseWithHyphenatedDebianRevision(t *testing.T) {
	got, err := Parse("0.4.9.11-0+deb13u1-0")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if got.Base != "0.4.9.11-0+deb13u1" || got.SeqName != "" || got.SeqNum != 0 {
		t.Fatalf("Parse returned %#v", got)
	}
}

func TestParseBaseWithHyphenatedDebianRevisionAndNamedSequence(t *testing.T) {
	got, err := Parse("0.4.9.11-0+deb13u1-tor.4")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if got.Base != "0.4.9.11-0+deb13u1" || got.SeqName != "tor" || got.SeqNum != 4 {
		t.Fatalf("Parse returned %#v", got)
	}
}
