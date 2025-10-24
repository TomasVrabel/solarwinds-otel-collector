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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
	"go.uber.org/zap"

	"github.com/solarwinds/solarwinds-otel-collector/extension/swjobengineextension/internal"
	"go.opentelemetry.io/collector/pdata/plog"
)

// DiscoveryLogReceiver defines the interface for components that can receive discovery logs
type DiscoveryLogReceiver interface {
	ProcessDiscoveryLog(logRecord plog.LogRecord)
}

type SwJobEngineExtension struct {
	logger *zap.Logger
	config *internal.Config
	client *internal.JobEngineClient

	discoveryContext *DiscoveryContext
}

type DiscoveryContext struct {
	credentials map[int]internal.CredentialSnmpV2
}

func newExtension(ctx context.Context, set extension.Settings, cfg *internal.Config) (*SwJobEngineExtension, error) {
	set.Logger.Info("Creating Solarwinds Extension")
	set.Logger.Info("Config", zap.Any("config", cfg))

	e := &SwJobEngineExtension{
		logger: set.Logger,
		config: cfg,
		discoveryContext: &DiscoveryContext{
			credentials: make(map[int]internal.CredentialSnmpV2),
		},
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
	e.logger.Info("All jobs deleted")

	e.runDiscoveryJobs()

	job_definitions, err := internal.ReadJobDefinitions(e.config.JobDefinitionsFilePath)

	if err != nil {
		e.logger.Error("Failed to read poller definitions", zap.Error(err))
	}

	e.createPollers(job_definitions)

	count, err := e.client.ListJobs()
	e.logger.Info("Current JobEngine job count", zap.Int("count", count))

	return nil
}

func (e *SwJobEngineExtension) createPollers(job_definitions []internal.PollerJob) error {

	e.logger.Info("Processing poller definitions", zap.Int("count", len(job_definitions)))

	// map pollert type to template
	var jobTemplateMap = map[string]string{
		"N.Cpu.SNMP.CiscoGen3":        "core_job_snmp_cpu.json",
		"N.Memory.SNMP.CiscoAsr":      "core_job_snmp_memory.json",
		"N.Memory.SNMP.CiscoGen3":     "core_job_snmp_memory.json",
		"N.Details.SNMP.Generic":      "core_job_inventory.json",
		"PCU.Statistics.SNMP.Generic": "core_job_pcu.json",

		"N.StatusAndResponseTime.ICMP.SendEcho": "core_job_icmp.json",
		"N.ResponseTime.ICMP.Native":            "core_job_icmp.json",
		"N.Status.ICMP.Native":                  "core_job_icmp.json",

		"I.Rediscovery.SNMP.IfTable":         "interface_job_snmp.json",
		"I.StatisticsErrors32.SNMP.IfTable":  "interface_job_snmp.json",
		"I.StatisticsTraffic.SNMP.Universal": "interface_job_snmp.json",
		"I.Status.SNMP.IfTable":              "interface_job_snmp.json",

		"N.Uptime.SNMP.Generic": "core_job_snmp_uptime.json",
	}

	for _, job_definition := range job_definitions {
		e.logger.Debug("Job Definition", zap.Any("job_definition", job_definition))

		variables := make(map[string]string)
		for _, variable := range job_definition.Variables {
			variables[variable.Name] = variable.Value
		}

		// Get the handler for this poller type
		templateId, exists := jobTemplateMap[job_definition.PollerType]
		if !exists {
			e.logger.Warn("Unknown poller type", zap.String("type", job_definition.PollerType))
			continue
		}

		// Create the job
		uid, err := e.client.CreateJob(templateId, job_definition, variables)
		if err != nil {
			e.logger.Error("Failed to create job",
				zap.String("type", job_definition.PollerType),
				zap.Error(err))
			continue
		}

		e.logger.Debug(fmt.Sprintf("Creating %s job", job_definition.PollerType),
			zap.String("uid", uid))
	}

	return nil
}

func (e *SwJobEngineExtension) runDiscoveryJobs() error {
	discovery_definitions, err := internal.ReadDiscoveryDefinitions(e.config.DiscoveryDefinitionsFilePath)

	if err != nil {
		e.logger.Error("Failed to read discovery definitions", zap.Error(err))
	}

	e.logger.Info("Discovery definitions", zap.Int("count", len(discovery_definitions)))

	// Run discovery jobs
	for _, discovery_definition := range discovery_definitions {
		e.logger.Info("Job Definition", zap.Any("discovery_definition", discovery_definition))

		// Cache credentials
		for _, credential := range discovery_definition.CredentialSnmpV2 {
			e.discoveryContext.credentials[credential.Id] = credential
		}

		// Create the job
		uid, err := e.client.CreateJob_CoreDiscovery(discovery_definition)
		if err != nil {
			e.logger.Error("Failed to create discovery job",
				zap.String("name", discovery_definition.Name),
				zap.Error(err))
			continue
		}

		e.logger.Info(fmt.Sprintf("Discovery job created"), zap.String("uid", uid))
	}

	return nil
}

// ProcessDiscoveryLog processes a discovery log record received from the discovery processor
func (e *SwJobEngineExtension) ProcessDiscoveryLog(logRecord plog.LogRecord) {
	var pollerJobs []internal.PollerJob

	e.logger.Info("Received discovery log",
		zap.String("body", logRecord.Body().AsString()),
		zap.Any("attributes", logRecord.Attributes().AsRaw()),
		zap.Time("timestamp", logRecord.Timestamp().AsTime()))

	logBody := logRecord.Body().AsString()

	// Deserialize JSON body as list of map[string]string
	var discoveryData []map[string]string
	if err := json.Unmarshal([]byte(logBody), &discoveryData); err != nil {
		e.logger.Error("Failed to deserialize discovery log JSON",
			zap.Error(err),
			zap.String("logBody", logBody))
		return
	}

	e.logger.Info("Successfully deserialized discovery log",
		zap.Int("itemCount", len(discoveryData)))

	// Process each discovery item
	for i, item := range discoveryData {
		e.logger.Debug("Processing discovery item",
			zap.Int("index", i),
			zap.Any("item", item))

		// Extract common fields if they exist
		ip := item["ip"]
		pollerType := item["pollerType"]
		credentialIdStr := item["credentialId"]

		credentialId, _ := strconv.Atoi(credentialIdStr)

		credential, credentialExists := e.discoveryContext.credentials[credentialId]

		if !credentialExists {
			e.logger.Warn("No credential found for node",
				zap.String("ip", ip),
				zap.Int("credentialId", credentialId),
				zap.String("pollerName", pollerType))

			// ICMP jobs don't need credential, use dummy credential
			credential = internal.CredentialSnmpV2{
				Community: "credential-not-available",
			}
		}

		// for now process only nodes, ignore interfaces
		if strings.HasPrefix(pollerType, "I.") {

			ifIndex := item["ifIndex"]

			jobContext := internal.JobContext{
				Type: internal.JobTypePoll,
				Entity: internal.EntityContext{
					EntityType: "NetworkInterface",
					EntityId: map[string]string{
						"sw.collector.Interfaces.Uri": "cloudId:" + ip + "-" + ifIndex,
						"sw.collector.Nodes.Category": "1",
					},
					EntityAttributes: map[string]string{},
					Relations: []internal.Relation{
						internal.Relation{
							RelationType: "Has",
							Entity: internal.EntityContext{
								EntityType: "NetworkDevice",
								EntityId: map[string]string{
									"sw.collector.Nodes.Uri":      "cloudId:" + ip,
									"sw.collector.Nodes.Category": "1",
								},
								EntityAttributes: map[string]string{
									"sw.collector.Nodes.IPAddress": ip,
								},
							},
						},
					},
				},
			}

			pollerJob := internal.PollerJob{
				ID:         "job-" + pollerType + "-" + ip + "-" + ifIndex,
				PollerType: pollerType,
				State:      jobContext,
				Frequency:  e.config.DefaultJobFrequency,
				Variables: []internal.Variable{
					{Name: "Community", Value: credential.Community},
					{Name: "IP", Value: ip},
					{Name: "InterfaceIndex", Value: ifIndex},
					{Name: "NetObjectId", Value: "0"},
				},
			}

			pollerJobs = append(pollerJobs, pollerJob)

			continue
		}

		// create Job State, TBD: create proper job state
		jobContext := internal.JobContext{
			Type: internal.JobTypePoll,
			Entity: internal.EntityContext{
				EntityType: "NetworkDevice",
				EntityId: map[string]string{
					"sw.collector.Nodes.Uri":      "cloudId:" + ip,
					"sw.collector.Nodes.Category": "1",
				},
				EntityAttributes: map[string]string{
					"sw.collector.Nodes.IPAddress": ip,
				},
			},
		}

		pollerJob := internal.PollerJob{
			ID:         "job-" + pollerType + "-" + ip,
			PollerType: pollerType,
			State:      jobContext,
			Frequency:  e.config.DefaultJobFrequency,
			Variables: []internal.Variable{
				{Name: "Community", Value: credential.Community},
				{Name: "IP", Value: ip},
				{Name: "NetObjectId", Value: "0"},
			},
		}

		pollerJobs = append(pollerJobs, pollerJob)
	}

	// creating pollers asychronously
	go func() {
		time.Sleep(time.Second * 2) // Delay to ensure all jobs are processed
		e.createPollers(pollerJobs)
	}()

	e.logger.Info("Finished processing discovery job results", zap.String("job_id", "111111111111111111111"))
}

func (e *SwJobEngineExtension) Shutdown(ctx context.Context) error {
	e.logger.Info("Shutting down SolarWinds JobEngine Extension")

	defer e.client.Cancel()

	// Everything must be shut down, regardless of the failure.
	//return e.heartbeat.Shutdown(ctx)

	return nil
}
