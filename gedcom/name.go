package gedcom

import (
	"errors"
	"github.com/jochenboesmans/gedcom-parser/util"
	"strings"
)

type Name struct {
	GivenName string
	Surname   string
	Primary   bool
}

func NameStructure(nameLines []*Line) (*Name, error) {
	rootLevel, err := nameLines[0].Level()
	if err != nil {
		return nil, errors.New("failed to parse root level of name structure")
	}

	name := Name{}
	if nameParts := strings.Split(nameLines[0].Value(), "/"); nameParts[0] != "" || len(nameParts) > 1 && nameParts[1] != "" {
		name.GivenName = strings.TrimSpace(nameParts[0])
		if len(nameParts) > 1 {
			name.Surname = strings.TrimSpace(nameParts[1])
		}
	}
	for _, nameLine := range nameLines[1:] {
		level, err := nameLine.Level()
		if err != nil {
			continue // continue searching
		}
		if level <= rootLevel {
			break // end  of name structure
		}

		tag, err := nameLine.Tag()
		if err != nil {
			continue
		}
		switch tag {
		case "GIVN":
			name.GivenName = nameLine.Value()
		case "SURN":
			name.Surname = nameLine.Value()
		case "_PRIM":
			name.Primary = util.PrimaryBoolByValue[strings.ToUpper(nameLine.Value())]
		}
	}
	return &name, nil
}
