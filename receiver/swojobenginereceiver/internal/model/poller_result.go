package model

type PollerAssignment struct {
	NetObjectType string `json:"NetObjectType"`
	NetObjectID   int    `json:"NetObjectID"`
	PollerType    string `json:"PollerType"`
	Enabled       bool   `json:"Enabled"`
}

type PCUObject struct {
	SerialNumber       string  `json:"SerialNumber"`
	FirmwareVersion    string  `json:"FirmwareVersion"`
	BatteryPackCount   int64   `json:"BatteryPackCount"`
	BatteryCapacity    float64 `json:"BatteryCapacity"`
	BatteryTemperature float64 `json:"BatteryTemperature"`
	TimeOnBattery      float64 `json:"TimeOnBattery"`
	ReplaceIndicator   int64   `json:"ReplaceIndicator"`
	BasicBatteryStatus int64   `json:"BasicBatteryStatus"`
	OutputStatus       int64   `json:"OutputStatus"`
	Status             int64   `json:"Status"`
	OutputPercentLoad  float64 `json:"OutputPercentLoad"`
	RunTimeRemaining   float64 `json:"RunTimeRemaining"`
	LastFailCause      int64   `json:"LastFailCause"`
	Model              string  `json:"Model"`
}

type PollerResult struct {
	Type      string    `json:"$type"`
	PCUObject PCUObject `json:"PCUObject"`
	Outcome   string    `json:"Outcome"`
}

type Result struct {
	PollerAssignment PollerAssignment `json:"PollerAssignment"`
	PollerResult     PollerResult     `json:"PollerResult"`
	ResultType       string           `json:"ResultType"`
	Outcome          string           `json:"Outcome"`
}

type PollerJobOutput struct {
	Results []Result `json:"Results"`
}
