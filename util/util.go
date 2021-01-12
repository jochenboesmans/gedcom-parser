package util

var MonthAbbrByInt = map[string]string{
	"01": "JAN",
	"02": "FEB",
	"03": "MAR",
	"04": "APR",
	"05": "MAY",
	"06": "JUN",
	"07": "JUL",
	"08": "AUG",
	"09": "SEP",
	"10": "OCT",
	"11": "NOV",
	"12": "DEC",
}

var MonthIntByAbbr = invertStringStringMap(MonthAbbrByInt)

var PrimaryValueByBool = map[bool]string{
	true:  "Y",
	false: "N",
}

var PrimaryBoolByValue = invertBoolStringMap(PrimaryValueByBool)

var GenderLetterByFull = map[string]string{
	"MALE":   "M",
	"FEMALE": "F",
}

var GenderFullByLetter = invertStringStringMap(GenderLetterByFull)

func invertBoolStringMap(m map[bool]string) map[string]bool {
	r := map[string]bool{}
	for k, v := range m {
		r[v] = k
	}
	return r
}
func invertStringStringMap(m map[string]string) map[string]string {
	r := map[string]string{}
	for k, v := range m {
		r[v] = k
	}
	return r
}

// TODO: Start using below generic function once there's better support for generics in Go.
//func invert(type Key, Value) (m map[Key]Value) map[Value]Key {
//	r := map[Value]Key{}
//	for k, v := range m {
//		r[v] = k
//	}
//	return r
//}
