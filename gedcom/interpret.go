package gedcom

import (
	"log"
	"sync"
)

// InterpretHeader interprets all gedcom metadata.
// It must be executed before any other record interpretation.
func (g *ConcurrencySafeGedcom) InterpretHeader(headerLines []*Line) error {
	h := &Gedcom_HeaderType{}
	for i, headerLine := range headerLines {
		tag, err := headerLine.Tag()
		if err != nil || tag != "HEAD" {
			continue // search lines until HEAD is found
		}
		for _, deepHeaderLine := range headerLines[i+1:] {
			tag, err := deepHeaderLine.Tag()
			if err != nil {
				continue
			}
			switch tag {
			case "SOUR":
				h.Source = deepHeaderLine.Value()
			}
		}
		break
	}
	g.lock()
	g.Header = h
	g.unlock()
	return nil
}

// InterpretRecord concurrently interprets a top-level GEDCOM record
// and puts the interpreted data in the ConcurrencySafeGedcom
//
// Records can be one of:
//
// * FAM_RECORD (FAM)
//
// * INDIVIDUAL_RECORD (INDI)
//
// * MULTIMEDIA_RECORD (OBJE)
//
// * NOTE_RECORD (NOTE)
//
// * REPOSITORY_RECORD (REPO)
//
// * SOURCE_RECORD (SOUR)
//
// * SUBMITTER_RECORD (SUBN)
//
func (g *ConcurrencySafeGedcom) InterpretRecord(recordLines []*Line, waitGroup *sync.WaitGroup) {
	tag, err := recordLines[0].Tag()
	if err != nil {
		return
	}
	switch tag {
	case "FAM":
		g.interpretFamilyRecord(recordLines)
	case "INDI":
		g.interpretIndividualRecord(recordLines)
	case "OBJE":
		// TODO
	case "NOTE":
		// TODO
	case "REPO":
		// TODO
	case "SOUR":
		// TODO
	case "SUBN":
		// TODO
	}
	waitGroup.Done()
}

func (g *ConcurrencySafeGedcom) interpretIndividualRecord(recordLines []*Line) {
	individualXRefID := recordLines[0].XRefID()
	individualInstance := Gedcom_Individual{
		Id: individualXRefID,
	}
	rootLevel, err := recordLines[0].Level()
	if err != nil {
		return
	}
	for i, line := range recordLines[1:] {
		level, err := line.Level()
		if err != nil {
			continue
		}
		if level <= rootLevel {
			break // end of record
		}

		tag, err := line.Tag()
		if err != nil {
			continue
		}
		switch tag {
		case "NAME":
			g.interpretIndividualName(recordLines[1+i:], &individualInstance)
		case "SEX":
			g.interpretIndividualSex(recordLines[1+i:], &individualInstance)
		case "BIRT":
			g.interpretIndividualEvent(recordLines[1+i:], &individualInstance, "BIRT")
		case "DEAT":
			g.interpretIndividualEvent(recordLines[1+i:], &individualInstance, "DEAT")
		}
	}
	g.lock()
	g.Gedcom.Individuals = append(g.Gedcom.Individuals, &individualInstance)
	g.unlock()

}

func (g *ConcurrencySafeGedcom) interpretIndividualSex(recordLines []*Line, individualInstance *Gedcom_Individual) {
	genderFull, err := interpretSexStructure(recordLines[0])
	if err != nil {
		logError(recordLines[0], "sex", err)
		return
	}
	individualInstance.Gender = genderFull
}

func (g *ConcurrencySafeGedcom) interpretIndividualEvent(recordLines []*Line, individualInstance *Gedcom_Individual, kind string) {
	event, err := interpretEventStructure(recordLines)
	if err != nil {
		logError(recordLines[0], "event", err)
		return
	}

	gedcomIndividualEvent := event.toGedcomIndividualEvent()
	switch kind {
	case "BIRT":
		individualInstance.BirthEvents = append(individualInstance.BirthEvents, &gedcomIndividualEvent)
	case "DEAT":
		individualInstance.DeathEvents = append(individualInstance.DeathEvents, &gedcomIndividualEvent)
	}
}

func (g *ConcurrencySafeGedcom) interpretIndividualName(recordLines []*Line, individualInstance *Gedcom_Individual) {
	name, err := interpretNameStructure(recordLines)
	if err != nil || name.IsEmpty() {
		logError(recordLines[0], "name", err)
		return
	}

	gedcomIndividualName := name.toGedcomIndividualName()
	individualInstance.Names = append(individualInstance.Names, &gedcomIndividualName)
}

func (g *ConcurrencySafeGedcom) interpretFamilyRecord(recordLines []*Line) {
	familyId := recordLines[0].XRefID()
	familyInstance := Gedcom_Family{
		Id: familyId,
	}
	rootLevel, err := recordLines[0].Level()
	if err != nil {
		return
	}
	for _, line := range recordLines[1:] {
		level, err := line.Level()
		if err != nil {
			continue
		}
		if level <= rootLevel {
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

func logError(firstLine *Line, structureKind string, err error) {
	l, toStringErr := firstLine.ToString()
	if toStringErr != nil {
		log.Printf("failed to interpret %s structure with error: %s\n", structureKind, err)
	} else {
		log.Printf("failed to interpret %s structure starting with %s with error: %s\n", structureKind, l, err)
	}
}
