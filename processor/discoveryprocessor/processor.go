// Copyright 2025 SolarWinds Worldwide, LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package discoveryprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/processor/processorhelper"
	"go.uber.org/zap"
)

const (
	discoveryAttribute = "sw.discovery"
)

// DiscoveryLogReceiver defines the interface for extensions that can receive discovery logs
type DiscoveryLogReceiver interface {
	ProcessDiscoveryLog(logRecord plog.LogRecord)
}

type discoveryProcessor struct {
	logger                 *zap.Logger
	config                 *Config
	telemetrySettings      component.TelemetrySettings
	jobEngineExtensionName string
	jobEngineExtension     DiscoveryLogReceiver
}

// processLogs processes all log records and removes discovery logs from the pipeline
func (dp *discoveryProcessor) processLogs(_ context.Context, ld plog.Logs) (plog.Logs, error) {
	resourceLogs := ld.ResourceLogs()

	// Process each resource log in reverse order to safely remove items
	for i := resourceLogs.Len() - 1; i >= 0; i-- {
		rl := resourceLogs.At(i)
		scopeLogs := rl.ScopeLogs()

		// Process each scope log in reverse order
		for j := scopeLogs.Len() - 1; j >= 0; j-- {
			sl := scopeLogs.At(j)
			logRecords := sl.LogRecords()

			// Process each log record in reverse order
			for k := logRecords.Len() - 1; k >= 0; k-- {
				lr := logRecords.At(k)

				if dp.isDiscoveryLog(lr.Attributes()) {
					// Send to job engine extension
					if dp.jobEngineExtension != nil {
						dp.jobEngineExtension.ProcessDiscoveryLog(lr)
						dp.logger.Debug("Processed and removed discovery log from pipeline")
					} else {
						dp.logger.Warn("Job engine extension not available for discovery log processing")
					}

					// Remove the discovery log record
					logRecords.RemoveIf(func(record plog.LogRecord) bool {
						return record == lr
					})
				}
			}

			// Remove empty scope logs
			if sl.LogRecords().Len() == 0 {
				scopeLogs.RemoveIf(func(scope plog.ScopeLogs) bool {
					return scope == sl
				})
			}
		}

		// Remove empty resource logs
		if rl.ScopeLogs().Len() == 0 {
			resourceLogs.RemoveIf(func(resource plog.ResourceLogs) bool {
				return resource == rl
			})
		}
	}

	// If all logs were removed, skip processing
	if ld.ResourceLogs().Len() == 0 {
		dp.logger.Debug("All logs were discovery logs, skipping further processing")
		return ld, processorhelper.ErrSkipProcessingData
	}

	return ld, nil
}

// isDiscoveryLog checks if a log record has the discovery=true attribute
func (dp *discoveryProcessor) isDiscoveryLog(attributes pcommon.Map) bool {
	discoveryValue, exists := attributes.Get(discoveryAttribute)
	if !exists {
		return false
	}

	// Check if the value is true (as boolean or string)
	switch discoveryValue.Type() {
	case pcommon.ValueTypeBool:
		return discoveryValue.Bool()
	case pcommon.ValueTypeStr:
		return discoveryValue.Str() == "true"
	default:
		return false
	}
}

// Start initializes the processor and finds the job engine extension
func (dp *discoveryProcessor) Start(_ context.Context, host component.Host) error {
	dp.logger.Info("Starting discovery processor")

	// Find the job engine extension
	extensions := host.GetExtensions()
	for id, ext := range extensions {
		dp.logger.Info("Checking extension", zap.String("extension_id", id.String()))
		if id.String() == dp.jobEngineExtensionName {
			if discoveryReceiver, ok := ext.(DiscoveryLogReceiver); ok {
				dp.jobEngineExtension = discoveryReceiver
				dp.logger.Info("Found job engine extension for discovery processing",
					zap.String("extension_id", id.String()))
				break
			} else {
				dp.logger.Warn("Extension does not implement DiscoveryLogReceiver interface",
					zap.String("extension_id", id.String()))
			}
		}
	}

	if dp.jobEngineExtension == nil {
		dp.logger.Warn("Job engine extension not found, discovery logs will be filtered but not processed",
			zap.String("expected_extension_name", dp.jobEngineExtensionName))
	}

	return nil
}

// Shutdown cleans up the processor
func (dp *discoveryProcessor) Shutdown(_ context.Context) error {
	dp.logger.Info("Shutting down discovery processor")
	return nil
}
