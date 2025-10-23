package internal

import (
	"encoding/json"
	"os"
)

type Variable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// JobType represents the type of job being executed
type JobType string

const (
	JobTypeOther     JobType = "other"
	JobTypePoll      JobType = "poll"
	JobTypeInventory JobType = "inventory"
	JobTypeDiscovery JobType = "discovery"
)

// JobContext represents the main job execution context
type JobContext struct {
	Type             JobType          `json:"type" validate:"required"`
	Entity           EntityContext    `json:"entity,omitempty"`
	DiscoveryContext DiscoveryContext `json:"discoveryContext,omitempty"`
}

// EntityContext represents the entity being processed
type EntityContext struct {
	EntityType       string            `json:"entityType" validate:"required"`
	EntityId         map[string]string `json:"entityId" validate:"required"`
	EntityAttributes map[string]string `json:"entityAttributes" validate:"required"`
	Relations        []Relation        `json:"relations,omitempty"`
}

// Relation represents a relationship between entities
type Relation struct {
	RelationType string        `json:"relationType" validate:"required"`
	Entity       EntityContext `json:"entity" validate:"required"`
}

type DiscoveryContext struct {
	DiscoveryId    string `json:"discoveryId" validate:"required"`
	DiscoveryRunId string `json:"discoveryRunId" validate:"required"`
}

type PollerJob struct {
	ID          string     `json:"id"`
	PollerType  string     `json:"pollerType"`
	State       JobContext `json:"state"`
	Frequency   uint       `json:"frequency"`
	InitialWait uint       `json:"initialWait"`
	Variables   []Variable `json:"variables"`
}

func parseJobDefinitions(jsonStr string) ([]PollerJob, error) {
	var pollers []PollerJob
	err := json.Unmarshal([]byte(jsonStr), &pollers)
	if err != nil {
		return nil, err
	}
	return pollers, nil
}

func ReadJobDefinitions(filePath string) ([]PollerJob, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	pollers, err := parseJobDefinitions(string(data))
	if err != nil {
		return nil, err
	}

	return pollers, nil
}

func JobStateToJSON(jc JobContext) (string, error) {
	data, err := json.MarshalIndent(jc, "", "  ")
	if err != nil {
		return nil, err
	}
	return string(data), nil
}
