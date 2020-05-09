package model

type Source struct {
	ApprovedSystemId string
	Version          string
	ProductName      string
	Corporation      SourceCorporation
	Data             SourceData
}

type SourceCorporation struct {
	Name       string
	Address    *Address
	WebsiteURL string
}

type SourceData struct {
	Name            string
	PublicationDate string
	Copyright       string
}

type Address struct {
	MainLine string
	City     string
	PostCode string
	Country  string
}
