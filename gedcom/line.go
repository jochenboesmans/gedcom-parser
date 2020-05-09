package gedcom

import (
	"github.com/jochenboesmans/gedcom-parser/util"
	"strconv"
	"strings"
)

type Line struct {
	originalLine *string
	levelMemo    *uint8
	xRefIDMemo   *string
	tagMemo      *string
	valueMemo    *string
}

func NewLine(gedcomLinePtr *string) *Line {
	return &Line{
		originalLine: gedcomLinePtr,
	}
}

func (gedcomLine *Line) Level() *uint8 {
	if gedcomLine.levelMemo != nil {
		return gedcomLine.levelMemo
	}
	parts := strings.SplitN(*gedcomLine.originalLine, " ", 2)
	level, err := strconv.Atoi(parts[0])
	util.MaybePanic(err)
	levelUint8 := uint8(level)
	result := &levelUint8
	gedcomLine.levelMemo = result
	return result
}
func (gedcomLine *Line) XRefID() *string {
	if gedcomLine.xRefIDMemo != nil {
		return gedcomLine.xRefIDMemo
	}
	parts := strings.SplitN(*gedcomLine.originalLine, " ", 3)
	var result *string = nil
	if len(parts) >= 2 && parts[1][0] == '@' {
		result = &parts[1]
	}
	gedcomLine.xRefIDMemo = result
	return result
}

func (gedcomLine *Line) Tag() *string {
	if gedcomLine.tagMemo != nil {
		return gedcomLine.tagMemo
	}
	parts := strings.SplitN(*gedcomLine.originalLine, " ", 4)
	var result *string = nil
	if len(parts) >= 2 && parts[1][0] != '@' {
		result = &parts[1]
	}
	if len(parts) >= 3 && parts[1][0] == '@' {
		result = &parts[2]
	}
	gedcomLine.tagMemo = result
	return result
}

func (gedcomLine *Line) Value() *string {
	if gedcomLine.valueMemo != nil {
		return gedcomLine.valueMemo
	}
	parts := strings.SplitN(*gedcomLine.originalLine, " ", 4)
	var result *string = nil
	if len(parts) == 3 && parts[1][0] != '@' {
		result = &parts[2]
	}
	if len(parts) == 4 {
		if parts[1][0] == '@' {
			result = &parts[3]
		} else {
			lastParts := parts[2] + " " + parts[3]
			result = &lastParts
		}
	}
	gedcomLine.valueMemo = result
	return result
}
