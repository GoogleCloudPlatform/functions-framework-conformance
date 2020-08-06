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
    /path/to/functions-framework-conformance/client/client -cmd "go run ." -type http
    ```

    For example, to test a .NET CloudEvent function, you would invoke:

    ```sh
    /path/to/functions-framework-conformance/client/client -cmd "dotnet run MyFunction" -type cloudevent
    ```

    For example, to test a Node.js legacy event function, you would invoke:

    ```sh
    /path/to/functions-framework-conformance/client/client -cmd "npx @google-cloud/functions-framework --target MyFunction --signature-type=event" -type legacyevent
    ```

## Usage

```
Usage of client:

  -cmd string
        command or container image to run a Functions Framework server at localhost:8080
  -type string
        type of function to validate (must be 'http', 'cloudevent', or 'legacyevent' (default "http")
  -validate-mapping
        whether to validate mapping from legacy->cloud events and vice versa (as applicable) (default true)
  -run-in-container
        whether to run a Functions Framework server in a container (default false)
```
