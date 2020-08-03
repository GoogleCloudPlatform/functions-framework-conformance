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
	"github.com/google/go-cmp/cmp"
)

// ValidateEvent validates that a particular function output matches the expected contents.
func ValidateEvent(name string, t EventType, got []byte) error {
	want := OutputData(name, t)
	if want == nil {
		// List available event types for debugging.
		available := []string{}
		for name, data := range Events {
			switch t {
			case LegacyEvent:
				if data.Output.LegacyEvent != nil {
					available = append(available, name)
				}
			case CloudEvent:
				if data.Output.CloudEvent != nil {
					available = append(available, name)
				}
			}
		}
		// Include the possibilities in the error.
		return fmt.Errorf("no expected output value found for %q. Available event types: %v", name, available)
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
		return fmt.Errorf("unmarshalling function-received version of legacy event %q: %v", name, err)
	}

	want := make(map[string]interface{})
	err = json.Unmarshal(wantBytes, &want)
	if err != nil {
		return fmt.Errorf("unmarshalling expected contents of legacy event %q: %v", name, err)
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
		return fmt.Errorf("unmarshalling function-received version of cloud event %q: %v", name, err)
	}

	want := &cloudevents.Event{}
	err = json.Unmarshal(wantBytes, want)
	if err != nil {
		return fmt.Errorf("unmarshalling expected contents of cloud event %q: %v", name, err)
	}

	gotContext := got.Context.AsV1()
	wantContext := want.Context.AsV1()

	fields := []struct {
		name      string
		gotValue  interface{}
		wantValue interface{}
	}{
		{
			name:      "ID",
			gotValue:  gotContext.ID,
			wantValue: wantContext.ID,
		},
		{
			name:      "source",
			gotValue:  gotContext.Source,
			wantValue: wantContext.Source,
		},
		{
			name:      "type",
			gotValue:  gotContext.Type,
			wantValue: wantContext.Type,
		},
		{
			name:      "datacontenttype",
			gotValue:  *gotContext.DataContentType,
			wantValue: *wantContext.DataContentType,
		},
		{
			name:      "data",
			gotValue:  got.DataEncoded,
			wantValue: want.DataEncoded,
		},
	}
	for _, field := range fields {
		if !cmp.Equal(field.gotValue, field.wantValue) {
			return fmt.Errorf("unexpected %q field in %q: got %v, want %v", field.name, name, field.gotValue, field.wantValue)
		}
	}

	return nil
}
