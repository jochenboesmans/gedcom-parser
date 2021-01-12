package gedcom

import (
	"testing"
)

func TestInterpretDateStructure(t *testing.T) {
	testDateValues := []string{
		"0 DATE 02 OCT 1822",
		"0 DATE 1822",
		"0 DATE OCT 1822",
	}
	expectedValues := []Date{
		{
			Year:  "1822",
			Month: "10",
			Day:   "02",
		},
		{
			Year:  "1822",
			Month: "",
			Day:   "",
		},
		{
			Year:  "1822",
			Month: "10",
			Day:   "",
		},
	}
	for i, testDateValue := range testDateValues {
		l := NewLine(&testDateValue)
		result := interpretDateStructure(l)
		if result != expectedValues[i] {
			t.Errorf("result date %+v does not match expected date %+v", result, expectedValues[i])
		}
	}
}
