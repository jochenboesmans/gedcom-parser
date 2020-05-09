package model

import "sync"

type Gedcom struct {
	Lock                sync.RWMutex
	Persons             []Person
	Familys             []Family
	Childs              []Child
	SourceRepos         []string
	MasterSources       []Source
	Medias              []string
	FactTypes           []string
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
}

func NewGedcom() *Gedcom {
	return &Gedcom{
		Persons:       []Person{},
		Familys:       []Family{},
		Childs:        []Child{},
		SourceRepos:   []string{},
		MasterSources: []Source{},
		Medias:        []string{},
		FactTypes:     []string{},
	}
}
