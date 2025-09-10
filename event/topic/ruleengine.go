package topic

import "github.com/nguyencobap/onvif/xsd"

type RuleEngine struct {
	Topic                *xsd.Boolean          `xml:"topic,attr"`
	MotionRegionDetector *MotionRegionDetector `json:"," xml:","`
	CellMotionDetector   *CellMotionDetector   `json:"," xml:","`
	TamperDetector       *TamperDetector       `json:"," xml:","`
	Recognition          *Recognition          `json:"," xml:","`
	CountAggregation     *CountAggregation     `json:"," xml:","`
}

type MotionRegionDetector struct {
	Topic  *xsd.Boolean `xml:"topic,attr"`
	Motion *Motion      `json:"Motion" xml:"Motion"`
}

type CellMotionDetector struct {
	Topic  *xsd.Boolean `xml:"topic,attr"`
	Motion *Motion
}

type Motion struct {
	Topic              *xsd.Boolean        `xml:"topic,attr"`
	MessageDescription *MessageDescription `json:"," xml:","`
}

type TamperDetector struct {
	Topic  *xsd.Boolean `xml:"topic,attr"`
	Tamper *Tamper
}

type Tamper struct {
	Topic              *xsd.Boolean        `xml:"topic,attr"`
	MessageDescription *MessageDescription `json:"," xml:","`
}

type Recognition struct {
	Topic        *xsd.Boolean `xml:"topic,attr"`
	Face         *Face        `json:"," xml:","`
	LicensePlate *Face        `json:"," xml:","`
}

type Face struct {
	Topic              *xsd.Boolean        `xml:"topic,attr"`
	MessageDescription *MessageDescription `json:"," xml:","`
}

type LicensePlate struct {
	Topic              *xsd.Boolean        `xml:"topic,attr"`
	MessageDescription *MessageDescription `json:"," xml:","`
}

type CountAggregation struct {
	Topic   *xsd.Boolean `xml:"topic,attr"`
	Counter *Counter     `json:"," xml:","`
}

type Counter struct {
	Topic              *xsd.Boolean        `xml:"topic,attr"`
	MessageDescription *MessageDescription `json:"," xml:","`
}
