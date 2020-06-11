1. Create a set of locally-runnable test functions (one for each type). Each
   function write its inputs to the file `function_output.json`.

   - The HTTP function should write the request body.
   - The CloudEvent function should serialize the CloudEvent parameter to JSON
     and write the resulting string.
   - The legacy event function should serialize the data and context parameters
     to JSON in the format of `{"data": ...data..., "context": ...context...}`
     and write the resulting string.

1. Build the test client: `cd framework_validation/client && go build`. This
   will create a `client` binary in `framework_validation`.

1. Invoke the client binary with the command to run your function server and the
   type of the function.

   For example, to test a Go HTTP function, you would invoke:

   ```sh
   /path/to/framework_validation/client/client -cmd "go run ." -type http
   ```

   For example, to test a .NET CloudEvent function, you would invoke:

   ```sh
   /path/to/framework_validation/client/client -cmd "dotnet run MyFunction" -type cloudevent
   ```
