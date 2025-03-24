package model

import (
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func appendMeasuremnt(metricName string, metricValue string, value float64, scopeMetrics pmetric.MetricSlice) {
	rttMetric := scopeMetrics.AppendEmpty()
	rttMetric.SetName(metricName)
	rttMetric.SetUnit(metricValue)
	rttMetricDataPoints := rttMetric.SetEmptyGauge().DataPoints()

	appendDataPoint(rttMetricDataPoints, value)
}

func appendDataPoint(metricDataPoints pmetric.NumberDataPointSlice, value float64) {
	dp := metricDataPoints.AppendEmpty()
	dp.SetDoubleValue(value)
	dp.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
	//dp.Attributes().PutStr(ATTR_PEER_IP, "10.10.10.10")
}
