# Step 1: Build the Go binary
FROM golang:1.23.2-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy only go.mod and go.sum files to cache dependencies
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Install sqlc
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate database code
WORKDIR /app/internal/database
RUN sqlc generate

# Return to root directory and build the binary
WORKDIR /app
RUN go build -o nestnet .

# Step 2: Create the runtime image
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /app/nestnet /usr/local/bin/nestnet

# Set the entrypoint to the binary
ENTRYPOINT ["/usr/local/bin/nestnet"]
