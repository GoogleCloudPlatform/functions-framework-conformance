# Functions framework conformance test javascript action

This action runs the Functions Framework conformance tests with the specified
parameters.

Requires Go to be installed prior to running (e.g. actions/setup-go).

## Inputs

### `cmd`

The command to be run.

### `functionType`

The type of function to validate. Can be `http`, `legacyevent`, or `cloudevent`.
Default `http`.

### `validateMapping`

Whether or not to validate legacy->cloudevent mapping and vice versa. Default
`true`.

## Example usage

```
uses: actions/setup-go@v1
uses: GoogleCloudPlatform/functions-framework-conformance/actions@v1
with:
  cmd: '<cmd to run function>'
  functionType: 'http'
  validateMapping: false
```
