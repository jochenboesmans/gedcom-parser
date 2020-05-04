package model

import (
	"github.com/jochenboesmans/gedcom-parser/util"
	"strconv"
)

type Child struct {
	Id                   uint32
	FamilyId             uint32
	ChildId              uint32
	RelationshipToFather uint8
	RelationshipToMother uint8
}

func NewChild(familyRecordId string, ithChildInFamily int, childId uint32) Child {
	childPersonId, err := util.Hash("CHILD-" + strconv.Itoa(ithChildInFamily) + "-" + familyRecordId)
	util.MaybePanic(err)
	familyId, err := util.Hash(familyRecordId)
	util.MaybePanic(err)
	util.MaybePanic(err)
	return Child{
		Id:       childPersonId,
		FamilyId: familyId,
		ChildId:  childId,
	}
}
