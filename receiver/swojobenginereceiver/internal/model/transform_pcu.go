package model

import (
	"encoding/json"
	"fmt"
	"strconv"

	"go.opentelemetry.io/collector/pdata/pmetric"

	jobEngineEvents "github.com/solarwinds/solarwinds-otel-collector/receiver/swojobenginereceiver/internal/job-engine-events"
)

func formatUri(entityType string, entityId int) string {
	return entityType + ":" + strconv.Itoa(entityId)
}

func Transform_PCU_toMetrics(in *jobEngineEvents.NotifyJobFinishedRequest) (*pmetric.Metrics, error) {
	metrics := pmetric.NewMetrics()

	for _, job := range in.FinishedJobs {
		var outputStr = job.GetResult().GetOutput()

		var root PollerJobOutput
		err := json.Unmarshal([]byte(outputStr), &root)
		if err != nil {
			fmt.Errorf("Error deserializing JSON:", err)
			return nil, err
		}

		// Add a ResourceMetrics to the Metrics object
		rm := metrics.ResourceMetrics().AppendEmpty()

		// Set resource attributes
		resource := rm.Resource()
		resource.Attributes().PutStr("sw.collector.EntityType", "sw.collector.PowerControlUnit")
		resource.Attributes().PutStr("sw.collector.Nodes.Category", "1")
		resource.Attributes().PutStr("sw.collector.Nodes.Uri", formatUri("NetworkDevice", root.Results[0].PollerAssignment.NetObjectID))
		resource.Attributes().PutStr("sw.collector.PowerControlUnit.Uri", formatUri("PowerControlUnit", root.Results[0].PollerAssignment.NetObjectID))

		resource.Attributes().PutInt("sw.collector.PowerControlUnit.BasicBatteryStatus", root.Results[0].PollerResult.PCUObject.BasicBatteryStatus)
		resource.Attributes().PutInt("sw.collector.PowerControlUnit.BatteryPackCount", root.Results[0].PollerResult.PCUObject.BatteryPackCount)
		resource.Attributes().PutStr("sw.collector.PowerControlUnit.DisplayName", "PCU")
		resource.Attributes().PutStr("sw.collector.PowerControlUnit.FirmwareVersion", root.Results[0].PollerResult.PCUObject.FirmwareVersion)
		resource.Attributes().PutInt("sw.collector.PowerControlUnit.LastFailCause", root.Results[0].PollerResult.PCUObject.LastFailCause)
		resource.Attributes().PutInt("sw.collector.PowerControlUnit.ReplaceIndicator", root.Results[0].PollerResult.PCUObject.ReplaceIndicator)
		resource.Attributes().PutStr("sw.collector.PowerControlUnit.SerialNumber", root.Results[0].PollerResult.PCUObject.SerialNumber)
		resource.Attributes().PutInt("sw.collector.PowerControlUnit.Status", root.Results[0].PollerResult.PCUObject.Status)

		scopeMetrics := rm.ScopeMetrics().AppendEmpty().Metrics()

		appendMeasuremnt("sw.collector.PowerControlUnit.Metrics.AvgBatteryCapacity", "percentage", root.Results[0].PollerResult.PCUObject.BatteryCapacity, scopeMetrics)
		appendMeasuremnt("sw.collector.PowerControlUnit.Metrics.AvgBatteryTemperature", "°C", root.Results[0].PollerResult.PCUObject.BatteryTemperature, scopeMetrics)
		appendMeasuremnt("sw.collector.PowerControlUnit.Metrics.AvgTimeOnBattery", "minutes", root.Results[0].PollerResult.PCUObject.TimeOnBattery, scopeMetrics)
		appendMeasuremnt("sw.collector.PowerControlUnit.Metrics.AvgOutputPercentLoad", "percentage", root.Results[0].PollerResult.PCUObject.OutputPercentLoad, scopeMetrics)
		appendMeasuremnt("sw.collector.PowerControlUnit.Metrics.AvgRunTimeRemaining", "minutes", root.Results[0].PollerResult.PCUObject.RunTimeRemaining, scopeMetrics)
	}
	return &metrics, nil
}
