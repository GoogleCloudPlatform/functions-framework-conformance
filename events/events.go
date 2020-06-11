// Package events contains the validation logic for different types of events.
package events

import (
	"log"
	"time"

	"cloud.google.com/go/functions/metadata"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type legacyEventBuilder func(e event) map[string]interface{}

type event struct {
	name   string
	ceType string
	data   map[string]interface{}
	meta   metadata.Metadata
	source string
	params map[string]interface{}
	// Builders is the list of functions usable to build this event.
	Builders []legacyEventBuilder
}

var (
	events = []event{}
)

func init() {
	timestamp, err := time.Parse(time.RFC3339, "2020-01-02T15:04:05Z")
	if err != nil {
		log.Fatalf("Parsing timestamp: %v", err)
	}
	// Test pubsub, storage, firestore, firebase realtime database, firebase auth, and firebase analytics events.
	events = []event{
		{
			name:   "pubsub",
			ceType: "com.google.cloud.pubsub.topic.publish.v0",
			source: "http://pubsub.googleapis.com/projects/sample-project/topics/gcf-test",
			meta: metadata.Metadata{
				EventID:   "1144231683168617",
				Timestamp: timestamp,
				EventType: "google.pubsub.topic.publish",
				Resource: &metadata.Resource{
					Service: "pubsub.googleapis.com",
					Name:    "projects/sample-project/topics/gcf-test",
					Type:    "type.googleapis.com/google.pubsub.v1.PubsubMessage",
					RawPath: "projects/sample-project/topics/gcf-test",
				},
			},
			data: map[string]interface{}{
				"@type": "type.googleapis.com/google.pubsub.v1.PubsubMessage",
				"attributes": map[string]interface{}{
					"attr1": "attr1-value",
				},
				"data": "dGVzdCBtZXNzYWdlIDM=",
			},
			Builders: []legacyEventBuilder{buildLegacyEvent, buildLegacyEventWithContext},
		},
		{
			name:   "storage",
			ceType: "com.google.cloud.storage.object.finalize.v0",
			source: "http://storage.googleapis.com/projects/_/buckets/sample-bucket/objects/MyFile#1588778055917163",
			meta: metadata.Metadata{
				EventID:   "1147091835525187",
				Timestamp: timestamp,
				EventType: "google.storage.object.finalize",
				Resource: &metadata.Resource{
					Service: "storage.googleapis.com",
					Name:    "projects/_/buckets/some-bucket/objects/Test.cs",
					Type:    "storage#object",
					RawPath: "projects/_/buckets/sample-bucket/objects/MyFile#1588778055917163",
				},
			},
			data: map[string]interface{}{
				"bucket":                  "some-bucket",
				"contentType":             "text/plain",
				"crc32c":                  "rTVTeQ==",
				"etag":                    "CNHZkbuF/ugCEAE=",
				"generation":              "1587627537231057",
				"id":                      "some-bucket/Test.cs/1587627537231057",
				"kind":                    "storage#object",
				"md5Hash":                 "kF8MuJ5+CTJxvyhHS1xzRg==",
				"mediaLink":               "https://www.googleapis.com/download/storage/v1/b/some-bucket/o/Test.cs?generation=1587627537231057\u0026alt=media",
				"metageneration":          "1",
				"name":                    "Test.cs",
				"selfLink":                "https://www.googleapis.com/storage/v1/b/some-bucket/o/Test.cs",
				"size":                    "352",
				"storageClass":            "MULTI_REGIONAL",
				"timeCreated":             "2020-04-23T07:38:57.230Z",
				"timeStorageClassUpdated": "2020-04-23T07:38:57.230Z",
				"updated":                 "2020-04-23T07:38:57.230Z",
			},
			Builders: []legacyEventBuilder{buildLegacyEvent, buildLegacyEventWithContext},
		},
		{
			name:   "firestore",
			ceType: "com.google.cloud.firestore.document.write.v0",
			source: "http://firestore.googleapis.com/projects/project-id/databases/(default)/documents/gcf-test/2Vm2mI1d0wIaK2Waj5to",
			meta: metadata.Metadata{
				EventID:   "7b8f1804-d38b-4b68-b37d-e2fb5d12d5a0-0",
				EventType: "providers/cloud.firestore/eventTypes/document.write",
				Resource: &metadata.Resource{
					RawPath: "projects/project-id/databases/(default)/documents/gcf-test/2Vm2mI1d0wIaK2Waj5to",
				},
				Timestamp: timestamp,
			},
			params: map[string]interface{}{
				"doc": "2Vm2mI1d0wIaK2Waj5to",
			},
			data: map[string]interface{}{
				"oldValue": map[string]interface{}{
					"createTime": "2020-04-23T09:58:53.211035Z",
					"fields": map[string]interface{}{
						"another test": map[string]interface{}{
							"stringValue": "asd",
						},
					},
					"name":       "projects/project-id/databases/(default)/documents/gcf-test/2Vm2mI1d0wIaK2Waj5to",
					"updateTime": "2020-04-23T12:00:27.247187Z",
				},
				"updateMask": map[string]interface{}{
					"fieldPaths": []interface{}{
						"count",
					},
				},
				"value": map[string]interface{}{
					"createTime": "2020-04-23T09:58:53.211035Z",
					"fields": map[string]interface{}{
						"another test": map[string]interface{}{
							"stringValue": "asd",
						},
					},
					"name":       "projects/project-id/databases/(default)/documents/gcf-test/2Vm2mI1d0wIaK2Waj5to",
					"updateTime": "2020-04-23T12:00:27.247187Z",
				},
			},
			Builders: []legacyEventBuilder{buildLegacyEvent},
		},
		{
			name:   "firebase-realtime-database",
			ceType: "com.google.cloud.firebase.database.write.v0",
			source: "http://firebase.googleapis.com/projects/_/instances/my-project-id/refs/gcf-test/xyz",
			meta: metadata.Metadata{
				EventType: "providers/google.firebase.database/eventTypes/ref.write",
				Resource: &metadata.Resource{
					RawPath: "projects/_/instances/my-project-id/refs/gcf-test/xyz",
				},
				Timestamp: timestamp,
				EventID:   "/SnHth9OSlzK1Puj85kk4tDbF90=",
			},
			params: map[string]interface{}{
				"child": "xyz",
			},
			data: map[string]interface{}{
				"data": nil,
				"delta": map[string]interface{}{
					"grandchild": "other",
				},
			},
			Builders: []legacyEventBuilder{buildLegacyEvent},
		},
		{
			name:   "firebase-auth",
			ceType: "com.google.cloud.firebase.auth.user.create.v0",
			source: "http://firebase.googleapis.com/projects/my-project-id",
			meta: metadata.Metadata{
				EventID:   "4423b4fa-c39b-4f79-b338-977a018e9b55",
				EventType: "providers/firebase.auth/eventTypes/user.create",
				Resource: &metadata.Resource{
					RawPath: "projects/my-project-id"},
				Timestamp: timestamp,
			},
			data: map[string]interface{}{
				"email": "test@nowhere.com",
				"metadata": map[string]interface{}{
					"createdAt": "2020-05-26T10:42:27Z",
				},
				"providerData": []interface{}{
					map[string]interface{}{
						"email":      "test@nowhere.com",
						"providerId": "password",
						"uid":        "test@nowhere.com",
					},
				},
				"uid": "UUpby3s4spZre6kHsgVSPetzQ8l2",
			},
			Builders: []legacyEventBuilder{buildLegacyEvent},
		},
		// {
		// 	name:   "firebase-analytics",
		// 	ceType: "com.google.cloud.firebase.analytics.log.v0",
		// 	Builders: []legacyEventBuilder{buildLegacyEvent},
		// },
	}
}

// BuildCloudEvent creates a CloudEvent based on the given event.
func BuildCloudEvent(e event) *cloudevents.Event {
	ce := cloudevents.NewEvent()
	ce.SetData("application/json", e.data)
	ce.Context.SetID(e.meta.EventID)
	ce.Context.SetType(e.ceType)
	ce.Context.SetTime(e.meta.Timestamp)
	ce.Context.SetSource(e.source)
	return &ce
}

func buildLegacyEventWithContext(e event) map[string]interface{} {
	le := map[string]interface{}{
		"data":   e.data,
		"params": e.params,
		"context": map[string]interface{}{
			"eventId":   e.meta.EventID,
			"timestamp": e.meta.Timestamp.Format(time.RFC3339),
			"eventType": e.meta.EventType,
			"resource": map[string]interface{}{
				"service": e.meta.Resource.Service,
				"name":    e.meta.Resource.Name,
				"type":    e.meta.Resource.Type,
			},
		},
	}
	return le
}

func buildLegacyEvent(e event) map[string]interface{} {
	le := map[string]interface{}{
		"data":      e.data,
		"params":    e.params,
		"eventId":   e.meta.EventID,
		"timestamp": e.meta.Timestamp.Format(time.RFC3339),
		"eventType": e.meta.EventType,
		"resource":  e.meta.Resource.RawPath,
	}
	return le
}

// AllEvents returns a list of all events to test.
func AllEvents() ([]event, error) {
	all := []event{}
	for _, e := range events {
		all = append(all, e)
	}
	return all, nil
}
