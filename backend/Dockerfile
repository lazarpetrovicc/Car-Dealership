# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.22.3-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Copy .env file
COPY .env ./

# Build the Go app
RUN go build -o main .

# Test stage
FROM golang:1.22.3-alpine AS tester
WORKDIR /app

# Copy the source, go mod and .env files
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY .env ./

# Run the tests
CMD ["go", "test", "-v", "./tests/..."]

# Run stage
FROM alpine:3.18 AS runner
WORKDIR /app

# Copy the executable and .env from the build stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env ./

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./main"]