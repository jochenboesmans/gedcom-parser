package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jochenboesmans/gedcom-parser/model/child"
	"github.com/jochenboesmans/gedcom-parser/model/family"
	"github.com/jochenboesmans/gedcom-parser/model/header"
	"github.com/jochenboesmans/gedcom-parser/model/note"
	"github.com/jochenboesmans/gedcom-parser/model/person"
	"github.com/jochenboesmans/gedcom-parser/model/repository"
	"github.com/jochenboesmans/gedcom-parser/model/shared"
	"github.com/jochenboesmans/gedcom-parser/model/source"
	"github.com/jochenboesmans/gedcom-parser/model/submission"
	"github.com/jochenboesmans/gedcom-parser/model/submitter"
	"strconv"

	//"github.com/jochenboesmans/gedcom-parser/model/repository"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	gedcomSpec "github.com/jochenboesmans/gedcom-parser/gedcom"
	"github.com/jochenboesmans/gedcom-parser/model"
	pb "github.com/jochenboesmans/gedcom-parser/proto"
	"github.com/jochenboesmans/gedcom-parser/util"
	"github.com/pquerna/ffjson/ffjson"
)

type OutputGedcom struct {
	Header      *header.Header
	Submission  *submission.Submission
	Persons     []*person.Person
	Familys     []*family.Family
	Childs      []*child.Child
	Notes       []*note.Note
	Repositorys []*repository.Repository
	Sources     []*source.Source
	Submitters  []*submitter.Submitter
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
	beginTime := time.Now()
	pathToGedcomFile := flag.String("pathToGedcomFile", "./test-input/ITIS.ged", "relative path to input gedcom file (with .ged extension if present)")
	pathToJsonFile := flag.String("pathToJsonFile", "./artifacts/ITIS.json", "relative path to output json file (with .json extension if wanted)")
	useProtobuf := flag.Bool("useProtobuf", false, "whether to use protobuf instead of json as serialization format")
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
	writeTime := time.Now()

	gedcomWithoutLock := OutputGedcom{
		Header:      gedcom.Header,
		Submission:  gedcom.Submission,
		Persons:     gedcom.Persons,
		Familys:     gedcom.Familys,
		Childs:      gedcom.Childs,
		Notes:       gedcom.Notes,
		Repositorys: gedcom.Repositorys,
		Sources:     gedcom.Sources,
		Submitters:  gedcom.Submitters,
	}

	if !*useProtobuf {
		gedcomJson, err := ffjson.Marshal(gedcomWithoutLock)
		writeFile, err := os.Create(*pathToJsonFile)
		writer := bufio.NewWriter(writeFile)
		_, err = writer.Write(gedcomJson)
		util.MaybePanic(err)
		err = writer.Flush()
		util.MaybePanic(err)
	} else {
		// WIP: needs full gedcom protobuf structure to be built
		pbPerson := &pb.Person{
			Id:        gedcom.Persons[0].Id,
			PersonRef: gedcom.Persons[0].PersonRef,
			IsLiving:  gedcom.Persons[0].IsLiving,
		}

		personProto, err := proto.Marshal(pbPerson)
		personWriteFile, err := os.Create("./artifacts/personproto")

		personWriter := bufio.NewWriter(personWriteFile)
		_, err = personWriter.Write(personProto)
		util.MaybePanic(err)
		err = personWriter.Flush()
		util.MaybePanic(err)
	}

	fmt.Printf("wrote to file in %f second.\n", float64(time.Since(writeTime))*math.Pow10(-9))
	fmt.Printf("total time taken: %f second.\n", float64(time.Since(beginTime))*math.Pow10(-9))
}

func interpretRecord(gedcom *model.Gedcom, currentRecordDeepLines []*gedcomSpec.Line, currentRecordLine *gedcomSpec.Line, waitGroup *sync.WaitGroup) {
	switch *currentRecordLine.Tag() {
	case "HEAD":
		interpretHeadRecord(gedcom, currentRecordDeepLines, currentRecordLine)
	case "INDI":
		interpretPersonRecord(gedcom, currentRecordDeepLines, currentRecordLine)
	case "FAM":
		interpretFamilyRecord(gedcom, currentRecordDeepLines, currentRecordLine)
	case "NOTE":
		interpretNoteRecord(gedcom, currentRecordDeepLines)
	case "REPO":
		interpretRepoRecord(gedcom, currentRecordDeepLines)
	case "SOUR":
		//interpretSourceRecord(gedcom, currentRecordDeepLines, currentRecordLine)
	case "SUBN":
		interpretSubmitterRecord(gedcom, currentRecordDeepLines)
	case "SUBM":
		interpretSubmissionRecord(gedcom, currentRecordDeepLines)
		//case "TRLR": nothing really to do here except maybe validate?
	}
	waitGroup.Done()
}

