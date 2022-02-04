# Functions Framework conformance test action

This action runs the Functions Framework conformance tests with the specified
parameters.

Requires Go to be installed prior to running (e.g. actions/setup-go).

## Inputs

### `version`

The version of conformance tests to run. Default to latest release if unspecified.

### `cmd`

The command to be run.

### `functionType`

The type of function to validate. Can be `http`, `legacyevent`, or `cloudevent`.
Default `http`.

### `validateMapping`

Whether or not to validate legacy->cloudevent mapping and vice versa. Default
`true`.

### `validateConcurrency`

Whether or not to validate concurrent requests are handled. Default `false`.

## Example usage

```yaml
uses: actions/setup-go@v1
uses: GoogleCloudPlatform/functions-framework-conformance/actions@v1
with:
  version: 'v1.0.0'
  functionType: 'http'
  validateMapping: false
  validateConcurrency: true
  source: 'testdata/testfunc.go'
  target: 'HTTP'
  runtime: 'go113'
```
