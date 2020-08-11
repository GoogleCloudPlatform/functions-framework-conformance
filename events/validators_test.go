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
	"fmt"
	"testing"
)

func TestValidateLegacyEvent(t *testing.T) {
	testName := "firebase-auth"
	data := OutputData(testName, LegacyEvent)
	if data == nil {
		t.Fatalf("no legacy event data")
	}
	// Validate the event output against itself.
	if vi := ValidateEvent(testName, LegacyEvent, data); vi.Errs != nil {
		t.Errorf("validating legacy event: %v", vi.Errs)
	}
}

func TestValidateCloudEvent(t *testing.T) {
	testName := "firebase-auth"
	data := OutputData(testName, CloudEvent)
	if data == nil {
		t.Fatalf("no cloudevent data")
	}
	// Validate the event output against itself.
	if vi := ValidateEvent(testName, CloudEvent, data); vi.Errs != nil {
		t.Errorf("validating cloudevent: %v", vi.Errs)
	}
}

func TestPrintValidationInfos(t *testing.T) {
	vis := []*ValidationInfo{
		&ValidationInfo{
			Name: "with error",
			Errs: []error{
				fmt.Errorf("first error"),
			},
		},
		&ValidationInfo{
			Name: "with multiple errors",
			Errs: []error{
				fmt.Errorf("first error"),
				fmt.Errorf("second error"),
			},
		},
		&ValidationInfo{
			Name: "passed",
		},
		&ValidationInfo{
			Name:          "skipped",
			SkippedReason: "skipping",
		},
	}

	wantLog := `Events tried:
	- with error (FAILED)
	- with multiple errors (FAILED)
	- passed (PASSED)
	- skipped (SKIPPED: skipping)`

	wantErr := fmt.Errorf(`Validation errors:
	- with error:
		- first error
	- with multiple errors:
		- first error
		- second error`)

	gotLog, gotErr := PrintValidationInfos(vis)
	if gotLog != wantLog {
		t.Errorf("PrintValidationInfos log: got %s, want %s", gotLog, wantLog)
	}
	if gotErr.Error() != wantErr.Error() {
		t.Errorf("PrintValidationInfos error: got %v, want %v", gotErr, wantErr)
	}

	passedVIs := []*ValidationInfo{
		&ValidationInfo{
			Name: "passed",
		},
		&ValidationInfo{
			Name:          "skipped",
			SkippedReason: "skipping",
		},
	}

	wantLog = `Events tried:
	- passed (PASSED)
	- skipped (SKIPPED: skipping)`

	gotLog, gotErr = PrintValidationInfos(passedVIs)
	if gotLog != wantLog {
		t.Errorf("PrintValidationInfos log: got %s, want %s", gotLog, wantLog)
	}
	if gotErr != nil {
		t.Errorf("PrintValidationInfos error: got %v, want nil", gotErr)
	}
}
