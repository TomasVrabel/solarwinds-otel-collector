package internal

import (
	"encoding/xml"
	"strings"
)

// OrionDiscoveryJobResult represents the root XML structure
type OrionDiscoveryJobResult struct {
	XMLName        xml.Name      `xml:"OrionDiscoveryJobResult"`
	EngineID       int           `xml:"EngineId"`
	ProfileID      int           `xml:"ProfileID"`
	PluginResults  PluginResults `xml:"PluginResults"`
	CanceledByUser bool          `xml:"CanceledByUser"`
}

// PluginResults represents the plugin results section
type PluginResults struct {
	KnownTypes KnownTypes `xml:"knownTypes"`
	PluginItem PluginItem `xml:"pluginItem"`
}

// KnownTypes represents the known types section
type KnownTypes struct {
	ArrayOfString ArrayOfString `xml:"ArrayOfstring"`
}

// ArrayOfString represents an array of strings
type ArrayOfString struct {
	Strings []string `xml:"string"`
}

// PluginItem represents a plugin item
type PluginItem struct {
	ArrayOfDiscoveryPluginResultBase ArrayOfDiscoveryPluginResultBase `xml:"ArrayOfDiscoveryPluginResultBase"`
}

// ArrayOfDiscoveryPluginResultBase represents an array of discovery plugin results
type ArrayOfDiscoveryPluginResultBase struct {
	DiscoveryPluginResultBase DiscoveryPluginResultBase `xml:"DiscoveryPluginResultBase"`
}

// DiscoveryPluginResultBase represents a discovery plugin result
type DiscoveryPluginResultBase struct {
	PluginTypeName              string                 `xml:"PluginTypeName"`
	ProfileID                   string                 `xml:"ProfileId"`
	AllowCrossEngineNodeUpdates bool                   `xml:"AllowCrossEngineNodeUpdates"`
	DiscoveredMACAddresses      DiscoveredMACAddresses `xml:"DiscoveredMACAddresses"`
	DiscoveredNodes             DiscoveredNodes        `xml:"DiscoveredNodes"`
	DiscoveredPollers           DiscoveredPollers      `xml:"DiscoveredPollers"`
}

// DiscoveredMACAddresses represents a collection of discovered MAC addresses
type DiscoveredMACAddresses struct {
	MACAddresses []MACAddress `xml:"DM"`
}

// MACAddress represents a discovered MAC address
type MACAddress struct {
	IsSelected bool   `xml:"IS"`
	MAC        string `xml:"MA"`
	NodeID     int    `xml:"NID"`
}

// DiscoveredNodes represents a collection of discovered nodes
type DiscoveredNodes struct {
	Nodes []Node `xml:"DN"`
}

// Node represents a discovered node
type Node struct {
	IsSelected   bool   `xml:"IS"`
	Contact      string `xml:"C"`
	IsExternal   bool   `xml:"EXT"`
	Hostname     string `xml:"HN"`
	ID           int    `xml:"ID"`
	IP           string `xml:"IP"`
	Location     string `xml:"L"`
	OID          string `xml:"OID"`
	ProfileID    int    `xml:"PID"`
	Description  string `xml:"SD"`
	Name         string `xml:"SN"`
	Status       string `xml:"SS"`
	Type         string `xml:"ST"`
	SNMPVersion  string `xml:"SV"`
	CredentialID int    `xml:"CID"`
}

// DiscoveredPollers represents a collection of discovered pollers
type DiscoveredPollers struct {
	Pollers []Poller `xml:"DP"`
}

// Poller represents a discovered poller
type Poller struct {
	IsSelected bool   `xml:"IS"`
	NodeID     int    `xml:"NID"`
	ObjectType string `xml:"OT"`
	PollerType string `xml:"PT"`
}

// ParseDiscoveryJobResult parses an XML string into an OrionDiscoveryJobResult struct
func ParseDiscoveryJobResult(xmlData []byte) (*OrionDiscoveryJobResult, error) {
	var result OrionDiscoveryJobResult

	// Fix encoding issues - sometimes XML may have escaped quotes
	xmlString := string(xmlData)
	xmlString = strings.ReplaceAll(xmlString, "\\\"", "\"")

	err := xml.Unmarshal([]byte(xmlString), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
