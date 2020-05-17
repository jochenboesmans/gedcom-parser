package submission

import "github.com/jochenboesmans/gedcom-parser/model/shared"

type Submission struct {
	SubmitterId              uint32
	Id                       uint32
	NameOfFamilyFile         string
	TempleCode               string
	GenerationsOfAncestors   uint32
	GenerationsOfDescendants uint32
	OrdinanceProcessFlag     bool
	AutomatedRecordId        string
	Notes                    []*shared.Note
	ChangeDate               *shared.ChangeDate
}
