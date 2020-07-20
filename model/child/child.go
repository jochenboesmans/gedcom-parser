package child

type Child struct {
	FamilyId             *string
	ChildId              *string
	RelationshipToFather bool
	RelationshipToMother bool
}

func NewChild(familyId *string, identificationString *string) Child {
	return Child{
		FamilyId: familyId,
		ChildId:  identificationString,
	}
}
