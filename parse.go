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
	Persons             []*model.Person
	Familys             []*model.Family
	Childs              []*model.Child
	SourceRepos         []string
	MasterSources       []*model.Source
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
	ContentDescription  *string
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

var personTime time.Duration
var familyTime time.Duration
var headTime time.Duration

func main() {
	pathToGedcomFile := flag.String("pathToGedcomFile", "./test-input/hugetree.ged", "relative path to input gedcom file (with .ged extension if present)")
	pathToJsonFile := flag.String("pathToJsonFile", "./artifacts/hugetree.json", "relative path to output json file (with .json extension if wanted)")
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
		gedcomLine := gedcomSpec.NewLine(&line)

		// interpret record once it's fully read
		if currentRecordLine != nil && *gedcomLine.Level() == 0 {
			waitGroup.Add(1)
			go interpretRecord(gedcom, currentRecordDeepLines, currentRecordLine, waitGroup)
			currentRecordLine = nil
			currentRecordDeepLines = []*gedcomSpec.Line{}
		}
		if *gedcomLine.Level() == 0 {
			currentRecordLine = gedcomLine
		}
		if currentRecordLine != nil {
			currentRecordDeepLines = append(currentRecordDeepLines, gedcomLine)
		}
		i++
	}

	waitGroup.Wait()
	fmt.Printf("interpreted head in %f second.\n", float64(headTime)*math.Pow10(-9))
	fmt.Printf("interpreted persons in %f second.\n", float64(personTime)*math.Pow10(-9))
	fmt.Printf("interpreted familys in %f second.\n", float64(familyTime)*math.Pow10(-9))
	writeJsonTime := time.Now()

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

	fmt.Printf("wrote json in %f second.\n", float64(time.Since(writeJsonTime))*math.Pow10(-9))
}

