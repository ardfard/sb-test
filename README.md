# SB-Test

[![Unit Test](https://github.com/ardfard/sb-test/actions/workflows/go.yml/badge.svg)](https://github.com/ardfard/sb-test/actions/workflows/go.yml)
[![Smoke Test](https://github.com/ardfard/sb-test/actions/workflows/smoke-test.yml/badge.svg)](https://github.com/ardfard/sb-test/actions/workflows/smoke-test.yml)
[![codecov](https://codecov.io/gh/ardfard/sb-test/graph/badge.svg)](https://codecov.io/gh/ardfard/sb-test)
[![Go Report Card](https://goreportcard.com/badge/github.com/ardfard/sb-test)](https://goreportcard.com/report/github.com/ardfard/sb-test)

This is an audio processing service that demonstrates a basic API for handling audio file uploads and format conversions. This service provides the following endpoints:

- POST /audio/user/{user_id}/phrase/{phrase_id} (Upload)
- GET /audio/user/{user_id}/phrase/{phrase_id}/{format} (Download)
- POST /users (Create a basic user)
- POST /users/{user_id}/phrases (Create a basic phrase for the user)

## Running the service and simulate common use cases
Assuming you have docker installed, you can run the service and simulate common use cases with the following script:

```bash
./scripts/run.sh
```

This will:
- Create a user
- Create a phrase for the user
- Upload an audio file for the phrase
- Download the audio file in both m4a and flac formats
- Run the service at http://localhost:8080

[![asciicast](https://asciinema.org/a/ei7JOoIGszahzomC5OxXlQIpB.svg)](https://asciinema.org/a/ei7JOoIGszahzomC5OxXlQIpB)

## Assumptions

There are few assumptions that were made when designing the service:

- The service is more focused on the audio processing and storage than the other user functionalities. Hence why there is only rudimentary user management and no authentication or authorization.
- The service is designed to be run on a local machine with simplicity in mind for the sake of demo purposes.
- The service exposes a RESTful API that is used to upload and download the audio files.
- The audio file is always converted to wav format for long term storage in the object storage. After the conversion, the original audio file is deleted from the object storage.
- When downloading the audio file, the service will convert the audio file to the requested format from the WAV format or from the original format if it's not yet deleted from the object storage.
- The service uses a background processing to convert the audio files to the wav format when uploading for better response time and reliability.

## Using the service

The service exposes a RESTful API that is used to upload and download the audio files. Before you can upload an audio file, you need to create a user and a phrase.

### Creating a user

```bash
curl -X POST http://localhost:8080/users -H 'Content-Type: application/json' -d '{"name": "John Musou"}'
```

### Creating a phrase

```bash
curl -X POST http://localhost:8080/users/{user_id}/phrases -H 'Content-Type: application/json' -d '{"text": "Hello, world!"}'
```

### Uploading an audio file

```bash
curl -X POST http://localhost:8080/audio/user/{user_id}/phrase/{phrase_id} -H 'Content-Type: multipart/form-data' -F 'audio_file=@path/to/your/audio/file'
```

### Downloading an audio file

```bash
curl -X GET http://localhost:8080/audio/user/{user_id}/phrase/{phrase_id}/{format}
```

Currently the service supports the following formats:
- m4a
- flac
- mp3
- wav

## Running the tests

```bash
make test
```

## Architecture

The technologies chosen are aimed for simplicity and ease of running on a local machine for the sake of demo purposes. The service is built using Go with layered architecture inspired by Clean Architecture so we can switch the implementation details without changing the business logic. It uses the following technologies:
- Go 1.23
- FFmpeg for audio format conversion
- SQLite for database and message queue for background processing
- Docker (optional, for containerized deployment)
- Local directory for storage (can be replaced with S3 provided you have the credentials)

### Project Structure

```
sb-test/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
|   ├── delivery
│   │   └── http
│   ├── domain/
│   │   ├── converter
│   │   ├── entity
│   │   ├── queue
│   │   ├── repository
│   │   └── storage
│   └── infrastructure/
│       ├── database
│       ├── message_queue
│       └── storage
├── pkg/
│ └── utils/
├── tests/
├── Makefile
├── README.md
├── config.yaml
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── scripts/
├── .gitignore
└── .dockerignore
```

Below is explanation of most important files and directories:

* `cmd/server/` - The main command for running the server
* `internal/delivery/` - The delivery layer (sometimes called the "interface adapters" layer) is responsible for “delivering” requests from external sources (like HTTP) to application's core business logic
* `internal/domain/` - The domain layer containing the core business logic
* `internal/infrastructure/` - The infrastructure layer containing actual implementations of the domain logic like adapters for the database and message queue or external services like AWS S3
* `internal/usecase/` - The usecase directory encapsulates the core business logic of the application. It defines and implements the specific actions or workflows (use cases) that the service provides.
* `internal/worker/` - The worker directory contains the background worker that is responsible for the background processing of the audio files
* `pkg` - The pkg directory contains the utility functions that can be used in many contexts
* `tests/` - Files related to testing the service
* `Makefile` - The Makefile for managing the project, including building and running the project
* `config.yaml` - The configuration file for the server
* `scripts/` - Files related to running the service and simulating common use cases

### Database

The database is a simple SQLite database that is used to store the user, phrase and audio file data. By default the database is created in current working directory. There are only three tables in the database and you can find the schema in `internal/infrastructure/database/schema.sql`.

### Background Processing

The background processing is implemented using a message queue. The message queue is a simple sqlite database that is used to store the messages and managing state for the background processing. You can find the schema in `internal/infrastructure/message_queue/schema.sql`.

### Storage

The project uses S3 or a local directory to store the audio files. Like a typical layered architecture, the storage in the infrastructure layer is responsible for the actual storage of the data that can be configured in the `config.yaml` file.


## Future Improvements

The current implementation is a simple one and no way near production ready. There are many ways to improve it. Here are some of the improvements that can be made:
- Change the message queue to a more robust solution like Kafka or using something like [temporal](https://temporal.io/) for better reliability and visibility of the background processing
- Use a more production ready database like PostgreSQL or MySQL
- When converting the audio file, the service is doing IO operations to the local disk first for conversion using ffmpeg. It should be possible that we stream the audio file to the ffmpeg to the response and avoid the local disk IO. Or even better, we can use a cloud native solution like AWS Elastic Transcoder or GCP Transcoding API for the conversion and convert it from the object storage to object storage.
- Put a modern proxy like [Envoy](https://www.envoyproxy.io/) in front of the service for load balancing, caching, monitoring and logging. We can also utilize the new features and protocol of HTTP/3 supported by the proxy for better latency and reliability. Additionally, we can enable compression for more efficient data transfer with the trade-off of increased CPU usage.
- Use CDN for caching the audio file. We can preconvert the audio file to the most popular formats after the upload and cache them in the CDN.
- Add authentication and authorization
- For long term storage, use a more efficient lossless format like [Opus](https://opus-codec.org/) or [FLAC](https://xiph.org/flac/).

