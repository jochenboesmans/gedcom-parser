package gedcom

import (
	"fmt"
	"strconv"
	"strings"
)

// structure used for parsing gedcom lines
type Line struct {
	lineString *string
	//level        int8
	//xRefID       string
	//tag          string
	//value        string
}

type GedcomFields struct {
	level  int8
	xRefID string
	tag    string
	value  string
}

func NewLine(gedcomLineString string) *Line {
	withoutNewLine := strings.TrimSuffix(gedcomLineString, "\n")
	return &Line{
		lineString: &withoutNewLine,
		//level:        -1, // use -1 instead of 0 for unset field
	}
}

// required field, must have a value >= 0 in a valid line
func (gedcomLine *Line) Level() (int8, error) {
	//if gedcomLine.level != -1 {
	//	return gedcomLine.level, nil
	//}

	if gedcomLine.lineString == nil {
		return 0, fmt.Errorf("failed to parse level from line because of improperly initialized Line struct")
	}
	parts := strings.SplitN(*gedcomLine.lineString, " ", 2)
	level, err := strconv.Atoi(parts[0])
	if err != nil {
		//gedcomLine.level = -1
		return -1, err
	}

	result := int8(level)
	//gedcomLine.level = result
	return result, nil
}

// optional field, can be an empty string in a valid line
func (gedcomLine *Line) XRefID() string {
	//if gedcomLine.xRefID != "" {
	//	return gedcomLine.xRefID
	//}

	if gedcomLine.lineString == nil {
		return ""
	}
	parts := strings.SplitN(*gedcomLine.lineString, " ", 3)
	result := ""
	if len(parts) >= 2 && parts[1][0] == '@' {
		result = parts[1]
	}
	//gedcomLine.xRefID = result
	return result
}

// required field, can't be an empty string in a valid line
func (gedcomLine *Line) Tag() (string, error) {
	//if gedcomLine.tag != "" {
	//	return gedcomLine.tag, nil
	//}

	if gedcomLine.lineString == nil {
		return "", fmt.Errorf("failed to parse tag from line because of improperly initialized Line struct")
	}
	parts := strings.SplitN(*gedcomLine.lineString, " ", 4)
	var result string
	//var valueToMemo string
	if len(parts) >= 2 && parts[1][0] != '@' {
		result = parts[1]
	}
	if len(parts) >= 3 && parts[1][0] == '@' {
		result = parts[2]
	}
	//if len(parts) == 3 && parts[1][0] != '@' {
	//	valueToMemo = parts[2]
	//}
	//if len(parts) == 4 {
	//	if parts[1][0] == '@' {
	//		valueToMemo = parts[3]
	//	} else {
	//		lastParts := parts[2] + " " + parts[3]
	//		valueToMemo = lastParts
	//	}
	//}

	if result == "" {
		return "", fmt.Errorf("no value for required field 'tag' of gedcom line")
	}

	//gedcomLine.tag = result
	//gedcomLine.value = valueToMemo
	return result, nil
}

func (gedcomLine *Line) Value() string {
	//if gedcomLine.value != "" {
	//	return gedcomLine.value
	//}

	if gedcomLine.lineString == nil {
		return ""
	}
	parts := strings.SplitN(*gedcomLine.lineString, " ", 4)
	var result string
	//var tagToMemo string
	//if len(parts) >= 2 && parts[1][0] != '@' {
	//	tagToMemo = parts[1]
	//}
	//if len(parts) >= 3 && parts[1][0] == '@' {
	//	tagToMemo = parts[2]
	//}
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
	//gedcomLine.tag = tagToMemo
	//gedcomLine.value = result
	return result
}

func (gf *GedcomFields) ToLine() (string, error) {
	var sb strings.Builder

	// level
	l := gf.level
	sb.WriteString(strconv.Itoa(int(l)))

	// xRefID
	x := gf.xRefID
	if x != "" {
		sb.WriteString(" ")
		sb.WriteString(x)
	}

	// tag
	t := gf.tag
	sb.WriteString(" ")
	sb.WriteString(strings.ToUpper(t))

	// value
	v := gf.value
	if v != "" {
		sb.WriteString(" ")
		sb.WriteString(v)
	}

	sb.WriteString("\n")

	return sb.String(), nil
}
