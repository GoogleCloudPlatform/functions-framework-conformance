// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/GoogleCloudPlatform/functions-framework-conformance/events"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

var (
	defaultStdoutFile = path.Join(os.TempDir(), "serverlog_stdout.txt")
	defaultStderrFile = path.Join(os.TempDir(), "serverlog_stderr.txt")
)

type functionServer interface {
	Start(stdoutFile, stderrFile, functionOutputFile string) (func(), error)
	OutputFile() ([]byte, error)
}

func send(url string, t events.EventType, data []byte) error {
	switch t {
	case events.LegacyEvent:
		return sendHTTP(url, data)
	case events.CloudEvent:
		ce, err := events.BuildCloudEvent(data)
		if err != nil {
			return fmt.Errorf("building cloudevent: %v", err)
		}
		return sendCE(url, *ce)
	}
	return nil
}

func sendHTTP(url string, data []byte) error {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if err != nil {
			return fmt.Errorf("reading HTTP response body: %v", err)
		}
		return fmt.Errorf("validation failed with exit code %d: %v", resp.StatusCode, string(body))
	}
	return nil
}

func sendCE(url string, e cloudevents.Event) error {
	ctx := cloudevents.ContextWithTarget(context.Background(), url)

	p, err := cloudevents.NewHTTP()
	if err != nil {
		return fmt.Errorf("failed to create protocol: %v", err)
	}

	c, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		return fmt.Errorf("failed to create client, %v", err)
	}

	res := c.Send(ctx, e)
	if !cloudevents.IsACK(res) {
		return fmt.Errorf("failed to send CloudEvent: %v", res)
	}
	return nil
}
