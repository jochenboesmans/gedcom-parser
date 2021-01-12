package gedcom

import (
	"github.com/jochenboesmans/gedcom-parser/util"
	"strings"
)

type Date struct {
	Year  string
	Month string
	Day   string
}

func interpretDateStructure(line *Line) Date {
	dateParts := strings.SplitN(line.Value(), " ", 3)
	date := Date{}
	if len(dateParts) > 0 {
		date.Year = dateParts[len(dateParts)-1]
	}
	if len(dateParts) > 1 {
		if month, ok := util.MonthIntByAbbr[strings.ToUpper(dateParts[len(dateParts)-2])]; ok {
			date.Month = month
		}
	}
	if len(dateParts) > 2 {
		date.Day = dateParts[len(dateParts)-3]
	}
	return date
}

func (date *Date) toGedcomIndividualDate() Gedcom_Individual_Date {
	return Gedcom_Individual_Date{
		Year:  date.Year,
		Month: date.Month,
		Day:   date.Day,
	}
}
