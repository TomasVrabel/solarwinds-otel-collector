package swojobenginereceiver

import (
	"context"

	"go.opentelemetry.io/collector/receiver/receiverhelper"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

// let's take inspiration from https://github.com/open-telemetry/opentelemetry-collector/tree/main/receiver/otlpreceiver
// and also https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/statsdreceiver

var (
	typeStr = component.MustNewType("swojobengine")
)

const (
	defaultPort = 8555
)

// Creates factory capable of creating swohostmetrics receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithMetrics(createObsMetricsReceiver, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		EndpointPort: defaultPort,
	}
}

func createObsMetricsReceiver(
	ctx context.Context,
	settings receiver.Settings,
	config component.Config,
	metrics consumer.Metrics,
) (receiver.Metrics, error) {
	receiverConfig := config.(*Config)

	r := &jobengineGRPCReceiver{
		conf:         receiverConfig,
		nextConsumer: metrics,
		settings:     settings,
	}

	var err error
	r.obsrepGRPC, err = receiverhelper.NewObsReport(receiverhelper.ObsReportSettings{
		ReceiverID:             settings.ID,
		Transport:              "grpc",
		ReceiverCreateSettings: settings,
	})
	if err != nil {
		return nil, err
	}

	return r, nil
}
