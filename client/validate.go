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
	"log"

	"github.com/GoogleCloudPlatform/functions-framework-conformance/events"
	"github.com/google/go-cmp/cmp"
)

type validatorParams struct {
	useBuildpacks   bool
	validateMapping bool
	runCmd          string
	outputFile      string
	source          string
	target          string
	runtime         string
	tag             string
	functionType    string
}

type validator struct {
	funcServer      functionServer
	validateMapping bool
	functionType    string
}

func newValidator(params validatorParams) *validator {
	v := validator{
		validateMapping: params.validateMapping,
		functionType:    params.functionType,
	}

	if !params.useBuildpacks {
		v.funcServer = &localFunctionServer{
			output: params.outputFile,
			cmd:    params.runCmd,
		}
		return &v
	}

	if params.functionType == "legacyevent" {
		params.functionType = "event"
	}

	v.funcServer = &buildpacksFunctionServer{
		output:   params.outputFile,
		source:   params.source,
		target:   params.target,
		runtime:  params.runtime,
		tag:      params.tag,
		funcType: params.functionType,
	}
	return &v
}

func (v validator) runValidation() error {
	log.Printf("Validating for %s...", *functionType)

	shutdown, err := v.funcServer.Start()
	if shutdown != nil {
		defer shutdown()
	}

	if err != nil {
		return fmt.Errorf("unable to start server: %v", err)
	}

	if err := v.validate("http://localhost:8080"); err != nil {
		return fmt.Errorf("Validation failure: %v", err)
	}

	log.Printf("All validation passed!")
	return nil
}

// The HTTP function should copy the contents of the request into the response.
func (v validator) validateHTTP(url string) error {
	want := map[string]string{
		"res": "PASS",
	}

	req, err := json.Marshal(want)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %v", err)
	}

	if err := sendHTTP(url, req); err != nil {
		return fmt.Errorf("failed to get response from HTTP function: %v", err)
	}
	output, err := v.funcServer.OutputFile()
	if err != nil {
		return fmt.Errorf("reading output file from HTTP function: %v", err)
	}

	got := make(map[string]string)
	if err = json.Unmarshal(output, &got); err != nil {
		return fmt.Errorf("failed to unmarshal json: %v", err)
	}

	if !cmp.Equal(got, want) {
		return fmt.Errorf("unexpected HTTP output data: got %v, want %v", got, want)
	}
	return nil
}

func (v validator) validateEvents(url string, inputType, outputType events.EventType) error {
	eventNames, err := events.EventNames(inputType)
	if err != nil {
		return err
	}

	vis := []*events.ValidationInfo{}
	for _, name := range eventNames {
		input := events.InputData(name, inputType)
		if input == nil {
			return fmt.Errorf("no input data for event %q", name)
		}
		err = send(url, inputType, input)
		if err != nil {
			return fmt.Errorf("failed to get response from function for %q: %v", name, err)
		}
		output, err := v.funcServer.OutputFile()
		if err != nil {
			return fmt.Errorf("reading output file from function for %q: %v", name, err)
		}
		if vi := events.ValidateEvent(name, inputType, outputType, output); vi != nil {
			vis = append(vis, vi)
		}
	}

	logStr, err := events.PrintValidationInfos(vis)
	log.Printf(logStr)
	return err
}

func (v validator) validate(url string) error {
	switch v.functionType {
	case "http":
		// Validate HTTP signature, if provided
		log.Printf("HTTP validation started...")
		if err := v.validateHTTP(url); err != nil {
			return err
		}
		log.Printf("HTTP validation passed!")
		return nil
	case "cloudevent":
		// Validate CloudEvent signature, if provided
		log.Printf("CloudEvent validation with CloudEvent requests...")
		if err := v.validateEvents(url, events.CloudEvent, events.CloudEvent); err != nil {
			return err
		}
		if v.validateMapping {
			log.Printf("CloudEvent validation with legacy event requests...")
			if err := v.validateEvents(url, events.LegacyEvent, events.CloudEvent); err != nil {
				return err
			}
		}
		log.Printf("CloudEvent validation passed!")
		return nil
	case "legacyevent":
		// Validate legacy event signature, if provided
		log.Printf("Legacy event validation with legacy event requests...")
		if err := v.validateEvents(url, events.LegacyEvent, events.LegacyEvent); err != nil {
			return err
		}
		if v.validateMapping {
			log.Printf("Legacy event validation with CloudEvent requests...")
			if err := v.validateEvents(url, events.CloudEvent, events.LegacyEvent); err != nil {
				return err
			}
		}
		log.Printf("Legacy event validation passed!")
		return nil
	}
	return fmt.Errorf("Expected type to be one of 'http', 'cloudevent', or 'legacyevent', got %s", v.functionType)
}
