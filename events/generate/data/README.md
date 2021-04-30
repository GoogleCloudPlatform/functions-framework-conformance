# Test data

The test data in this directory consists of events in various representations.
Each type of event (e.g. "firebase-auth") is a particular event that can be
represented in the Google Cloud Functions background function format or in a
format compliant with the CloudEvent [spec](https://cloudevents.io/).

## Adding test cases

To add a new test case or to adjust an existing test case, add or change the
appropriate files in this directory.

A new test case should have an expected legacy event input and output and an
expected CloudEvent input and output between which each Functions Frameworks
should be able to convert. These files should be named:

-   `newevent-legacy-input.json`
-   `newevent-legacy-output.json`
-   `newevent-cloudevent-input.json`
-   `newevent-cloudevent-output.json`

The conformance test suite will interpret one set of such files as one
validation test case.

Generally, try to stay consistent across tests for the following fields (arbitrary values, but consistent ones):

- EventID: "aaaaaa-1111-bbbb-2222-cccccccccccc"
- Timestamp: "2020-09-29T11:32:00.000Z"
- Project ID: my-project-id

The `output-converted.json` suffix can be used to override the expected output
when converting between event types. For example, consider the following files:

-   `newevent-legacy-output.json`
-   `newevent-legacy-output-converted.json`

The `newevent-legacy-output.json` file will be used to validate legacy events,
and `newevent-legacy-output-converted.json` will be used to validate converting
between cloud events and legacy events.

Once you have the input and output data, generate the test cases to embed them
in the binary. Run the following:

`go generate ./...`

This will update the `events/event_data.go` file with your changes. Include
changes to this file in your commit.
