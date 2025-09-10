package analytics

//go:generate python3 ../python/gen_commands.py

import (
	"github.com/nguyencobap/onvif/xsd"
	"github.com/nguyencobap/onvif/xsd/onvif"
)

// GetSupportedAnalyticsModules and its properties are defined in the Onvif specification:
// https://www.onvif.org/ver20/analytics/wsdl/analytics.wsdl#op.GetSupportedAnalyticsModules
type GetSupportedAnalyticsModules struct {
	XMLName            string               `xml:"tan:GetSupportedAnalyticsModules"`
	ConfigurationToken onvif.ReferenceToken `xml:"tan:ConfigurationToken"`
}

type GetSupportedAnalyticsModulesResponse struct {
	SupportedAnalyticsModules SupportedAnalyticsModules
}

type SupportedAnalyticsModules struct {
	Limit                                *xsd.Int                     `json:","`
	AnalyticsModuleContentSchemaLocation *xsd.String                  `json:","`
	AnalyticsModuleDescription           []AnalyticsModuleDescription `json:","`
}

type AnalyticsModuleDescription struct {
	Name         string      `xml:"Name,attr"`
	Fixed        bool        `xml:"fixed,attr"`
	MaxInstances int         `xml:"maxInstances,attr"`
	Parameters   *Parameters `json:","`
	Messages     *Messages   `json:","`
}

// CreateAnalyticsModules and its properties are defined in the Onvif specification:
// https://www.onvif.org/ver20/analytics/wsdl/analytics.wsdl#op.CreateAnalyticsModules
type CreateAnalyticsModules struct {
	XMLName            string                `xml:"tev:CreateAnalyticsModules"`
	ConfigurationToken onvif.ReferenceToken  `xml:"tan:ConfigurationToken"`
	AnalyticsModule    []onvif.ConfigRequest `xml:"tan:AnalyticsModule"`
}

type CreateAnalyticsModulesResponse struct{}

// DeleteAnalyticsModules and its properties are defined in the Onvif specification:
// https://www.onvif.org/ver20/analytics/wsdl/analytics.wsdl#op.DeleteAnalyticsModules
type DeleteAnalyticsModules struct {
	XMLName             string               `xml:"tan:DeleteAnalyticsModules"`
	ConfigurationToken  onvif.ReferenceToken `xml:"tan:ConfigurationToken"`
	AnalyticsModuleName []xsd.String         `xml:"tan:AnalyticsModuleName"`
}

type DeleteAnalyticsModulesResponse struct{}

// GetAnalyticsModules and its properties are defined in the Onvif specification:
// https://www.onvif.org/ver20/analytics/wsdl/analytics.wsdl#op.GetAnalyticsModules
type GetAnalyticsModules struct {
	XMLName            string               `xml:"tan:GetAnalyticsModules"`
	ConfigurationToken onvif.ReferenceToken `xml:"tan:ConfigurationToken"`
}

type GetAnalyticsModulesResponse struct {
	AnalyticsModule []onvif.Config
}

// GetAnalyticsModuleOptions and its properties are defined in the Onvif specification:
// https://www.onvif.org/ver20/analytics/wsdl/analytics.wsdl#op.GetAnalyticsModuleOptions
type GetAnalyticsModuleOptions struct {
	XMLName            string               `xml:"tan:GetAnalyticsModuleOptions"`
	Type               xsd.QName            `xml:"tan:Type,"`
	ConfigurationToken onvif.ReferenceToken `xml:"tan:ConfigurationToken"`
}

type GetAnalyticsModuleOptionsResponse struct {
	Options []AnalyticsModuleOptions
}

type AnalyticsModuleOptions struct {
	RuleType        string       `json:"," xml:",attr"`
	Name            string       `json:"," xml:",attr"`
	Type            string       `json:"," xml:",attr"`
	AnalyticsModule string       `json:"," xml:",attr"`
	IntRange        *IntRange    `json:","`
	StringItems     *StringItems `json:","`
}

type IntRange struct {
	Min int
	Max int
}

type StringItems struct {
	Item []string
}

// ModifyAnalyticsModules and its properties are defined in the Onvif specification:
// https://www.onvif.org/ver20/analytics/wsdl/analytics.wsdl#op.ModifyAnalyticsModules
type ModifyAnalyticsModules struct {
	XMLName            string                `xml:"tan:ModifyAnalyticsModules"`
	ConfigurationToken onvif.ReferenceToken  `xml:"tan:ConfigurationToken"`
	AnalyticsModule    []onvif.ConfigRequest `xml:"tan:AnalyticsModule"`
}

