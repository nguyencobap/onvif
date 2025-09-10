package topic

import "github.com/nguyencobap/onvif/xsd"

type MessageDescription struct {
	IsProperty xsd.Boolean `xml:"IsProperty,attr"`
	Source     Source      `json:"," xml:","`
	Data       Data        `json:"," xml:","`
}

type Source struct {
	SimpleItemDescription []SimpleItemDescription `json:"," xml:","`
}

type Data struct {
	SimpleItemDescription []SimpleItemDescription `json:"," xml:","`
}

type SimpleItemDescription struct {
	Name xsd.AnyType `xml:"Name,attr"`
	Type xsd.AnyType `xml:"Type,attr"`
}
