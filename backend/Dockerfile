FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Define build arguments for target architecture with default value
ARG GOARCH=amd64

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -o backend .

# Use a small alpine image for the final image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/backend .

# Expose the port the server listens on
EXPOSE 8080

# Run the binary
CMD ["./backend"]
