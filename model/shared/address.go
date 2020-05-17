package shared

type Address struct {
	PhysicalAddress *PhysicalAddress
	PhoneNumber     []*string
	Email           []*string
	Fax             []*string
	WebPage         []*string
}

type PhysicalAddress struct {
	MainLine string
	Line1    string
	Line2    string
	Line3    string
	City     string
	State    string
	PostCode string
	Country  string
}
