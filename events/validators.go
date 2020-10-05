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

// ValidationInfo contains information about a particular validation step, including a reason why
// the validation for this event and type was skipped or the relevant error.
type ValidationInfo struct {
	Name          string
	Errs          []error
	SkippedReason string
}

// PrintValidationInfos takes a list of ValidationInfos and collapses them into a single error and
// a single log line recording which events were validation, which skipped, and why.
func PrintValidationInfos(vis []*ValidationInfo) (string, error) {
	errStr := "Validation errors:"
	logStr := "Events tried:"

	errsOccurred := false
	for _, vi := range vis {
		// Collect errors into one string.
		if vi.Errs != nil {
			errsOccurred = true
			viErrStr := fmt.Sprintf("%s:", vi.Name)
			for _, err := range vi.Errs {
				viErrStr = fmt.Sprintf("%s\n\t\t- %v", viErrStr, err)
			}
			errStr = fmt.Sprintf("%s\n\t- %s", errStr, viErrStr)
			logStr = fmt.Sprintf("%s\n\t- %s (FAILED)", logStr, vi.Name)
			continue
		}

		// Collect events run and skipped into one string.
		if vi.SkippedReason != "" {
			logStr = fmt.Sprintf("%s\n\t- %s (SKIPPED: %s)", logStr, vi.Name, vi.SkippedReason)
		} else {
			logStr = fmt.Sprintf("%s\n\t- %s (PASSED)", logStr, vi.Name)
		}
	}

	if errsOccurred {
		return logStr, fmt.Errorf(errStr)
	}
	return logStr, nil
}

// ValidateEvent validates that a particular function output matches the expected contents.
func ValidateEvent(name string, t EventType, got []byte) *ValidationInfo {
	want := OutputData(name, t)
	if want == nil {
		// Include the possibilities in the error.
		return &ValidationInfo{
			Name:          name,
			SkippedReason: fmt.Sprintf("no expected output value of type %s", t),
		}
	}

	switch t {
	case LegacyEvent:
		return validateLegacyEvent(name, got, want)
	case CloudEvent:
		return validateCloudEvent(name, got, want)
	}

	// Should be unreachable.
	return nil
}

func validateLegacyEvent(name string, gotBytes, wantBytes []byte) *ValidationInfo {
	vi := &ValidationInfo{
		Name: name,
	}
	got := make(map[string]interface{})
	err := json.Unmarshal(gotBytes, &got)
	if err != nil {
		vi.Errs = append(vi.Errs, fmt.Errorf("unmarshalling function-received version of legacy event %q: %v", name, err))
	}

	want := make(map[string]interface{})
	err = json.Unmarshal(wantBytes, &want)
	if err != nil {
		vi.Errs = append(vi.Errs, fmt.Errorf("unmarshalling expected contents of legacy event %q: %v", name, err))
	}

	// If there were issues extracting the data, bail early.
	if vi.Errs != nil {
		return vi
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
			vi.Errs = append(vi.Errs, fmt.Errorf("unexpected %q in event %q:\ngot %+v,\nwant %+v", field.name, name, field.gotValue, field.wantValue))
		}
	}

	return vi
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

func validateCloudEvent(name string, gotBytes, wantBytes []byte) *ValidationInfo {
	vi := &ValidationInfo{
		Name: name,
	}

	got := &cloudevents.Event{}
	err := json.Unmarshal(gotBytes, got)
	if err != nil {
		vi.Errs = append(vi.Errs, fmt.Errorf("unmarshalling function-received version of cloud event %q: %v", name, err))
	}

	want := &cloudevents.Event{}
	err = json.Unmarshal(wantBytes, want)
	if err != nil {
		vi.Errs = append(vi.Errs, fmt.Errorf("unmarshalling expected contents of cloud event %q: %v", name, err))
	}

	// If there were issues extracting the data, bail early.
	if vi.Errs != nil {
		return vi
	}

	fields := []struct {
		name      string
		gotValue  interface{}
		wantValue interface{}
	}{
		{
			name:      "ID",
			gotValue:  got.ID(),
			wantValue: want.ID(),
		},
		{
			name:      "source",
			gotValue:  got.Source(),
			wantValue: want.Source(),
		},
		{
			name:      "type",
			gotValue:  got.Type(),
			wantValue: want.Type(),
		},
		{
			name:      "datacontenttype",
			gotValue:  got.DataContentType(),
			wantValue: want.DataContentType(),
		},
		{
			name:      "data",
			gotValue:  unmarshalMap(got.Data(), vi),
			wantValue: unmarshalMap(want.Data(), vi),
		},
	}
	for _, field := range fields {
		if !cmp.Equal(field.gotValue, field.wantValue) {
			vi.Errs = append(vi.Errs, fmt.Errorf("unexpected %q field in %q: got %v, want %v", field.name, name, field.gotValue, field.wantValue))
		}
	}

	return vi
}

func unmarshalMap(data []byte, vi *ValidationInfo) (dataMap map[string]interface{}) {
	if err := json.Unmarshal(data, &dataMap); err != nil {
		vi.Errs = append(vi.Errs, fmt.Errorf("could not parse CloudEvent data as map: %v", err))
	}
	return
}