func interpretRecord(gedcom *model.Gedcom, currentRecordDeepLines []*gedcomSpec.Line, currentRecordLine *gedcomSpec.Line, waitGroup *sync.WaitGroup) {
	switch *currentRecordLine.Tag() {
	case "HEAD":
		interpretHeadRecord(gedcom, currentRecordDeepLines, currentRecordLine)
	case "INDI":
		interpretPersonRecord(gedcom, currentRecordDeepLines, currentRecordLine)
	case "FAM":
		interpretFamilyRecord(gedcom, currentRecordDeepLines, currentRecordLine)
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
	startTime := time.Now()
	for i, line := range currentRecordDeepLines {
		if i != 0 && *line.Level() == 0 {
			break
		}
		if *line.Level() == 1 {
			switch *line.Tag() {
			case "SOUR":
				source := model.Source{
					ApprovedSystemId: *line.Value(),
				}
				for j, sourceLine := range currentRecordDeepLines[i+1:] {
					if *sourceLine.Level() < 2 {
						break
					}
					switch *sourceLine.Tag() {
					case "VERS":
						if sourceLine.Value() != nil {
							source.Version = *sourceLine.Value()
						}
					case "NAME":
						if sourceLine.Value() != nil {
							source.ProductName = *sourceLine.Value()
						}
					case "CORP":
						if sourceLine.Value() != nil {
							corporation := model.SourceCorporation{
								Name: *sourceLine.Value(),
							}
							for k, corpLine := range currentRecordDeepLines[i+1+j+1:] {
								if *corpLine.Level() < 3 {
									break
								}
								switch *corpLine.Tag() {
								case "ADDR":
									if corpLine.Value() != nil {
										address := model.Address{
											MainLine: *corpLine.Value(),
										}
										for _, addrLine := range currentRecordDeepLines[i+1+j+1+k+1:] {
											if *addrLine.Level() < 4 {
												break
											}
											switch *addrLine.Tag() {
											case "CITY":
												if addrLine.Value() != nil {
													address.City = *addrLine.Value()
												}
											case "POST":
												if addrLine.Value() != nil {
													address.PostCode = *addrLine.Value()
												}
											case "CTRY":
												if addrLine.Value() != nil {
													address.Country = *addrLine.Value()
												}
											}
										}
										corporation.Address = &address
									}
								case "WWW":
									if corpLine.Value() != nil {
										corporation.WebsiteURL = *corpLine.Value()
									}
								}
							}
							source.Corporation = corporation
						}
					}
				}
				gedcom.Lock.Lock()
				gedcom.MasterSources = append(gedcom.MasterSources, &source)
				gedcom.Lock.Unlock()
			case "DATE":
				if line.Value() != nil {
					date := ""
					dateParts := strings.SplitN(*line.Value(), " ", 3)
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
					if timeLine.Value() != nil {
						date += "T" + *timeLine.Value()
					}

					gedcom.Lock.Lock()
					gedcom.TransmissionDate = date
					gedcom.Lock.Unlock()
				}
			case "DEST":
				if line.Value() != nil {
					gedcom.Lock.Lock()
					gedcom.ReceivingSystemName = *line.Value()
					gedcom.Lock.Unlock()
				}
			case "SUBM":
				if line.Value() != nil {
					gedcom.Lock.Lock()
					// TODO: ID-ify xrefid (hash or whatever)
					gedcom.SubmitterRecordId = *line.Value()
					gedcom.Lock.Unlock()
				}
			case "SUBN":
				if line.Value() != nil {
					gedcom.Lock.Lock()
					// TODO: ID-ify xrefid (hash or whatever)
					gedcom.SubmissionRecordId = *line.Value()
					gedcom.Lock.Unlock()
				}
			case "FILE":
				if line.Value() != nil {
					gedcom.Lock.Lock()
					gedcom.FileName = *line.Value()
					gedcom.Lock.Unlock()
				}
			case "COPR":
				if line.Value() != nil {
					gedcom.Lock.Lock()
					gedcom.Copyright = *line.Value()
					gedcom.Lock.Unlock()
				}
			case "GEDC":
				metadata := model.GedcomMetadata{}
				for _, gedcLine := range currentRecordDeepLines[i+1:] {
					if *gedcLine.Level() < 2 {
						break
					}
					switch *gedcLine.Tag() {
					case "VERS":
						metadata.Version = *gedcLine.Value()
					case "FORM":
						metadata.Form = *gedcLine.Value()
					}
				}
				gedcom.Lock.Lock()
				gedcom.Metadata = metadata
				gedcom.Lock.Unlock()
			case "CHAR":
				if line.Value() != nil {
					characterSet := model.CharacterSet{
						Value: *line.Value(),
					}
					if len(currentRecordDeepLines) > i+1 {
						characterSet.Version = *currentRecordDeepLines[i+1].Value()
					}
					gedcom.Lock.Lock()
					gedcom.CharacterSet = characterSet
					gedcom.Lock.Unlock()
				}
			case "LANG":
				if line.Value() != nil {
					gedcom.Lock.Lock()
					gedcom.Language = *line.Value()
					gedcom.Lock.Unlock()
				}
			case "PLAC":
				if line.Value() != nil {
					gedcom.Lock.Lock()
					gedcom.PlaceHierarchy = *line.Value()
					gedcom.Lock.Unlock()
				}
			case "NOTE":
				if line.Value() != nil {
					note := *line.Value()
					for _, noteLine := range currentRecordDeepLines[i+1:] {
						switch *noteLine.Tag() {
						case "CONT":
						case "CONC":
							if noteLine.Value() != nil {
								note += " " + *noteLine.Value()
							}
						}
					}
					gedcom.Lock.Lock()
					gedcom.ContentDescription = &note
					gedcom.Lock.Unlock()
				}
			}
		}
	}
	headTime += time.Since(startTime)
}

func interpretPersonRecord(gedcom *model.Gedcom, currentRecordDeepLines []*gedcomSpec.Line, currentRecordLine *gedcomSpec.Line) {
	startTime := time.Now()
	personPtr := model.NewPerson(currentRecordLine.XRefID())
	for i, line := range currentRecordDeepLines {
		if i != 0 && *line.Level() == 0 {
			break
		}
		if *line.Level() == 1 {
			switch *line.Tag() {
			case "NAME":
				name := model.PersonName{
					FactTypeId: 100,
				}
				for _, nameLine := range currentRecordDeepLines[i+1:] {
					if *nameLine.Level() < 2 {
						break
					}
					switch *nameLine.Tag() {
					case "GIVN":
						name.GivenNames = *nameLine.Value()
					case "SURN":
						name.Surnames = *nameLine.Value()
					}
				}
				if name.GivenNames != "" || name.Surnames != "" {
					personPtr.Names = append(personPtr.Names, &name)
				}
			case "BIRT":
				birthFact := model.PersonFact{
					FactTypeId: 405,
				}
				for _, birthFactLine := range currentRecordDeepLines[i+1:] {
					if *birthFactLine.Level() < 2 {
						break
					}
					switch *birthFactLine.Tag() {
					case "_PRIM":
						if birthFactLine.Value() != nil {
							switch *birthFactLine.Value() {
							case "Y":
								birthFact.Preferred = true
							case "N":
								birthFact.Preferred = false
							}
						}
					case "DATE":
						if birthFactLine.Value() != nil {
							birthFact.DateDetail = *birthFactLine.Value()
						}
					case "PLAC":
						place := model.PersonPlace{}
						if birthFactLine.Value() != nil {
							place.PlaceName = *birthFactLine.Value()
						}
						birthFact.Place = place
					}
				}
				personPtr.Facts = append(personPtr.Facts, &birthFact)
			//TODO: case "DEAT":
			case "SEX":
				if line.Value() != nil {
					switch *line.Value() {
					case "M":
						personPtr.Gender = 1
					case "F":
						personPtr.Gender = 2
					}
				}
			case "CHAN":
				date := ""
				for _, chanLine := range currentRecordDeepLines[i+1:] {
					if *chanLine.Level() < 2 {
						break
					}
					switch *chanLine.Tag() {
					case "DATE":
						if chanLine.Value() != nil {
							dateParts := strings.SplitN(*chanLine.Value(), " ", 3)
							if len(dateParts) >= 1 {
								date = dateParts[0]
							}
							if len(dateParts) >= 2 {
								date = monthNumberByAbbreviation[dateParts[1]] + "-" + date
							}
							if len(dateParts) >= 3 {
								date = dateParts[2] + "-" + date
							}
						}
					case "TIME":
						date += "T" + *chanLine.Value()
					}
				}
				personPtr.DateCreated = date
			case "_UID":
				if line.Value() != nil {
					personPtr.PersonRef = *line.Value()
				}
			}
		}
	}
	gedcom.Lock.Lock()
	gedcom.Persons = append(gedcom.Persons, personPtr)
	gedcom.Lock.Unlock()
	personTime += time.Since(startTime)
}

func interpretFamilyRecord(gedcom *model.Gedcom, currentRecordDeepLines []*gedcomSpec.Line, currentRecordLine *gedcomSpec.Line) {
	startTime := time.Now()
	family := model.NewFamily(currentRecordLine.XRefID())
	for i, line := range currentRecordDeepLines {
		if i != 0 && *line.Level() == 0 {
			break
		}
		switch *line.Tag() {
		case "HUSB":
			if line.Value() != nil {
				fatherId, err := util.Hash(*line.Value())
				util.MaybePanic(err)
				family.FatherId = fatherId
			}
		case "WIFE":
			if line.Value() != nil {
				motherId, err := util.Hash(*line.Value())
				util.MaybePanic(err)
				family.MotherId = motherId
			}
		case "CHIL":
			if line.Value() != nil {
				childId, err := util.Hash(*line.Value())
				util.MaybePanic(err)
				family.ChildIds = append(family.ChildIds, childId)
			}
		case "CHAN":
			date := ""
			for _, chanLine := range currentRecordDeepLines[i+1:] {
				if *chanLine.Level() < 2 {
					break
				}
				switch *chanLine.Tag() {
				case "DATE":
					if chanLine.Value() != nil {
						dateParts := strings.SplitN(*chanLine.Value(), " ", 3)
						if len(dateParts) >= 1 {
							date = dateParts[0]
						}
						if len(dateParts) >= 2 {
							date = monthNumberByAbbreviation[dateParts[1]] + "-" + date
						}
						if len(dateParts) >= 3 {
							date = dateParts[2] + "-" + date
						}
					}
				case "TIME":
					if chanLine.Value() != nil {
						date += "T" + *chanLine.Value()
					}
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
		gedcom.Childs = append(gedcom.Childs, &child)
		gedcom.Lock.Unlock()

	}

	gedcom.Lock.Lock()
	gedcom.Familys = append(gedcom.Familys, &family)
	gedcom.Lock.Unlock()
	familyTime += time.Since(startTime)
}
