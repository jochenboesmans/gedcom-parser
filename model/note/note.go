package note

type Note struct {
	Id                uint32
	SubmitterText     string
	UserReferences    []*UserReference
	AutomatedRecordId string
	SourceCitation    []*string
	ChangeDate        string
}

type UserReference struct {
	Number string
	Type   string
}
