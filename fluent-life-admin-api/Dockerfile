# Stage 1: Build the Go application
# go.mod requires go >= 1.24.6
FROM golang:1.24.6-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
# The output binary will be named 'server' and placed in the current directory (/app)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Stage 2: Create the final lightweight image
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/server .

# Copy the configuration file
COPY configs/ ./configs/

# Expose the port the application runs on
EXPOSE 8082

# Command to run the application
CMD ["./server"]
