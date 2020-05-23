package source

import "github.com/jochenboesmans/gedcom-parser/model/shared"

type Source struct {
	Id               uint32
	Data             *Data
	Originator       string
	DescriptiveTitle string
	FiledByEntry     string
	PublicationFacts string
	TextFromSource   string
	//RepositoryCitations []*shared.RepositoryCitation
	UserReferences    []*shared.UserReference
	AutomatedRecordId string
	ChangeDate        *shared.ChangeDate
	Notes             []*shared.Note
	//MultimediaLinks     []*shared.MultimediaLink
}

type Data struct {
	ResponsibleAgency string
	Notes             []*shared.Note
	Events            []*Event
}

type Event struct {
	AttributeTypes    string
	DatePeriod        string
	JurisdictionPlace string
}
