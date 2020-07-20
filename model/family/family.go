package family

type Family struct {
	Id       *string
	FatherId *string
	MotherId *string
	ChildIds []*string
}

type NoPointerFamily struct {
	Id       string
	FatherId string
	MotherId string
	ChildIds []string
}

func NewFamily(identificationString *string) Family {
	return Family{
		Id:       identificationString,
		ChildIds: []*string{},
	}
}
