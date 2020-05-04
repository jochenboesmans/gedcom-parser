package util

func MaybePanic(err error) {
	if err != nil {
		panic(err)
	}
}
