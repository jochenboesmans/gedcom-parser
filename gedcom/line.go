package gedcom

import (
	"github.com/jochenboesmans/gedcom-parser/util"
	"strconv"
	"strings"
)

type Line struct {
	level  uint8
	xRefId string
	tag    string
	value  string
}

func NewLine(gedcomLine string) *Line {
	bomTrimmedGedcomLine := strings.Trim(gedcomLine, "\uFEFF")
	parts := strings.SplitN(bomTrimmedGedcomLine, " ", 4)
	level, err := strconv.ParseUint(parts[0], 10, 64)
	util.MaybePanic(err)
	line := Line{
		level: uint8(level),
	}

	if len(parts) >= 2 {
		if parts[1][0] == '@' {
			line.xRefId = parts[1]
		} else {
			line.tag = parts[1]
		}
	}

	if len(parts) >= 3 {
		if line.xRefId != "" {
			line.tag = parts[2]
		} else {
			line.value = parts[2]
		}
	}

	if len(parts) >= 4 {
		if line.xRefId != "" {
			line.value = parts[3]
		} else {
			line.value = parts[2] + " " + parts[3]
		}
	}

	return &line
}

func (gedcomLine *Line) Level() uint8 {
	return gedcomLine.level
}
func (gedcomLine *Line) XRefID() string {
	return gedcomLine.xRefId
}
func (gedcomLine *Line) Tag() string {
	return gedcomLine.tag
}
func (gedcomLine *Line) Value() string {
	return gedcomLine.value
}
