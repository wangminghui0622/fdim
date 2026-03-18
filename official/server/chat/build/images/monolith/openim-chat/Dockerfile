# Use Go 1.22 Alpine as the base image for building the application
FROM golang:1.22-alpine AS builder

# Define the base directory for the application as an environment variable
ENV SERVER_DIR=/openim-chat

# Set the working directory inside the container based on the environment variable
WORKDIR $SERVER_DIR

# Set the Go proxy to improve dependency resolution speed (optional)
# ENV GOPROXY=https://goproxy.io,direct

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy all files from the current directory into the container
COPY . .

# Install Mage to use for building the application
RUN go install github.com/magefile/mage@v1.15.0

# Build all OpenIM Chat services
RUN mage build || true

# Using Alpine Linux with Go environment for the final image
FROM golang:1.22-alpine

# Install necessary packages
RUN apk add --no-cache bash

# Set the environment and work directory
ENV SERVER_DIR=/openim-chat
WORKDIR $SERVER_DIR

# Copy the compiled binaries from the builder image to the final image
COPY --from=builder $SERVER_DIR/_output $SERVER_DIR/_output
COPY --from=builder $SERVER_DIR/config $SERVER_DIR/config
COPY --from=builder /go/bin/mage /usr/local/bin/mage
COPY --from=builder $SERVER_DIR/magefile_windows.go $SERVER_DIR/
COPY --from=builder $SERVER_DIR/magefile_unix.go $SERVER_DIR/
COPY --from=builder $SERVER_DIR/magefile.go $SERVER_DIR/
COPY --from=builder $SERVER_DIR/go.mod $SERVER_DIR/
COPY --from=builder $SERVER_DIR/go.sum $SERVER_DIR/

# Install all dependencies including gomake for mage start
RUN go mod download
RUN go get github.com/openimsdk/gomake@v0.0.15-alpha.11

# Copy start-config.yml file required by mage start
COPY start-config.yml $SERVER_DIR/start-config.yml

# Set the command to run when the container starts
# Ensure dependencies are available before running mage start
ENTRYPOINT ["sh", "-c", "cd /openim-chat && go mod download && mage start && tail -f /dev/null"]
