package main

import (
	"bufio"
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
	"github.com/pquerna/ffjson/ffjson"
)

type OutputGedcom struct {
	Persons             []model.Person
	Familys             []model.Family
	Childs              []model.Child
	SourceRepos         []string
	MasterSources       []model.Source
	Medias              []string
	FactTypes           []string
	ReceivingSystemName string
	TransmissionDate    string
	SubmitterRecordId   string
	FileName            string
	Copyright           string
	Metadata            model.GedcomMetadata
	CharacterSet        model.CharacterSet
	Language            string
	PlaceHierarchy      string
	ContentDescription  string
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

	pathToGedcomFile := flag.String("pathToGedcomFile", "./test-input/hugetree.ged", "relative path to input gedcom file (with .ged extension if present)")
	pathToJsonFile := flag.String("pathToJsonFile", "./artifacts/actual-hugetree.json", "relative path to output json file (with .json extension if wanted)")
	flag.Parse()

	file, err := os.Open(*pathToGedcomFile)
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

	gedcomJson, err := ffjson.Marshal(gedcomWithoutLock)
	writeFile, err := os.Create(*pathToJsonFile)

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
	case "HEAD":
		interpretHeadRecord(gedcom, currentRecordDeepLines, currentRecordLine)
		// case "NOTE":
		// case "REPO":
		// case "SOUR":
		// case "SUBN":
		// case "SUBM":
		// case "TRLR":
	}
	waitGroup.Done()
}

func interpretHeadRecord(gedcom *model.Gedcom, currentRecordDeepLines []*gedcomSpec.Line, currentRecordLine *gedcomSpec.Line) {
	for i, line := range currentRecordDeepLines {
		if i != 0 && line.Level() == 0 {
			break
		}
		if line.Level() == 1 {
			switch line.Tag() {
			case "SOUR":
				source := model.Source{
					ApprovedSystemId: line.Value(),
				}
				for j, sourceLine := range currentRecordDeepLines[i+1:] {
					if sourceLine.Level() < 2 {
						break
					}
					switch sourceLine.Tag() {
					case "VERS":
						source.Version = sourceLine.Value()
					case "NAME":
						source.ProductName = sourceLine.Value()
					case "CORP":
						corporation := model.SourceCorporation{
							Name: sourceLine.Value(),
						}
						for k, corpLine := range currentRecordDeepLines[i+1+j+1:] {
							if corpLine.Level() < 3 {
								break
							}
							switch corpLine.Tag() {
							case "ADDR":
								address := model.Address{
									MainLine: sourceLine.Value(),
								}
								for _, addrLine := range currentRecordDeepLines[i+1+j+1+k+1:] {
									if addrLine.Level() < 4 {
										break
									}
									switch addrLine.Tag() {
									case "CITY":
										address.City = addrLine.Value()
									case "POST":
										address.PostCode = addrLine.Value()
									case "CTRY":
										address.Country = addrLine.Value()
									}
								}
								corporation.Address = address
							case "WWW":
								corporation.WebsiteURL = corpLine.Value()
							}
						}
						source.Corporation = corporation
					}
				}
				gedcom.Lock.Lock()
				gedcom.MasterSources = append(gedcom.MasterSources, source)
				gedcom.Lock.Unlock()
			case "DATE":
				date := ""
				dateParts := strings.SplitN(line.Value(), " ", 3)
				if len(dateParts) >= 1 {
					date = dateParts[0]
				}
				if len(dateParts) >= 2 {
					date = monthNumberByAbbreviation[dateParts[1]] + "-" + date
				}
				if len(dateParts) >= 3 {
					date = dateParts[2] + "-" + date
				}

				timeLine := currentRecordDeepLines[i+1]
				date += "T" + timeLine.Value()

				gedcom.Lock.Lock()
				gedcom.TransmissionDate = date
				gedcom.Lock.Unlock()
			case "DEST":
				gedcom.Lock.Lock()
				gedcom.ReceivingSystemName = line.Value()
				gedcom.Lock.Unlock()
			case "SUBM":
				gedcom.Lock.Lock()
				// TODO: ID-ify xrefid (hash or whatever)
				gedcom.SubmitterRecordId = line.Value()
				gedcom.Lock.Unlock()
			case "SUBN":
				gedcom.Lock.Lock()
				// TODO: ID-ify xrefid (hash or whatever)
				gedcom.SubmissionRecordId = line.Value()
				gedcom.Lock.Unlock()
			case "FILE":
				gedcom.Lock.Lock()
				gedcom.FileName = line.Value()
				gedcom.Lock.Unlock()
			case "COPR":
				gedcom.Lock.Lock()
				gedcom.Copyright = line.Value()
				gedcom.Lock.Unlock()
			case "GEDC":
				metadata := model.GedcomMetadata{}
				for _, gedcLine := range currentRecordDeepLines[i+1:] {
					if gedcLine.Level() < 2 {
						break
					}
					switch gedcLine.Value() {
					case "VERS":
						metadata.Version = gedcLine.Value()
					case "FORM":
						metadata.Form = gedcLine.Value()
					}
				}
				gedcom.Lock.Lock()
				gedcom.Metadata = metadata
				gedcom.Lock.Unlock()
			case "CHAR":
				characterSet := model.CharacterSet{
					Value: line.Value(),
				}
				if len(currentRecordDeepLines) > i+1 {
					characterSet.Version = currentRecordDeepLines[i+1].Value()
				}
				gedcom.Lock.Lock()
				gedcom.CharacterSet = characterSet
				gedcom.Lock.Unlock()
			case "LANG":
				gedcom.Lock.Lock()
				gedcom.Language = line.Value()
				gedcom.Lock.Unlock()
			case "PLAC":
				gedcom.Lock.Lock()
				gedcom.PlaceHierarchy = line.Value()
				gedcom.Lock.Unlock()
			case "NOTE":
				note := line.Value()
				for _, noteLine := range currentRecordDeepLines[i+1:] {
					switch noteLine.Tag() {
					case "CONT":
					case "CONC":
						note += " " + noteLine.Value()
					}
				}
				gedcom.Lock.Lock()
				gedcom.ContentDescription = note
				gedcom.Lock.Unlock()
			}
		}
	}
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
					if nameLine.Level() < 2 {
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
