package submitter

import "github.com/jochenboesmans/gedcom-parser/model/shared"

type Submitter struct {
	Id      uint32
	Name    string
	Address *shared.Address
	//MultimediaLinks        []*shared.MultimediaLink
	LanguagePreference     []string
	SubmitterRegisteredRFN string
	AutomatedRecordId      string
	Notes                  []*shared.Note
	ChangeDate             *shared.ChangeDate
}