func interpretSubmitterRecord(gedcom *model.Gedcom, currentRecordDeepLines []*gedcomSpec.Line) {
	baseLevel := *currentRecordDeepLines[0].Level()
	idString := *currentRecordDeepLines[0].XRefID()
	id, err := util.Hash(idString)
	util.MaybePanic(err)
	s := submitter.Submitter{
		Id: id,
	}
	for _, line := range currentRecordDeepLines[1:] {
		if *line.Level() <= baseLevel {
			break
		}
		switch *line.Tag() {
		case "NAME":
			s.Name = *line.Value()
		//case "ADDR":
		// interpretAddressStructure
		//case "OBJE":
		// interpretMultimediaLink
		case "LANG":
			s.LanguagePreference = append(s.LanguagePreference, *line.Value())
		case "RFN":
			s.SubmitterRegisteredRFN = *line.Value()
		case "RIN":
			s.AutomatedRecordId = *line.Value()
			//case "NOTE":
			//interpretNoteStructure
			//case "CHAN":
			// interpretChangeDate
		}
	}
	gedcom.Lock.Lock()
	gedcom.Submitters = append(gedcom.Submitters, &s)
	gedcom.Lock.Unlock()
}

func interpretSubmissionRecord(gedcom *model.Gedcom, currentRecordDeepLines []*gedcomSpec.Line) {
	baseLevel := *currentRecordDeepLines[0].Level()
	idString := *currentRecordDeepLines[0].XRefID()
	id, err := util.Hash(idString)
	util.MaybePanic(err)
	s := submission.Submission{
		Id: id,
	}
	for _, line := range currentRecordDeepLines[1:] {
		if *line.Level() <= baseLevel {
			break
		}
		switch *line.Tag() {
		case "SUBM":
			submitterIdString := *line.Value()
			submitterId, err := util.Hash(submitterIdString)
			util.MaybePanic(err)
			s.SubmitterId = submitterId
		case "FAMF":
			s.NameOfFamilyFile = *line.Value()
		case "TEMP":
			s.TempleCode = *line.Value()
		case "ANCE":
			gensString := *line.Value()
			gensInt, err := strconv.ParseUint(gensString, 10, 32)
			util.MaybePanic(err)
			s.GenerationsOfAncestors = uint32(gensInt)
		case "DESC":
			gensString := *line.Value()
			gensInt, err := strconv.ParseUint(gensString, 10, 32)
			util.MaybePanic(err)
			s.GenerationsOfDescendants = uint32(gensInt)
		case "ORDI":
			switch strings.ToUpper(*line.Value()) {
			case "YES":
				s.OrdinanceProcessFlag = true
			case "NO":
				s.OrdinanceProcessFlag = false
			}
		case "RIN":
			s.AutomatedRecordId = *line.Value()
			//case "NOTE":
			//interpretNoteStructure
			//case "CHAN":
			// interpretChangeDate
		}
	}
	gedcom.Lock.Lock()
	gedcom.Submission = &s
	gedcom.Lock.Unlock()
}

func interpretRepoRecord(gedcom *model.Gedcom, currentRecordDeepLines []*gedcomSpec.Line) {
	idString := *currentRecordDeepLines[0].XRefID()
	id, err := util.Hash(idString)
	util.MaybePanic(err)
	r := repository.Repository{
		Id:      id,
		Address: &shared.Address{},
	}
	for i, line := range currentRecordDeepLines[1:] {
		if *line.Level() == 0 {
			break
		}
		if *line.Level() == 1 {
			switch *line.Tag() {
			case "NAME":
				r.Name = *line.Value()
			case "ADDR":
				if line.Value() != nil {
					physicalAddress := shared.PhysicalAddress{
						MainLine: *line.Value(),
					}
					for _, addrLine := range currentRecordDeepLines[i+1:] {
						if *addrLine.Level() < 2 {
							break
						}
						switch *addrLine.Tag() {
						case "CONT":
							physicalAddress.MainLine = physicalAddress.MainLine + " " + *addrLine.Value()
						case "ADR1":
							physicalAddress.Line1 = *addrLine.Value()
						case "ADR2":
							physicalAddress.Line2 = *addrLine.Value()
						case "ADR3":
							physicalAddress.Line3 = *addrLine.Value()
						case "CITY":
							physicalAddress.City = *addrLine.Value()
						case "POST":
							physicalAddress.PostCode = *addrLine.Value()
						case "CTRY":
							physicalAddress.Country = *addrLine.Value()
						}
					}
					r.Address.PhysicalAddress = &physicalAddress
				}
			case "PHON":
				r.Address.PhoneNumber = append(r.Address.PhoneNumber, line.Value())
			case "EMAIL":
				r.Address.Email = append(r.Address.Email, line.Value())
			case "FAX":
				r.Address.Fax = append(r.Address.Fax, line.Value())
			case "WWW":
				r.Address.WebPage = append(r.Address.WebPage, line.Value())
			}
		}
	}
	gedcom.Repositorys = append(gedcom.Repositorys, &r)
}

