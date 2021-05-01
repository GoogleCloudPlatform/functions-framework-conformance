# Functions Framework conformance test action

This action runs the Functions Framework conformance tests with the specified
parameters.

Requires Go to be installed prior to running (e.g. actions/setup-go).

## Inputs

<!--BEGIN GENERATED DOCS-->

### `outputFile`

The output file from the function. Default value: `'function_output.json'`.


### `functionType`

The type of function to validate. Can be `http`, `legacyevent`, or `cloudevent`. Default value: `'http'`.


### `validateMapping`

Whether to validate mapping from legacy->cloud event and vice versa. Default value: `true`.


### `source`

Function source code. Default value: `''`.


### `target`

Function target. Default value: `''`.


### `runtime`

Function runtime (e.g. nodejs10, go113). Default value: `''`.


### `tag`

GCR tag to use for builder image. Default value: `'latest'`.


### `useBuildpacks`

GCR tag to use for builder image. Default value: `true`.


### `cmd`

Command to run a Functions Framework server at localhost:8080. Ignored if -buildpacks=true.. Default value: `''`.


### `startDelay`

GCR tag to use for builder image. Default value: `1`.


### `workingDirectory`

The subdirectory in which the conformance tests should run. Default value: `''`.

<!--END GENERATED DOCS-->

## Example usage

```yaml
uses: actions/setup-go@v1
uses: GoogleCloudPlatform/functions-framework-conformance/actions@v1
with:
  functionType: 'http'
  validateMapping: false
  source: 'testdata/testfunc.go'
  target: 'HTTP'
  runtime: 'go113'
```
