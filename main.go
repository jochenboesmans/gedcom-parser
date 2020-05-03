package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"strings"
)

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

type Gedcom struct {
	Persons       []Person
	Familys       []Family
	Childs        []Child
	SourceRepos   []string
	MasterSources []string
	Medias        []string
	FactTypes     []string
}

type Family struct {
	Id          string
	FatherId    string
	MotherId    string
	ChildIds    []string
	DateCreated string
}

type Child struct {
	Id                   string
	FamilyId             string
	ChildId              string
	RelationshipToFather int8
	RelationshipToMother int8
}

type Person struct {
	Id          string
	PersonRef   string
	IsLiving    bool
	Gender      int8
	DateCreated string
	Names       []Name
	Facts       []Fact
}

type Name struct {
	FactTypeId int16
	GivenNames string
	Surnames   string
}

type Fact struct {
	FactTypeId int16
	DateDetail string
	Place      Place
	Preferred  bool
}

type Place struct {
	PlaceName string
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
	scenarioFlagPtr := flag.String("scenario", "sibling", "which gedcom to parse")
	flag.Parse()
	file, err := os.Open("./test-input/" + *scenarioFlagPtr + ".ged")
	maybePanic(err)
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	gedcom := Gedcom{
		Persons:       []Person{},
		Familys:       []Family{},
		Childs:        []Child{},
		SourceRepos:   []string{},
		MasterSources: []string{},
		Medias:        []string{},
		FactTypes:     []string{},
	}
	currentLines := []string{}
	currentRecord := []string{}
	for fileScanner.Scan() {
		line := fileScanner.Text()
		words := strings.SplitN(line, " ", 3)

		// interpret record once it's fully read
		if len(currentLines) > 0 && words[0] == "0" {
			interpretRecord(&gedcom, currentLines, currentRecord)
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

	gedcomJson, err := json.MarshalIndent(gedcom, "", "  ")
	writeFile, err := os.Create("./artifacts/actual-" + *scenarioFlagPtr + ".json")

	writer := bufio.NewWriter(writeFile)
	_, err = writer.Write(gedcomJson)
	maybePanic(err)
	err = writer.Flush()
	maybePanic(err)
}

func interpretRecord(gedcom *Gedcom, recordLines []string, currentRecord []string) {
	if len(currentRecord) >= 3 && currentRecord[2] == "INDI" {
		interpretPersonRecord(gedcom, recordLines, currentRecord)
	} else if len(currentRecord) >= 3 && currentRecord[2] == "FAM" {
		interpretFamilyRecord(gedcom, recordLines, currentRecord)
	}
}

func interpretPersonRecord(gedcom *Gedcom, recordLines []string, currentRecord []string) {
	// person is assumed living unless proven to be dead
	person := Person{
		Id:       currentRecord[1],
		IsLiving: true,
		Facts:    []Fact{},
		Names:    []Name{},
	}
	for i, line := range recordLines {
		words := strings.SplitN(line, " ", 3)
		if i != 0 && words[0] == "0" {
			break
		}
		if words[0] == "1" && words[1] == "NAME" {
			name := Name{
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
			birthFact := Fact{
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
					birthFact.Place = Place{birthFactWords[2]}
				}
			}
			person.Facts = append(person.Facts, birthFact)
		}
		if words[0] == "1" && words[1] == "DEAT" {
			// TODO: Actually add death facts
			person.IsLiving = false
		}
		if words[0] == "1" && words[1] == "SEX" {
			if words[2] == "M" {
				person.Gender = 1
			} else if words[2] == "F" {
				person.Gender = 2
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
	gedcom.Persons = append(gedcom.Persons, person)
}

func interpretFamilyRecord(gedcom *Gedcom, recordLines []string, currentRecord []string) {
	family := Family{
		Id:          currentRecord[1],
		FatherId:    "",
		MotherId:    "",
		ChildIds:    []string{},
		DateCreated: "",
	}
	for i, line := range recordLines {
		words := strings.SplitN(line, " ", 3)
		if words[1] == "HUSB" {
			family.FatherId = words[2]
		}
		if words[1] == "WIFE" {
			family.MotherId = words[2]
		}
		if words[1] == "CHIL" {
			family.ChildIds = append(family.ChildIds, words[2])
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
		child := Child{
			Id:       "CHILD-" + strconv.Itoa(i) + "-" + currentRecord[1],
			FamilyId: currentRecord[1],
			ChildId:  childId,
		}
		if family.MotherId != "" {
			child.RelationshipToMother = 1
		}
		if family.FatherId != "" {
			child.RelationshipToFather = 1
		}
		gedcom.Childs = append(gedcom.Childs, child)

	}

	gedcom.Familys = append(gedcom.Familys, family)
}
