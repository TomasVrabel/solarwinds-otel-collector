package model

import (
	"time"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

type Buffers struct {
	BufferNoMem  float64 `json:"BufferNoMem"`
	BufferSmMiss float64 `json:"BufferSmMiss"`
	BufferMdMiss float64 `json:"BufferMdMiss"`
	BufferBgMiss float64 `json:"BufferBgMiss"`
	BufferLgMiss float64 `json:"BufferLgMiss"`
	BufferHgMiss float64 `json:"BufferHgMiss"`
}

type ResponseTime struct {
	ResponseTimeMilliseconds int `json:"ResponseTimeMilliseconds"`
}

type Status struct {
	PolledStatus string `json:"PolledStatus"`
}

type CPUPollerResult struct {
	Type            string         `json:"$type"`
	CpuLoadPerIndex map[string]int `json:"CpuLoadPerIndex,omitempty"`
	Outcome         string         `json:"Outcome"`
}

type MemoryPollerResult struct {
	Type        string  `json:"$type"`
	Buffers     Buffers `json:"Buffers"`
	TotalMemory float64 `json:"TotalMemory"`
	UsedMemory  float64 `json:"UsedMemory"`
	Outcome     string  `json:"Outcome"`
}
type EchoPollerResult struct {
	Type            string       `json:"$type"`
	IPAddress       string       `json:"IPAddress"`
	IsDynamic       bool         `json:"IsDynamic"`
	DNS             *string      `json:"DNS"`
	ResponseTime    ResponseTime `json:"ResponseTime"`
	Status          Status       `json:"Status"`
	ObservationTime time.Time    `json:"ObservationTime"`
	FastPollState   string       `json:"FastPollState"`
	Outcome         string       `json:"Outcome"`
}

type NodeDetailsPollerResult struct {
	Type        string   `json:"$type"`
	VendorOID   string   `json:"VendorOID"`
	Contact     string   `json:"Contact"`
	Location    string   `json:"Location"`
	Description string   `json:"Description"`
	SysName     string   `json:"SysName"`
	SysServices int      `json:"SysServices"`
	MAC         []string `json:"MAC"`
	Outcome     string   `json:"Outcome"`
}

func addResult_CoreInventory(rm *pmetric.ResourceMetrics, result *NodeDetailsPollerResult, assignment *PollerAssignment) error {
	// Set resource attributes
	resource := rm.Resource()
	resource.Attributes().PutStr("sw.collector.Nodes.Category", "1")
	resource.Attributes().PutStr("sw.collector.Nodes.Uri", formatUri("NetworkDevice", assignment.NetObjectID))

	resource.Attributes().PutStr("sw.collector.Nodes.Location", result.Location)
	resource.Attributes().PutStr("sw.collector.Nodes.Contact", result.Contact)
	resource.Attributes().PutStr("sw.collector.Nodes.SysDescription", "TBD")
	resource.Attributes().PutStr("sw.collector.Nodes.Description", result.Description)
	resource.Attributes().PutStr("sw.collector.Nodes.DisplayName", result.SysName)
	resource.Attributes().PutStr("sw.collector.Nodes.Vendor", result.VendorOID)

	return nil
}

func averageMapValues(data map[string]int) float64 {
	if len(data) == 0 {
		return 0
	}

	sum := 0
	for _, value := range data {
		sum += value
	}

	return float64(sum) / float64(len(data))
}

func addResult_CPU(rm *pmetric.ResourceMetrics, result *CPUPollerResult, assignment *PollerAssignment) error {
	// Set resource attributes
	resource := rm.Resource()
	resource.Attributes().PutStr("sw.collector.Nodes.Category", "1")
	resource.Attributes().PutStr("sw.collector.Nodes.Uri", formatUri("NetworkDevice", assignment.NetObjectID))

	scopeMetrics := rm.ScopeMetrics().AppendEmpty().Metrics()

	cpuLoad := averageMapValues(result.CpuLoadPerIndex)
	appendMeasuremnt("sw.collector.CPULoad.AvgLoad", "percentage", "CPU Utilization", cpuLoad, scopeMetrics)

	return nil
}

func addResult_Memory(rm *pmetric.ResourceMetrics, result *MemoryPollerResult, assignment *PollerAssignment) error {

	// Set resource attributes
	resource := rm.Resource()
	resource.Attributes().PutStr("sw.collector.Nodes.Category", "1")
	resource.Attributes().PutStr("sw.collector.Nodes.Uri", formatUri("NetworkDevice", assignment.NetObjectID))

	scopeMetrics := rm.ScopeMetrics().AppendEmpty().Metrics()

	appendMeasuremnt("sw.collector.CPULoad.AvgMemoryUsed", "By", "Memory used", result.UsedMemory, scopeMetrics)
	appendMeasuremnt("sw.collector.CPULoad.TotalMemory", "By", "Total memory", result.TotalMemory, scopeMetrics)

	return nil
}

func addResult_Echo(rm *pmetric.ResourceMetrics, result *EchoPollerResult, assignment *PollerAssignment) error {

	// Set resource attributes
	resource := rm.Resource()
	resource.Attributes().PutStr("sw.collector.Nodes.Category", "1")
	resource.Attributes().PutStr("sw.collector.Nodes.Uri", formatUri("NetworkDevice", assignment.NetObjectID))

	scopeMetrics := rm.ScopeMetrics().AppendEmpty().Metrics()

	if result.Status.PolledStatus == "Up" {
		appendMeasuremnt("sw.collector.ResponseTime.AvgResponseTime", "ms", "Response time", float64(result.ResponseTime.ResponseTimeMilliseconds), scopeMetrics)
		appendMeasuremnt("sw.collector.ResponseTime.Availability", "percentage", "Average availability", float64(100), scopeMetrics)
		appendMeasuremnt("sw.collector.ResponseTime.PercentLoss", "percentage", "Percentage of lost packets", float64(0), scopeMetrics)
	} else {
		appendMeasuremnt("sw.collector.ResponseTime.Availability", "percentage", "Average availability", float64(0), scopeMetrics)
		appendMeasuremnt("sw.collector.ResponseTime.PercentLoss", "percentage", "Percentage of lost packets", float64(100), scopeMetrics)

	}
	return nil
}
