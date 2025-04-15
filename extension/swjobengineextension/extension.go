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

package swjobengineextension

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
	"go.uber.org/zap"

	"github.com/solarwinds/solarwinds-otel-collector/extension/swjobengineextension/internal"
)

type SwJobEngineExtension struct {
	logger *zap.Logger
	config *internal.Config
	client *internal.JobEngineClient
}

func newExtension(ctx context.Context, set extension.Settings, cfg *internal.Config) (*SwJobEngineExtension, error) {
	set.Logger.Info("Creating Solarwinds Extension")
	set.Logger.Info("Config", zap.Any("config", cfg))

	e := &SwJobEngineExtension{
		logger: set.Logger,
		config: cfg,
	}
	/*
		var err error
		e.heartbeat, err = internal.NewHeartbeat(ctx, set, cfg)
		if err != nil {
			return nil, err
		}
	*/
	return e, nil
}

func (e *SwJobEngineExtension) Start(ctx context.Context, host component.Host) error {
	e.logger.Info("Starting SolarWinds JobEngine Extension")

	client, err := internal.NewGrpcClient(e.config, e.logger)

	if err != nil {
		e.logger.Error("Failed to create gRPC client", zap.Error(err))
	}

	e.client = client

	err = e.client.DeleteJobs()
	e.logger.Info("Deleting all jobs")

	job_definitions, err := internal.ReadJobDefinitions(e.config.JobDefinitionsFilePath)
	e.logger.Info("Job definitions", zap.Int("count", len(job_definitions)))

	var uid string
	for _, job_definition := range job_definitions {
		e.logger.Info("Job Definition", zap.Any("job_definition", job_definition))

		variables := make(map[string]string)
		for _, variable := range job_definition.Variables {
			variables[variable.Name] = variable.Value
		}

		if job_definition.PollerType == "ICMP" {
			uid, err = e.client.CreateJob_ICMP(variables)
			e.logger.Info("Creating ICMP job", zap.String("uid", uid))
		}
		if job_definition.PollerType == "SNMP" {
			uid, err = e.client.CreateJob_SNMP(variables)
			e.logger.Info("Creating SNMP job", zap.String("uid", uid))
		}
		if job_definition.PollerType == "PCU" {
			uid, err = e.client.CreateJob_PCU(variables)
			e.logger.Info("Creating PCU job", zap.String("uid", uid))
		}
		if job_definition.PollerType == "CoreInventory" {
			uid, err = e.client.CreateJob_CoreInventory(variables)
			e.logger.Info("Creating CoreInventory job", zap.String("uid", uid))
		}
	}

	count, err := e.client.ListJobs()
	e.logger.Info("Current JobEngine job count", zap.Int("count", count))

	return nil
}

func (e *SwJobEngineExtension) Shutdown(ctx context.Context) error {
	e.logger.Info("Shutting down SolarWinds JobEngine Extension")

	defer e.client.Cancel()

	// Everything must be shut down, regardless of the failure.
	//return e.heartbeat.Shutdown(ctx)

	return nil
}
