package header

type Source struct {
	ApprovedSystemId string
	Version          string
	ProductName      string
	Corporation      *Corporation
	Data             *Data
}

type Corporation struct {
	Name       string
	Address    *Address
	WebsiteURL string
}

type Data struct {
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
