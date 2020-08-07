package util

var MonthAbbrByInt = map[int]string{
	1:  "JAN",
	2:  "FEB",
	3:  "MAR",
	4:  "APR",
	5:  "MAY",
	6:  "JUN",
	7:  "JUL",
	8:  "AUG",
	9:  "SEP",
	10: "OCT",
	11: "NOV",
	12: "DEC",
}

var MonthIntByAbbr = invertIntStringMap(MonthAbbrByInt)

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

func invertIntStringMap(m map[int]string) map[string]int {
	r := map[string]int{}
	for k, v := range m {
		r[v] = k
	}
	return r
}
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
