package events

import (
	"encoding/json"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func TestValidateLegacyEvent(t *testing.T) {
	tcs := []struct {
		name  string
		event map[string]interface{}
	}{
		{
			name: "resource service defined",
			event: map[string]interface{}{
				"data": map[string]interface{}{
					"@type": "type.googleapis.com/google.pubsub.v1.PubsubMessage",
					"attributes": map[string]interface{}{
						"attr1": "attr1-value",
					},
					"data": "dGVzdCBtZXNzYWdlIDM=",
				},
				"context": map[string]interface{}{
					"eventId":   "1144231683168617",
					"timestamp": "2020-01-02T15:04:05Z",
					"eventType": "google.pubsub.topic.publish",
					"resource": map[string]interface{}{
						"service": "pubsub.googleapis.com",
						"name":    "projects/sample-project/topics/gcf-test",
						"type":    "type.googleapis.com/google.pubsub.v1.PubsubMessage",
					},
				},
			},
		},
		{
			name: "resource path defined",
			event: map[string]interface{}{
				"data": map[string]interface{}{
					"@type": "type.googleapis.com/google.pubsub.v1.PubsubMessage",
					"attributes": map[string]interface{}{
						"attr1": "attr1-value",
					},
					"data": "dGVzdCBtZXNzYWdlIDM=",
				},
				"context": map[string]interface{}{
					"eventId":   "1144231683168617",
					"timestamp": "2020-01-02T15:04:05Z",
					"eventType": "google.pubsub.topic.publish",
					"resource":  "projects/sample-project/topics/gcf-test",
				},
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			e, err := json.Marshal(tc.event)
			if err != nil {
				t.Fatalf("marshalling legacy event: %v", err)
			}
			// First event in events is the pubsub event representation.
			if err := ValidateLegacyEvent(string(e), events[0]); err != nil {
				t.Errorf("validating legacy event: %v", err)
			}
		})
	}
}

func TestValidateCloudEvent(t *testing.T) {
	timestamp, err := time.Parse(time.RFC3339, "2020-01-02T15:04:05Z")
	if err != nil {
		t.Fatalf("Parsing timestamp: %v", err)
	}
	ce := cloudevents.NewEvent()
	ce.Context.SetID("1144231683168617")
	ce.Context.SetType("com.google.cloud.pubsub.topic.publish.v0")
	ce.Context.SetTime(timestamp)
	ce.Context.SetSource("http://pubsub.googleapis.com/projects/sample-project/topics/gcf-test")
	ce.SetData("application/json", map[string]interface{}{
		"@type": "type.googleapis.com/google.pubsub.v1.PubsubMessage",
		"attributes": map[string]interface{}{
			"attr1": "attr1-value",
		},
		"data": "dGVzdCBtZXNzYWdlIDM=",
	})

	bytes, err := json.Marshal(ce)
	if err != nil {
		t.Fatalf("marshalling cloud event: %v", err)
	}
	// First event in events is the pubsub event representation.
	if err := ValidateCloudEvent(string(bytes), events[0]); err != nil {
		t.Errorf("validating cloud event: %v", err)
	}
}
