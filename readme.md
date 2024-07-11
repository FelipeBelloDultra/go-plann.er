# go-plann.er

## What is it?

Rest API developed in GoLang for trips scheduling with sending emails.

## API specifications

Full API specifications for the backend of the go-plann.er application built during Rocketseat's NLW Journey.

## Version: 1.0.0

### /trips/{tripId}/confirm

#### GET

##### Summary:

Confirm a trip and send e-mail invitations.

##### Parameters

| Name   | Located in | Description | Required | Schema        |
| ------ | ---------- | ----------- | -------- | ------------- |
| tripId | path       |             | Yes      | string (uuid) |

##### Responses

| Code | Description      |
| ---- | ---------------- |
| 204  | Default Response |
| 400  | Bad request      |

### /participants/{participantId}/confirm

#### PATCH

##### Summary:

Confirms a participant on a trip.

##### Parameters

| Name          | Located in | Description | Required | Schema        |
| ------------- | ---------- | ----------- | -------- | ------------- |
| participantId | path       |             | Yes      | string (uuid) |

##### Responses

| Code | Description      |
| ---- | ---------------- |
| 204  | Default Response |
| 400  | Bad request      |

### /trips/{tripId}/invites

#### POST

##### Summary:

Invite someone to the trip.

##### Parameters

| Name   | Located in | Description | Required | Schema        |
| ------ | ---------- | ----------- | -------- | ------------- |
| tripId | path       |             | Yes      | string (uuid) |

##### Responses

| Code | Description      |
| ---- | ---------------- |
| 201  | Default Response |
| 400  | Bad request      |

### /trips/{tripId}/activities

#### POST

##### Summary:

Create a trip activity.

##### Parameters

| Name   | Located in | Description | Required | Schema        |
| ------ | ---------- | ----------- | -------- | ------------- |
| tripId | path       |             | Yes      | string (uuid) |

##### Responses

| Code | Description      |
| ---- | ---------------- |
| 201  | Default Response |
| 400  | Bad request      |

#### GET

##### Summary:

Get a trip activities.

##### Description:

This route will return all the dates between the trip starts_at and ends_at dates, even those without activities.

##### Parameters

| Name   | Located in | Description | Required | Schema        |
| ------ | ---------- | ----------- | -------- | ------------- |
| tripId | path       |             | Yes      | string (uuid) |

##### Responses

| Code | Description      |
| ---- | ---------------- |
| 200  | Default Response |
| 400  | Bad request      |

### /trips/{tripId}/links

#### POST

##### Summary:

Create a trip link.

##### Parameters

| Name   | Located in | Description | Required | Schema        |
| ------ | ---------- | ----------- | -------- | ------------- |
| tripId | path       |             | Yes      | string (uuid) |

##### Responses

| Code | Description      |
| ---- | ---------------- |
| 201  | Default Response |
| 400  | Bad request      |

#### GET

##### Summary:

Get a trip links.

##### Parameters

| Name   | Located in | Description | Required | Schema        |
| ------ | ---------- | ----------- | -------- | ------------- |
| tripId | path       |             | Yes      | string (uuid) |

##### Responses

| Code | Description      |
| ---- | ---------------- |
| 200  | Default Response |
| 400  | Bad request      |

### /trips

#### POST

##### Summary:

Create a new trip

##### Responses

| Code | Description      |
| ---- | ---------------- |
| 201  | Default Response |
| 400  | Bad request      |

### /trips/{tripId}

#### GET

##### Summary:

Get a trip details.

##### Parameters

| Name   | Located in | Description | Required | Schema        |
| ------ | ---------- | ----------- | -------- | ------------- |
| tripId | path       |             | Yes      | string (uuid) |

##### Responses

| Code | Description      |
| ---- | ---------------- |
| 200  | Default Response |
| 400  | Bad request      |

#### PUT

##### Summary:

Update a trip.

##### Parameters

| Name   | Located in | Description | Required | Schema        |
| ------ | ---------- | ----------- | -------- | ------------- |
| tripId | path       |             | Yes      | string (uuid) |

##### Responses

| Code | Description      |
| ---- | ---------------- |
| 204  | Default Response |
| 400  | Bad request      |

### /trips/{tripId}/participants

#### GET

##### Summary:

Get a trip participants.

##### Parameters

| Name   | Located in | Description | Required | Schema        |
| ------ | ---------- | ----------- | -------- | ------------- |
| tripId | path       |             | Yes      | string (uuid) |

##### Responses

| Code | Description      |
| ---- | ---------------- |
| 200  | Default Response |
| 400  | Bad request      |
