package main

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
)

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

type Gedcom struct {
	Persons       []Person
	Familys       []string
	Childs        []string
	SourceRepos   []string
	MasterSources []string
	Medias        []string
	FactTypes     []string
}

type Person struct {
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
	file, err := os.Open("one-node.ged")
	maybePanic(err)
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	gedcom := Gedcom{
		Persons:       []Person{},
		Familys:       []string{},
		Childs:        []string{},
		SourceRepos:   []string{},
		MasterSources: []string{},
		Medias:        []string{},
		FactTypes:     []string{},
	}
	// person is assumed living unless proven to be dead
	person := Person{
		IsLiving: true,
	}
	currentLines := []string{}
	openStack := false
	for fileScanner.Scan() {
		line := fileScanner.Text()
		words := strings.SplitN(line, " ", 3)
		if words[0] == "1" && (words[1] == "NAME" || words[1] == "BIRT" || words[1] == "DEATH" || words[1] == "CHAN" || words[1] == "SEX") {
			openStack = true
		}
		if openStack && words[0] == "0" {
			openStack = false
			for i, line := range currentLines {
				words := strings.SplitN(line, " ", 3)
				if words[0] == "1" && words[1] == "NAME" {
					name := Name{
						FactTypeId: 100,
					}
					for _, nameLine := range currentLines[i+1 : len(currentLines)] {
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
					for _, birthFactLine := range currentLines[i+1 : len(currentLines)] {
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
					for _, chanLine := range currentLines[i+1 : len(currentLines)] {
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
			}
		}
		if openStack {
			currentLines = append(currentLines, line)
		}
	}
	gedcom.Persons = append(gedcom.Persons, person)

	gedcomJson, err := json.MarshalIndent(gedcom, "", "  ")
	writeFile, err := os.Create("actual.json")

	writer := bufio.NewWriter(writeFile)
	_, err = writer.Write(gedcomJson)
	maybePanic(err)
	err = writer.Flush()
	maybePanic(err)
}
