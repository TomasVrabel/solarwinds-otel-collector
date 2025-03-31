package model

import (
	"encoding/json"
)

type PollerAssignment struct {
	NetObjectType string `json:"NetObjectType"`
	NetObjectID   int    `json:"NetObjectID"`
	PollerType    string `json:"PollerType"`
	Enabled       bool   `json:"Enabled"`
}

type Result struct {
	PollerAssignment PollerAssignment `json:"PollerAssignment"`
	PollerResult     json.RawMessage  `json:"PollerResult"`
	ResultType       string           `json:"ResultType"`
	Outcome          string           `json:"Outcome"`
}

type PollerJobOutput struct {
	Results []Result `json:"Results"`
}
