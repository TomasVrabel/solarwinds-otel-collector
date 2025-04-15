package internal

import (
	"encoding/json"
	"os"
)

type Variable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Poller struct {
	ID         string     `json:"id"`
	PollerType string     `json:"pollerType"`
	Variables  []Variable `json:"variables"`
}

func parseJobDefinitions(jsonStr string) ([]Poller, error) {
	var pollers []Poller
	err := json.Unmarshal([]byte(jsonStr), &pollers)
	if err != nil {
		return nil, err
	}
	return pollers, nil
}

func ReadJobDefinitions(filePath string) ([]Poller, error) {
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
