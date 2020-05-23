package note

import "github.com/jochenboesmans/gedcom-parser/model/shared"

type Note struct {
	Id                uint32
	SubmitterText     string
	UserReferences    []*shared.UserReference
	AutomatedRecordId string
	//SourceCitation    []*shared.SourceCitation
	ChangeDate *shared.ChangeDate
}
