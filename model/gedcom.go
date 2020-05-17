package model

import (
	"github.com/jochenboesmans/gedcom-parser/model/child"
	"github.com/jochenboesmans/gedcom-parser/model/family"
	"github.com/jochenboesmans/gedcom-parser/model/header"
	"github.com/jochenboesmans/gedcom-parser/model/person"
	"sync"
)

type Gedcom struct {
	Lock          sync.RWMutex
	Persons       []*person.Person
	Familys       []*family.Family
	Childs        []*child.Child
	SourceRepos   []string
	MasterSources []*header.Source
	Medias        []string
	FactTypes     []string
	Header        header.Header
}

func NewGedcom() *Gedcom {
	return &Gedcom{
		Persons:       []*person.Person{},
		Familys:       []*family.Family{},
		Childs:        []*child.Child{},
		SourceRepos:   []string{},
		MasterSources: []*header.Source{},
		Medias:        []string{},
		FactTypes:     []string{},
		Header:        header.Header{},
	}
}
