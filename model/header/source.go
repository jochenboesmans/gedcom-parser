package header

import "github.com/jochenboesmans/gedcom-parser/model/shared"

type Source struct {
	ApprovedSystemId string
	Version          string
	ProductName      string
	Corporation      *Corporation
	Data             *Data
}

type Corporation struct {
	Name       string
	Address    *shared.PhysicalAddress
	WebsiteURL string
}

type Data struct {
	Name            string
	PublicationDate string
	Copyright       string
}
