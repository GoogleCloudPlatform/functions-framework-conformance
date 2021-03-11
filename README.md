# Functions Framework Conformance

This project contains tooling for conformance and validation of Function
Frameworks to the Functions Framework contract.

## Quickstart

1. In your Functions Framework repo:
   - Create a set of locally-runnable test functions – one function for each signature type).

   Each function write its inputs to the file `function_output.json`.

   - The `http` functions should write the request body.
   - The `cloudevent` function should serialize the CloudEvent parameter to
     JSON and write the resulting string.
   - The `event` (legacy event) function should serialize the data and context
     parameters to JSON in the format:
       `{"data": ...data..., "context": ...context...}`
     and write the resulting string.

1.  Build the test client:

    ```sh
    cd functions-framework-conformance/client && \
      go build
    ```

    This will create a `client` binary in `functions-framework-conformance/client` directory.
    You can use this binary to test the conformance of Function Frameworks.

1.  Invoke the client binary with the following command to run your function server and
    the type of the function.

    - **Go _HTTP_** function Example:

        ```sh
        $HOME/functions-framework-conformance/client/client \
          -cmd "go run ." \
          -type http \
          -builder-source testdata \
          -buildpacks=false
        ```

    - **.NET _CloudEvent_** function example:

        ```sh
        $HOME/functions-framework-conformance/client/client \
          -cmd "dotnet run MyFunction" \
          -type cloudevent \
          -buildpacks=false
        ```

    - **Node.js _legacy event_** function example:

        ```sh
        $HOME/functions-framework-conformance/client/client \
          -cmd "npx @google-cloud/functions-framework --target MyFunction --signature-type=event" \
          -type legacyevent \
          -buildpacks=false
        ```

    If there are validation errors, an error will be logged in the output, causing your conformance test to fail.

## Usage

<nobr>

| configuration flag | type | default | description |
| --- | --- | --- | --- |
| `-cmd` | string | `""` | command to run a Functions Framework server at localhost:8080. Ignored if `-buildpacks=true`. |
| `-type` | string | `"http"` | type of function to validate (must be "http", "cloudevent", or "legacyevent") |
| `-validate-mapping` | boolean | `true` | whether to validate mapping from legacy->cloud events and vice versa (as applicable) |
| `-output-file` | string | `"function_output.json"` | name of file output by function |
| `-buildpacks` | boolean | `true` | whether to use the current release of buildpacks to run the validation. If `true`, `-cmd` is ignored and `--builder-*` flags must be set. |
| `-builder-source` | string | `""` | function source directory to use in building. Required if `-buildpacks=true` |
| `-builder-target` | string | `""` | function target to use in building. Required if `-buildpacks=true` |
| `-builder-runtime` | string | `""` | runtime to use in building. Required if `-buildpacks=true` |
| `-builder-tag` | string | `"latest"` | builder image tag to use in building |
| `-start-delay` | uint | `1` | Seconds to wait before sending HTTP request to command process |

</nobr>

If `-buildpacks` is `true`, you must specify the following flags:

- `-builder-runtime`
- `-builder-source`
- `-builder-target`