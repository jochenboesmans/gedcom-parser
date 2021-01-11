package gedcom

type Place string

func (place *Place) toString() string {
	return string(*place)
}
