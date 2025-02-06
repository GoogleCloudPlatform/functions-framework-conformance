# Mapping between Google Cloud Function events and CloudEvents

One of the responsibilities of Functions Frameworks that is tested
by the conformance tests in this repository is that of event
conversion.

There are three representations of events that are relevant:

- The HTTP request representation of events as they are currently
  delivered to Google Cloud Functions. This document refers to this
  as the "GCF HTTP" representation.
- The "event, context" in-language representation as they are
  delivered by Functions Frameworks to classic background events,  
  [as described here](https://cloud.google.com/functions/docs/writing/background).  
- [CNCF CloudEvents](https://github.com/cloudevents/spec/blob/v1.0/spec.md) -
  either in one of the [HTTP
  representations](https://github.com/cloudevents/spec/blob/v1.0/http-protocol-binding.md)
  or an in-language representation using a standard CloudEvents SDK.

All Functions Framework provide conversions from the GCF HTTP
format to an in-language CloudEvent format. This allows Functions
Framework end users to write functions with a CloudEvent signature.

Some Functions Frameworks (Java, Go, Node, Python) also provide
conversions from the CloudEvent format to the "event, context"
in-language representation. This allows functions using the "event,
context" signature to operate in without changes an environment
where CloudEvents are the native representation in HTTP requests.

This document describes both kinds of conversion. For each kind,
there's a general approach, followed by event-specific additional
requirements.

# GCF HTTP to CloudEvent conversion

## General flow

All relevant information in the GCF HTTP representation is within
the HTTP body, as a JSON object. Note that context information can
come from either the root of the JSON, or within a `context`
property. The description below uses paths within the JSON to
represent nested information. For example, the path `context/eventType`
refers to "the `eventType` property within the `context` property of the root
object".

First obtain the following information from the JSON:

- *gcf_event_type*: from `context/eventType` or `eventType`
- *gcf_service*: from `context/resource/service`, or mapped from the
  GCF event type if absent (see below).
- *gcf_resource_name*: from `context/resource/name` or `resource`
- *gcf_event_id*: from `context/eventId` or `eventId`
- *gcf_timestamp*: from `context/timestamp` or `timestamp`
- *gcf_data*: the whole JSON object in `data`

Where multiple sources are listed and information is present in more
than one location, the first matching location is used. (In reality
we don't expect to see this, but the rule provides consistency
between Functions frameworks.)

Some GCF HTTP requests do not contain a `context/resource/service`
value. In this case, the *gcf_service* is determined by a prefix
match of *gcf_event_type* against the table below.

|*gcf_event_type* prefix|*gcf_service*|
|-|-|
|providers/cloud.firestore/|firestore.googleapis.com|
|providers/google.firebase.analytics/|firebaseanalytics.googleapis.com|
|providers/firebase.auth/|firebaseauth.googleapis.com|
|providers/google.firebase.database/|firebasedatabase.googleapis.com|
|providers/cloud.pubsub/|pubsub.googleapis.com|
|providers/cloud.storage/|storage.googleapis.com|

Next, determine the *ce_type* based on an exact match of
*gcf_event_type* as given below:

|*gcf_event_type*|*ce_type*|
|-|-|
|google.pubsub.topic.publish|google.cloud.pubsub.topic.v1.messagePublished|
|providers/cloud.pubsub/eventTypes/topic.publish|google.cloud.pubsub.topic.v1.messagePublished|
|google.storage.object.finalize|google.cloud.storage.object.v1.finalized|
|google.storage.object.delete|google.cloud.storage.object.v1.deleted|
|google.storage.object.archive|google.cloud.storage.object.v1.archived|
|google.storage.object.metadataUpdate|google.cloud.storage.object.v1.metadataUpdated|
|providers/cloud.firestore/eventTypes/document.write|google.cloud.firestore.document.v1.written|
|providers/cloud.firestore/eventTypes/document.create|google.cloud.firestore.document.v1.created|
|providers/cloud.firestore/eventTypes/document.update|google.cloud.firestore.document.v1.updated|
|providers/cloud.firestore/eventTypes/document.delete|google.cloud.firestore.document.v1.deleted|
|providers/firebase.auth/eventTypes/user.create|google.firebase.auth.user.v1.created|
|providers/firebase.auth/eventTypes/user.delete|google.firebase.auth.user.v1.deleted|
|providers/firebase.remoteConfig/remoteconfig.update|google.firebase.remoteconfig.remoteConfig.v1.updated|
|providers/google.firebase.analytics/eventTypes/event.log|google.firebase.analytics.log.v1.written|
|providers/google.firebase.database/eventTypes/ref.create|google.firebase.database.ref.v1.created|
|providers/google.firebase.database/eventTypes/ref.write|google.firebase.database.ref.v1.written|
|providers/google.firebase.database/eventTypes/ref.update|google.firebase.database.ref.v1.updated|
|providers/google.firebase.database/eventTypes/ref.delete|google.firebase.database.ref.v1.deleted|

Finally, construct a new CloudEvent with the following attributes:

- `id`: *gcf_event_id*
- `source`: `//`*gcf_service*`/`*gcf_resource_name*
- `type`: *ce_type*
- `time`: *gcf_timestamp*
- `specversion`: `1.0`
- `datacontenttype`: `application/json`
- `data`: *gcf_data*

Note: the `subject` attribute is only present where the specific
event type defines it in the sections below.

## Event-specific requirements

### Cloud Storage events

The GCF HTTP representation of Cloud Storage events (events with a
*gcf_service* of `storage.googleapis.com`) specify a
*gcf_resource_name* which includes both the
[bucket name](https://cloud.google.com/storage/docs/key-terms#bucket-names) and the
[object
name](https://cloud.google.com/storage/docs/key-terms#object-names).
This sometimes including the *generation* of the object as well.

The conversion of the event to a CloudEvent separates these out: the
CloudEvent `source` attribute includes the resource name up to and
including the bucket name, and the CloudEvent `subject` attribute is
the remainder of the resource name, but not including the
generation. Note that the object name may include slashes (but the
bucket name will not).

For example, after performing the generic data extraction described
earlier, if the results include:

- *gcf_service*: `storage.googleapis.com`
- *gcf_resource_name*: `projects/_/buckets/sample-bucket/objects/folder/MyFile#1588778055917163`

This will lead to CloudEvent attributes of:

- `source`: `//storage.googleapis.com/projects/_/buckets/sample-bucket`
- `subject`: `objects/folder/MyFile`

### Cloud PubSub events

The GCF HTTP representation of Cloud PubSub events (events with a
*gcf_service* of `pubsub.googleapis.com`) contain a `data` property
which needs to be wrapped in an extra JSON object in the CloudEvent
`data` attribute, to conform with the [expected CloudEvent
representation](https://github.com/googleapis/google-cloudevents/blob/main/proto/google/events/cloud/pubsub/v1/data.proto).

So the CloudEvent should have:

`data`: `{ "message":` *gcf_data* `}`

(The `subscription` property of the CloudEvent `data` attribute
should not be populated at the moment. If the subscription name
becomes available in the GCF HTTP representation, this document will
be updated, along with the conformance tests.)

Additionally, two properties should be populated in the message,
based on the context:

- The `messageId` property in the `message` object should be set to *gcf_event_id*
- The `publishTime` property in the `message` object should be set to *gcf_timestamp*

The conversion should **not** parse *gcf_data* to ensure that only
expected properties are present. (For example, the GCF HTTP
representation usually includes a `@type` property, always with a
value of `"type.googleapis.com/google.pubsub.v1.PubsubMessage"`).
The inclusion or removal of the extra properties should make little
difference to users, it's simpler to write conformance tests if all
Functions Frameworks behave consistently.

### Firebase RTDB events

The `resource` in the GCF HTTP representation is of the form
`projects/_/instances/{instance-id}/refs/{ref-path}`. Additionally,
there is a top-level `domain` property, which is used to determine the `location` part
of the CloudEvent representation:

- The `domain` property must be present as a string; if it's missing, the conversion should fail.
- If the `domain` value is `firebaseio.com`, the location is `us-central1`.
- Otherwise, the location is the value of `domain` before the first period. For example,
  a `domain` value of `europe-west1.firebasedatabase.app` would lead to a location value
  of `europe-west1`.

In the CloudEvent representation, this information is split between
the `source` and the `subject`:

- `source`: `//firestore.googleapis.com/projects/_/locations/{location}/instances/{instance-id}`
- `subject: refs/{ref-path}`

### Firebase analytics events

The `resource` property in the GCF HTTP representation is of the
form `projects/{project-id}/events/{event-name}`. As part of
conversion to a CloudEvent, this is split between the `subject` and
the `source` attributes:

- `source`: `//firebaseanalytics.googleapis.com/projects/{project-id}/apps/{app-id}`
- `subject`: `events/{event-name}`

The `app-id` part it obtained from the data within the GCF HTTP representation, from
a path of `userDim.appInfo.appId` (both the `userDim` and `appInfo` properties are
expected to have object values; the `appId` property is expected to have a string value.)

### Firebase auth events

The `subject` attribute of the CloudEvent is of the form `users/{uid}` where the
`uid` value is taken from the `uid` property within the original
`data` property. This value is still retained within the CloudEvent data as well.

Additionally, two properties within the CloudEvent have different names
to those in the GCF HTTP event. Within the `metadata` top-level
property, there are two timestamps. In GCF HTTP events these have
names of `createdAt` and `lastSignedInAt`; in the CloudEvent
representation the names are `createTime` and `lastSignInTime`
respectively.

### Firestore document events

The `resource` in the GCF HTTP representation is of the form
`projects/{project-id}/databases/{database-id}/documents/{path-to-document}`.
In the CloudEvent representation, this information is split between
the `source` and the `subject`:

- `source`: `//firestore.googleapis.com/projects/{project-id}/databases/{database-id}`
- `subject: documents/{document-id}`

# CloudEvent to "event, context" representation

Here "event" represents the payload of the event, either as a string
of raw JSON or a deserialized in-lanuage representation. The
"context" contains the following information:

- `eventId`: A unique ID for the event. (String)
- `timestamp`: The date/time this event was created. (String)
- `eventType`: The type of the event.
  For example: "google.pubsub.topic.publish". (String) 
- `resource`: The resource that emitted the event. (The format of this
  varies by language.)

The Java Functions Framework additionally includes an `attributes`
map from string keys to string values.

## General flow

The "event" corresponds directly to the CloudEvent `data` attribute,
unless otherwise specified.

The "context" is populated as follows:

- `eventId` from the CloudEvent `id` attribute
- `timestamp` from the CloudEvent `timestamp` attribute
- `eventType` from the CloudEvent `type` attribute, with a reverse
  mapping of the *gcf_event_type* to *ce_type* table earlier applied.
- `resource`: TBD

### Cloud Storage events

TBD

### Cloud PubSub events

TBD
