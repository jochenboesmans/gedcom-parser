package gedcom

import (
	"testing"
)

var lines = [][]string{
	{
		"1 NAME William_Lee\n",
	},
	{
		"1 NAME /Parry/\n",
	},
	{
		"1 NAME William Lee /Parry/\n",
	},
	{
		"10 NAME William_Lee\n",
	},
	{
		"1 NAME\n",
		"2 GIVN William Lee\n",
		"2 SURN Parry\n",
		"2 _PRIM Y\n",
	},
	{
		"1 NAME\n",
		"2 GIVN William Lee\n",
		"2 SURN Parry\n",
		"2 _PRIM N\n",
	},
	{
		"1 NAME\n",
		"2 GIVN William Lee\n",
		"2 SURN Parry\n",
	},
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
	{
		GivenName: "William_Lee",
		Surname:   "",
		Primary:   false,
	},
	{
		GivenName: "William Lee",
		Surname:   "Parry",
		Primary:   true,
	},
	{
		GivenName: "William Lee",
		Surname:   "Parry",
		Primary:   false,
	},
	{
		GivenName: "William Lee",
		Surname:   "Parry",
		Primary:   false,
	},
}

func recordLines(multiLines []string) []*Line {
	r := []*Line{}
	for i := range multiLines {
		r = append(r, NewLine(&multiLines[i]))
	}
	return r
}

var testNames = [][]*Line{
	recordLines(lines[0]),
	recordLines(lines[1]),
	recordLines(lines[2]),
	recordLines(lines[3]),
	recordLines(lines[4]),
	recordLines(lines[5]),
	recordLines(lines[6]),
}

func TestNameStructure(t *testing.T) {
	for i := range testNames {
		result, err := interpretNameStructure(testNames[i])
		if err != nil {
			t.Errorf("failed to interpret %s as name structure with error: %s", lines[i], err)
		}

		if *result != expectedResults[i] {
			t.Errorf("result name does not equal expected; result: %+v, expected %+v", *result, expectedResults[i])
		}
	}
}
