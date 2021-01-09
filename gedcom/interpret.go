package gedcom

import (
	"github.com/jochenboesmans/gedcom-parser/util"
	"strconv"
	"strings"
	"sync"
)

func (g *ConcurrencySafeGedcom) InterpretRecord(recordLines []*Line, waitGroup *sync.WaitGroup) {
	tag, err := recordLines[0].Tag()
	if err != nil {
		return
	}
	switch tag {
	case "INDI":
		g.interpretIndividualRecord(recordLines)
	case "FAM":
		g.interpretFamilyRecord(recordLines)
	}
	waitGroup.Done()
}

func (g *ConcurrencySafeGedcom) interpretIndividualRecord(recordLines []*Line) {
	individualXRefID := recordLines[0].XRefID()
	individualInstance := Gedcom_Individual{
		Id: individualXRefID,
	}
	for i, line := range recordLines[1:] {
		level, err := line.Level()
		if err != nil {
			continue
		}
		if level < 1 {
			break // end of record
		}

		tag, err := line.Tag()
		if err != nil {
			continue
		}
		switch tag {
		case "NAME":
			g.interpretName(line, recordLines, i, &individualInstance)
		case "SEX":
			g.interpretSexLine(line, &individualInstance)
		case "BIRT":
			g.interpretIndividualEvent(recordLines, i, &individualInstance, "BIRT")
		case "DEAT":
			g.interpretIndividualEvent(recordLines, i, &individualInstance, "DEAT")
		}
	}
	g.lock()
	g.Gedcom.Individuals = append(g.Gedcom.Individuals, &individualInstance)
	g.unlock()

}

func (g *ConcurrencySafeGedcom) interpretIndividualEvent(recordLines []*Line, i int, individualInstance *Gedcom_Individual, kind string) {
	e := Gedcom_Individual_Event{
		Date:    nil,
		Place:   "",
		Primary: false,
	}
	for _, eventLine := range recordLines[i+1:] {
		level, err := eventLine.Level()
		if err != nil {
			continue
		}
		if level < 2 {
			break // end of event structure
		}

		tag, err := eventLine.Tag()
		if err != nil {
			continue
		}
		if tag == "DATE" {
			birthDate := parseDate(eventLine)
			e.Date = birthDate
		}
		if tag == "PLAC" {
			e.Place = eventLine.Value()
		}
		if tag == "_PRIM" {
			if primaryBool, ok := util.PrimaryBoolByValue[eventLine.Value()]; ok {
				e.Primary = primaryBool
			}
		}
	}

	switch kind {
	case "BIRT":
		individualInstance.BirthEvents = append(individualInstance.BirthEvents, &e)
	case "DEAT":
		individualInstance.DeathEvents = append(individualInstance.DeathEvents, &e)
	}
}

func (g *ConcurrencySafeGedcom) interpretName(line *Line, recordLines []*Line, i int, individualInstance *Gedcom_Individual) {
	name := Gedcom_Individual_Name{}
	if nameParts := strings.Split(line.Value(), "/"); nameParts[0] != "" || nameParts[1] != "" {
		name.GivenName = nameParts[0]
		name.Surname = nameParts[1]
	}
	for _, nameLine := range recordLines[i+1:] {
		level, err := nameLine.Level()
		if err != nil {
			continue
		}
		if level < 2 {
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

	if name.GivenName != "" || name.Surname != "" {
		individualInstance.Names = append(individualInstance.Names, &name)
	}

}

func (g *ConcurrencySafeGedcom) interpretSexLine(line *Line, individualInstance *Gedcom_Individual) {
	if genderFull, ok := util.GenderFullByLetter[line.Value()]; ok {
		individualInstance.Gender = genderFull
	}
}

func parseDate(line *Line) *Gedcom_Individual_Date {
	dateParts := strings.SplitN(line.Value(), " ", 3)
	date := &Gedcom_Individual_Date{}
	if len(dateParts) > 0 {
		if year, err := strconv.Atoi(dateParts[0]); err == nil {
			date.Year = uint32(year)
		}
	}
	if len(dateParts) > 1 {
		if month, ok := util.MonthIntByAbbr[strings.ToUpper(dateParts[1])]; ok {
			date.Month = uint32(month)
		}
	}
	if len(dateParts) > 2 {
		if day, err := strconv.Atoi(dateParts[2]); err == nil {
			date.Day = uint32(day)
		}
	}
	return date
}

func (g *ConcurrencySafeGedcom) interpretFamilyRecord(recordLines []*Line) {
	familyId := recordLines[0].XRefID()
	familyInstance := Gedcom_Family{Id: familyId}
	for _, line := range recordLines {
		level, err := line.Level()
		if err != nil {
			continue
		}
		if level < 1 {
			break // end of record
		}

		tag, err := line.Tag()
		if err != nil {
			continue
		}
		switch tag {
		case "HUSB":
			familyInstance.FatherId = line.Value()
		case "WIFE":
			familyInstance.MotherId = line.Value()
		case "CHIL":
			familyInstance.ChildIds = append(familyInstance.ChildIds, line.Value())
		}
	}
	g.lock()
	g.Gedcom.Families = append(g.Gedcom.Families, &familyInstance)
	g.unlock()
}
