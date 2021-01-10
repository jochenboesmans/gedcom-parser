package gedcom

import (
	"testing"
)

var lines = []string{
	"1 NAME William_Lee",
	"1 NAME /Parry/",
	"1 NAME William Lee /Parry/",
}
var testNames = [][]*Line{
	{NewLine(&lines[0])},
	{NewLine(&lines[1])},
	{NewLine(&lines[2])},
}
var expectedResults = []Name{
	{
		GivenName: "William_Lee",
		Surname:   "",
		Primary:   false,
	},
	{
		GivenName: "",
		Surname:   "Parry",
		Primary:   false,
	},
	{
		GivenName: "William Lee",
		Surname:   "Parry",
		Primary:   false,
	},
}

func TestNameStructure(t *testing.T) {
	for i := range testNames {
		result, err := NameStructure(testNames[i])
		if err != nil {
			t.Errorf("failed to interpret %s as name structure with error: %s", lines[i], err)
		}

		if *result != expectedResults[i] {
			t.Errorf("result name does not equal expected; result: %v, expected %v", *result, expectedResults[i])
		}
	}
}
