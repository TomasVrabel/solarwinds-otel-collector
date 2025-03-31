package model

import (
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type PCUPollerResult struct {
	Type      string    `json:"$type"`
	PCUObject PCUObject `json:"PCUObject"`
	Outcome   string    `json:"Outcome"`
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

func addResult_PCU(rm *pmetric.ResourceMetrics, result *PCUPollerResult, assignment *PollerAssignment) error {
	// Set resource attributes
	resource := rm.Resource()
	resource.Attributes().PutStr("sw.collector.EntityType", "sw.collector.PowerControlUnit")
	resource.Attributes().PutStr("sw.collector.Nodes.Category", "1")
	resource.Attributes().PutStr("sw.collector.Nodes.Uri", formatUri("NetworkDevice", assignment.NetObjectID))
	resource.Attributes().PutStr("sw.collector.PowerControlUnit.Uri", formatUri("PowerControlUnit", assignment.NetObjectID))

	resource.Attributes().PutInt("sw.collector.PowerControlUnit.BasicBatteryStatus", result.PCUObject.BasicBatteryStatus)
	resource.Attributes().PutInt("sw.collector.PowerControlUnit.BatteryPackCount", result.PCUObject.BatteryPackCount)
	resource.Attributes().PutStr("sw.collector.PowerControlUnit.DisplayName", "PCU")
	resource.Attributes().PutStr("sw.collector.PowerControlUnit.FirmwareVersion", result.PCUObject.FirmwareVersion)
	resource.Attributes().PutInt("sw.collector.PowerControlUnit.LastFailCause", result.PCUObject.LastFailCause)
	resource.Attributes().PutInt("sw.collector.PowerControlUnit.ReplaceIndicator", result.PCUObject.ReplaceIndicator)
	resource.Attributes().PutStr("sw.collector.PowerControlUnit.SerialNumber", result.PCUObject.SerialNumber)
	resource.Attributes().PutInt("sw.collector.PowerControlUnit.Status", result.PCUObject.Status)

	scopeMetrics := rm.ScopeMetrics().AppendEmpty().Metrics()

	appendMeasuremnt("sw.collector.PowerControlUnit.Metrics.AvgBatteryCapacity", "percentage", "Battery capacity", result.PCUObject.BatteryCapacity, scopeMetrics)
	appendMeasuremnt("sw.collector.PowerControlUnit.Metrics.AvgBatteryTemperature", "°C", "Battery temperature", result.PCUObject.BatteryTemperature, scopeMetrics)
	appendMeasuremnt("sw.collector.PowerControlUnit.Metrics.AvgTimeOnBattery", "minutes", "Time on Battery", result.PCUObject.TimeOnBattery, scopeMetrics)
	appendMeasuremnt("sw.collector.PowerControlUnit.Metrics.AvgOutputPercentLoad", "percentage", "Output Percent Load", result.PCUObject.OutputPercentLoad, scopeMetrics)
	appendMeasuremnt("sw.collector.PowerControlUnit.Metrics.AvgRunTimeRemaining", "minutes", "Run Time Remaining on battery", result.PCUObject.RunTimeRemaining, scopeMetrics)

	return nil
}
