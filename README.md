# Functions Framework Conformance

This project contains tooling for conformance and validation of Function
Frameworks to the Functions Framework contract.

## Quickstart

1. In your Functions Framework repo:
   - Create a set of locally-runnable test functions – one function for each signature type.

   Each function write its inputs to the file `function_output.json`.

   - The `http` functions should write the request body.
   - The `cloudevent` function should serialize the CloudEvent parameter to
     JSON and write the resulting string.
   - The `event` (legacy event) function should serialize the data and context
     parameters to JSON in the format:
       `{"data": ...data..., "context": ...context...}`
     and write the resulting string.
    - The `typed` function should accept a JSON request object and should echo
      the request object back in the "payload" field of the response. I.e. if
      the request is `{"a":"b"}` the response should be `{"payload":{"a":
      "b"}}`.

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
          -cmd="go run ." \
          -type=http \
          -builder-source=testdata \
          -buildpacks=false
        ```

    - **.NET _CloudEvent_** function example:

        ```sh
        $HOME/functions-framework-conformance/client/client \
          -cmd="dotnet run MyFunction" \
          -type=cloudevent \
          -buildpacks=false
        ```

    - **Node.js _legacy event_** function example:

        ```sh
        $HOME/functions-framework-conformance/client/client \
          -cmd="npx @google-cloud/functions-framework --target MyFunction --signature-type=event" \
          -type=legacyevent \
          -buildpacks=false
        ```
    - **Ruby __typed http__** function Example:

        ```sh
        $HOME/functions-framework-conformance/client/client \
          -cmd="bundle exec functions-framework-ruby --source test/conformance/app.rb --target typed_func --signature-type http" \
          --type=http \
          --declarative-type=typed \
          -buildpacks=false
        ```

    If there are validation errors, an error will be logged in the output, causing your conformance test to fail.

## Usage

<nobr>

| Configuration flag | Type | Default | Description |
| --- | --- | --- | --- |
| `-cmd` | string | `"''"` | A string with the command to run a Functions Framework server at `localhost:8080`. Must be wrapped in quotes. Ignored if `-buildpacks=true`. |
| `-type` | string | `"http"` | The function signature to use (must be `"http"`, `"cloudevent"`, or `"legacyevent"`). |
| `-declarative-type` | string | `""` | The declarative signature type of the function (must be 'http', 'cloudevent', 'legacyevent', or 'typed'), default matches -type |
| `-validate-mapping` | boolean | `true` | Whether to validate mapping from legacy->cloud events and vice versa (as applicable). |
| `-output-file` | string | `"function_output.json"` | Name of file output by function. |
| `-buildpacks` | boolean | `true` | Whether to use the current release of buildpacks to run the validation. If `true`, `-cmd` is ignored and `--builder-*` flags must be set. |
| `-builder-source` | string | `""` | Function source directory to use in building. Required if `-buildpacks=true`. |
| `-builder-target` | string | `""` | Function target to use in building. Required if `-buildpacks=true`. |
| `-builder-runtime` | string | `""` | Runtime to use in building. Required if `-buildpacks=true`. |
| `-builder-runtime-version` | string | `""` | Runtime version used while building. Buildpack will use the latest version if flag is not specified. |
| `-builder-tag` | string | `"latest"` | Builder image tag to use in building. Ignored if `-builder-url` is specified. |
| `-builder-url` | string | `""` | Builder image url to use in building including tag. Client defaults to `gcr.io/gae-runtimes/buildpacks/<language>/builder:<builder-tag>` if none is specified. |
| `-start-delay` | uint | `1` | Seconds to wait before sending HTTP request to command process. |
| `-envs` | string | `""` | A comma separated string of additional runtime environment variables. |

</nobr>

If `-buildpacks` is `true`, you must specify the following flags:

- `-builder-runtime`
- `-builder-source`
- `-builder-target`

