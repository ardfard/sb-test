# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies and certificates
RUN apk add --no-cache gcc musl-dev make

WORKDIR /app

# Copy go mod files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make build

# Final stage
FROM alpine:latest

# Install ffmpeg
RUN apk add --no-cache ffmpeg ca-certificates


# Copy the binary
COPY --from=builder /app/bin/sb-test /app/bin/sb-test

COPY config.yaml /etc/sb-test/config.yaml

ENTRYPOINT ["/app/bin/sb-test"]

CMD ["-c", "/etc/sb-test/config.yaml"]
