package gedcom

import (
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
				name := Gedcom_Individual_Name{}
				nameParts := strings.Split(*line.Value(), "/")
				if nameParts[0] != "" || nameParts[1] != "" {
					name.GivenName = nameParts[0]
					name.Surname = nameParts[1]
				} else {
					for _, nameLine := range recordLines[i+1:] {
						if *nameLine.Level() < 2 {
							break
						}
						switch *nameLine.Tag() {
						case "GIVN":
							name.GivenName = *nameLine.Value()
						case "SURN":
							name.Surname = *nameLine.Value()
						}
					}
				}
				if name.GivenName != "" || name.Surname != "" {
					individualInstance.Names = append(individualInstance.Names, &name)
				}
			case "SEX":
				if line.Value() != nil {
					switch *line.Value() {
					case "M":
						individualInstance.Gender = "MALE"
					case "F":
						individualInstance.Gender = "FEMALE"
					}
				}
			case "BIRT":
				for _, birthLine := range recordLines[i+1:] {
					if *birthLine.Level() < 2 {
						break
					}
					if *recordLines[i+1].Tag() == "DATE" {
						birthDate := parseDate(recordLines[i+1])
						individualInstance.BirthDate = birthDate
					}
				}
			case "DEAT":
				for _, deathLine := range recordLines[i+1:] {
					if *deathLine.Level() < 2 {
						break
					}
					if *recordLines[i+1].Tag() == "DATE" {
						deathDate := parseDate(recordLines[i+1])
						individualInstance.DeathDate = deathDate
					}
				}
			}
		}
	}
	g.lock()
	g.Gedcom.Individuals = append(g.Gedcom.Individuals, &individualInstance)
	g.unlock()
}

func parseDate(line *Line) *Gedcom_Individual_Date {
	monthIntByAbbr := map[string]int{
		"JAN": 1,
		"FEB": 2,
		"MAR": 3,
		"APR": 4,
		"MAY": 5,
		"JUN": 6,
		"JUL": 7,
		"AUG": 8,
		"SEP": 9,
		"OCT": 10,
		"NOV": 11,
		"DEC": 12,
	}
	dateParts := strings.SplitN(*line.Value(), " ", 3)
	date := &Gedcom_Individual_Date{}
	if len(dateParts) > 0 {
		if year, err := strconv.Atoi(dateParts[0]); err == nil {
			date.Year = uint32(year)
		}
	}
	if len(dateParts) > 1 {
		if month, ok := monthIntByAbbr[strings.ToUpper(dateParts[1])]; ok {
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

	for i, f := range g.Families {
		if _, ok := indexedIndividuals[f.MotherId]; !ok {
			g.removeFamilyAt(i)
		}

		if _, ok := indexedIndividuals[f.FatherId]; !ok {
			g.removeFamilyAt(i)
		}
		for _, childId := range f.ChildIds {
			if _, ok := indexedIndividuals[childId]; !ok {
				g.removeFamilyAt(i)
			}
		}
	}
}

func (g *ConcurrencySafeGedcom) removeFamilyAt(i int) {
	g.rwlock.Lock()
	g.Families = withoutFamily(g.Families, i)
	g.rwlock.Unlock()
}

func withoutFamily(families []*Gedcom_Family, index int) []*Gedcom_Family {
	families[len(families)-1], families[index] = families[index], families[len(families)-1]
	return families[:len(families)-1]
}
