package internal

import (
	"encoding/json"
	"os"
)

type Variable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PollerJob struct {
	ID          string     `json:"id"`
	PollerType  string     `json:"pollerType"`
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
