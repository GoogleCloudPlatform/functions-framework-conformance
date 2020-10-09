# Functions Framework Conformance

This project contains tooling for conformance and validation of Function
Frameworks to the Functions Framework contract.

## Quickstart

1.  Create a set of locally-runnable test functions (one for each type). Each
    function write its inputs to the file `function_output.json`.

    -   The HTTP function should write the request body.
    -   The CloudEvent function should serialize the CloudEvent parameter to
        JSON and write the resulting string.
    -   The legacy event function should serialize the data and context
        parameters to JSON in the format of `{"data": ...data..., "context":
        ...context...}` and write the resulting string.

1.  Build the test client: `cd functions-framework-conformance/client && go
    build`. This will create a `client` binary in
    `functions-framework-conformance/client`.

1.  Invoke the client binary with the command to run your function server and
    the type of the function.

    For example, to test a Go HTTP function, you would invoke:

    ```sh
    /path/to/functions-framework-conformance/client/client -cmd "go run ." -type http -buildpacks false
    ```

    For example, to test a .NET CloudEvent function, you would invoke:

    ```sh
    /path/to/functions-framework-conformance/client/client -cmd "dotnet run MyFunction" -type cloudevent -buildpacks false
    ```

    For example, to test a Node.js legacy event function, you would invoke:

    ```sh
    /path/to/functions-framework-conformance/client/client -cmd "npx @google-cloud/functions-framework --target MyFunction --signature-type=event" -type legacyevent -buildpacks false
    ```

## Usage

```
Usage of client:

  -builder-runtime string
        runtime to use in building. Required if -buildpacks=true
  -builder-source string
        function source directory to use in building. Required if -buildpacks=true
  -builder-tag string
        builder image tag to use in building (default "latest")
  -builder-target string
        function target to use in building. Required if -buildpacks=true
  -buildpacks
        whether to use the current release of buildpacks to run the validation. If true, -cmd is ignored and --builder-* flags must be set. (default true)
  -cmd string
        command to run a Functions Framework server at localhost:8080. Ignored if -buildpacks=true.
  -output-file string
        name of file output by function (default "function_output.json")
  -type string
        type of function to validate (must be 'http', 'cloudevent', or 'legacyevent' (default "http")
  -validate-mapping
        whether to validate mapping from legacy->cloud events and vice versa (as applicable) (default true)
  -start-delay int
        number of seconds to wait after command process startup before sending HTTP request
```
