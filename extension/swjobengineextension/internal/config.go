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

package internal

import (
	"errors"

	"go.opentelemetry.io/collector/component"
)

// Config represents a Solarwinds Extension configuration.
type Config struct {
	// JobEngineServiceEndpoint is the endpoint for the Job Engine service.
	JobEngineServiceEndpoint     string `mapstructure:"job_engine_service_endpoint"`
	TLS_ServerNamme              string `mapstructure:"tls_server_name"`
	TLS_PublicKey                string `mapstructure:"tls_public_key"`
	TLS_PrivateKey               string `mapstructure:"tls_private_key"`
	JobDefinitionsFilePath       string `mapstructure:"job_definitions_file_path"`
	DiscoveryDefinitionsFilePath string `mapstructure:"discovery_definitions_file_path"`
	DefaultJobFrequency          uint   `mapstructure:"default_job_frequency"`
	DefaultJobInitialWait        uint   `mapstructure:"default_job_initial_wait"`

	EndpointPort int `mapstructure:"endpoint_port"`
}

var (
	ErrMissingServiceEndpoint = errors.New("invalid configuration: 'job_engine_service_endpoint' must be set")
)

// NewDefaultConfig creates a new default configuration.
//
// Warning: it doesn't define mandatory `Token` and `DataCenter`
// fields that need to be explicitly provided.
func NewDefaultConfig() component.Config {
	return &Config{
		DefaultJobFrequency:   120,
		DefaultJobInitialWait: 1,
	}
}

// Validate checks the configuration for its validity.
func (cfg *Config) Validate() error {
	if cfg.JobEngineServiceEndpoint == "" {
		return ErrMissingServiceEndpoint
	}

	return nil
}
