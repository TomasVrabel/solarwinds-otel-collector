module github.com/solarwinds/solarwinds-otel-collector/receiver/swojobenginereceiver

go 1.23.4

require (
	github.com/solarwinds/solarwinds-otel-collector/pkg/job-engine-events v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/collector/component v0.119.0
	go.opentelemetry.io/collector/component/componentstatus v0.119.0
	go.opentelemetry.io/collector/consumer v1.25.0
	go.opentelemetry.io/collector/pdata v1.25.0
	go.opentelemetry.io/collector/receiver v0.119.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.70.0

)

require (
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.119.0 // indirect
	go.opentelemetry.io/collector/pipeline v0.119.0 // indirect
	go.opentelemetry.io/otel v1.34.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.34.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.36.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)

replace github.com/solarwinds/solarwinds-otel-collector/pkg/job-engine-events => ../../pkg/job-engine-events
