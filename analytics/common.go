package analytics

import "github.com/nguyencobap/onvif/xsd"

type Parameters struct {
	SimpleItemDescription  []SimpleItemDescription  `json:","`
	ElementItemDescription []ElementItemDescription `json:","`
	Extension              *xsd.String              `json:","`
}

type SimpleItemDescription struct {
	Name  string `json:"," xml:",attr"`
	Type  string `json:"," xml:",attr"`
	Value string `json:"," xml:",attr"`
}

type ElementItemDescription struct {
	Name  string `json:"," xml:",attr"`
	Value string `json:"," xml:",attr"`
}

type Messages struct {
	IsProperty  *xsd.Boolean `json:"," xml:",attr"`
	Source      *Source      `json:","`
	Key         *Key         `json:","`
	Data        *Data        `json:","`
	Extension   *xsd.String  `json:","`
	ParentTopic *xsd.String  `json:","`
}

type Source struct {
	SimpleItemDescription  []SimpleItemDescription  `json:","`
	ElementItemDescription []ElementItemDescription `json:","`
	Extension              *xsd.String              `json:","`
}

type Key struct {
	SimpleItemDescription  []SimpleItemDescription  `json:","`
	ElementItemDescription []ElementItemDescription `json:","`
	Extension              *xsd.String              `json:","`
}

type Data struct {
	SimpleItemDescription  []SimpleItemDescription  `json:","`
	ElementItemDescription []ElementItemDescription `json:","`
	Extension              *xsd.String              `json:","`
}
