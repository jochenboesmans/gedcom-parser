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

func (g *ConcurrencySafeGedcom) ToSerializedGedcom() (*bytes.Buffer, error) {
	gedcom := g.Gedcom
	buf := bytes.NewBuffer([]byte{})
	lineCounter := 0

	rootLevel := 0
	headerLine := &Line{
		level: int8(rootLevel),
		tag:   "HEAD",
	}
	err := writeLine(headerLine, buf, &lineCounter)
	if err != nil {
		return nil, err
	}

	for _, i := range gedcom.Individuals {
		firstLine := &Line{
			level:  int8(rootLevel),
			xRefID: i.Id,
			tag:    "INDI",
		}
		err := writeLine(firstLine, buf, &lineCounter)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, n := range i.Names {
			nameLevel := rootLevel + 1
			nameLine := &Line{
				level: int8(nameLevel),
				tag:   "NAME",
			}
			err := writeLine(nameLine, buf, &lineCounter)
			if err != nil {
				log.Println(err)
				continue
			}

			if n.GivenName != "" {
				givenNameLine := &Line{
					level: int8(nameLevel + 1),
					tag:   "GIVN",
					value: n.GivenName,
				}
				err := writeLine(givenNameLine, buf, &lineCounter)
				if err != nil {
					log.Println(err)
				}
			}
			if n.Surname != "" {
				surnameLine := &Line{
					level: int8(nameLevel + 1),
					tag:   "SURN",
					value: n.Surname,
				}
				err := writeLine(surnameLine, buf, &lineCounter)
				if err != nil {
					log.Println(err)
				}
			}
			if primValue, ok := util.PrimaryValueByBool[n.Primary]; ok {
				primLine := &Line{
					level: int8(nameLevel + 1),
					tag:   "_PRIM",
					value: primValue,
				}
				err := writeLine(primLine, buf, &lineCounter)
				if err != nil {
					log.Println(err)
				}
			}
		}

		for _, b := range i.BirthEvents {
			firstLine := fmt.Sprintf("1 BIRT\n")
			buf.WriteString(firstLine)

			var secondLine string
			if b.Date.Year != "" && b.Date.Month != "" && b.Date.Day != "" {
				secondLine = fmt.Sprintf("2 DATE %s %s %s\n", b.Date.Day, util.MonthAbbrByInt[b.Date.Month], b.Date.Year)
			} else if b.Date.Year != "" && b.Date.Month != "" {
				secondLine = fmt.Sprintf("2 DATE %s %s\n", util.MonthAbbrByInt[b.Date.Month], b.Date.Year)
			} else if b.Date.Year != "" {
				secondLine = fmt.Sprintf("2 DATE %s\n", b.Date.Year)
			}
			if secondLine != "" {
				buf.WriteString(firstLine)
			}

			if b.Place != "" {
				placeLine := fmt.Sprintf("2 PLAC %s\n", b.Place)
				buf.WriteString(placeLine)
			}

			primaryLine := fmt.Sprintf("2 _PRIM %s\n", util.PrimaryValueByBool[b.Primary])
			buf.WriteString(primaryLine)
		}

		for _, d := range i.DeathEvents {
			firstLine := fmt.Sprintf("1 DEAT\n")
			buf.WriteString(firstLine)

			var secondLine string
			if d.Date.Year != "" && d.Date.Month != "" && d.Date.Day != "" {
				secondLine = fmt.Sprintf("2 DATE %s %s %s\n", d.Date.Day, util.MonthAbbrByInt[d.Date.Month], d.Date.Year)
			} else if d.Date.Year != "" && d.Date.Month != "" {
				secondLine = fmt.Sprintf("2 DATE %s %s\n", util.MonthAbbrByInt[d.Date.Month], d.Date.Year)
			} else if d.Date.Year != "" {
				secondLine = fmt.Sprintf("2 DATE %s\n", d.Date.Year)
			}
			if secondLine != "" {
				buf.WriteString(firstLine)
			}

			if d.Place != "" {
				placeLine := fmt.Sprintf("2 PLAC %s\n", d.Place)
				buf.WriteString(placeLine)
			}

			if primaryValue, hit := util.PrimaryValueByBool[d.Primary]; hit {
				primaryLine := fmt.Sprintf("2 _PRIM %s\n", primaryValue)
				buf.WriteString(primaryLine)
			}
		}

		if genderLetter, hit := util.GenderLetterByFull[i.Gender]; hit {
			genderLine := fmt.Sprintf("1 SEX %s\n", genderLetter)
			buf.WriteString(genderLine)
		}
	}

	for _, f := range gedcom.Families {
		firstLine := fmt.Sprintf("0 %s FAM\n", f.Id)
		buf.WriteString(firstLine)

		if f.FatherId != "" {
			fatherLine := fmt.Sprintf("1 HUSB %s\n", f.FatherId)
			buf.WriteString(fatherLine)
		}
		if f.MotherId != "" {
			motherLine := fmt.Sprintf("1 WIFE %s\n", f.MotherId)
			buf.WriteString(motherLine)
		}

		for _, childId := range f.ChildIds {
			childLine := fmt.Sprintf("1 CHIL %s\n", childId)
			buf.WriteString(childLine)
		}
	}

	trailer := "0 TRLR\n"
	buf.WriteString(trailer)

	return buf, nil
}
