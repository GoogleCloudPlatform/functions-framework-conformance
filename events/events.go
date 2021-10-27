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

//go:generate go run generate/events_generate.go

// Package events contains the validation logic for different types of events.
package events

import (
	"encoding/json"
	"fmt"
	"sort"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// EventType is the type of event to validate.
type EventType int

const (
	// LegacyEvent represents a legacy event type.
	LegacyEvent EventType = iota
	// CloudEvent represents a CloudEvent type.
	CloudEvent
)

func (e EventType) String() string {
	switch e {
	case LegacyEvent:
		return "legacy event"
	case CloudEvent:
		return "cloud event"
	}
	return ""
}

// EventNames returns a list of event names to use as inputs for a particular event type.
func EventNames(t EventType) ([]string, error) {
	eventNames := []string{}
	for name, data := range Events {
		switch t {
		case LegacyEvent:
			if data.Input.LegacyEvent != nil {
				eventNames = append(eventNames, name)
			}
		case CloudEvent:
			if data.Input.CloudEvent != nil {
				eventNames = append(eventNames, name)
			}
		}
	}

	// Sort the event names for deterministic output.
	sort.Strings(eventNames)

	return eventNames, nil
}

// InputData returns the contents of the input event for a particular event name and type.
func InputData(name string, t EventType) []byte {
	switch t {
	case LegacyEvent:
		return Events[name].Input.LegacyEvent
	case CloudEvent:
		return Events[name].Input.CloudEvent
	}
	return nil
}

// OutputData returns the contents of the output event for a particular event name and type.
func OutputData(name string, t EventType, isConversion bool) []byte {
	switch t {
	case LegacyEvent:
		if isConversion && Events[name].ConvertedOutput.LegacyEvent != nil {
			return Events[name].ConvertedOutput.LegacyEvent
		}
		return Events[name].Output.LegacyEvent
	case CloudEvent:
		if isConversion && Events[name].ConvertedOutput.CloudEvent != nil {
			return Events[name].ConvertedOutput.CloudEvent
		}
		return Events[name].Output.CloudEvent
	}
	return nil
}

// BuildCloudEvent creates a CloudEvent from a byte slice.
func BuildCloudEvent(data []byte) (*cloudevents.Event, error) {
	ce := &cloudevents.Event{}
	err := json.Unmarshal(data, ce)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling cloud event: %v", err)
	}
	return ce, nil
}
