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

	gedcomSpec "github.com/jochenboesmans/gedcom-parser/gedcom"
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

	currentRecordDeepLines := []*gedcomSpec.Line{}
	var currentRecordLine *gedcomSpec.Line
	waitGroup := &sync.WaitGroup{}

	gedcom := model.NewGedcom()

	i := 0
	for fileScanner.Scan() {
		line := ""
		if i == 0 {
			line = strings.TrimPrefix(fileScanner.Text(), "\uFEFF")
		} else {
			line = fileScanner.Text()
		}
		gedcomLine := gedcomSpec.NewLine(line)

		// interpret record once it's fully read
		if currentRecordLine != nil && gedcomLine.Level() == 0 {
			waitGroup.Add(1)
			go interpretRecord(gedcom, currentRecordDeepLines, currentRecordLine, waitGroup)
			currentRecordLine = nil
			currentRecordDeepLines = []*gedcomSpec.Line{}
		}
		if gedcomLine.Level() == 0 {
			currentRecordLine = gedcomLine
		}
		if currentRecordLine != nil {
			currentRecordDeepLines = append(currentRecordDeepLines, gedcomLine)
		}
		i++
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

func interpretRecord(gedcom *model.Gedcom, currentRecordDeepLines []*gedcomSpec.Line, currentRecordLine *gedcomSpec.Line, waitGroup *sync.WaitGroup) {
	switch currentRecordLine.Tag() {
	case "INDI":
		interpretPersonRecord(gedcom, currentRecordDeepLines, currentRecordLine)
	case "FAM":
		interpretFamilyRecord(gedcom, currentRecordDeepLines, currentRecordLine)
		// case "HEAD":
		// case "NOTE":
		// case "REPO":
		// case "SOUR":
		// case "SUBN":
		// case "SUBM":
		// case "TRLR":
	}
	waitGroup.Done()
}

func interpretPersonRecord(gedcom *model.Gedcom, currentRecordDeepLines []*gedcomSpec.Line, currentRecordLine *gedcomSpec.Line) {
	person := model.NewPerson(currentRecordLine.XRefID())
	for i, line := range currentRecordDeepLines {
		if i != 0 && line.Level() == 0 {
			break
		}
		if line.Level() == 1 {
			switch line.Tag() {
			case "NAME":
				name := model.PersonName{
					FactTypeId: 100,
				}
				for _, nameLine := range currentRecordDeepLines[i+1:] {
					if nameLine.Level() == 0 {
						break
					}
					switch nameLine.Tag() {
					case "GIVN":
						name.GivenNames = nameLine.Value()
					case "SURN":
						name.Surnames = nameLine.Value()
					}
				}
				if name.GivenNames != "" || name.Surnames != "" {
					person.Names = append(person.Names, name)
				}
			case "BIRT":
				birthFact := model.PersonFact{
					FactTypeId: 405,
				}
				for _, birthFactLine := range currentRecordDeepLines[i+1:] {
					if birthFactLine.Level() < 2 {
						break
					}
					switch birthFactLine.Tag() {
					case "_PRIM":
						switch birthFactLine.Value() {
						case "Y":
							birthFact.Preferred = true
						case "N":
							birthFact.Preferred = false
						}
					case "DATE":
						birthFact.DateDetail = birthFactLine.Value()
					case "PLAC":
						birthFact.Place = model.PersonPlace{PlaceName: birthFactLine.Value()}
					}
				}
				person.Facts = append(person.Facts, birthFact)
			//TODO: case "DEAT":
			case "SEX":
				switch line.Value() {
				case "M":
					person.Gender = 1
				case "F":
					person.Gender = 2
				}
			case "CHAN":
				date := ""
				for _, chanLine := range currentRecordDeepLines[i+1:] {
					if chanLine.Level() < 2 {
						break
					}
					switch chanLine.Tag() {
					case "DATE":
						dateParts := strings.SplitN(chanLine.Value(), " ", 3)
						if len(dateParts) >= 1 {
							date = dateParts[0]
						}
						if len(dateParts) >= 2 {
							date = monthNumberByAbbreviation[dateParts[1]] + "-" + date
						}
						if len(dateParts) >= 3 {
							date = dateParts[2] + "-" + date
						}
					case "TIME":
						date += "T" + chanLine.Value()
					}
				}
				if date != "" {
					person.DateCreated = date
				}
			case "_UID":
				person.PersonRef = line.Value()
			}
		}
	}
	gedcom.Lock.Lock()
	gedcom.Persons = append(gedcom.Persons, *person)
	gedcom.Lock.Unlock()
}

func interpretFamilyRecord(gedcom *model.Gedcom, currentRecordDeepLines []*gedcomSpec.Line, currentRecordLine *gedcomSpec.Line) {
	family := model.NewFamily(currentRecordLine.XRefID())
	for i, line := range currentRecordDeepLines {
		if i != 0 && line.Level() == 0 {
			break
		}
		switch line.Tag() {
		case "HUSB":
			fatherId, err := util.Hash(line.Value())
			util.MaybePanic(err)
			family.FatherId = fatherId
		case "WIFE":
			motherId, err := util.Hash(line.Value())
			util.MaybePanic(err)
			family.MotherId = motherId
		case "CHIL":
			childId, err := util.Hash(line.Value())
			util.MaybePanic(err)
			family.ChildIds = append(family.ChildIds, childId)
		case "CHAN":
			date := ""
			for _, chanLine := range currentRecordDeepLines[i+1:] {
				if chanLine.Level() < 2 {
					break
				}
				switch chanLine.Tag() {
				case "DATE":
					dateParts := strings.SplitN(chanLine.Value(), " ", 3)
					if len(dateParts) >= 1 {
						date = dateParts[0]
					}
					if len(dateParts) >= 2 {
						date = monthNumberByAbbreviation[dateParts[1]] + "-" + date
					}
					if len(dateParts) >= 3 {
						date = dateParts[2] + "-" + date
					}
				case "TIME":
					date += "T" + chanLine.Value()
				}
			}
			if date != "" {
				family.DateCreated = date
			}
		}
	}

	for i, childId := range family.ChildIds {
		child := model.NewChild(currentRecordLine.XRefID(), i, childId)
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
