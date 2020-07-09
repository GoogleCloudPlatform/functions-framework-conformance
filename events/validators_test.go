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
	"testing"
)

func TestValidateLegacyEvent(t *testing.T) {
	testName := "firebase-auth"
	data := OutputData(testName, LegacyEvent)
	if data == nil {
		t.Fatalf("no legacy event data")
	}
	// Validate the event output against itself.
	if err := ValidateEvent(testName, LegacyEvent, data); err != nil {
		t.Errorf("validating legacy event: %v", err)
	}
}

func TestValidateCloudEvent(t *testing.T) {
	testName := "firebase-auth"
	data := OutputData(testName, CloudEvent)
	if data == nil {
		t.Fatalf("no cloudevent data")
	}
	// Validate the event output against itself.
	if err := ValidateEvent(testName, CloudEvent, data); err != nil {
		t.Errorf("validating cloudevent: %v", err)
	}
}
