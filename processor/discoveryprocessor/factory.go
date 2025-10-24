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
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

const (
	typeStr   = "discovery"
	stability = component.StabilityLevelBeta
)

var (
	consumerCapabilities = consumer.Capabilities{MutatesData: true}
)

// NewFactory creates a new processor factory for the discovery processor
func NewFactory() processor.Factory {
	return processor.NewFactory(
		component.MustNewType(typeStr),
		func() component.Config { return createDefaultConfig() },
		processor.WithLogs(createLogsProcessor, stability),
	)
}

// createLogsProcessor creates a new logs processor instance
func createLogsProcessor(
	ctx context.Context,
	params processor.Settings,
	cfg component.Config,
	nextLogsConsumer consumer.Logs,
) (processor.Logs, error) {
	discoveryConfig := cfg.(*Config)

	dp := &discoveryProcessor{
		logger:                 params.Logger,
		config:                 discoveryConfig,
		telemetrySettings:      params.TelemetrySettings,
		jobEngineExtensionName: discoveryConfig.JobEngineExtensionName,
	}

	return processorhelper.NewLogs(
		ctx,
		params,
		cfg,
		nextLogsConsumer,
		dp.processLogs,
		processorhelper.WithCapabilities(consumerCapabilities),
		processorhelper.WithStart(dp.Start),
		processorhelper.WithShutdown(dp.Shutdown),
	)
}
