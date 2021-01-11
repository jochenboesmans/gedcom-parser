package gedcom

import (
	"fmt"
	"github.com/jochenboesmans/gedcom-parser/util"
	"strings"
)

type Name struct {
	GivenName string
	Surname   string
	Primary   bool
}

func interpretNameStructure(nameLines []*Line) (*Name, error) {
	rootLevel, err := nameLines[0].Level()
	if err != nil {
		return nil, fmt.Errorf("failed to parse root level of name structure: %s", err)
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

func (name *Name) IsEmpty() bool {
	return name.GivenName == "" && name.Surname == ""
}

func (name *Name) toGedcomIndividualName() Gedcom_Individual_Name {
	return Gedcom_Individual_Name{
		GivenName: name.GivenName,
		Surname:   name.Surname,
		Primary:   name.Primary,
	}
}
