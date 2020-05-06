package gedcom

import (
	"github.com/jochenboesmans/gedcom-parser/util"
	"strconv"
	"strings"
)

type Line struct {
	originalLine string
}

func NewLine(gedcomLine string) *Line {
	return &Line{
		originalLine: gedcomLine,
	}
}

func (gedcomLine *Line) Level() uint8 {
	line := gedcomLine.originalLine
	parts := strings.SplitN(line, " ", 2)
	level, err := strconv.Atoi(parts[0])
	util.MaybePanic(err)
	return uint8(level)
}
func (gedcomLine *Line) XRefID() string {
	line := gedcomLine.originalLine
	parts := strings.SplitN(line, " ", 3)
	if len(parts) >= 2 && parts[1][0] == '@' {
		return parts[1]
	}
	return ""
}

func (gedcomLine *Line) Tag() string {
	line := gedcomLine.originalLine
	parts := strings.SplitN(line, " ", 4)
	if len(parts) >= 2 && parts[1][0] != '@' {
		return parts[1]
	}
	if len(parts) >= 3 && parts[1][0] == '@' {
		return parts[2]
	}
	return ""
}

func (gedcomLine *Line) Value() string {
	line := gedcomLine.originalLine
	parts := strings.SplitN(line, " ", 4)
	if len(parts) == 3 && parts[1][0] != '@' {
		return parts[2]
	}
	if len(parts) == 4 {
		if parts[1][0] == '@' {
			return parts[3]
		} else {
			lastParts := parts[2] + " " + parts[3]
			return lastParts
		}
	}
	return ""
}
