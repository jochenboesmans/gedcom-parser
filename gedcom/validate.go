package gedcom

import (
	"fmt"
	"log"
)

func (g *ConcurrencySafeGedcom) ValidateIdUniqueness() {
	alternativeId := 999999 // TODO: assign definitely unique alternative ids
	ids := map[string]bool{}
	for i, indi := range g.Individuals {
		if alreadyInMap := ids[indi.Id]; alreadyInMap {
			g.lock()
			g.Individuals[i].Id = fmt.Sprintf("@I%d@", alternativeId)
			g.unlock()
			alternativeId--
		}
	}
}

func contains(submitters []*Gedcom_Submitter, xRefId string) bool {
	for _, s := range submitters {
		if s.Id == xRefId {
			return true
		}
	}
	return false
}

func (g *ConcurrencySafeGedcom) ValidateHeaderXRefIntegrity() {
	submitterXRefId := g.Header.Submitter
	if submitterXRefId == "" {
		return
	}
	if !contains(g.Submitters, submitterXRefId) {
		if len(g.Submitters) > 0 {
			alternative := g.Submitters[0].Id
			log.Printf("invalid submitter xRefId in header (%s), defaulting to %s", submitterXRefId, alternative)
		} else {
			log.Printf("invalid submitter xRefId in header (%s), no alternative found, removing submitter xRefId from header", submitterXRefId)
		}
	}

}

// ValidateFamilyRecordXRefIdIntegrity ensures integrity of cross references to indi records in family records
// COST WARNING: O(f*c) where f is the amount of family records and c is the amount of children in a family
func (g *ConcurrencySafeGedcom) ValidateFamilyRecordXRefIdIntegrity() {
	indexedIndividuals := g.IndividualsByIds()

	for i, f := range g.Families {
		if _, ok := indexedIndividuals[f.MotherId]; !ok {
			g.lock()
			g.Families[i].MotherId = ""
			g.unlock()
		}
		if _, ok := indexedIndividuals[f.FatherId]; !ok {
			g.lock()
			g.Families[i].FatherId = ""
			g.unlock()
		}
		for j, childId := range f.ChildIds {
			if _, ok := indexedIndividuals[childId]; !ok {
				g.lock()
				g.Families[i].ChildIds[j] = ""
				g.unlock()
			}
		}
	}
}

func (g *ConcurrencySafeGedcom) Validate() {
	g.ValidateIdUniqueness()
	g.ValidateHeaderXRefIntegrity()
	g.ValidateFamilyRecordXRefIdIntegrity()
}
