package header

type Header struct {
	ReceivingSystemName string
	TransmissionDate    string
	SubmitterRecordId   string
	SubmissionRecordId  string
	FileName            string
	Copyright           string
	Metadata            GedcomMetadata
	CharacterSet        CharacterSet
	Language            string
	PlaceHierarchy      string
	ContentDescription  string
	Source              *Source
}
