# Use the official Golang image as the base image for building the app
FROM golang:latest AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies (using Go modules)
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o gambler

# Use a minimal image for running the app
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /app/gambler .

# Copy the .env file to the container
COPY .env .env

# Expose port (if your app listens on a specific port)
EXPOSE 8080

# Command to run the executable with environment variables loaded
CMD ["sh", "-c", "export $(cat .env | xargs) && ./gambler"]
