module github.com/solarwinds/solarwinds-otel-collector/processor/discoveryprocessor

go 1.22.0

require (
	go.opentelemetry.io/collector/component v0.116.0
	go.opentelemetry.io/collector/consumer v0.116.0
	go.opentelemetry.io/collector/pdata v1.22.0
	go.opentelemetry.io/collector/processor v0.116.0
	go.opentelemetry.io/collector/processor/processorhelper v0.116.0
	go.uber.org/zap v1.27.0
)

replace github.com/solarwinds/solarwinds-otel-collector/extension/swjobengineextension => ../../extension/swjobengineextension