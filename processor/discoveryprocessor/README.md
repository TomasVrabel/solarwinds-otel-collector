# Discovery Processor

The Discovery Processor filters logs that contain the `discovery=true` attribute and forwards them to the SolarWinds Job Engine Extension for processing. After forwarding discovery logs to the extension, these logs are removed from the pipeline to prevent them from being sent to downstream consumers.

## Configuration

The processor accepts the following configuration parameters:

- `job_engine_extension_name` (string, default: "swjobengine"): The name of the job engine extension to send discovery data to.

## Example Configuration

```yaml
processors:
  discovery:
    job_engine_extension_name: "swjobengine"
```

## How it works

1. The processor examines all incoming log records
2. For each log record, it checks if the `discovery` attribute is present and set to `true`
3. If a discovery log is found:
   - The entire log record (with all attributes) is sent to the configured job engine extension
   - The log record is removed from the pipeline
4. Non-discovery logs continue through the pipeline normally

## Discovery Log Format

Discovery logs should have the `discovery` attribute set to `true` (either as a boolean or string value). For example:

```json
{
  "timestamp": "2025-01-23T10:00:00Z",
  "body": "Discovery data content",
  "attributes": {
    "discovery": true,
    "source": "network_scanner",
    "target": "10.0.0.1"
  }
}
```

## Dependencies

This processor requires the SolarWinds Job Engine Extension to be configured and running in order to process discovery logs. If the extension is not found, discovery logs will still be filtered out but will not be processed.

## Supported pipeline types

- logs

## Stability

This processor is in development and may change in future releases.