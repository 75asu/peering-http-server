# Use a minimal base image with Go support
FROM golang:1.17-alpine AS build

# Set the working directory
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod ./

# Download Go dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 go build -o server

# Create a lightweight production image
FROM alpine:3.15

# Install required packages (curl, ping, telnet, bash)
RUN apk add --no-cache curl iputils bind-tools bash

# Copy the built Go binary from the build stage
COPY --from=build /app/server /usr/local/bin/server

# Set a non-root user for running the application
RUN adduser -D appuser
USER appuser

# Set the entrypoint command
ENTRYPOINT ["/usr/local/bin/server"]

# Expose the server's port
EXPOSE 8080
