package model

import (
	"encoding/json"
	"strconv"
	"time"

	jobEngineEvents "github.com/solarwinds/solarwinds-otel-collector/pkg/job-engine-events"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type JobResultContext struct {
	State string
}

func appendMeasuremnt(metricName string, metricValue string, metricDesc string, value float64, scopeMetrics pmetric.MetricSlice) {
	rttMetric := scopeMetrics.AppendEmpty()
	rttMetric.SetName(metricName)
	rttMetric.SetUnit(metricValue)
	rttMetric.SetDescription(metricDesc)
	rttMetricDataPoints := rttMetric.SetEmptyGauge().DataPoints()

	appendDataPoint(rttMetricDataPoints, value)
}

func appendDataPoint(metricDataPoints pmetric.NumberDataPointSlice, value float64) {
	dp := metricDataPoints.AppendEmpty()
	dp.SetDoubleValue(value)
	dp.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
	//dp.Attributes().PutStr(ATTR_PEER_IP, "10.10.10.10")
}

func formatUri(entityType string, entityId int) string {
	return entityType + ":" + strconv.Itoa(entityId)
}

func createTransformFunc[T any](processFunc func(*JobResultContext, *pmetric.ResourceMetrics, *T, *PollerAssignment) error) func([]byte, *JobResultContext, *pmetric.ResourceMetrics, *PollerAssignment) error {
	return func(pollerResultData []byte, context *JobResultContext, rm *pmetric.ResourceMetrics, assignment *PollerAssignment) error {
		var pollerResult T
		if err := json.Unmarshal(pollerResultData, &pollerResult); err != nil {
			return err
		}
		return processFunc(context, rm, &pollerResult, assignment)
	}
}

var transformMap = map[string]func([]byte, *JobResultContext, *pmetric.ResourceMetrics, *PollerAssignment) error{
	"MultiCoreCpuLoadResult":          createTransformFunc(addResult_CPU),
	"CiscoMemoryResult":               createTransformFunc(addResult_Memory),
	"NodeDetailsResult":               createTransformFunc(addResult_CoreInventory),
	"DeclarativePollerResultBase":     createTransformFunc(addResult_PCU),
	"NodeStatusAndResponseTimeResult": createTransformFunc(addResult_Echo),
}

func Transform_toMetrics(in *jobEngineEvents.NotifyJobFinishedRequest, logger *zap.Logger) (*pmetric.Metrics, error) {
	metrics := pmetric.NewMetrics()

	for _, job := range in.FinishedJobs {
		var state = job.GetState()
		var outputStr = job.GetResult().GetOutput()
		var pollerError = job.GetResult().GetError()

		if pollerError != "" {
			logger.Warn("Poller job finished with error", zap.String("pollerError", pollerError), zap.String("state", state), zap.String("job_id", job.GetScheduledJobId()))
			continue
		}

		var jobResultContext = &JobResultContext{
			State: state,
		}

		var root PollerJobOutput
		err := json.Unmarshal([]byte(outputStr), &root)
		if err != nil {
			logger.Error("Error deserializing JSON", zap.String("input", string(outputStr)), zap.String("pollerError", pollerError), zap.String("state", state), zap.Error(err))
			continue
		}

		for _, result := range root.Results {
			// Add a ResourceMetrics to the Metrics object
			rm := metrics.ResourceMetrics().AppendEmpty()

			transformFunction := transformMap[result.ResultType]

			if transformFunction == nil {
				message := "unknown result type"
				logger.Error(message, zap.String("ResultType", result.ResultType), zap.String("result", string(outputStr)), zap.Error(err))
				continue
			}

			err = transformFunction(result.PollerResult, jobResultContext, &rm, &result.PollerAssignment)

			if err != nil {
				message := "Error processing PollerResult"
				logger.Error(message, zap.String("ResultType", result.ResultType), zap.Error(err))
				continue
			}
		}
	}
	return &metrics, nil
}
