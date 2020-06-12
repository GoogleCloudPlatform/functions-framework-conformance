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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"example.com/event-validation/events"
)

const (
	outputFile = "function_output.json"
)

// The HTTP function should copy the contents of the request into the response.
func validateHTTP(url string) error {
	req := `{"res":"PASS"}`
	err := sendHTTP(url, req)
	if err != nil {
		return fmt.Errorf("failed to get response: %v", err)
	}
	output, err := ioutil.ReadFile(outputFile)
	if err != nil {
		return fmt.Errorf("reading output file: %v", err)
	}
	if string(output) != req {
		return fmt.Errorf("unexpected HTTP data: got %s, want %s", output, req)
	}
	return nil
}

func validateLegacyEvents(url string) error {
	allEvents, err := events.AllEvents()
	if err != nil {
		return err
	}

	for _, le := range allEvents {
		for _, build := range le.Builders {
			leJSON, err := json.Marshal(build(le))
			if err != nil {
				return fmt.Errorf("encoding event: %v", err)
			}
			err = sendHTTP(url, string(leJSON))
			if err != nil {
				return fmt.Errorf("response failed: %v", err)
			}
			output, err := ioutil.ReadFile(outputFile)
			if err != nil {
				return fmt.Errorf("reading output file: %v", err)
			}
			if err := events.ValidateLegacyEvent(string(output), le); err != nil {
				return fmt.Errorf("unexpected legacy event: %v", err)
			}
		}
	}

	return nil
}

func validateCloudEvents(url string) error {
	allEvents, err := events.AllEvents()
	if err != nil {
		return err
	}

	for _, ce := range allEvents {
		err := sendCE(url, *events.BuildCloudEvent(ce))
		if err != nil {
			return fmt.Errorf("response failed: %v", err)
		}
		output, err := ioutil.ReadFile(outputFile)
		if err != nil {
			return fmt.Errorf("reading output file: %v", err)
		}
		if err := events.ValidateCloudEvent(string(output), ce); err != nil {
			return fmt.Errorf("unexpected cloud event: %v", err)
		}
	}
	return nil
}

func validate(url, functionType string) error {
	switch functionType {
	case "http":
		// Validate HTTP signature, if provided
		log.Printf("HTTP validation started...")
		if err := validateHTTP(url); err != nil {
			return err
		}
		log.Printf("HTTP validation passed!")
		return nil
	case "cloudevent":
	case "ce":
		// Validate CloudEvent signature, if provided
		log.Printf("CloudEvent validation started...")
		if err := validateCloudEvents(url); err != nil {
			return err
		}
		log.Printf("CloudEvent validation passed!")
		return nil
	case "legacyevent":
	case "le":
		// Validate legacy event signature, if provided
		log.Printf("Legacy event validation started...")
		if err := validateLegacyEvents(url); err != nil {
			return err
		}
		log.Printf("Legacy event validation passed!")
		return nil
	}
	return fmt.Errorf("Expected type to be one of 'http', 'cloudevent', or 'legacyevent', got %s", functionType)
}
