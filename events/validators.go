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

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// ValidateEvent validates that a particular function output matches the expected contents.
func ValidateEvent(name string, t EventType, got []byte) error {
	want := OutputData(name, t)
	if want == nil {
		return fmt.Errorf("no output found for %q", name)
	}

	switch t {
	case LegacyEvent:
		return validateLegacyEvent(name, got, want)
	case CloudEvent:
		return validateCloudEvent(name, got, want)
	}
	return nil
}

func validateLegacyEvent(name string, gotBytes, wantBytes []byte) error {
	got := make(map[string]interface{})
	err := json.Unmarshal(gotBytes, &got)
	if err != nil {
		return fmt.Errorf("unmarshalling legacy event %q: %v", name, err)
	}

	want := make(map[string]interface{})
	err = json.Unmarshal(wantBytes, &want)
	if err != nil {
		return fmt.Errorf("unmarshalling expected legacy event %q: %v", name, err)
	}

	if !reflect.DeepEqual(got["data"], want["data"]) {
		return fmt.Errorf("unexpected data in event %q:\ngot %v,\nwant %v", name, got["data"], want["data"])
	}

	gotContext := got["context"].(map[string]interface{})
	wantContext := want["context"].(map[string]interface{})

	// For some fields in the context, they can be written in more than one way. Check all.
	type eventFields struct {
		name      string
		gotValue  interface{}
		wantValue interface{}
	}
	fields := []eventFields{
		{
			name:      "ID",
			gotValue:  getMaybeSnakeCaseField(gotContext, "eventId"),
			wantValue: wantContext["eventId"],
		},
		{
			name:      "type",
			gotValue:  getMaybeSnakeCaseField(gotContext, "eventType"),
			wantValue: wantContext["eventType"],
		},
		{
			name:      "timestamp",
			gotValue:  gotContext["timestamp"],
			wantValue: wantContext["timestamp"],
		},
		{
			name:      "resource",
			gotValue:  gotContext["resource"],
			wantValue: wantContext["resource"],
		},
		{
			name:      "data",
			gotValue:  got["data"],
			wantValue: want["data"],
		},
	}

	for _, field := range fields {
		if !reflect.DeepEqual(field.gotValue, field.wantValue) {
			return fmt.Errorf("unexpected %q in event %q:\ngot %+v,\nwant %+v", field.name, name, field.gotValue, field.wantValue)
		}
	}

	return nil
}

// Some fields can present with either a CamelCase or a snake_case key. Both are acceptable.
func getMaybeSnakeCaseField(gotContext map[string]interface{}, field string) interface{} {
	if gotVal, ok := gotContext[field]; ok {
		return gotVal
	}

	var lowerField string
	if field == "eventId" {
		lowerField = "event_id"
	}
	if field == "eventType" {
		lowerField = "event_type"
	}

	if gotVal, ok := gotContext[lowerField]; lowerField != "" && ok {
		return gotVal
	}

	return nil
}

func validateCloudEvent(name string, gotBytes, wantBytes []byte) error {
	got := &cloudevents.Event{}
	err := json.Unmarshal(gotBytes, got)
	if err != nil {
		return fmt.Errorf("unmarshalling cloud event %q: %v", name, err)
	}

	want := &cloudevents.Event{}
	err = json.Unmarshal(wantBytes, want)
	if err != nil {
		return fmt.Errorf("unmarshalling expected cloud event %q: %v", name, err)
	}

	if want.String() != got.String() {
		return fmt.Errorf("unexpected event %q: got %s, want %s", name, got.String(), want.String())
	}
	return nil
}
