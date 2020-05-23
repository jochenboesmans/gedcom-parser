package source

import "github.com/jochenboesmans/gedcom-parser/model/shared"

type Source struct {
	Id               uint32
	Events           []*Event
	Originator       string
	DescriptiveTitle string
	FiledByEntry     string
	PublicationFacts string
	TextFromSource   string
	//RepositoryCitations []*shared.RepositoryCitation
	UserReference     *shared.UserReference
	AutomatedRecordId string
	ChangeDate        *shared.ChangeDate
	Notes             []*shared.Note
	//MultimediaLinks     []*shared.MultimediaLink
}

type Event struct {
	DatePeriod        string
	JurisdictionPlace string
	ResponsibleAgency string
	Notes             []*shared.Note
}
