package job_engine_events

import (
	"encoding/json"
	"fmt"
)

const (
	JOB_STATE_NODES_CATEGORY_ATTRIBUTE       = "sw.collector.Nodes.Category"
	JOB_STATE_NODES_URI_ATTRIBUTE            = "sw.collector.Nodes.Uri"
	JOB_STATE_NODES_IP_ADDRESS_ATTRIBUTE     = "sw.collector.Nodes.IPAddress"
	JOB_STATE_NODES_POLLING_METHOD_ATTRIBUTE = "sw.collector.Nodes.PollingMethod"
)

// PollerJobState represents a map of string keys and string values
type PollerJobState map[string]string

func (m PollerJobState) serializeJobState() ([]byte, error) {
	return json.Marshal(m)
}

// SerializeJobStateToString converts the PollerJobState to a JSON string
func (m PollerJobState) SerializeJobStateToString() (string, error) {
	bytes, err := m.serializeJobState()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// DeserializeJobState converts JSON bytes to PollerJobState
func DeserializeJobState(data []byte) (PollerJobState, error) {
	var result PollerJobState
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JobState: %w", err)
	}
	return result, nil
}

// DeserializeJobStateFromString converts a JSON string to PollerJobState
func DeserializeJobStateFromString(data string) (PollerJobState, error) {
	return DeserializeJobState([]byte(data))
}

// NewJobState creates a new empty PollerJobState
func NewJobState() PollerJobState {
	return make(PollerJobState)
}
