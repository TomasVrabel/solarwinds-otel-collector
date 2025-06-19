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
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componentstatus"
	"go.opentelemetry.io/collector/extension"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/solarwinds/solarwinds-otel-collector/extension/swjobengineextension/internal"
	jobEngineEvents "github.com/solarwinds/solarwinds-otel-collector/pkg/job-engine-events"
	job_engine_events "github.com/solarwinds/solarwinds-otel-collector/pkg/job-engine-events"
)

type SwJobEngineExtension struct {
	logger *zap.Logger
	config *internal.Config
	client *internal.JobEngineClient

	discoveryContext *DiscoveryContext

	serverGRPC *grpc.Server
	server     *server
}

type DiscoveryContext struct {
	credentials map[int]internal.CredentialSnmpV2
}

// server is used to implement Job Engine evetns
type server struct {
	extension *SwJobEngineExtension
	logger    *zap.Logger

	jobEngineEvents.UnimplementedJobEngineEventsServer
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

	err = e.startEventEndpoint(ctx, host)
	if err != nil {
		e.logger.Error("Failed to start event endpoint", zap.Error(err))
	}

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
		"N.Details.SNMP.Generic":      "core_job_inventory.json",
		"PCU.Statistics.SNMP.Generic": "core_job_pcu.json",

		"N.StatusAndResponseTime.ICMP.SendEcho": "core_job_icmp.json",
		"N.ResponseTime.ICMP.Native":            "core_job_icmp.json",
		"N.Status.ICMP.Native":                  "core_job_icmp.json",

		"N.Uptime.SNMP.Generic": "core_job_snmp_uptime.json",
	}

	for _, job_definition := range job_definitions {
		e.logger.Info("Job Definition", zap.Any("job_definition", job_definition))

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

		e.logger.Info(fmt.Sprintf("Creating %s job", job_definition.PollerType),
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

func (e *SwJobEngineExtension) startEventEndpoint(ctx context.Context, host component.Host) error {
	var err error

	e.serverGRPC = grpc.NewServer()
	if err != nil {
		return fmt.Errorf("failed create grpc server error: %w", err)
	}

	e.server = &server{
		extension: e,
		logger:    e.logger,
	}

	jobEngineEvents.RegisterJobEngineEventsServer(e.serverGRPC, e.server)

	err = e.startGRPCServer(ctx, host)
	if err != nil {
		return fmt.Errorf("failed to start grpc server error: %w", err)
	}

	return err
}

// OnJobFinished implements jobEngineEvents.OnJobFinished
func (r *server) NotifyJobFinished(_ context.Context, in *jobEngineEvents.NotifyJobFinishedRequest) (*jobEngineEvents.NotifyJobFinishedResponse, error) {
	r.logger.Info("Received job finished notification")

	var pollerJobs []internal.PollerJob
	var deviceMap = make(map[int]internal.Node)

	for _, job := range in.FinishedJobs {
		r.logger.Info("Job finished",
			zap.String("scheduled_job_id", job.ScheduledJobId),
			zap.String("job_id", job.Result.JobId),
			//zap.String("output", string(job.GetResult().GetOutput())),
			zap.String("job_state", job.State))

		output := job.GetResult().GetOutput()
		if len(output) > 0 {
			// Try to parse as discovery job result
			discoveryResult, err := internal.ParseDiscoveryJobResult(output)
			if err != nil {
				r.logger.Error("Failed to parse discovery job result", zap.Error(err))
			} else {
				r.logger.Info("Parsed discovery job result",
					zap.Int("engineId", discoveryResult.EngineID),
					zap.Int("profileId", discoveryResult.ProfileID),
					zap.Int("nodeCount", len(discoveryResult.PluginResults.PluginItem.ArrayOfDiscoveryPluginResultBase.DiscoveryPluginResultBase.DiscoveredNodes.Nodes)),
					zap.String("base.pluginTypeName", discoveryResult.PluginResults.PluginItem.ArrayOfDiscoveryPluginResultBase.DiscoveryPluginResultBase.PluginTypeName),
					zap.String("base.profileId", discoveryResult.PluginResults.PluginItem.ArrayOfDiscoveryPluginResultBase.DiscoveryPluginResultBase.ProfileID),
					zap.Bool("base.allowCrossEngineNodeUpdates", discoveryResult.PluginResults.PluginItem.ArrayOfDiscoveryPluginResultBase.DiscoveryPluginResultBase.AllowCrossEngineNodeUpdates))

				discoveryPluginResultBase := discoveryResult.PluginResults.PluginItem.ArrayOfDiscoveryPluginResultBase.DiscoveryPluginResultBase

				// Log discovered nodes
				for _, node := range discoveryPluginResultBase.DiscoveredNodes.Nodes {
					r.logger.Info("Discovered node",
						zap.Int("id", node.ID),
						zap.String("ip", node.IP),
						zap.String("name", node.Name),
						zap.String("type", node.Type),
						zap.String("description", node.Description),
						zap.Int("profileId", node.ProfileID),
						zap.String("status", node.Status),
						zap.String("location", node.Location),
						zap.String("hostname", node.Hostname),
						zap.String("contact", node.Contact),
						zap.Bool("isExternal", node.IsExternal),
						zap.String("oid", node.OID),
						zap.Bool("isSelected", node.IsSelected),
						zap.Int("credentialId", node.CredentialID))

					deviceMap[node.ID] = node
				}

				// Log discovered pollers
				r.logger.Info("Discovered pollers",
					zap.Int("count", len(discoveryPluginResultBase.DiscoveredPollers.Pollers)))

				for _, poller := range discoveryPluginResultBase.DiscoveredPollers.Pollers {
					r.logger.Info("Discovered poller",
						zap.Int("nodeId", poller.NodeID),
						zap.String("type", poller.PollerType),
						zap.String("objectType", poller.ObjectType))

					node := deviceMap[poller.NodeID]

					credential, exists := r.extension.discoveryContext.credentials[node.CredentialID]

					if !exists {
						r.logger.Warn("No credential found for node",
							zap.Int("nodeId", poller.NodeID),
							zap.Int("credentialId", node.CredentialID))
						continue
					}

					// create Job State
					jobState := job_engine_events.NewJobState()
					jobState[jobEngineEvents.JOB_STATE_NODES_URI_ATTRIBUTE] = "networkDevice-" + node.IP
					jobState[jobEngineEvents.JOB_STATE_NODES_CATEGORY_ATTRIBUTE] = "1" // Network Device
					jobStateString, err := jobState.SerializeJobStateToString()

					if err != nil {
						r.logger.Error("Failed to serialize job state", zap.Error(err))
					}

					pollerJob := internal.PollerJob{
						ID:         "job-" + poller.PollerType + "-" + strconv.Itoa(poller.NodeID),
						PollerType: poller.PollerType,
						State:      jobStateString,
						Frequency:  r.extension.config.DefaultJobFrequency,
						Variables: []internal.Variable{
							{Name: "Community", Value: credential.Community},
							{Name: "IP", Value: node.IP},
							{Name: "NetObjectId", Value: strconv.Itoa(node.ID)},
						},
					}

					pollerJobs = append(pollerJobs, pollerJob)
				}

				// creating pollers asychronously
				go func() {
					time.Sleep(time.Second * 2) // Delay to ensure all jobs are processed
					r.extension.createPollers(pollerJobs)
				}()

				r.logger.Info("Finshed processing discovery job results", zap.String("job_id", job.Result.JobId))
			}
		} else {
			r.logger.Info("No output for job", zap.String("job_id", job.Result.JobId))
		}
	}

	return &jobEngineEvents.NotifyJobFinishedResponse{}, nil
}

func (r *SwJobEngineExtension) startGRPCServer(ctx context.Context, host component.Host) error {
	r.logger.Info("Starting GRPC server", zap.Int("endpoint.port", r.config.EndpointPort))
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", r.config.EndpointPort))
	if err != nil {
		return err
	}

	go func() {
		if errGRPC := r.serverGRPC.Serve(listener); !errors.Is(errGRPC, grpc.ErrServerStopped) && errGRPC != nil {
			componentstatus.ReportStatus(host, componentstatus.NewFatalErrorEvent(errGRPC))
		}
	}()
	return nil
}

func (e *SwJobEngineExtension) Shutdown(ctx context.Context) error {
	e.logger.Info("Shutting down SolarWinds JobEngine Extension")

	defer e.client.Cancel()

	if e.serverGRPC != nil {
		e.serverGRPC.GracefulStop()
	}

	// Everything must be shut down, regardless of the failure.
	//return e.heartbeat.Shutdown(ctx)

	return nil
}
