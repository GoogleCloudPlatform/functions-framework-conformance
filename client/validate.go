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
	"fmt"
	"io/ioutil"
	"log"

	"github.com/GoogleCloudPlatform/functions-framework-conformance/events"
)

const (
	outputFile = "function_output.json"
)

// The HTTP function should copy the contents of the request into the response.
func validateHTTP(url string) error {
	req := []byte(`{"res":"PASS"}`)
	err := sendHTTP(url, req)
	if err != nil {
		return fmt.Errorf("failed to get response: %v", err)
	}
	output, err := ioutil.ReadFile(outputFile)
	if err != nil {
		return fmt.Errorf("reading output file: %v", err)
	}
	if string(output) != string(req) {
		return fmt.Errorf("unexpected HTTP data: got %s, want %s", output, req)
	}
	return nil
}

func validateEvents(url string, inputType, outputType events.EventType) error {
	eventNames, err := events.EventNames(inputType)
	if err != nil {
		return err
	}

	for _, name := range eventNames {
		input := events.InputData(name, inputType)
		if input == nil {
			return fmt.Errorf("no input data for %q", name)
		}
		err = send(url, inputType, input)
		if err != nil {
			return fmt.Errorf("response failed for %q: %v", name, err)
		}
		output, err := ioutil.ReadFile(outputFile)
		if err != nil {
			return fmt.Errorf("reading function output file for %q: %v", name, err)
		}
		if err := events.ValidateEvent(name, outputType, output); err != nil {
			return fmt.Errorf("unexpected output for %q: %v", name, err)
		}
	}

	return nil
}

func validate(url, functionType string, validateMapping bool) error {
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
		// Validate CloudEvent signature, if provided
		log.Printf("CloudEvent validation started...")
		if err := validateEvents(url, events.CloudEvent, events.CloudEvent); err != nil {
			return err
		}
		if validateMapping {
			if err := validateEvents(url, events.LegacyEvent, events.CloudEvent); err != nil {
				return err
			}
		}
		log.Printf("CloudEvent validation passed!")
		return nil
	case "legacyevent":
		// Validate legacy event signature, if provided
		log.Printf("Legacy event validation started...")
		if err := validateEvents(url, events.LegacyEvent, events.LegacyEvent); err != nil {
			return err
		}
		if validateMapping {
			if err := validateEvents(url, events.CloudEvent, events.LegacyEvent); err != nil {
				return err
			}
		}
		log.Printf("Legacy event validation passed!")
		return nil
	}
	return fmt.Errorf("Expected type to be one of 'http', 'cloudevent', or 'legacyevent', got %s", functionType)
}