type ModifyAnalyticsModulesResponse struct{}

// GetSupportedRules and its properties are defined in the Onvif specification:
// https://www.onvif.org/ver20/analytics/wsdl/analytics.wsdl#op.GetSupportedRules
type GetSupportedRules struct {
	XMLName            string               `xml:"tan:GetSupportedRules"`
	ConfigurationToken onvif.ReferenceToken `xml:"tan:ConfigurationToken"`
}

type GetSupportedRulesResponse struct {
	SupportedRules SupportedRules
}

type SupportedRules struct {
	Limit                     *xsd.Int    `json:","`
	RuleContentSchemaLocation *xsd.String `json:","`
	RuleDescription           []RuleDescription
}

type RuleDescription struct {
	Name         *xsd.String  `json:"," xml:",attr"`
	Fixed        *xsd.Boolean `json:"," xml:"fixed,attr"`
	MaxInstances *xsd.Int     `json:"," xml:"maxInstances,attr"`
	Parameters   Parameters
	Messages     Messages `json:","`
}

// GetRules and its properties are defined in the Onvif specification:
// https://www.onvif.org/ver20/analytics/wsdl/analytics.wsdl#op.GetRules
type GetRules struct {
	XMLName            string               `xml:"tan:GetRules"`
	ConfigurationToken onvif.ReferenceToken `xml:"tan:ConfigurationToken"`
}

type GetRulesResponse struct {
	Rule []onvif.Config
}

// CreateRules and its properties are defined in the Onvif specification:
// https://www.onvif.org/ver20/analytics/wsdl/analytics.wsdl#op.CreateRules
type CreateRules struct {
	XMLName            string                `xml:"tan:CreateRules"`
	ConfigurationToken onvif.ReferenceToken  `xml:"tan:ConfigurationToken"`
	Rule               []onvif.ConfigRequest `xml:"tan:Rule"`
}

type ItemListExtension xsd.AnyType

type CreateRulesResponse struct{}

// DeleteRules and its properties are defined in the Onvif specification:
// https://www.onvif.org/ver20/analytics/wsdl/analytics.wsdl#op.DeleteRules
type DeleteRules struct {
	XMLName            string               `xml:"tan:DeleteRules"`
	ConfigurationToken onvif.ReferenceToken `xml:"tan:ConfigurationToken"`
	RuleName           []xsd.String         `xml:"tan:RuleName"`
}

type DeleteRulesResponse struct{}

// GetRuleOptions and its properties are defined in the Onvif specification:
// https://www.onvif.org/ver20/analytics/wsdl/analytics.wsdl#op.GetRuleOptions
type GetRuleOptions struct {
	XMLName            string               `xml:"tan:GetRuleOptions"`
	RuleType           xsd.QName            `xml:"tan:RuleType"`
	ConfigurationToken onvif.ReferenceToken `xml:"tan:ConfigurationToken"`
}

type GetRuleOptionsResponse struct {
	RuleOptions []RuleOptions
}

type RuleOptions struct {
	RuleType                  *xsd.String                `json:","`
	Name                      *xsd.String                `json:"," xml:",attr"`
	Type                      *xsd.String                `json:"," xml:",attr"`
	MinOccurs                 *xsd.String                `json:"," xml:"minOccurs,attr"`
	MaxOccurs                 *xsd.String                `json:"," xml:"maxOccurs,attr"`
	AnalyticsModule           *xsd.String                `json:","`
	IntRange                  *IntRange                  `json:","`
	StringItems               *StringItems               `json:","`
	PolygonOptions            *PolygonOptions            `json:","`
	MotionRegionConfigOptions *MotionRegionConfigOptions `json:","`
	StringList                *xsd.String                `json:","`
}

type PolygonOptions struct {
	VertexLimits VertexLimits
}

type VertexLimits struct {
	Min int
	Max int
}

type MotionRegionConfigOptions struct {
	DisarmSupport  bool
	PolygonSupport bool
	PolygonLimits  VertexLimits
}

// ModifyRules and its properties are defined in the Onvif specification:
// https://www.onvif.org/ver20/analytics/wsdl/analytics.wsdl#op.ModifyRules
type ModifyRules struct {
	XMLName            string                `xml:"tan:ModifyRules"`
	ConfigurationToken onvif.ReferenceToken  `xml:"tan:ConfigurationToken"`
	Rule               []onvif.ConfigRequest `xml:"tan:Rule"`
}

type ModifyRulesResponse struct{}

type GetServiceCapabilities struct {
	XMLName string `xml:"tan:GetServiceCapabilities"`
}
