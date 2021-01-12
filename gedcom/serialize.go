package gedcom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jochenboesmans/gedcom-parser/util"
	"log"
)

func (gedcom *ConcurrencySafeGedcom) ToJson() (*[]byte, error) {
	gedcomJson, err := json.Marshal(&gedcom.Gedcom)
	if err != nil {
		return nil, err
	}
	return &gedcomJson, nil
}

func (gedcom *ConcurrencySafeGedcom) ToProto() (*[]byte, error) {
	gedcomProto, err := proto.Marshal(&gedcom.Gedcom)
	if err != nil {
		return nil, err
	}
	return &gedcomProto, nil
}

func writeLine(line *Line, buf *bytes.Buffer, lineCounter *int) error {
	lineString, err := line.ToString()
	if err != nil {
		return fmt.Errorf("failed to serialize line %d with error: %s", *lineCounter, err)
	}
	buf.WriteString(lineString)
	*lineCounter++
	return nil
}

func createAndWriteLine(level int, xRefID string, tag string, value string, lineCounter *int, buf *bytes.Buffer) error {
	line := &Line{
		level:  int8(level),
		tag:    tag,
		xRefID: xRefID,
		value:  value,
	}
	err := writeLine(line, buf, lineCounter)
	return err
}

func (g *ConcurrencySafeGedcom) ToSerializedGedcom() (*bytes.Buffer, error) {
	gedcom := g.Gedcom
	buf := bytes.NewBuffer([]byte{})
	lineCounter := 0
	rootLevel := 0

	err := createAndWriteLine(rootLevel, "", "HEAD", "", &lineCounter, buf)
	if err != nil {
		// completely fail write if header write fails
		return nil, err
	}

	for _, i := range gedcom.Individuals {
		indiLevel := rootLevel
		err := createAndWriteLine(indiLevel, i.Id, "INDI", "", &lineCounter, buf)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, n := range i.Names {
			nameLevel := indiLevel + 1
			err := createAndWriteLine(nameLevel, "", "NAME", "", &lineCounter, buf)
			if err != nil {
				log.Println(err)
				continue
			}

			if n.GivenName != "" {
				givenNameLevel := nameLevel + 1
				err := createAndWriteLine(givenNameLevel, "", "GIVN", n.GivenName, &lineCounter, buf)
				if err != nil {
					log.Println(err)
				}
			}
			if n.Surname != "" {
				surnameLevel := nameLevel + 1
				err := createAndWriteLine(surnameLevel, "", "SURN", n.Surname, &lineCounter, buf)
				if err != nil {
					log.Println(err)
				}
			}
			if primValue, ok := util.PrimaryValueByBool[n.Primary]; ok {
				primLevel := nameLevel + 1
				err := createAndWriteLine(primLevel, "", "_PRIM", primValue, &lineCounter, buf)
				if err != nil {
					log.Println(err)
				}
			}
		}

		for _, b := range i.BirthEvents {
			eventLevel := indiLevel + 1
			err := createAndWriteLine(eventLevel, "", "BIRT", "", &lineCounter, buf)
			if err != nil {
				log.Println(err)
				continue
			}
			createAndWriteDeepEventLines(b, eventLevel, &lineCounter, buf)
		}

		for _, d := range i.DeathEvents {
			eventLevel := indiLevel + 1
			err := createAndWriteLine(eventLevel, "", "DEAT", "", &lineCounter, buf)
			if err != nil {
				log.Println(err)
				continue
			}
			createAndWriteDeepEventLines(d, eventLevel, &lineCounter, buf)
		}

		if genderLetter, hit := util.GenderLetterByFull[i.Gender]; hit {
			genderLevel := indiLevel + 1
			err := createAndWriteLine(genderLevel, "", "SEX", genderLetter, &lineCounter, buf)
			if err != nil {
				log.Println(err)
			}
		}
	}

	for _, f := range gedcom.Families {
		familyLevel := rootLevel
		err := createAndWriteLine(familyLevel, f.Id, "FAM", "", &lineCounter, buf)
		if err != nil {
			log.Println(err)
			continue
		}

		if f.FatherId != "" {
			fatherLevel := familyLevel + 1
			err := createAndWriteLine(fatherLevel, "", "HUSB", f.FatherId, &lineCounter, buf)
			if err != nil {
				log.Println(err)
			}
		}
		if f.MotherId != "" {
			motherLevel := familyLevel + 1
			err := createAndWriteLine(motherLevel, "", "WIFE", f.MotherId, &lineCounter, buf)
			if err != nil {
				log.Println(err)
			}
		}

		for _, childId := range f.ChildIds {
			childLevel := familyLevel + 1
			err := createAndWriteLine(childLevel, "", "CHIL", childId, &lineCounter, buf)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}

	err = createAndWriteLine(rootLevel, "", "TRLR", "", &lineCounter, buf)
	if err != nil {
		log.Println(err)
	}

	return buf, nil
}

func createAndWriteDeepEventLines(event *Gedcom_Individual_Event, eventLevel int, lineCounter *int, buf *bytes.Buffer) {
	var dateValue string
	if event.Date.Year != "" && event.Date.Month != "" && event.Date.Day != "" {
		dateValue = fmt.Sprintf("%s %s %s", event.Date.Day, util.MonthAbbrByInt[event.Date.Month], event.Date.Year)
	} else if event.Date.Year != "" && event.Date.Month != "" {
		dateValue = fmt.Sprintf("%s %s", util.MonthAbbrByInt[event.Date.Month], event.Date.Year)
	} else if event.Date.Year != "" {
		dateValue = fmt.Sprintf("%s", event.Date.Year)
	}
	if dateValue != "" {
		dateLevel := eventLevel + 1
		err := createAndWriteLine(dateLevel, "", "DATE", dateValue, lineCounter, buf)
		if err != nil {
			log.Println(err)
		}
	}

	if event.Place != "" {
		placeLevel := eventLevel + 1
		err := createAndWriteLine(placeLevel, "", "PLAC", event.Place, lineCounter, buf)
		if err != nil {
			log.Println(err)
		}
	}

	if primValue, ok := util.PrimaryValueByBool[event.Primary]; ok {
		primLevel := eventLevel + 1
		err := createAndWriteLine(primLevel, "", "_PRIM", primValue, lineCounter, buf)
		if err != nil {
			log.Println(err)
		}
	}

}
