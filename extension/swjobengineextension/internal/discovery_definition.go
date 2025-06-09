package internal

import (
	"encoding/json"
	"os"
)

type DiscoveryJob struct {
	Name             string             `json:"name"`
	IpAddresses      []string           `json:"ipAddresses"`
	Subnets          []Subnet           `json:"subnets"`
	IpRanges         []IpRange          `json:"ipRanges"`
	CredentialSnmpV2 []CredentialSnmpV2 `json:"credentials"`
}

type CredentialSnmpV2 struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Community string `json:"community"`
}

type Subnet struct {
	IpAddress string `json:"ipAddress"`
	Mask      string `json:"mask"`
}

type IpRange struct {
	FirstIpAddress string `json:"firstIpAddress"`
	LastIpAddress  string `json:"lastIpAddress"`
}

func parseDiscoveryDefinitions(jsonStr string) ([]DiscoveryJob, error) {
	var jobs []DiscoveryJob
	err := json.Unmarshal([]byte(jsonStr), &jobs)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func ReadDiscoveryDefinitions(filePath string) ([]DiscoveryJob, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	jobs, err := parseDiscoveryDefinitions(string(data))
	if err != nil {
		return nil, err
	}

	return jobs, nil
}
