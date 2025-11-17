# Step 1: Build the Go application
FROM golang:1.23-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code to the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/main.go

# Step 2: Create a small image with only the built binary
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/main .


# Expose the port the app runs on
EXPOSE 5000

# Command to run the binary
CMD ["./main"]