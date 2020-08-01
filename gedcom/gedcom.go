package gedcom

import (
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
			}
		}
	}
	g.lock()
	g.Gedcom.Individuals = append(g.Gedcom.Individuals, &individualInstance)
	g.unlock()
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
