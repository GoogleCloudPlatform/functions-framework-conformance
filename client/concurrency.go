// Copyright 2022 Google LLC
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
	"log"
	"sync"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-conformance/events"
)

func timeExecution(fn func() error) (time.Duration, error) {
	start := time.Now()
	err := fn()
	return time.Since(start), err
}

// validateConcurrency validates a server can handle concurrent requests by
// valdating that the response time for a single request does not increase
// linearly with n concurrent requests, given a function that:
// 1. Is not CPU-bound (e.g. sleeps)
// 2. Executes for at least 1s to ensure non-trivial measurement differences
func validateConcurrency(url string, functionType string) error {
	log.Printf("%s validation with concurrent requests...", functionType)
	var sendFn func() error
	switch functionType {
	case "http":
		sendFn = func() error {
			return sendHTTP(url, []byte(`{"data": "hello"}`))
		}
	case "cloudevent":
		// Arbitrary payload that conforms to CloudEvent schema
		sendFn = func() error {
			return send(url, events.CloudEvent, []byte(`{
			"specversion": "1.0",
			"type": "google.firebase.auth.user.v1.created",
			"source": "//firebaseauth.googleapis.com/projects/my-project-id",
			"subject": "users/UUpby3s4spZre6kHsgVSPetzQ8l2",
			"id": "aaaaaa-1111-bbbb-2222-cccccccccccc",
			"time": "2020-09-29T11:32:00.123Z",
			"datacontenttype": "application/json",
			"data": {
			  "email": "test@nowhere.com",
			  "metadata": {
				"createTime": "2020-05-26T10:42:27Z",
				"lastSignInTime": "2020-10-24T11:00:00Z"
			  },
			  "providerData": [
				{
				  "email": "test@nowhere.com",
				  "providerId": "password",
				  "uid": "test@nowhere.com"
				}
			  ],
			  "uid": "UUpby3s4spZre6kHsgVSPetzQ8l2"
			}
		  }`))
		}
	case "legacyevent":
		// Arbitrary payload that conforms to Background event schema
		sendFn = func() error {
			return send(url, events.LegacyEvent, []byte(`{
			"data": {
			  "email": "test@nowhere.com",
			  "metadata": {
				"createdAt": "2020-05-26T10:42:27Z",
				"lastSignedInAt": "2020-10-24T11:00:00Z"
			  },
			  "providerData": [
				{
				  "email": "test@nowhere.com",
				  "providerId": "password",
				  "uid": "test@nowhere.com"
				}
			  ],
			  "uid": "UUpby3s4spZre6kHsgVSPetzQ8l2"
			},
			"eventId": "aaaaaa-1111-bbbb-2222-cccccccccccc",
			"eventType": "providers/firebase.auth/eventTypes/user.create",
			"notSupported": {
			},
			"resource": "projects/my-project-id",
			"timestamp": "2020-09-29T11:32:00.123Z"
		  }`))
		}
	default:
		return fmt.Errorf("expected type to be one of 'http', 'cloudevent', or 'legacyevent', got %s", functionType)
	}
	if err := sendConcurrentRequests(sendFn); err != nil {
		return err
	}
	log.Printf("Concurrency validation passed!")
	return nil
}

func sendConcurrentRequests(sendFn func() error) error {
	// Get a benchmark for the time it takes for a single request
	singleReqTime, singleReqErr := timeExecution(func() error {
		return sendFn()
	})
	if singleReqErr != nil {
		return fmt.Errorf("concurrent validation unable to send single request to benchmark response time: %v", singleReqErr)
	}

	minWait := 1 * time.Second
	if singleReqTime < minWait {
		return fmt.Errorf("concurrent validation requires a function that waits at least %s before responding, function responded in %s", minWait, singleReqTime)
	}
	log.Printf("Single request response time benchmarked, took %s for 1 request", singleReqTime)

	// Get a benchmark for the time it takes for concurrent requests
	const numConReqs = 10
	log.Printf("Starting %d concurrent workers to send requests", numConReqs)

	type workerResponse struct {
		id  int
		err error
	}
	var wg sync.WaitGroup
	respCh := make(chan workerResponse, numConReqs)
	conReqTime, _ := timeExecution(func() error {
		for i := 0; i < numConReqs; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				err := sendFn()
				respCh <- workerResponse{id: id, err: err}
			}(i)
		}

		wg.Wait()
		return nil
	})

	maybeErrMessage := ""
	for i := 0; i < numConReqs; i++ {
		resp := <-respCh
		if resp.err != nil {
			maybeErrMessage += fmt.Sprintf("error #%d: %v\n", i, resp.err)
		} else {
			log.Printf("Worker #%d done", resp.id)
		}
	}

	if maybeErrMessage != "" {
		return fmt.Errorf("at least one concurrent request failed:\n%s", maybeErrMessage)
	}

	// Validate that the concurrent requests were handled faster than if all
	// the requests were handled serially, using the single request time
	// as a benchmark. The concurrent time should be less than half of the
	// time it would have taken to execute all requests serially.
	if conReqTime > 2*singleReqTime {
		return fmt.Errorf("function took too long to complete %d concurrent requests. %d concurrent request time: %s, single request time: %s", numConReqs, numConReqs, conReqTime, singleReqTime)
	}
	log.Printf("Concurrent request response time benchmarked, took %s for %d requests", conReqTime, numConReqs)
	return nil
}
