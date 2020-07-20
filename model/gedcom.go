package model

import (
	"github.com/jochenboesmans/gedcom-parser/model/family"
	"github.com/jochenboesmans/gedcom-parser/model/individual"
	"sync"
)

type ConcurrencySafeGedcom struct {
	Gedcom
	Lock sync.RWMutex
}

type Gedcom struct {
	Individuals []*individual.Individual
	Families    []*family.Family
}

type NoPointerGedcom struct {
	Individuals []individual.NoPointerIndividual
	Families    []family.NoPointerFamily
}
