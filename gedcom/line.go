package gedcom

import (
	"fmt"
	"log"
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
	var result *uint8 = nil
	if err != nil {
		fmt.Printf("%s", *gedcomLine.originalLine)
	} else {
		levelUint8 := uint8(level)
		result = &levelUint8
		if result == nil {
			log.Println("no value for required field 'level' of gedcom line.")
		}
	}
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
	var valueToMemo string
	if len(parts) >= 2 && parts[1][0] != '@' {
		result = &parts[1]
	}
	if len(parts) >= 3 && parts[1][0] == '@' {
		result = &parts[2]
	}
	if len(parts) == 3 && parts[1][0] != '@' {
		valueToMemo = parts[2]
	}
	if len(parts) == 4 {
		if parts[1][0] == '@' {
			valueToMemo = parts[3]
		} else {
			lastParts := parts[2] + " " + parts[3]
			valueToMemo = lastParts
		}
	}
	if result == nil {
		log.Printf("no value for required field 'tag' of gedcom line.")
	}
	gedcomLine.tagMemo = result
	safeValueMemo := strconv.QuoteToASCII(valueToMemo)
	gedcomLine.valueMemo = &safeValueMemo
	return result
}

func (gedcomLine *Line) Value() *string {
	if gedcomLine.valueMemo != nil {
		return gedcomLine.valueMemo
	}
	parts := strings.SplitN(*gedcomLine.originalLine, " ", 4)
	var result string
	var tagToMemo *string
	if len(parts) >= 2 && parts[1][0] != '@' {
		tagToMemo = &parts[1]
	}
	if len(parts) >= 3 && parts[1][0] == '@' {
		tagToMemo = &parts[2]
	}
	if len(parts) == 3 && parts[1][0] != '@' {
		result = parts[2]
	}
	if len(parts) == 4 {
		if parts[1][0] == '@' {
			result = parts[3]
		} else {
			lastParts := parts[2] + " " + parts[3]
			result = lastParts
		}
	}
	gedcomLine.tagMemo = tagToMemo
	safeResult := strconv.QuoteToASCII(result)
	gedcomLine.valueMemo = &safeResult
	return &safeResult
}
