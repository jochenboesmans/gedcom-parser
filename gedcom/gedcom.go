package gedcom

import (
	"strconv"
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

func (g *ConcurrencySafeGedcom) RemoveInvalidFamilies() {
	indexedIndividuals := g.IndividualsByIds()

	familyIndicesToRemove := []int{}
familiesLoop:
	for i, f := range g.Families {
		if _, ok := indexedIndividuals[f.MotherId]; !ok {
			familyIndicesToRemove = append(familyIndicesToRemove, i)
			continue
		}
		if _, ok := indexedIndividuals[f.FatherId]; !ok {
			familyIndicesToRemove = append(familyIndicesToRemove, i)
			continue
		}
		for _, childId := range f.ChildIds {
			if _, ok := indexedIndividuals[childId]; !ok {
				familyIndicesToRemove = append(familyIndicesToRemove, i)
				continue familiesLoop
			}
		}
	}

	g.removeFamiliesAt(familyIndicesToRemove)
}

func (g *ConcurrencySafeGedcom) removeFamiliesAt(i []int) {
	g.lock()
	for _, index := range i {
		g.Families = withoutFamily(g.Families, index)
	}
	g.unlock()
}

func withoutFamily(families []*Gedcom_Family, index int) []*Gedcom_Family {
	families[len(families)-1], families[index] = families[index], families[len(families)-1]
	return families[:len(families)-1]
}

// ensures any non-utf8 chars that were encoded during parsing of original gedcom are decoded again
func (g *ConcurrencySafeGedcom) DecodeUnicodeFields() error {
	for _, i := range g.Gedcom.Individuals {
		for _, n := range i.Names {
			decodedGivenName, err := strconv.Unquote(n.GivenName)
			if err != nil {
				return err
			}
			n.GivenName = decodedGivenName

			decodedSurname, err := strconv.Unquote(n.Surname)
			if err != nil {
				return err
			}
			n.Surname = decodedSurname
		}
	}
	return nil
}
