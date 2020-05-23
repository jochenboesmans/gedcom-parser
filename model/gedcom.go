package model

import (
	"github.com/jochenboesmans/gedcom-parser/model/child"
	"github.com/jochenboesmans/gedcom-parser/model/family"
	"github.com/jochenboesmans/gedcom-parser/model/header"
	"github.com/jochenboesmans/gedcom-parser/model/multimedia"
	"github.com/jochenboesmans/gedcom-parser/model/note"
	"github.com/jochenboesmans/gedcom-parser/model/person"
	"github.com/jochenboesmans/gedcom-parser/model/repository"
	"github.com/jochenboesmans/gedcom-parser/model/source"
	"github.com/jochenboesmans/gedcom-parser/model/submission"
	"github.com/jochenboesmans/gedcom-parser/model/submitter"
	"sync"
)

type Gedcom struct {
	Lock        sync.RWMutex
	Header      *header.Header
	Submission  *submission.Submission
	Persons     []*person.Person
	Familys     []*family.Family
	Childs      []*child.Child
	Notes       []*note.Note
	Repositorys []*repository.Repository
	Sources     []*source.Source
	Submitters  []*submitter.Submitter
	Multimedias []*multimedia.Multimedia
}

func NewGedcom() *Gedcom {
	return &Gedcom{
		Persons:     []*person.Person{},
		Familys:     []*family.Family{},
		Childs:      []*child.Child{},
		Notes:       []*note.Note{},
		Repositorys: []*repository.Repository{},
		Sources:     []*source.Source{},
		Submitters:  []*submitter.Submitter{},
		Multimedias: []*multimedia.Multimedia{},
	}
}
