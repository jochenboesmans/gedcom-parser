package gedcom

import (
	"fmt"
	"strconv"
	"strings"
)

// structure used for parsing gedcom lines
// holds a ref to the original line as well as memos for each part, allowing for lazy parsing
type Line struct {
	originalLine *string
	levelMemo    int8
	xRefIDMemo   string
	tagMemo      string
	valueMemo    string
}

func NewLine(gedcomLinePtr *string) *Line {
	return &Line{
		originalLine: gedcomLinePtr,
		levelMemo:    -1, // use -1 instead of 0 for unset field
	}
}

// required field, must have a value >= 0 in a valid line
func (gedcomLine *Line) Level() (int8, error) {
	if gedcomLine.levelMemo != -1 {
		return gedcomLine.levelMemo, nil
	}

	parts := strings.SplitN(*gedcomLine.originalLine, " ", 2)
	level, err := strconv.Atoi(parts[0])
	if err != nil {
		gedcomLine.levelMemo = -1
		return -1, err
	}

	result := int8(level)
	gedcomLine.levelMemo = result
	return result, nil
}

// optional field, can be an empty string in a valid line
func (gedcomLine *Line) XRefID() string {
	if gedcomLine.xRefIDMemo != "" {
		return gedcomLine.xRefIDMemo
	}

	parts := strings.SplitN(*gedcomLine.originalLine, " ", 3)
	result := ""
	if len(parts) >= 2 && parts[1][0] == '@' {
		result = parts[1]
	}
	gedcomLine.xRefIDMemo = result
	return result
}

// required field, can't be an empty string in a valid line
func (gedcomLine *Line) Tag() (string, error) {
	if gedcomLine.tagMemo != "" {
		return gedcomLine.tagMemo, nil
	}

	parts := strings.SplitN(*gedcomLine.originalLine, " ", 4)
	var result string
	var valueToMemo string
	if len(parts) >= 2 && parts[1][0] != '@' {
		result = parts[1]
	}
	if len(parts) >= 3 && parts[1][0] == '@' {
		result = parts[2]
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

	if result == "" {
		return "", fmt.Errorf("no value for required field 'tag' of gedcom line")
	}

	gedcomLine.tagMemo = result
	safeValueMemo := strconv.QuoteToASCII(valueToMemo) // encode unicode chars in line values
	gedcomLine.valueMemo = safeValueMemo
	return result, nil
}

func (gedcomLine *Line) Value() string {
	if gedcomLine.valueMemo != "" {
		return gedcomLine.valueMemo
	}

	parts := strings.SplitN(*gedcomLine.originalLine, " ", 4)
	var result string
	var tagToMemo string
	if len(parts) >= 2 && parts[1][0] != '@' {
		tagToMemo = parts[1]
	}
	if len(parts) >= 3 && parts[1][0] == '@' {
		tagToMemo = parts[2]
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
	safeResult := strconv.QuoteToASCII(result) // encode unicode chars in line values
	gedcomLine.valueMemo = safeResult
	return safeResult
}
