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
		for j, deepHeaderLine := range headerLines[i+1:] {
			tag, err := deepHeaderLine.Tag()
			if err != nil {
				continue
			}
			switch tag {
			case "SOUR":
				h.Source = deepHeaderLine.Value()
			case "SUBM":
				h.Submitter = deepHeaderLine.Value()
			case "GEDC": // TODO: extract to function
				gedcomMetaData := &Gedcom_HeaderType_GedcomMetaDataType{
					VersionNumber: deepHeaderLine.Value(),
				}
				for _, deepGedcomMetaDataLine := range headerLines[j+1:] {
					tag, err := deepGedcomMetaDataLine.Tag()
					if err != nil {
						continue
					}
					switch tag {
					case "VERS":
						gedcomMetaData.VersionNumber = deepGedcomMetaDataLine.Value()
					case "FORM":
						gedcomMetaData.GedcomForm = deepGedcomMetaDataLine.Value()
					}
				}
				h.GedcomMetaData = gedcomMetaData
			case "CHAR":
				h.CharacterSet = deepHeaderLine.Value()
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
// * SUBMITTER_RECORD (SUBM)
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
		g.interpretNoteRecord(recordLines)
	case "REPO":
		g.interpretRepositoryRecord(recordLines)
	case "SOUR":
		g.interpretSourceRecord(recordLines)
	case "SUBM":
		g.interpretSubmitterRecord(recordLines)
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

func (g *ConcurrencySafeGedcom) interpretNoteRecord(recordLines []*Line) {
	xRefID, submitterText := recordLines[0].XRefID(), recordLines[0].Value()
	note := Gedcom_Note{
		Id:            xRefID,
		SubmitterText: submitterText,
	}
	g.lock()
	g.Gedcom.Notes = append(g.Gedcom.Notes, &note)
	g.unlock()
}

func (g *ConcurrencySafeGedcom) interpretMultimediaRecord(recordLines []*Line) {
	xRefID := recordLines[0].XRefID()
	multimedia := Gedcom_Multimedia{
		Id:    xRefID,
		Files: []*Gedcom_Multimedia_File{},
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
		case "FILE":
			// TODO: extract to function
			reference := line.Value()
			file := Gedcom_Multimedia_File{
				Reference: reference,
			}
			for _, fileLine := range recordLines[i+1:] {
				fileLevel, err := fileLine.Level()
				if err != nil {
					continue
				}
				if fileLevel <= rootLevel+1 {
					break // end of record
				}
				fileTag, err := fileLine.Tag()
				if err != nil {
					continue
				}
				switch fileTag {
				case "FORM":
					file.Format = fileLine.Value()
				}
			}
			multimedia.Files = append(multimedia.Files, &file)
		}
	}
	g.lock()
	g.Multimedias = append(g.Multimedias, &multimedia)
	g.unlock()
}

func (g *ConcurrencySafeGedcom) interpretRepositoryRecord(recordLines []*Line) {
	xRefID := recordLines[0].XRefID()
	repository := Gedcom_Repository{
		Id: xRefID,
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
		case "NAME":
			repository.Name = line.Value()
		}
	}
	g.lock()
	g.Gedcom.Repositories = append(g.Gedcom.Repositories, &repository)
	g.unlock()
}

func (g *ConcurrencySafeGedcom) interpretSourceRecord(recordLines []*Line) {
	xRefID := recordLines[0].XRefID()
	source := Gedcom_Source{
		Id: xRefID,
	}
	g.lock()
	g.Gedcom.Sources = append(g.Gedcom.Sources, &source)
	g.unlock()
}

func (g *ConcurrencySafeGedcom) interpretSubmitterRecord(recordLines []*Line) {
	xRefID := recordLines[0].XRefID()
	submitterInstance := Gedcom_Submitter{
		Id: xRefID,
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
		case "NAME":
			submitterInstance.Name = line.Value()
		}
	}
	g.lock()
	g.Gedcom.Submitters = append(g.Gedcom.Submitters, &submitterInstance)
	g.unlock()

}

func logError(firstLine *Line, structureKind string, err error) {
	l := firstLine.lineString
	log.Printf("failed to interpret %s structure starting with %s with error: %s\n", structureKind, *l, err)
}