func interpretNoteRecord(gedcom *model.Gedcom, currentRecordDeepLines []*gedcomSpec.Line) {
	idString := *currentRecordDeepLines[0].XRefID()
	id, err := util.Hash(idString)
	util.MaybePanic(err)
	n := note.Note{
		Id:             id,
		SubmitterText:  *currentRecordDeepLines[0].Value(),
		UserReferences: []*shared.UserReference{},
	}
	for i, line := range currentRecordDeepLines[1:] {
		if *line.Level() == 0 {
			break
		}
		if *line.Level() == 1 {
			switch *line.Tag() {
			case "CONC":
			case "CONT":
				n.SubmitterText += *line.Value()
			case "REFN":
				reference := shared.UserReference{
					Number: *line.Value(),
				}
				for _, noteLine := range currentRecordDeepLines[i+1:] {
					if *noteLine.Level() < 2 {
						break
					}
					switch *noteLine.Tag() {
					case "TYPE":
						reference.Type = *line.Value()
					}
				}
				n.UserReferences = append(n.UserReferences, &reference)
			case "RIN":
				n.AutomatedRecordId = *line.Value()
				// TODO: sourcecitation and changedate
			}
		}
	}
	gedcom.Lock.Lock()
	gedcom.Notes = append(gedcom.Notes, &n)
	gedcom.Lock.Unlock()
}

func interpretHeadRecord(gedcom *model.Gedcom, currentRecordDeepLines []*gedcomSpec.Line, currentRecordLine *gedcomSpec.Line) {
	startTime := time.Now()
	h := header.Header{}
	for i, line := range currentRecordDeepLines {
		if i != 0 && *line.Level() == 0 {
			break
		}
		if *line.Level() == 1 {
			switch *line.Tag() {
			case "SOUR":
				source := header.Source{
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
							corporation := header.Corporation{
								Name: *sourceLine.Value(),
							}
							for k, corpLine := range currentRecordDeepLines[i+1+j+1:] {
								if *corpLine.Level() < 3 {
									break
								}
								switch *corpLine.Tag() {
								case "ADDR":
									if corpLine.Value() != nil {
										address := shared.PhysicalAddress{
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
							source.Corporation = &corporation
						}
					}
				}
				h.Source = &source
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

					h.TransmissionDate = date
				}
			case "DEST":
				if line.Value() != nil {
					h.ReceivingSystemName = *line.Value()
				}
			case "SUBM":
				if line.Value() != nil {
					h.SubmitterRecordId = *line.Value()
					// TODO: ID-ify xrefid (hash or whatever)
				}
			case "SUBN":
				if line.Value() != nil {
					h.SubmissionRecordId = *line.Value()
				}
			case "FILE":
				if line.Value() != nil {
					h.FileName = *line.Value()
				}
			case "COPR":
				if line.Value() != nil {
					h.Copyright = *line.Value()
				}
			case "GEDC":
				metadata := header.GedcomMetadata{}
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
				h.Metadata = metadata
			case "CHAR":
				if line.Value() != nil {
					characterSet := header.CharacterSet{
						Value: *line.Value(),
					}
					if len(currentRecordDeepLines) > i+1 {
						characterSet.Version = *currentRecordDeepLines[i+1].Value()
					}
					h.CharacterSet = characterSet
				}
			case "LANG":
				if line.Value() != nil {
					h.Language = *line.Value()
				}
			case "PLAC":
				if line.Value() != nil {
					h.PlaceHierarchy = *line.Value()
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
					h.ContentDescription = note
				}
			}
		}
	}
	headTime += time.Since(startTime)

	gedcom.Lock.Lock()
	gedcom.Header = &h
	gedcom.Lock.Unlock()
}

func interpretPersonRecord(gedcom *model.Gedcom, currentRecordDeepLines []*gedcomSpec.Line, currentRecordLine *gedcomSpec.Line) {
	startTime := time.Now()
	personPtr := person.NewPerson(currentRecordLine.XRefID())
	for i, line := range currentRecordDeepLines {
		if i != 0 && *line.Level() == 0 {
			break
		}
		if *line.Level() == 1 {
			switch *line.Tag() {
			case "NAME":
				name := person.PersonName{
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
				birthFact := person.PersonFact{
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
						place := person.PersonPlace{}
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
	family := family.NewFamily(currentRecordLine.XRefID())
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
		child := child.NewChild(currentRecordLine.XRefID(), i, childId)
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
