package gedcom

import (
	"testing"
)

var gedcomLines = []string{
	"0 HEAD",
	"0 @1@ INDI",
	"1 NAME Robert Eugene/Williams/",
	"1 SEX M",
	"1 BIRT",
	"2 DATE 02 OCT 1822",
	"1 FAMC @4@",
}

func TestLine_Level(t *testing.T) {
	expectedLevels := []int8{0, 0, 1, 1, 1, 2, 1}
	for i, line := range gedcomLines {
		l := NewLine(&line)
		if level, err := l.Level(); err != nil || level != expectedLevels[i] {
			t.Errorf("unexpected level at index %d, expected: %d, actual: %d", i, expectedLevels[i], level)
		}
	}
}

func TestLine_XRefID(t *testing.T) {
	expectedValues := []string{"", "@1@", "", "", "", "", ""}
	for i, line := range gedcomLines {
		l := NewLine(&line)
		if l.XRefID() != "" && l.XRefID() != expectedValues[i] {
			t.Errorf("unexpected XRefId at index %d, expected: %s, actual: %s", i, expectedValues[i], l.XRefID())
		} else if l.XRefID() == "" && expectedValues[i] != "" {
			t.Errorf("unexpected XRefId at index %d, expected: %s, actual: %s", i, expectedValues[i], l.XRefID())
		}
	}
}

func TestLine_Tag(t *testing.T) {
	expectedValues := []string{"HEAD", "INDI", "NAME", "SEX", "BIRT", "DATE", "FAMC"}
	for i, line := range gedcomLines {
		l := NewLine(&line)
		if tag, err := l.Tag(); err != nil || tag != expectedValues[i] {
			t.Errorf("unexpected tag at index %d, expected: %s, actual: %s", i, expectedValues[i], tag)
		}
	}
}

func TestLine_Value(t *testing.T) {
	expectedValues := []string{
		"",
		"",
		"Robert Eugene/Williams/",
		"M",
		"",
		"02 OCT 1822",
		"@4@",
	}
	for i, line := range gedcomLines {
		l := NewLine(&line)
		if l.Value() != "" && l.Value() != expectedValues[i] {
			t.Errorf("unexpected value at index %d, expected: %s, actual: %s", i, expectedValues[i], l.Value())
		} else if l.Value() == "" && expectedValues[i] != "" {
			t.Errorf("unexpected value at index %d, expected: %s, actual: %s", i, expectedValues[i], l.Value())
		}
	}
}
