package model

import (
	"github.com/jochenboesmans/gedcom-parser/util"
)

type Person struct {
	Id          uint32
	PersonRef   string
	IsLiving    bool
	Gender      uint8
	DateCreated string
	Names       []PersonName
	Facts       []PersonFact
}

type PersonName struct {
	FactTypeId uint16
	GivenNames string
	Surnames   string
}

type PersonFact struct {
	FactTypeId uint16
	DateDetail string
	Place      PersonPlace
	Preferred  bool
}

type PersonPlace struct {
	PlaceName string
}

func NewPerson(identificationString string) *Person {
	id, err := util.Hash(identificationString)
	util.MaybePanic(err)

	// person is assumed living unless proven to be dead
	return &Person{
		Id:       id,
		IsLiving: true,
		Facts:    []PersonFact{},
		Names:    []PersonName{},
	}
}
