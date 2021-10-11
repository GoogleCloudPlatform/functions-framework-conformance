# Functions Framework conformance test action

This action runs the Functions Framework conformance tests with the specified
parameters.

Requires Go 1.16+ to be installed prior to running (e.g. actions/setup-go).

## Inputs

### `version`

The version of conformance tests to run. Default `latest`.

### `cmd`

The command to be run.

### `functionType`

The type of function to validate. Can be `http`, `legacyevent`, or `cloudevent`.
Default `http`.

### `validateMapping`

Whether or not to validate legacy->cloudevent mapping and vice versa. Default
`true`.

## Example usage

```yaml
uses: actions/setup-go@v1
uses: GoogleCloudPlatform/functions-framework-conformance/actions@v1
with:
  version: 'v1.0.0'
  functionType: 'http'
  validateMapping: false
  source: 'testdata/testfunc.go'
  target: 'HTTP'
  runtime: 'go113'
```
