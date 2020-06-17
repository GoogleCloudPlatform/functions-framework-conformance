// This binary contains a valiation framework for functions frameworks.
package main

import (
	"encoding/json"
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
	req := "PASS"
	err := sendHTTP(url, req)
	if err != nil {
		return fmt.Errorf("failed to get response: %v", err)
	}
	output, err := ioutil.ReadFile(outputFile)
	if err != nil {
		return fmt.Errorf("reading output file: %v", err)
	}
	if string(output) != "PASS" {
		return fmt.Errorf("unexpected HTTP data: got %q, want 'PASS'", output)
	}
	return nil
}

func validateEvents(url, functionType string) error {
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
		if err := validateEvents(url, "cloudevent"); err != nil {
			return err
		}
		log.Printf("CloudEvent validation passed!")
	case "legacyevent":
	case "le":
		// Validate legacy event signature, if provided
		log.Printf("Legacy event validation started...")
		if err := validateEvents(url, "legacyevent"); err != nil {
			return err
		}
		log.Printf("Legacy event validation passed!")
	}
	return fmt.Errorf("Expected type to be one of 'http', 'cloudevent', or 'legacyevent', got %s", functionType)
}
