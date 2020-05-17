package repository

import "github.com/jochenboesmans/gedcom-parser/model/shared"

type Repository struct {
	Id      uint32
	Name    string
	Address *shared.Address
	//Note *shared.Note
	//UserReference *shared.UserReference
	AutomatedRecordId string
	//ChangeDate *shared.ChangeDate
}
