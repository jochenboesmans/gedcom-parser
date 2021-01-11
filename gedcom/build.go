package gedcom

import (
	"bytes"
	"fmt"
	"github.com/jochenboesmans/gedcom-parser/util"
)

// TODO: replace with method in gedcom/serialize.go
func WritableGedcom(concSafeGedcom *ConcurrencySafeGedcom) *bytes.Buffer {
	// try to decode non-utf8 fields, keep encoded version if it fails
	_ = concSafeGedcom.DecodeUnicodeFields()

	gedcom := concSafeGedcom.Gedcom
	buf := bytes.NewBuffer([]byte{})

	header := "0 HEAD\n"
	buf.WriteString(header)

	for _, i := range gedcom.Individuals {
		firstLine := fmt.Sprintf("0 %s INDI\n", i.Id)
		buf.WriteString(firstLine)

		for _, n := range i.Names {
			nameLine := fmt.Sprintf("1 NAME %s/%s/\n", n.GivenName, n.Surname)
			buf.WriteString(nameLine)

			primaryLine := fmt.Sprintf("2 _PRIM %s\n", util.PrimaryValueByBool[n.Primary])
			buf.WriteString(primaryLine)
		}

		for _, b := range i.BirthEvents {
			firstLine := fmt.Sprintf("1 BIRT\n")
			buf.WriteString(firstLine)

			var secondLine string
			if b.Date.Year != 0 && b.Date.Month != 0 && b.Date.Day != 0 {
				secondLine = fmt.Sprintf("2 DATE %d %s %d\n", b.Date.Day, util.MonthAbbrByInt[int(b.Date.Month)], b.Date.Year)
			} else if b.Date.Year != 0 && b.Date.Month != 0 {
				secondLine = fmt.Sprintf("2 DATE %s %d\n", util.MonthAbbrByInt[int(b.Date.Month)], b.Date.Year)
			} else if b.Date.Year != 0 {
				secondLine = fmt.Sprintf("2 DATE %d\n", b.Date.Year)
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
			if d.Date.Year != 0 && d.Date.Month != 0 && d.Date.Day != 0 {
				secondLine = fmt.Sprintf("2 DATE %d %s %d\n", d.Date.Day, util.MonthAbbrByInt[int(d.Date.Month)], d.Date.Year)
			} else if d.Date.Year != 0 && d.Date.Month != 0 {
				secondLine = fmt.Sprintf("2 DATE %s %d\n", util.MonthAbbrByInt[int(d.Date.Month)], d.Date.Year)
			} else if d.Date.Year != 0 {
				secondLine = fmt.Sprintf("2 DATE %d\n", d.Date.Year)
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

	return buf
}
