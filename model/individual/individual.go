package individual

type Individual struct {
	Id     *string
	Names  []*Name
	Gender string
	//PersonalNameStructures []*PersonalNameStructure
	//SexValue string
	//IndividualEventStructures []*IndividualEventStructure
	//ChildToFamilyLinks []*ChildToFamilyLink
	//SpouseToFamilyLinks []*SpouseToFamilyLink
	//AssociationStructures []*AssociationStructure
}

type NoPointerIndividual struct {
	Id     string
	Names  []Name
	Gender string
}

type Name struct {
	GivenName string
	Surname   string
}

func NewIndividual(identificationString *string) Individual {
	return Individual{
		Id: identificationString,
	}
}
