package multimedia

import "github.com/jochenboesmans/gedcom-parser/model/shared"

type Multimedia struct {
	Id                uint32
	Files             []*File
	UserReferences    []*shared.UserReference
	AutomatedRecordId string
	// Notes []*shared.Note
	// Sources []*shared.Source
	// ChangeDate *shared.ChangeDate
}

type File struct {
	Reference string
	Format    string
	Type      string
	Title     string
}
