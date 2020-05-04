package model

import "sync"

type Gedcom struct {
	Lock          sync.RWMutex
	Persons       []Person
	Familys       []Family
	Childs        []Child
	SourceRepos   []string
	MasterSources []string
	Medias        []string
	FactTypes     []string
}

func NewGedcom() *Gedcom {
	return &Gedcom{
		Persons:       []Person{},
		Familys:       []Family{},
		Childs:        []Child{},
		SourceRepos:   []string{},
		MasterSources: []string{},
		Medias:        []string{},
		FactTypes:     []string{},
	}
}
