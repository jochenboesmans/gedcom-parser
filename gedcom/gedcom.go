package gedcom

import (
	"sync"
)

type ConcurrencySafeGedcom struct {
	Gedcom
	rwlock sync.RWMutex
}

func NewConcurrencySafeGedcom() *ConcurrencySafeGedcom {
	return &ConcurrencySafeGedcom{
		Gedcom: Gedcom{},
		rwlock: sync.RWMutex{},
	}
}

func (g *ConcurrencySafeGedcom) lock() {
	g.rwlock.Lock()
}

func (g *ConcurrencySafeGedcom) unlock() {
	g.rwlock.Unlock()
}

func (g *ConcurrencySafeGedcom) IndividualsByIds() map[string]*Gedcom_Individual {
	result := map[string]*Gedcom_Individual{}
	for _, i := range g.Individuals {
		result[i.Id] = i
	}
	return result
}

func (g *ConcurrencySafeGedcom) FamiliesByIds() map[string]*Gedcom_Family {
	result := map[string]*Gedcom_Family{}
	for _, i := range g.Families {
		result[i.Id] = i
	}
	return result
}
