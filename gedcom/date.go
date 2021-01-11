package gedcom

import (
	"github.com/jochenboesmans/gedcom-parser/util"
	"strconv"
	"strings"
)

type Date struct {
	Year  uint32
	Month uint32
	Day   uint32
}

func interpretDateStructure(line *Line) Date {
	dateParts := strings.SplitN(line.Value(), " ", 3)
	date := Date{}
	if len(dateParts) > 0 {
		if year, err := strconv.Atoi(dateParts[0]); err == nil {
			date.Year = uint32(year)
		}
	}
	if len(dateParts) > 1 {
		if month, ok := util.MonthIntByAbbr[strings.ToUpper(dateParts[1])]; ok {
			date.Month = uint32(month)
		}
	}
	if len(dateParts) > 2 {
		if day, err := strconv.Atoi(dateParts[2]); err == nil {
			date.Day = uint32(day)
		}
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
