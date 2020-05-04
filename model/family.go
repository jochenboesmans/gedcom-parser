package model

import "github.com/jochenboesmans/gedcom-parser/util"

type Family struct {
	Id          uint32
	FatherId    uint32
	MotherId    uint32
	ChildIds    []uint32
	DateCreated string
}

func NewFamily(identificationString string) Family {
	id, err := util.Hash(identificationString)
	util.MaybePanic(err)

	return Family{
		Id:       id,
		ChildIds: []uint32{},
	}
}
