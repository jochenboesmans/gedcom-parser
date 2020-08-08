package gedcom

import (
	"github.com/jochenboesmans/gedcom-parser/util"
	"strconv"
	"strings"
	"sync"
)

type ConcurrencySafeGedcom struct {
	Gedcom
	rwlock sync.RWMutex
}

func NewConcurrencySafeGedcom() *ConcurrencySafeGedcom {
	return &ConcurrencySafeGedcom{
		Gedcom: Gedcom{},
		rwlock: sync.RWMutex{},
	}
}

func (g *ConcurrencySafeGedcom) InterpretRecord(recordLines []*Line, waitGroup *sync.WaitGroup) {
	switch *recordLines[0].Tag() {
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
		Id: *individualXRefID,
	}
	for i, line := range recordLines {
		if i != 0 && *line.Level() == 0 {
			break
		}
		if *line.Level() == 1 {
			switch *line.Tag() {
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
		if *eventLine.Level() < 2 {
			break
		}
		if *eventLine.Tag() == "DATE" {
			birthDate := parseDate(eventLine)
			e.Date = birthDate
		}
		if *eventLine.Tag() == "PLAC" {
			e.Place = *eventLine.Value()
		}
		if *eventLine.Tag() == "_PRIM" {
			if primaryValue, ok := util.PrimaryBoolByValue[*eventLine.Value()]; ok {
				e.Primary = primaryValue
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
	if line.Value() != nil {
		if nameParts := strings.Split(*line.Value(), "/"); nameParts[0] != "" || nameParts[1] != "" {
			name.GivenName = nameParts[0]
			name.Surname = nameParts[1]
		}
	} else {
		for _, nameLine := range recordLines[i+1:] {
			if *nameLine.Level() < 2 {
				break
			}
			switch *nameLine.Tag() {
			case "GIVN":
				if nameLine.Value() != nil {
					name.GivenName = *nameLine.Value()
				}
			case "SURN":
				if nameLine.Value() != nil {
					name.Surname = *nameLine.Value()
				}
			}
		}
	}
	for _, nameLine := range recordLines[i+1:] {
		if *nameLine.Level() < 2 {
			break
		}
		if *nameLine.Tag() == "_PRIM" {
			name.Primary = util.PrimaryBoolByValue[strings.ToUpper(*nameLine.Value())]
		}
	}
	if name.GivenName != "" || name.Surname != "" {
		individualInstance.Names = append(individualInstance.Names, &name)
	}

}

func (g *ConcurrencySafeGedcom) interpretSexLine(line *Line, individualInstance *Gedcom_Individual) {
	if line.Value() != nil {
		switch *line.Value() {
		case "M":
			individualInstance.Gender = "MALE"
		case "F":
			individualInstance.Gender = "FEMALE"
		}
	}
}

func parseDate(line *Line) *Gedcom_Individual_Date {
	dateParts := strings.SplitN(*line.Value(), " ", 3)
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
	familyInstance := Gedcom_Family{Id: *familyId}
	for i, line := range recordLines {
		if i != 0 && *line.Level() == 0 {
			break
		}
		switch *line.Tag() {
		case "HUSB":
			if line.Value() != nil {
				fatherId := line.Value()
				familyInstance.FatherId = *fatherId
			}
		case "WIFE":
			if line.Value() != nil {
				motherId := line.Value()
				familyInstance.MotherId = *motherId
			}
		case "CHIL":
			if line.Value() != nil {
				childId := line.Value()
				familyInstance.ChildIds = append(familyInstance.ChildIds, *childId)
			}
		}
	}
	g.lock()
	g.Gedcom.Families = append(g.Gedcom.Families, &familyInstance)
	g.unlock()
}

func (g *ConcurrencySafeGedcom) lock() {
	g.rwlock.Lock()
}

func (g *ConcurrencySafeGedcom) unlock() {
	g.rwlock.Unlock()
}

func (g *ConcurrencySafeGedcom) IndividualsByIds() map[string]*Gedcom_Individual {
	result := map[string]*Gedcom_Individual{}
	for _, i := range g.Individuals {
		result[i.Id] = i
	}
	return result
}

func (g *ConcurrencySafeGedcom) RemoveInvalidFamilies() {
	indexedIndividuals := g.IndividualsByIds()

	familyIndicesToRemove := []int{}
familiesLoop:
	for i, f := range g.Families {
		if _, ok := indexedIndividuals[f.MotherId]; !ok {
			familyIndicesToRemove = append(familyIndicesToRemove, i)
			continue
		}
		if _, ok := indexedIndividuals[f.FatherId]; !ok {
			familyIndicesToRemove = append(familyIndicesToRemove, i)
			continue
		}
		for _, childId := range f.ChildIds {
			if _, ok := indexedIndividuals[childId]; !ok {
				familyIndicesToRemove = append(familyIndicesToRemove, i)
				continue familiesLoop
			}
		}
	}

	g.removeFamiliesAt(familyIndicesToRemove)
}

func (g *ConcurrencySafeGedcom) removeFamiliesAt(i []int) {
	g.lock()
	for _, index := range i {
		g.Families = withoutFamily(g.Families, index)
	}
	g.unlock()
}

func withoutFamily(families []*Gedcom_Family, index int) []*Gedcom_Family {
	families[len(families)-1], families[index] = families[index], families[len(families)-1]
	return families[:len(families)-1]
}

// ensures any non-utf8 chars that were encoded during parsing of original gedcom are decoded again
func (g *ConcurrencySafeGedcom) DecodeUnicodeFields() error {
	for _, i := range g.Gedcom.Individuals {
		for _, n := range i.Names {
			decodedGivenName, err := strconv.Unquote(n.GivenName)
			if err != nil {
				return err
			}
			n.GivenName = decodedGivenName

			decodedSurname, err := strconv.Unquote(n.Surname)
			if err != nil {
				return err
			}
			n.Surname = decodedSurname
		}
	}
	return nil
}
