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

package events

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// ValidateLegacyEvent validates that a particular data and context matches the expected contents.
func ValidateLegacyEvent(data string, e event) error {
	got := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &got)
	if err != nil {
		return fmt.Errorf("unmarshalling legacy event: %v", err)
	}

	want := buildLegacyEventWithContext(e)
	if !reflect.DeepEqual(got["data"], want["data"]) {
		return fmt.Errorf("unexpected legacy event data:\ngot %v,\nwant %v", got["data"], want["data"])
	}

	gotCtx, ok := got["context"].(map[string]interface{})
	if !ok {
		// If the 'context' key doesn't exist, that's okay, we'll just check the root-level data.
		gotCtx = got
	}

	if gotCtx["eventId"] != e.meta.EventID {
		return fmt.Errorf("unexpected legacy event ID:\ngot %v,\nwant %v", gotCtx["eventId"], e.meta.EventID)
	}
	if gotCtx["eventType"] != e.meta.EventType {
		return fmt.Errorf("unexpected legacy event type:\ngot %v,\nwant %v", gotCtx["eventType"], e.meta.EventType)
	}
	if gotCtx["timestamp"] != e.meta.Timestamp.Format(time.RFC3339) {
		return fmt.Errorf("unexpected legacy event timestamp:\ngot %v,\nwant %v", gotCtx["timestamp"], e.meta.Timestamp.Format(time.RFC3339))
	}
	if gotCtx["resource"] != e.meta.Resource.RawPath {
		gotResource, ok := gotCtx["resource"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected legacy event resource path:\ngot %v,\nwant %v", gotCtx["resource"], e.meta.Resource.RawPath)
		}
		if gotResource["service"] != e.meta.Resource.Service {
			return fmt.Errorf("unexpected legacy event resource service:\ngot %v,\nwant %v", gotCtx["service"], e.meta.Resource.Service)
		}
		if gotResource["name"] != e.meta.Resource.Name {
			return fmt.Errorf("unexpected legacy event resource name:\ngot %v,\nwant %v", gotCtx["name"], e.meta.Resource.Name)
		}
		if gotResource["type"] != e.meta.Resource.Type {
			return fmt.Errorf("unexpected legacy event resource type:\ngot %v,\nwant %v", gotCtx["type"], e.meta.Resource.Type)
		}
	}

	return nil
}

// ValidateCloudEvent validates that a particular Cloud Event matches the expected one.
func ValidateCloudEvent(data string, e event) error {
	got := &cloudevents.Event{}
	err := json.Unmarshal([]byte(data), got)
	if err != nil {
		return fmt.Errorf("unmarshalling cloud event: %v", err)
	}

	want := BuildCloudEvent(e)
	if want.String() != got.String() {
		return fmt.Errorf("got %s, want %s", got.String(), want.String())
	}
	return nil
}
