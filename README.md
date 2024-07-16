# Table of contents
1. [Gommentary](#gommentary)
2. [How to run](#how-to-run)
   1. [Docker](#docker)
3. [Configuration](#configuration)
4. [Open-api](#open-api)
5. [Development](#development)
6. [Example responses](#example-responses)

# Gommentary
Service offering a comment REST API.   
The idea is to create comments to "subjects". Subject is anything that can be commented on including other comments. This allows us to create replies. Replies and edit counts of the comment are shown in the main CommentsPage to the subject.   
Authors are pulled out of the JWT `sub` field based on the application configuration.
Subject ids must be unique.

# How to run
If you have Go environment set up just clone the project and run `make run-dev` to run in dev mode.
## Docker
Build the docker image `make build-docker` and then run it with 
```sh 
docker run -p 8090:8090 docker.io/library/gommentary:latest
```  
If you want to preserve the database file between restarts mount a directory and point database file to it 
```sh 
docker run -p 8090:8090 -v ${PWD}/db:/db -e DB_FILE=/db/gommentary.db docker.io/library/gommentary:latest
```

# Configuration
Example configuration is available [here](./conf/api/conf-prod.yaml). To run the app without JWT support set `jwt.required` to `false`. This will allow you to run the app without Authors and every request can edit any comment. If this is set to true only authors can edit their comments. Authors are pulled out of the JWT `sub` field.

Main configuration struct is located [here](./internal/config/config.go) take a look there to see what else can be configured.

# Open-api
Open api definition can be found in [here](./open-api.yaml)

# Development
[Makefile](./Makefile) is your friend. Run `make help` to get info about targets.

# Example responses

Get comments to subject `random-subject-65465161`
```json
{
	"content": [
		{
			"author": "johndoe",
			"date": "2024-07-16T13:03:06Z",
			"edits": 0,
			"id": "fedcdf7d-f9cc-42f8-a9ff-0b9fade02f37",
			"position": 1,
			"replies": 1,
			"text": "comment one"
		},
		{
			"author": "johndoe",
			"date": "2024-07-16T13:03:10Z",
			"edits": 0,
			"id": "448ff5a7-4977-4190-a7fa-6f4ba3bf5fa9",
			"position": 2,
			"replies": 0,
			"text": "comment two"
		},
		{
			"author": "johndoe",
			"date": "2024-07-16T13:03:11Z",
			"edits": 1,
			"id": "48d3c22e-6603-43af-a910-c58c4af6341b",
			"position": 3,
			"replies": 0,
			"text": "comment three"
		}
	],
	"first": true,
	"last": true,
	"size": 3,
	"subject": {
		"id": "random-subject-65465161"
	},
	"totalElements": 3,
	"totalPages": 1
}
```

Get replies to comment `fedcdf7d-f9cc-42f8-a9ff-0b9fade02f37`
```json
{
	"content": [
		{
			"author": "johndoe",
			"date": "2024-07-16T13:04:53Z",
			"edits": 0,
			"id": "4037b1a5-37de-4e90-bc1f-91ccf26e720a",
			"position": 1,
			"replies": 0,
			"text": "reply to comment one"
		}
	],
	"first": true,
	"last": true,
	"size": 1,
	"subject": {
		"id": "fedcdf7d-f9cc-42f8-a9ff-0b9fade02f37"
	},
	"totalElements": 1,
	"totalPages": 1
}
```
