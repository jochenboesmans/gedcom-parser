package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jochenboesmans/gedcom-parser/model"
	"github.com/jochenboesmans/gedcom-parser/util"
)

type OutputGedcom struct {
	Persons       []model.Person
	Familys       []model.Family
	Childs        []model.Child
	SourceRepos   []string
	MasterSources []string
	Medias        []string
	FactTypes     []string
}

var monthNumberByAbbreviation = map[string]string{
	"JAN": "01",
	"FEB": "02",
	"MAR": "03",
	"APR": "04",
	"MAY": "05",
	"JUN": "06",
	"JUL": "07",
	"AUG": "08",
	"SEP": "09",
	"OCT": "10",
	"NOV": "11",
	"DEC": "12",
}

func main() {
	startTime := time.Now()

	scenarioFlagPtr := flag.String("scenario", "sibling", "which gedcom to parse")
	flag.Parse()

	file, err := os.Open("./test-input/" + *scenarioFlagPtr + ".ged")
	util.MaybePanic(err)
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	currentLines := []string{}
	currentRecord := []string{}
	waitGroup := &sync.WaitGroup{}

	gedcom := model.NewGedcom()

	for fileScanner.Scan() {
		line := fileScanner.Text()
		words := strings.SplitN(line, " ", 3)

		// interpret record once it's fully read
		if len(currentLines) > 0 && words[0] == "0" {
			waitGroup.Add(1)
			go interpretRecord(gedcom, currentLines, currentRecord, waitGroup)
			currentRecord = []string{}
			currentLines = []string{}
		}
		if words[0] == "0" && len(words) >= 3 && (words[2] == "INDI" || words[2] == "FAM") {
			currentRecord = words
		}
		if len(currentRecord) > 0 {
			currentLines = append(currentLines, line)
		}
	}

	waitGroup.Wait()

	gedcomWithoutLock := OutputGedcom{
		Persons:       gedcom.Persons,
		Childs:        gedcom.Childs,
		Familys:       gedcom.Familys,
		SourceRepos:   gedcom.SourceRepos,
		MasterSources: gedcom.MasterSources,
		Medias:        gedcom.Medias,
		FactTypes:     gedcom.FactTypes,
	}

	gedcomJson, err := json.MarshalIndent(gedcomWithoutLock, "", "  ")
	writeFile, err := os.Create("./artifacts/actual-" + *scenarioFlagPtr + ".json")

	writer := bufio.NewWriter(writeFile)
	_, err = writer.Write(gedcomJson)
	util.MaybePanic(err)
	err = writer.Flush()
	util.MaybePanic(err)

	fmt.Printf("done in %f seconds.", float64(time.Since(startTime))*math.Pow10(-9))
}

func interpretRecord(gedcom *model.Gedcom, recordLines []string, currentRecord []string, waitGroup *sync.WaitGroup) {
	if len(currentRecord) >= 3 && currentRecord[2] == "INDI" {
		interpretPersonRecord(gedcom, recordLines, currentRecord)
	} else if len(currentRecord) >= 3 && currentRecord[2] == "FAM" {
		interpretFamilyRecord(gedcom, recordLines, currentRecord)
	}
	waitGroup.Done()
}

