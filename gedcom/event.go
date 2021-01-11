package gedcom

import (
	"fmt"
	"github.com/jochenboesmans/gedcom-parser/util"
)

type Event struct {
	Date
	Place
	Primary bool
}

func interpretEventStructure(eventLines []*Line) (*Event, error) {
	rootLevel, err := eventLines[0].Level()
	if err != nil {
		return nil, fmt.Errorf("failed to parse root level of event structure: %s", err)
	}

	event := Event{}
	for _, eventLine := range eventLines[1:] {
		level, err := eventLine.Level()
		if err != nil {
			continue
		}
		if level <= rootLevel {
			break // end of event structure
		}

		tag, err := eventLine.Tag()
		if err != nil {
			continue
		}
		switch tag {
		case "DATE":
			event.Date = interpretDateStructure(eventLine)
		case "PLAC":
			event.Place = Place(eventLine.Value())
		case "_PRIM":
			if primaryBool, ok := util.PrimaryBoolByValue[eventLine.Value()]; ok {
				event.Primary = primaryBool
			}
		}
	}
	return &event, nil
}

func (event *Event) toGedcomIndividualEvent() Gedcom_Individual_Event {
	gedcomIndividualDate := event.Date.toGedcomIndividualDate()
	placeString := event.Place.toString()
	return Gedcom_Individual_Event{
		Date:    &gedcomIndividualDate,
		Place:   placeString,
		Primary: event.Primary,
	}
}
