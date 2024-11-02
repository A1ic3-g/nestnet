LABEL authors="Alice Gowler, Dom Binks"

# Step 1: Build the Go binary
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application binary
RUN go build -o nestpost .

# Step 2: Create the runtime image
FROM alpine:latest

# Install any needed runtime dependencies (optional)
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/nestpost /usr/local/bin/nestpost

# Set the entrypoint to the binary
ENTRYPOINT ["/usr/local/bin/nestpost"]