func interpretPersonRecord(gedcom *model.Gedcom, recordLines []string, currentRecord []string) {
	person := model.NewPerson(currentRecord[1])
	for i, line := range recordLines {
		words := strings.SplitN(line, " ", 3)
		if i != 0 && words[0] == "0" {
			break
		}
		if words[0] == "1" && words[1] == "NAME" {
			name := model.PersonName{
				FactTypeId: 100,
			}
			for _, nameLine := range recordLines[i+1 : len(recordLines)] {
				nameWords := strings.SplitN(nameLine, " ", 3)
				if nameWords[0] == "1" {
					break
				}
				if nameWords[1] == "GIVN" {
					name.GivenNames = nameWords[2]
				}
				if nameWords[1] == "SURN" {
					name.Surnames = nameWords[2]
				}
				//if birthFactWords[1] == "_PRIM" {
				//	if birthFactWords[2] == "Y" {
				//		birthFact.Preferred = true
				//	} else if birthFactWords[2] == "N" {
				//		birthFact.Preferred = false
				//  }
				//}
			}
			if name.GivenNames != "" || name.Surnames != "" {
				person.Names = append(person.Names, name)
			}
		}
		if words[0] == "1" && words[1] == "BIRT" {
			birthFact := model.PersonFact{
				FactTypeId: 405,
			}
			for _, birthFactLine := range recordLines[i+1 : len(recordLines)] {
				birthFactWords := strings.SplitN(birthFactLine, " ", 3)
				if birthFactWords[0] == "1" {
					break
				}
				if birthFactWords[1] == "_PRIM" {
					if birthFactWords[2] == "Y" {
						birthFact.Preferred = true
					} else if birthFactWords[2] == "N" {
						birthFact.Preferred = false
					}
				}
				if birthFactWords[1] == "DATE" {
					birthFact.DateDetail = birthFactWords[2]
				}
				if birthFactWords[1] == "PLAC" {
					birthFact.Place = model.PersonPlace{PlaceName: birthFactWords[2]}
				}
			}
			person.Facts = append(person.Facts, birthFact)
		}
		if words[0] == "1" && words[1] == "DEAT" {
			// TODO: Actually add death facts
			person.IsLiving = false
		}
		if words[0] == "1" && words[1] == "SEX" {
			if len(words) > 2 {
				if words[2] == "M" {
					person.Gender = 1
				} else if words[2] == "F" {
					person.Gender = 2
				}
			}
		}
		if words[0] == "1" && words[1] == "CHAN" {
			date := ""
			for _, chanLine := range recordLines[i+1 : len(recordLines)] {
				chanLineWords := strings.SplitN(chanLine, " ", 3)
				if chanLineWords[0] == "1" {
					break
				}
				if chanLineWords[1] == "DATE" {
					dateParts := strings.SplitN(chanLineWords[2], " ", 3)
					date += dateParts[2] + "-" + monthNumberByAbbreviation[dateParts[1]] + "-" + dateParts[0]
				}
				if chanLineWords[1] == "TIME" {
					date += "T" + chanLineWords[2]
				}
			}
			if date != "" {
				person.DateCreated = date
			}
		}
		if words[1] == "_UID" {
			person.PersonRef = words[2]
		}
	}
	gedcom.Lock.Lock()
	gedcom.Persons = append(gedcom.Persons, *person)
	gedcom.Lock.Unlock()
}

func interpretFamilyRecord(gedcom *model.Gedcom, recordLines []string, currentRecord []string) {
	family := model.NewFamily(currentRecord[1])
	for i, line := range recordLines {
		words := strings.SplitN(line, " ", 3)
		if words[1] == "HUSB" {
			fatherId, err := util.Hash(words[2])
			util.MaybePanic(err)
			family.FatherId = fatherId
		}
		if words[1] == "WIFE" {
			motherId, err := util.Hash(words[2])
			util.MaybePanic(err)
			family.MotherId = motherId
		}
		if words[1] == "CHIL" {
			childId, err := util.Hash(words[2])
			util.MaybePanic(err)
			family.ChildIds = append(family.ChildIds, childId)
		}
		if words[0] == "1" && words[1] == "CHAN" {
			date := ""
			for _, chanLine := range recordLines[i+1 : len(recordLines)] {
				chanLineWords := strings.SplitN(chanLine, " ", 3)
				if chanLineWords[0] == "1" {
					break
				}
				if chanLineWords[1] == "DATE" {
					dateParts := strings.SplitN(chanLineWords[2], " ", 3)
					date += dateParts[2] + "-" + monthNumberByAbbreviation[dateParts[1]] + "-" + dateParts[0]
				}
				if chanLineWords[1] == "TIME" {
					date += "T" + chanLineWords[2]
				}
			}
			if date != "" {
				family.DateCreated = date
			}
		}
	}

	for i, childId := range family.ChildIds {
		child := model.NewChild(currentRecord[1], i, childId)
		if family.MotherId != 0 {
			child.RelationshipToMother = 1
		}
		if family.FatherId != 0 {
			child.RelationshipToFather = 1
		}
		gedcom.Lock.Lock()
		gedcom.Childs = append(gedcom.Childs, child)
		gedcom.Lock.Unlock()

	}

	gedcom.Lock.Lock()
	gedcom.Familys = append(gedcom.Familys, family)
	gedcom.Lock.Unlock()
}
